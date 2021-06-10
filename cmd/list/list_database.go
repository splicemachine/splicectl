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
	c "github.com/splicemachine/splicectl/cmd"
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

		// databaseName, _ := cmd.Flags().GetString("database-name")
		out, err := c.GetDatabaseList()
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

func init() {
	listCmd.AddCommand(listDatabaseCmd)

	// getDatabaseCRCmd.Flags().String("database-name", "", "Specify the database name")

}