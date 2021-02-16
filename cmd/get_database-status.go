package cmd

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var getDatabaseStatus = &cobra.Command{
	Use:   "database-status",
	Short: "Get the status of database.",
	Long: `EXAMPLES
	splicectl get database-status --database-name "test"
`,
	Run: func(cmd *cobra.Command, args []string) {
		var dberr error
		databaseName, _ := cmd.Flags().GetString("database-name")
		if len(databaseName) == 0 {
			databaseName, dberr = promptForDatabaseName()
			if dberr != nil {
				logrus.Fatal("Could not get name of Database", dberr)
			}
		}

		out, err := getDatabaseStatusData( databaseName)
		if err != nil {
			logrus.WithError(err).Error("Error getting status of database ")
		}

	    fmt.Println(out)

	},
}

func getDatabaseStatusData(databaseName string) (string, error) {

	restClient := resty.New()

	uri := fmt.Sprintf("splicectl/v1/splicedb/splicedatabasestatus?database-name=%s", databaseName)
	resp, resperr := restClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Token-Bearer", authClient.GetTokenBearer()).
		SetHeader("X-Token-Session", authClient.GetSessionID()).
		Get(fmt.Sprintf("%s/%s", apiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error getting Database status info")
		return "", resperr
	}
	return string(resp.Body()[:]), nil
}

func init() {
	getCmd.AddCommand(getDatabaseStatus)

	getDatabaseStatus.Flags().StringP("database-name", "d", "", "Specify the database name")

	getDatabaseStatus.MarkFlagRequired("database-name")

}
