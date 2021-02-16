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

var restartDatabaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Restart a specific database in the Splice DB Cluster.",
	Long: `EXAMPLES
	splicectl list database
	splicectl restart database --database-name splicedb
`,
	Run: func(cmd *cobra.Command, args []string) {
		var dberr error
		databaseName, _ := cmd.Flags().GetString("database-name")
		forceRestart, _ := cmd.Flags().GetBool("force")
		if len(databaseName) == 0 {
			databaseName, dberr = promptForDatabaseName()
			if dberr != nil {
				logrus.Fatal("Could not get a list of Databases", dberr)
			}
		}
		out, err := restartDatabase(databaseName, forceRestart)
		if err != nil {
			logrus.WithError(err).Error("Error restarting database")
		}
		var asData objects.ActionStatus
		marshErr := json.Unmarshal([]byte(out), &asData)
		if marshErr != nil {
			logrus.Fatal("Could not unmarshall data", marshErr)
		}

		if !formatOverridden {
			outputFormat = "text"
		}

		switch strings.ToLower(outputFormat) {
		case "json":
			asData.ToJSON()
		case "gron":
			asData.ToGRON()
		case "yaml":
			asData.ToYAML()
		case "text", "table":
			asData.ToTEXT(noHeaders)
		}

	},
}

func restartDatabase(dbname string, force bool) (string, error) {
	restClient := resty.New()

	uri := fmt.Sprintf("splicectl/v1/splicedb/splicedatabaserestart?database-name=%s&force=%t", dbname, force)
	resp, resperr := restClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Token-Bearer", authClient.GetTokenBearer()).
		SetHeader("X-Token-Session", authClient.GetSessionID()).
		Post(fmt.Sprintf("%s/%s", apiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error restarting the database")
		return "", resperr
	}

	return string(resp.Body()[:]), nil

}

func init() {
	restartCmd.AddCommand(restartDatabaseCmd)

	restartDatabaseCmd.Flags().StringP("database-name", "d", "", "Specify the database name")
	restartDatabaseCmd.Flags().BoolP("force", "f", false, "Force the restart")

}
