package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/blang/semver/v4"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/splicemachine/splicectl/cmd/objects"
)

var pauseCmd = &cobra.Command{
	Use:   "pause",
	Args:  cobra.MinimumNArgs(0),
	Short: "Pause Cluster Databases",
	Long: `EXAMPLES
	splicectl list databases
	splicectl pause --database-name <database> --message "<message>"`,
	Run: func(cmd *cobra.Command, args []string) {
		var dberr error
		var sv semver.Version

		_, sv = versionDetail.RequirementMet("pause")

		message, _ := cmd.Flags().GetString("message")
		databaseName, _ := cmd.Flags().GetString("database-name")
		if len(databaseName) == 0 {
			databaseName, dberr = promptForDatabaseName()
			if dberr != nil {
				logrus.Fatal("Could not get a list of Databases", dberr)
			}
		}
		if isDatabaseActive(databaseName) {
			out, err := pauseDatabase(databaseName, message)
			if err != nil {
				logrus.Warn("Pausing database failed.")
			}

			if semverV1, err := semver.ParseRange(">=0.1.7"); err != nil {
				logrus.Fatal("Failed to parse SemVer")
			} else {
				if semverV1(sv) {
					displayPauseDatabaseV1(out)
				}
			}
		} else {
			logrus.Warn("The database is not listed as Active, not paused")
		}

	},
}

func displayPauseDatabaseV1(in string) {
	fmt.Println(in)
	os.Exit(0)
}

func isDatabaseActive(db string) bool {
	dbJSON, err := getDatabaseList()
	if err != nil {
		logrus.WithError(err).Fatal("Error retreiving ClusterId list")
	}
	var dbList objects.DatabaseList

	marshErr := json.Unmarshal([]byte(dbJSON), &dbList)
	if marshErr != nil {
		logrus.Fatal("Could not unmarshall database list for ClusterId", marshErr)
	}

	for _, v := range dbList.Clusters {
		if v.DcosAppId == db {
			if v.Status == "Active" {
				return true
			}
			return false
		}
	}
	return false
}

func pauseDatabase(db string, msg string) (string, error) {
	restClient := resty.New()

	uri := "splicectl/v1/splicedb/splicedatabasepause"

	var resp *resty.Response
	var resperr error

	reqJSON := fmt.Sprintf("{ \"appId\": \"%s\", \"message\": \"%s\" }", db, msg)
	resp, resperr = restClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Token-Bearer", authClient.GetTokenBearer()).
		SetHeader("X-Token-Session", authClient.GetSessionID()).
		SetBody(reqJSON).
		Post(fmt.Sprintf("%s/%s", apiServer, uri))
	if resperr != nil {
		logrus.WithError(resperr).Error("Error getting Default CR Info")
		return "", resperr
	}

	return string(resp.Body()[:]), nil

}

func init() {
	rootCmd.AddCommand(pauseCmd)
	pauseCmd.Flags().StringP("database-name", "d", "", "Specify the Splice Machine Database to Pause")
	pauseCmd.Flags().StringP("message", "m", "", "Add a message to the database log")
}
