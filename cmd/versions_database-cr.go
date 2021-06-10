package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/splicemachine/splicectl/common"

	"github.com/spf13/cobra"
)

var versionsDatabaseCRCmd = &cobra.Command{
	Use:   "database-cr",
	Short: "Retrieve the versions for a specific workspace CR in the cluster.",
	Long: `EXAMPLES
	splicectl list workspace
	splicectl versions workspace-cr --database-name splicedb

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

		_, sv = VersionDetail.RequirementMet("versions_database-cr")

		databaseName := common.DatabaseName(cmd)
		if len(databaseName) == 0 {
			databaseName, dberr = PromptForDatabaseName()
			if dberr != nil {
				logrus.Fatal("Could not get a list of workspaces", dberr)
			}
		}
		out, err := getDatabaseCRVersions(databaseName)
		if err != nil {
			logrus.WithError(err).Error("Error getting workspace CR versions")
		}

		if semverV1, err := semver.ParseRange(">=0.0.15 <0.0.17"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV1(sv) {
				displayVersionsDatabaseCRV1(out)
			}
		}

		if semverV2, err := semver.ParseRange(">=0.0.17"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV2(sv) {
				displayVersionsDatabaseCRV2(out)
			}
		}
	},
}

func displayVersionsDatabaseCRV1(in string) {
	fmt.Println(in)
	os.Exit(0)
}

func displayVersionsDatabaseCRV2(in string) {
	if strings.ToLower(OutputFormat) == "raw" {
		fmt.Println(in)
		os.Exit(0)
	}
	crData, cerr := common.RestructureVersions(in)
	if cerr != nil {
		logrus.Fatal("Vault Version JSON conversion failed.")
	}

	if !FormatOverridden {
		OutputFormat = "text"
	}

	switch strings.ToLower(OutputFormat) {
	case "json":
		crData.ToJSON()
	case "gron":
		crData.ToGRON()
	case "yaml":
		crData.ToYAML()
	case "text", "table":
		crData.ToTEXT(NoHeaders)
	}
}

func getDatabaseCRVersions(db string) (string, error) {
	restClient := resty.New()

	uri := fmt.Sprintf("splicectl/v1/vault/databasecrversions?database-name=%s", db)
	resp, resperr := restClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Token-Bearer", AuthClient.GetTokenBearer()).
		SetHeader("X-Token-Session", AuthClient.GetSessionID()).
		Get(fmt.Sprintf("%s/%s", ApiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error getting workspace CR Versions")
		return "", resperr
	}

	return string(resp.Body()[:]), nil

}

func init() {
	versionsCmd.AddCommand(versionsDatabaseCRCmd)

	// add database name and aliases
	versionsDatabaseCRCmd.Flags().StringP("database-name", "d", "", "Specify the database name")
	versionsDatabaseCRCmd.Flags().String("database", "", "Alias for database-name, prefer the use of -d and --database-name.")
	versionsDatabaseCRCmd.Flags().String("workspace", "", "Alias for database-name, prefer the use of -d and --database-name.")

	// versionsDatabaseCRCmd.Flags().String("output", "json", "Specify the output type")
	// versionsDatabaseCRCmd.MarkFlagRequired("database-name")
}
