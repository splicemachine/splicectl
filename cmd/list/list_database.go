package list

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/sirupsen/logrus"
	"github.com/splicemachine/splicectl/cmd/objects"

	"github.com/spf13/cobra"
)

var listDatabaseCmd = &cobra.Command{
	Use:     "workspace",
	Aliases: []string{"workspaces", "database", "databases"},
	Short:   "Retrieve a list of splice databases in the cluster.",
	Long: `EXAMPLES
	splicectl list workspace
`,
	Run: func(cmd *cobra.Command, args []string) {
		_, sv := c.VersionDetail.RequirementMet("list_database")

		// check active and paused flag values
		active, err := cmd.Flags().GetBool("active")
		if err != nil {
			logrus.WithError(err).Error("Error getting Database CR Info")
		}

		paused, err := cmd.Flags().GetBool("paused")
		if err != nil {
			logrus.WithError(err).Error("Error getting Database CR Info")
		}

		// databaseName, _ := cmd.Flags().GetString("database-name")
		out, err := getDatabaseListWithFlags(active, paused)
		if err != nil {
			logrus.WithError(err).Error("Error getting Database CR Info")
		}

		if semverV1, err := semver.ParseRange(">=0.0.14 <0.0.17"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV1(sv) {
				displayListDatabaseV1(out)
			}
		}

		if semverV2, err := semver.ParseRange(">=0.0.17"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV2(sv) {
				displayListDatabaseV2(out)
			}
		}
	},
}

func displayListDatabaseV1(in string) {
	fmt.Println(in)
	os.Exit(0)
}

func displayListDatabaseV2(in string) {
	if strings.ToLower(c.OutputFormat) == "raw" {
		fmt.Println(in)
		os.Exit(0)
	}
	var dbList objects.DatabaseList

	marshErr := json.Unmarshal([]byte(in), &dbList)
	if marshErr != nil {
		logrus.Fatal("Could not unmarshall data", marshErr)
	}

	c.OutputData(&dbList)

}

// getDatabaseList - simple wrapper around getDatabaseListWithFlags to prevent
// cascading issues with change in API.
func getDatabaseList() (string, error) {
	return getDatabaseListWithFlags(false, false)
}

// getDatabaseListWithFlags - gets list of databases and filters/orders them
// based on flags.
func getDatabaseListWithFlags(active, paused bool) (string, error) {
	data, err := c.GetDatabaseList()
	if err != nil {
		return "", err
	}

	// filter and order DBList returned from api call
	rawRespBody, dbList := data, &objects.DatabaseList{}
	if err := json.Unmarshal(rawRespBody, dbList); err != nil {
		return "", err
	}
	dbList = dbList.FilterByStatus(active, paused)
	filteredDBList, err := json.Marshal(dbList)
	// TODO: return either bytes or json, do not go back to strings
	return string(filteredDBList), err

}

func init() {
	listCmd.AddCommand(listDatabaseCmd)

	listDatabaseCmd.Flags().BoolP("active", "a", false, "Select if you want to get active databases.")
	listDatabaseCmd.Flags().BoolP("paused", "p", false, "Select if you want to get paused databases.")
}
