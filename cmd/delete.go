package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/splicemachine/splicectl/cmd/objects"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Args:  cobra.MinimumNArgs(0),
	Short: "Delete Cluster Databases",
	Long: `EXAMPLES
	splicectl list databases
	splicectl delete --database-name <database> --delete

	  * The '--delete' is required as a validation for the deletion request`,
	Run: func(cmd *cobra.Command, args []string) {
		var dberr error
		verifyDelete, _ := cmd.Flags().GetBool("delete")
		databaseName, _ := cmd.Flags().GetString("database-name")
		if len(databaseName) == 0 {
			databaseName, dberr = promptForDatabaseName()
			if dberr != nil {
				logrus.Fatal("Could not get a list of Databases", dberr)
			}
		}
		if verifyDelete {
			clusterID := getMatchingClusterID(databaseName)
			if len(clusterID) > 0 {
				out, err := deleteDatabase(clusterID)
				if err != nil {
					logrus.Warn("Deleting database failed.")
				}
				fmt.Println(out)
			} else {
				logrus.Fatal("Unable to determine ClusterId from Database Name")
			}
		} else {
			logrus.Fatal("You MUST specify --delete on the commandline to validate the deletion")
		}
	},
}

func getMatchingClusterID(db string) string {
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
			return v.ClusterId
		}
	}

	return ""
}

func deleteDatabase(cid string) (string, error) {
	restClient := resty.New()

	uri := fmt.Sprintf("splicectl/v1/splicedb/splicedatabasedelete?database-name=%s", cid)

	var resp *resty.Response
	var resperr error

	resp, resperr = restClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Token-Bearer", authClient.GetTokenBearer()).
		SetHeader("X-Token-Session", authClient.GetSessionID()).
		Delete(fmt.Sprintf("%s/%s", apiServer, uri))
	if resperr != nil {
		logrus.WithError(resperr).Error("Error getting Default CR Info")
		return "", resperr
	}

	return string(resp.Body()[:]), nil

}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringP("database-name", "d", "", "Specify the Splice Machine Database to Pause")
	deleteCmd.Flags().Bool("delete", false, "Verification parameter to perform the deletion")
}
