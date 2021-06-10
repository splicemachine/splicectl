package apply

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/sirupsen/logrus"
	c "github.com/splicemachine/splicectl/cmd"
	"github.com/splicemachine/splicectl/cmd/objects"
	"github.com/splicemachine/splicectl/common"

	"github.com/spf13/cobra"
)

var applyDatabaseCRCmd = &cobra.Command{
	Use:   "database-cr",
	Short: "Submit a new CR for a specified database in the cluster.",
	Long: `EXAMPLES
	splicectl list workspace
    splicectl get database-cr --database-name splicedb -o json > ~/tmp/splicedb.json
    # edit the file
    splicectl apply database-cr --database-name splicedb --file ~/tmp/splicedb.json

	Note: --database-name and -d are the preferred way to supply the database name.
	However, --database and --workspace can also be used as well. In the event that
	more than one of them is supplied database-name and d are preferred over all
	and workspace is preferred over database. The most preferred option that is
	supplied will be used and a message will be displayed letting you know which
	option was chosen if more than one were supplied.
`,
	Run: func(cmd *cobra.Command, args []string) {
		var dberr error
		_, sv := c.VersionDetail.RequirementMet("apply_database-cr")

		databaseName := common.DatabaseName(cmd)
		if len(databaseName) == 0 {
			databaseName, dberr = c.PromptForDatabaseName()
			if dberr != nil {
				logrus.Fatal("Could not get a list of Databases", dberr)
			}
		}
		filePath, _ := cmd.Flags().GetString("file")
		fileBytes, _ := ioutil.ReadFile(filePath)

		jsonBytes, cerr := common.WantJSON(fileBytes)
		if cerr != nil {
			logrus.Fatal("The input data MUST be in either JSON or YAML format")
		}

		out, err := setDatabaseCR(databaseName, jsonBytes)
		if err != nil {
			logrus.WithError(err).Error("Error setting Database CR Info")
		}

		if semverV1, err := semver.ParseRange(">=0.0.14 <0.0.17"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV1(sv) {
				displayApplyDatabaseCRV1(out)
			}
		}

		if semverV2, err := semver.ParseRange(">=0.0.17"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV2(sv) {
				displayApplyDatabaseCRV2(out)
			}
		}

	},
}

func displayApplyDatabaseCRV1(in string) {
	fmt.Println(in)
	os.Exit(0)
}

func displayApplyDatabaseCRV2(in string) {

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

func setDatabaseCR(dbname string, in []byte) (string, error) {
	uri := fmt.Sprintf("splicectl/v1/vault/databasecr?database-name=%s", dbname)
	resp, resperr := c.RestyWithHeaders().
		SetBody(in).
		Post(fmt.Sprintf("%s/%s", c.ApiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error setting Database CR Info")
		return "", resperr
	}

	return string(resp.Body()[:]), nil

}

func init() {
	applyCmd.AddCommand(applyDatabaseCRCmd)

	// add database name and aliases
	applyDatabaseCRCmd.Flags().StringP("database-name", "d", "", "Specify the database name")
	applyDatabaseCRCmd.Flags().String("database", "", "Alias for database-name, prefer the use of -d and --database-name.")
	applyDatabaseCRCmd.Flags().String("workspace", "", "Alias for database-name, prefer the use of -d and --database-name.")

	applyDatabaseCRCmd.Flags().StringP("file", "f", "", "Specify the input file")
	// applyDatabaseCRCmd.MarkFlagRequired("database-name")
	applyDatabaseCRCmd.MarkFlagRequired("file")
}
