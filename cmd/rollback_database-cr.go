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

var rollbackDatabaseCRCmd = &cobra.Command{
	Use:   "database-cr",
	Short: "Rollback the Database CR to a specific vault version",
	Long: `EXAMPLES
	splicectl list database
	splicectl versions database-cr --database-name splicedb
	splicectl rollback database-cr --database-name splicedb --version 2
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
		version, _ := cmd.Flags().GetInt("version")
		out, err := rollbackDatabaseCR(databaseName, version)
		if err != nil {
			logrus.WithError(err).Error("Error getting Database CR Info")
		}
		var vvData objects.VaultVersion
		marshErr := json.Unmarshal([]byte(out), &vvData)
		if marshErr != nil {
			logrus.Fatal("Could not unmarshall data", marshErr)
		}

		if !formatOverridden {
			outputFormat = "text"
		}

		switch strings.ToLower(outputFormat) {
		case "json":
			vvData.ToJSON()
		case "gron":
			vvData.ToGRON()
		case "yaml":
			vvData.ToYAML()
		case "text", "table":
			vvData.ToTEXT(noHeaders)
		}

	},
}

func rollbackDatabaseCR(dbname string, ver int) (string, error) {
	restClient := resty.New()

	uri := fmt.Sprintf("splicectl/v1/vault/rollbackdatabasecr?version=%d&database-name=%s", ver, dbname)
	resp, resperr := restClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Token-Bearer", authClient.GetTokenBearer()).
		SetHeader("X-Token-Session", authClient.GetSessionID()).
		Post(fmt.Sprintf("%s/%s", apiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error getting Database CR Info")
		return "", resperr
	}

	return string(resp.Body()[:]), nil

}

func init() {
	rollbackCmd.AddCommand(rollbackDatabaseCRCmd)

	rollbackDatabaseCRCmd.Flags().StringP("database-name", "d", "", "Specify the database name")
	// rollbackDatabaseCRCmd.Flags().String("output", "json", "Specify the output type")
	rollbackDatabaseCRCmd.Flags().IntP("version", "v", 0, "Specify the version to retrieve, default latest")
	// rollbackDatabaseCRCmd.MarkFlagRequired("database-name")
	rollbackDatabaseCRCmd.MarkFlagRequired("version")

}
