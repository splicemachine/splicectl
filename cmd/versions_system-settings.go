package cmd

import (
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/splicemachine/splicectl/common"

	"github.com/spf13/cobra"
)

var versionsSystemSettingsCmd = &cobra.Command{
	Use:   "system-settings",
	Short: "Retrieve the versions of the system settings in the cluster.",
	Long: `EXAMPLES
	splicectl versions system-settings
`,
	Run: func(cmd *cobra.Command, args []string) {
		out, err := getSystemSettingsVersions()
		if err != nil {
			logrus.WithError(err).Error("Error getting System Settings")
		}
		ssData, cerr := common.RestructureVersions(out)
		if cerr != nil {
			logrus.Fatal("Vault Version JSON conversion failed.")
		}

		if !formatOverridden {
			outputFormat = "text"
		}

		switch strings.ToLower(outputFormat) {
		case "json":
			ssData.ToJSON()
		case "gron":
			ssData.ToGRON()
		case "yaml":
			ssData.ToYAML()
		case "text", "table":
			ssData.ToTEXT(noHeaders)
		}

	},
}

func getSystemSettingsVersions() (string, error) {
	restClient := resty.New()

	uri := "splicectl/v1/vault/systemsettingsversions"
	resp, resperr := restClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Token-Bearer", authClient.GetTokenBearer()).
		SetHeader("X-Token-Session", authClient.GetSessionID()).
		Get(fmt.Sprintf("%s/%s", apiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error getting System Settings")
		return "", resperr
	}

	return string(resp.Body()[:]), nil

}

func init() {
	versionsCmd.AddCommand(versionsSystemSettingsCmd)

}
