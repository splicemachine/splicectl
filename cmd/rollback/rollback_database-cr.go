package rollback

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/sirupsen/logrus"
	"github.com/splicemachine/splicectl/cmd/objects"
	"github.com/splicemachine/splicectl/common"

	"github.com/spf13/cobra"
)

var rollbackDatabaseCRCmd = &cobra.Command{
	Use:   "database-cr",
	Short: "Rollback the workspace CR to a specific vault version",
	Long: `EXAMPLES
	splicectl list workspace
	splicectl versions workspace-cr --database-name splicedb
	splicectl rollback workspace-cr --database-name splicedb --version 2

	Note: --database-name and -d are the preferred way to supply the database name.
	However, --database and --workspace can also be used as well. In the event that
	more than one of them is supplied database-name and d are preferred over all
	and workspace is preferred over database. The most preferred option that is
	supplied will be used and a message will be displayed letting you know which
	option was chosen if more than one were supplied.
`,
	Run: func(cmd *cobra.Command, args []string) {
		var dberr error
		_, sv := c.VersionDetail.RequirementMet("rollback_database-cr")

		databaseName := common.DatabaseName(cmd)
		if len(databaseName) == 0 {
			databaseName, dberr = c.PromptForDatabaseName()
			if dberr != nil {
				logrus.Fatal("Could not get a list of workspaces", dberr)
			}
		}
		version, _ := cmd.Flags().GetInt("version")
		out, err := rollbackDatabaseCR(databaseName, version)
		if err != nil {
			logrus.WithError(err).Error("Error getting workspace CR Info")
		}

		if semverV1, err := semver.ParseRange(">=0.0.15 <0.0.17"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV1(sv) {
				displayRollbackDatabaseCRV1(out)
			}
		}

		if semverV2, err := semver.ParseRange(">=0.0.17"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV2(sv) {
				displayRollbackDatabaseCRV2(out)
			}
		}
	},
}

func displayRollbackDatabaseCRV1(in string) {
	fmt.Println(in)
	os.Exit(0)
}

func displayRollbackDatabaseCRV2(in string) {
	if strings.ToLower(c.OutputFormat) == "raw" {
		fmt.Println(in)
		os.Exit(0)
	}
	var vvData objects.VaultVersion
	marshErr := json.Unmarshal([]byte(in), &vvData)
	if marshErr != nil {
		logrus.Fatal("Could not unmarshall data", marshErr)
	}

	c.OutputData(&vvData)
}

func rollbackDatabaseCR(dbname string, ver int) (string, error) {
	uri := fmt.Sprintf("splicectl/v1/vault/rollbackdatabasecr?version=%d&database-name=%s", ver, dbname)
	resp, resperr := c.RestyWithHeaders().
		Post(fmt.Sprintf("%s/%s", c.ApiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error getting workspace CR Info")
		return "", resperr
	}

	return string(resp.Body()[:]), nil

}

func init() {
	rollbackCmd.AddCommand(rollbackDatabaseCRCmd)

	// add database name and aliases
	rollbackDatabaseCRCmd.Flags().StringP("database-name", "d", "", "Specify the database name")
	rollbackDatabaseCRCmd.Flags().String("database", "", "Alias for database-name, prefer the use of -d and --database-name.")
	rollbackDatabaseCRCmd.Flags().String("workspace", "", "Alias for database-name, prefer the use of -d and --database-name.")

	// rollbackDatabaseCRCmd.Flags().String("output", "json", "Specify the output type")
	rollbackDatabaseCRCmd.Flags().IntP("version", "v", 0, "Specify the version to retrieve, default latest")
	// rollbackDatabaseCRCmd.MarkFlagRequired("database-name")
	rollbackDatabaseCRCmd.MarkFlagRequired("version")

}
