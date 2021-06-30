package get

import (
	"fmt"
	"os"

	"github.com/blang/semver/v4"
	"github.com/sirupsen/logrus"
	"github.com/splicemachine/splicectl/common"

	"github.com/spf13/cobra"
)

var getDatabaseStatus = &cobra.Command{
	Use:   "database-status",
	Short: "Get the status of database.",
	Long: `EXAMPLES
	splicectl get database-status --database-name "test"

	Note: --database-name and -d are the preferred way to supply the database name.
	However, --database and --workspace can also be used as well. In the event that
	more than one of them is supplied database-name and d are preferred over all
	and workspace is preferred over database. The most preferred option that is
	supplied will be used and a message will be displayed letting you know which
	option was chosen if more than one were supplied.
`,
	Run: func(cmd *cobra.Command, args []string) {
		var dberr error
		_, sv := c.VersionDetail.RequirementMet("get_database-status")

		databaseName := common.DatabaseName(cmd)
		if len(databaseName) == 0 {
			databaseName, dberr = c.PromptForDatabaseName()
			if dberr != nil {
				logrus.Fatal("Could not get name of Database", dberr)
			}
		}

		out, err := getDatabaseStatusData(databaseName)
		if err != nil {
			logrus.WithError(err).Error("Error getting status of database ")
		}

		if semverV1, err := semver.ParseRange(">=0.1.6"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV1(sv) {
				displayGetDatabaseStatusV1(out)
			}
		}
	},
}

func displayGetDatabaseStatusV1(in string) {
	fmt.Println(in)
	os.Exit(0)
}

func getDatabaseStatusData(databaseName string) (string, error) {
	uri := fmt.Sprintf("splicectl/v1/splicedb/splicedatabasestatus?database-name=%s", databaseName)
	resp, resperr := c.RestyWithHeaders().
		Get(fmt.Sprintf("%s/%s", c.ApiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error getting Database status info")
		return "", resperr
	}
	return string(resp.Body()[:]), nil
}

func init() {
	getCmd.AddCommand(getDatabaseStatus)

	// add database name and aliases
	getDatabaseStatus.Flags().StringP("database-name", "d", "", "Specify the database name")
	getDatabaseStatus.Flags().String("database", "", "Alias for database-name, prefer the use of -d and --database-name.")
	getDatabaseStatus.Flags().String("workspace", "", "Alias for database-name, prefer the use of -d and --database-name.")
}
