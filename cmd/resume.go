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
	"github.com/splicemachine/splicectl/common"
)

var resumeCmd = &cobra.Command{
	Use:   "resume",
	Args:  cobra.MinimumNArgs(0),
	Short: "Resume Cluster workspaces",
	Long: `EXAMPLES
	splicectl list workspace
	splicectl resume --database-name <database> --message "<message>"
	
	Note: --database-name and -d are the preferred way to supply the database name.
	However, --database and --workspace can also be used as well. In the event that
	more than one of them is supplied database-name and d are preferred over all
	and workspace is preferred over database. The most preferred option that is
	supplied will be used and a message will be displayed letting you know which
	option was chosen if more than one were supplied.`,
	Run: func(cmd *cobra.Command, args []string) {
		var dberr error
		var sv semver.Version

		_, sv = VersionDetail.RequirementMet("resume")

		message, _ := cmd.Flags().GetString("message")
		databaseName := common.DatabaseName(cmd)
		if len(databaseName) == 0 {
			databaseName, dberr = PromptForDatabaseName()
			if dberr != nil {
				logrus.Fatal("Could not get a list of workspaces", dberr)
			}
		}
		if isDatabasePaused(databaseName) {
			out, err := resumeDatabase(databaseName, message)
			if err != nil {
				logrus.Warn("Resuming workspace failed.")
			}

			if semverV1, err := semver.ParseRange(">=0.1.7"); err != nil {
				logrus.Fatal("Failed to parse SemVer")
			} else {
				if semverV1(sv) {
					displayResumeDatabaseV1(out)
				}
			}
		} else {
			logrus.Warn("The workspace is not listed as Paused, not resuming")
		}

	},
}

func displayResumeDatabaseV1(in string) {
	fmt.Println(in)
	os.Exit(0)
}

func isDatabasePaused(db string) bool {
	dbJSON, err := getDatabaseList()
	if err != nil {
		logrus.WithError(err).Fatal("Error retreiving ClusterId list")
	}
	var dbList objects.DatabaseList

	marshErr := json.Unmarshal([]byte(dbJSON), &dbList)
	if marshErr != nil {
		logrus.Fatal("Could not unmarshall workspace list for ClusterId", marshErr)
	}

	for _, v := range dbList.Clusters {
		if v.DcosAppId == db {
			if v.Status == "Paused" {
				return true
			}
			return false
		}
	}
	return false
}

func resumeDatabase(db string, msg string) (string, error) {
	restClient := resty.New()

	uri := "splicectl/v1/splicedb/splicedatabaseresume"

	var resp *resty.Response
	var resperr error

	reqJSON := fmt.Sprintf("{ \"appId\": \"%s\", \"message\": \"%s\" }", db, msg)
	resp, resperr = restClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Token-Bearer", AuthClient.GetTokenBearer()).
		SetHeader("X-Token-Session", AuthClient.GetSessionID()).
		SetBody(reqJSON).
		Post(fmt.Sprintf("%s/%s", ApiServer, uri))
	if resperr != nil {
		logrus.WithError(resperr).Error("Error getting Default CR Info")
		return "", resperr
	}

	return string(resp.Body()[:]), nil

}

func init() {
	RootCmd.AddCommand(resumeCmd)

	// add database name and aliases
	resumeCmd.Flags().StringP("database-name", "d", "", "Specify the database name")
	resumeCmd.Flags().String("database", "", "Alias for database-name, prefer the use of -d and --database-name.")
	resumeCmd.Flags().String("workspace", "", "Alias for database-name, prefer the use of -d and --database-name.")

	resumeCmd.Flags().StringP("message", "m", "", "Add a message to the workspace log")
}
