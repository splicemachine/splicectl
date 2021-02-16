package cmd

import (
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/splicemachine/splicectl/common"

	"github.com/spf13/cobra"
)

var versionsDatabaseCRCmd = &cobra.Command{
	Use:   "database-cr",
	Short: "Retrieve the versions for a specific database CR in the cluster.",
	Long: `EXAMPLES
	splicectl list database
	splicectl versions database-cr --database-name splicedb
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
		out, err := getDatabaseCRVersions(databaseName)
		if err != nil {
			logrus.WithError(err).Error("Error getting database CR versions")
		}
		crData, cerr := common.RestructureVersions(out)
		if cerr != nil {
			logrus.Fatal("Vault Version JSON conversion failed.")
		}

		if !formatOverridden {
			outputFormat = "text"
		}

		switch strings.ToLower(outputFormat) {
		case "json":
			crData.ToJSON()
		case "gron":
			crData.ToGRON()
		case "yaml":
			crData.ToYAML()
		case "text", "table":
			crData.ToTEXT(noHeaders)
		}

	},
}

func getDatabaseCRVersions(db string) (string, error) {
	restClient := resty.New()

	uri := fmt.Sprintf("splicectl/v1/vault/databasecrversions?database-name=%s", db)
	resp, resperr := restClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Token-Bearer", authClient.GetTokenBearer()).
		SetHeader("X-Token-Session", authClient.GetSessionID()).
		Get(fmt.Sprintf("%s/%s", apiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error getting Database CR Versions")
		return "", resperr
	}

	return string(resp.Body()[:]), nil

}

func init() {
	versionsCmd.AddCommand(versionsDatabaseCRCmd)

	versionsDatabaseCRCmd.Flags().StringP("database-name", "d", "", "Specify the database name")
	// versionsDatabaseCRCmd.Flags().String("output", "json", "Specify the output type")
	// versionsDatabaseCRCmd.MarkFlagRequired("database-name")
}
