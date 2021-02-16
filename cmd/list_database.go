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

var listDatabaseCmd = &cobra.Command{
	Use:     "database",
	Aliases: []string{"databases"},
	Short:   "Retrieve a list of splice databases in the cluster.",
	Long: `EXAMPLES
	splicectl list database
`,
	Run: func(cmd *cobra.Command, args []string) {
		// databaseName, _ := cmd.Flags().GetString("database-name")
		out, err := getDatabaseList()
		if err != nil {
			logrus.WithError(err).Error("Error getting Database CR Info")
		}
		var dbList objects.DatabaseList

		marshErr := json.Unmarshal([]byte(out), &dbList)
		if marshErr != nil {
			logrus.Fatal("Could not unmarshall data", marshErr)
		}

		if !formatOverridden {
			outputFormat = "table"
		}

		switch strings.ToLower(outputFormat) {
		case "json":
			dbList.ToJSON()
		case "gron":
			dbList.ToGRON()
		case "yaml":
			dbList.ToYAML()
		case "text", "table":
			dbList.ToTEXT(noHeaders)
		}
	},
}

func getDatabaseList() (string, error) {
	restClient := resty.New()

	uri := "splicectl/v1/splicedb/splicedatabase"
	resp, resperr := restClient.R().
		// SetHeader("Content-Type", "application/json").
		// SetHeader("Accept", "application/json").
		SetHeader("X-Token-Bearer", authClient.GetTokenBearer()).
		SetHeader("X-Token-Session", authClient.GetSessionID()).
		Execute("LIST", fmt.Sprintf("%s/%s", apiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error getting Database List")
		return "", resperr
	}

	return string(resp.Body()[:]), nil

}

func init() {
	listCmd.AddCommand(listDatabaseCmd)

	// getDatabaseCRCmd.Flags().String("database-name", "", "Specify the database name")

}
