package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/splicemachine/splicectl/cmd/objects"

	"github.com/spf13/cobra"
)

var getAccountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "Get a list of Cloud Manager Accounts",
	Long: `EXAMPLES
	splicectl get accounts

	    * if no accounts are listed, you will need to logon to the Ops Center
`,
	Run: func(cmd *cobra.Command, args []string) {

		var sv semver.Version

		_, sv = VersionDetail.RequirementMet("get_accounts")

		out, err := getAccounts()
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
	if strings.ToLower(OutputFormat) == "raw" {
		fmt.Println(in)
		os.Exit(0)
	}

	var accounts objects.AccountList

	marshErr := json.Unmarshal([]byte(in), &accounts)
	if marshErr != nil {
		logrus.Fatal("Could not unmarshall data", marshErr)
	}

	if !FormatOverridden {
		OutputFormat = "text"
	}

	switch strings.ToLower(OutputFormat) {

	case "json":
		accounts.ToJSON()
	case "gron":
		accounts.ToGRON()
	case "yaml":
		accounts.ToYAML()
	case "text", "table":
		accounts.ToTEXT(NoHeaders)
	}

}

func getAccounts() (string, error) {
	restClient := resty.New()

	uri := "splicectl/v1/cm/accounts"
	resp, resperr := restClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Token-Bearer", AuthClient.GetTokenBearer()).
		SetHeader("X-Token-Session", AuthClient.GetSessionID()).
		Get(fmt.Sprintf("%s/%s", ApiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error getting Account List Info")
		return "", resperr
	}

	return string(resp.Body()[:]), nil

}

func init() {
	getCmd.AddCommand(getAccountsCmd)
}
