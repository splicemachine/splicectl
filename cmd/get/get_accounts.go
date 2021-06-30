package get

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/splicemachine/splicectl/cmd/objects"
)

var getAccountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "Get a list of Cloud Manager Accounts",
	Long: `EXAMPLES
	splicectl get accounts

	    * if no accounts are listed, you will need to logon to the Ops Center
`,
	Run: func(cmd *cobra.Command, args []string) {
		_, sv := c.VersionDetail.RequirementMet("get_accounts")

		out, err := c.GetAccounts()
		if err != nil {
			logrus.WithError(err).Error("Error getting Default CR Info")
		}

		if semverV1, err := semver.ParseRange(">=0.1.7"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV1(sv) {
				displayGetAccountsV1(out)
			}
		}
	},
}

func displayGetAccountsV1(in string) {
	if strings.ToLower(c.OutputFormat) == "raw" {
		fmt.Println(in)
		os.Exit(0)
	}

	var accounts objects.AccountList

	marshErr := json.Unmarshal([]byte(in), &accounts)
	if marshErr != nil {
		logrus.Fatal("Could not unmarshall data", marshErr)
	}

	c.OutputData(&accounts)
}

func init() {
	getCmd.AddCommand(getAccountsCmd)
}
