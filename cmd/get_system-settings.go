package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/splicemachine/splicectl/cmd/objects"
)

var getSystemSettingsCmd = &cobra.Command{
	Use:   "system-settings",
	Short: "Get the default system settings for the cluster.",
	Long: `EXAMPLES
	splicectl get system-settings -o json > ~/tmp/system-settings.json
`,
	Run: func(cmd *cobra.Command, args []string) {
		version, _ := cmd.Flags().GetInt("version")
		decode, _ := cmd.Flags().GetBool("decode-values")

		out, err := getSystemSettings(version)
		if err != nil {
			logrus.WithError(err).Error("Error getting System Settings")
		}

		var sessData objects.SystemSettings
		marshErr := json.Unmarshal([]byte(out), &sessData)
		if marshErr != nil {
			logrus.Fatal("Could not unmarshall data", marshErr)
		}

		if !formatOverridden {
			outputFormat = "yaml"
		}

		switch strings.ToLower(outputFormat) {
		case "json":
			sessData.ToJSON()
		case "gron":
			sessData.ToGRON()
		case "yaml":
			sessData.ToYAML()
		case "text", "table":
			sessData.ToTEXT(noHeaders, decode)
		}
	},
}

func getSystemSettings(ver int) (string, error) {
	restClient := resty.New()

	uri := fmt.Sprintf("splicectl/v1/vault/systemsettings?version=%d", ver)
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
	getCmd.AddCommand(getSystemSettingsCmd)

	getSystemSettingsCmd.Flags().Int("version", 0, "Specify the version to retrieve, default latest")
	getSystemSettingsCmd.Flags().BoolP("decode-values", "d", false, "Decode Base64 Encoded Values")
}
