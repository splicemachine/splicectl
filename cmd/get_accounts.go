package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

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
		out, err := getAccounts()
		if err != nil {
			logrus.WithError(err).Error("Error getting Default CR Info")
		}

		var accounts objects.AccountList

		marshErr := json.Unmarshal([]byte(out), &accounts)
		if marshErr != nil {
			logrus.Fatal("Could not unmarshall data", marshErr)
		}

		if !formatOverridden {
			outputFormat = "text"
		}

		switch strings.ToLower(outputFormat) {

		case "json":
			accounts.ToJSON()
		case "gron":
			accounts.ToGRON()
		case "yaml":
			accounts.ToYAML()
		case "text", "table":
			accounts.ToTEXT(noHeaders)
		}
	},
}

func getAccounts() (string, error) {
	restClient := resty.New()

	uri := "splicectl/v1/cm/accounts"
	resp, resperr := restClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Token-Bearer", authClient.GetTokenBearer()).
		SetHeader("X-Token-Session", authClient.GetSessionID()).
		Get(fmt.Sprintf("%s/%s", apiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error getting Account List Info")
		return "", resperr
	}

	return string(resp.Body()[:]), nil

}

func init() {
	getCmd.AddCommand(getAccountsCmd)
}
