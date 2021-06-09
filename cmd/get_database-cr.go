package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/splicemachine/splicectl/cmd/objects"
	"github.com/splicemachine/splicectl/common"

	"github.com/spf13/cobra"
)

var getDatabaseCRCmd = &cobra.Command{
	Use:   "database-cr",
	Short: "Get the CR for a specific workspace in the cluster.",
	Long: `EXAMPLES
	splicectl list workspace
	splicectl get database-cr --database-name splicedb -o json > ~/tmp/splicedb-cr.json

	Note: --database-name and -d are the preferred way to supply the database name.
	However, --database and --workspace can also be used as well. In the event that
	more than one of them is supplied database-name and d are preferred over all
	and workspace is preferred over database. The most preferred option that is
	supplied will be used and a message will be displayed letting you know which
	option was chosen if more than one were supplied.
`,
	Run: func(cmd *cobra.Command, args []string) {
		var dberr error
		var sv semver.Version

		_, sv = versionDetail.RequirementMet("get_database-cr")

		databaseName := common.DatabaseName(cmd)
		if len(databaseName) == 0 {
			databaseName, dberr = promptForDatabaseName()
			if dberr != nil {
				logrus.Fatal("Could not get a list of Databases", dberr)
			}
		}
		filePath, _ := cmd.Flags().GetString("file")
		version, _ := cmd.Flags().GetInt("version")

		out, err := getDatabaseCR(databaseName, version)
		if err != nil {
			logrus.WithError(err).Error("Error getting workspace CR Info")
		}

		if semverV1, err := semver.ParseRange(">=0.0.14 <0.0.17"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV1(sv) {
				displayGetDatabaseV1(out, filePath)
			}
		}

		if semverV2, err := semver.ParseRange(">=0.0.17"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV2(sv) {
				displayGetDatabaseV2(out, filePath)
			}
		}
	},
}

func displayGetDatabaseV1(in string, fp string) {
	if len(fp) == 0 {
		fmt.Println(in)
	} else {
		objects.WriteToFile(fp, in)
	}
	os.Exit(0)
}
func displayGetDatabaseV2(in string, fp string) {
	if strings.ToLower(outputFormat) == "raw" {
		fmt.Println(in)
		os.Exit(0)
	}
	var dbCR objects.DatabaseCR
	marshErr := json.Unmarshal([]byte(in), &dbCR)
	if marshErr != nil {
		logrus.Fatal("Could not unmarshall data", marshErr)
	}

	if !formatOverridden {
		outputFormat = "yaml"
	}

	switch strings.ToLower(outputFormat) {
	case "json":
		dbCR.ToJSON(fp)
	case "gron":
		dbCR.ToGRON(fp)
	case "yaml":
		dbCR.ToYAML(fp)
	case "text", "table":
		dbCR.ToTEXT(noHeaders)
	}

}

func getDatabaseCR(dbname string, ver int) (string, error) {
	restClient := resty.New()

	uri := fmt.Sprintf("splicectl/v1/vault/databasecr?version=%d&database-name=%s", ver, dbname)
	resp, resperr := restClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Token-Bearer", authClient.GetTokenBearer()).
		SetHeader("X-Token-Session", authClient.GetSessionID()).
		Get(fmt.Sprintf("%s/%s", apiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error getting workspace CR Info")
		return "", resperr
	}

	return string(resp.Body()[:]), nil

}

func init() {
	getCmd.AddCommand(getDatabaseCRCmd)

	// add database name and aliases
	getDatabaseCRCmd.Flags().StringP("database-name", "d", "", "Specify the database name")
	getDatabaseCRCmd.Flags().String("database", "", "Alias for database-name, prefer the use of -d and --database-name.")
	getDatabaseCRCmd.Flags().String("workspace", "", "Alias for database-name, prefer the use of -d and --database-name.")

	getDatabaseCRCmd.Flags().IntP("version", "v", 0, "Specify the version to retrieve, default latest")
	getDatabaseCRCmd.Flags().StringP("file", "f", "", "Specify an output file")
	// getDatabaseCRCmd.MarkFlagRequired("database-name")

}
