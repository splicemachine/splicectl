package del

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/blang/semver/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/splicemachine/splicectl/cmd/config"
	"github.com/splicemachine/splicectl/cmd/objects"
	"github.com/splicemachine/splicectl/common"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Args:  cobra.MinimumNArgs(0),
	Short: "Delete Cluster workspaces",
	Long: `EXAMPLES
	splicectl list workspace
	splicectl delete --database-name <database> --delete

	* The '--delete' is required as a validation for the deletion request
	  
	Note: --database-name and -d are the preferred way to supply the database name.
	However, --database and --workspace can also be used as well. In the event that
	more than one of them is supplied database-name and d are preferred over all
	and workspace is preferred over database. The most preferred option that is
	supplied will be used and a message will be displayed letting you know which
	option was chosen if more than one were supplied. `,
	Run: func(cmd *cobra.Command, args []string) {
		var dberr error

		_, sv := c.VersionDetail.RequirementMet("delete")

		verifyDelete, _ := cmd.Flags().GetBool("delete")
		databaseName := common.DatabaseName(cmd)
		if len(databaseName) == 0 {
			databaseName, dberr = c.PromptForDatabaseName()
			if dberr != nil {
				logrus.Fatal("Could not get a list of workspaces", dberr)
			}
		}
		if verifyDelete {
			clusterID := getMatchingClusterID(databaseName)
			if len(clusterID) > 0 {
				out, err := deleteDatabase(clusterID)
				if err != nil {
					logrus.Warn("Deleting workspace failed.")
				}
				if semverV1, err := semver.ParseRange(">=0.1.7"); err != nil {
					logrus.Fatal("Failed to parse SemVer")
				} else {
					if semverV1(sv) {
						displayDeleteV1(out)
					}
				}
			} else {
				logrus.Fatal("Unable to determine ClusterId from workspace Name")
			}
		} else {
			logrus.Fatal("You MUST specify --delete on the commandline to validate the deletion")
		}
	},
}

func displayDeleteV1(in string) {
	fmt.Println(in)
	os.Exit(0)
}

func getMatchingClusterID(db string) string {
	dbJSON, err := c.GetDatabaseList()
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
			return v.ClusterId
		}
	}

	return ""
}

func deleteDatabase(cid string) (string, error) {
	uri := fmt.Sprintf("splicectl/v1/splicedb/splicedatabasedelete?database-name=%s", cid)
	resp, resperr := c.RestyWithHeaders().
		Delete(fmt.Sprintf("%s/%s", c.ApiServer, uri))
	if resperr != nil {
		logrus.WithError(resperr).Error("Error getting Default CR Info")
		return "", resperr
	}

	return string(resp.Body()[:]), nil

}

var c *config.Config

func InitSubCommands(conf *config.Config) *cobra.Command {
	c = conf
	return deleteCmd
}

func init() {
	// add database name and aliases
	deleteCmd.Flags().StringP("database-name", "d", "", "Specify the database name")
	deleteCmd.Flags().String("database", "", "Alias for database-name, prefer the use of -d and --database-name.")
	deleteCmd.Flags().String("workspace", "", "Alias for database-name, prefer the use of -d and --database-name.")

	deleteCmd.Flags().Bool("delete", false, "Verification parameter to perform the deletion")
}
