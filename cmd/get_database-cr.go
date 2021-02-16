package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/splicemachine/splicectl/cmd/objects"

	"github.com/spf13/cobra"
)

var getDatabaseCRCmd = &cobra.Command{
	Use:   "database-cr",
	Short: "Get the CR for a specific database in the cluster.",
	Long: `EXAMPLES
	splicectl list database
	splicectl get database-cr --database-name splicedb -o json > ~/tmp/splicedb-cr.json
`,
	Run: func(cmd *cobra.Command, args []string) {
		var dberr error
		databaseName, _ := cmd.Flags().GetString("database-name")
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
			logrus.WithError(err).Error("Error getting Database CR Info")
		}

		var dbCR objects.DatabaseCR
		marshErr := json.Unmarshal([]byte(out), &dbCR)
		if marshErr != nil {
			logrus.Fatal("Could not unmarshall data", marshErr)
		}

		if !formatOverridden {
			outputFormat = "yaml"
		}

		switch strings.ToLower(outputFormat) {
		case "json":
			dbCR.ToJSON(filePath)
		case "gron":
			dbCR.ToGRON(filePath)
		case "yaml":
			dbCR.ToYAML(filePath)
		case "text", "table":
			dbCR.ToTEXT(noHeaders)
		}
	},
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
		logrus.WithError(resperr).Error("Error getting Database CR Info")
		return "", resperr
	}

	return string(resp.Body()[:]), nil

}

func init() {
	getCmd.AddCommand(getDatabaseCRCmd)

	getDatabaseCRCmd.Flags().StringP("database-name", "d", "", "Specify the database name")
	getDatabaseCRCmd.Flags().IntP("version", "v", 0, "Specify the version to retrieve, default latest")
	getDatabaseCRCmd.Flags().StringP("file", "f", "", "Specify an output file")
	// getDatabaseCRCmd.MarkFlagRequired("database-name")

}
