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

var rollbackSystemSettingsCmd = &cobra.Command{
	Use:   "system-settings",
	Short: "Rollback the system settings to a specific vault version",
	Long: `EXAMPLES
	splicectl versions system-settings
	splicectl rollback system-settings --version 2
`,
	Run: func(cmd *cobra.Command, args []string) {
		version, _ := cmd.Flags().GetInt("version")
		out, err := rollbackSystemSettings(version)
		if err != nil {
			logrus.WithError(err).Error("Error rolling back System Settings")
		}
		var vvData objects.VaultVersion
		marshErr := json.Unmarshal([]byte(out), &vvData)
		if marshErr != nil {
			logrus.Fatal("Could not unmarshall data", marshErr)
		}

		if !formatOverridden {
			outputFormat = "text"
		}

		switch strings.ToLower(outputFormat) {
		case "json":
			vvData.ToJSON()
		case "gron":
			vvData.ToGRON()
		case "yaml":
			vvData.ToYAML()
		case "text", "table":
			vvData.ToTEXT(noHeaders)
		}

	},
}

func rollbackSystemSettings(ver int) (string, error) {
	restClient := resty.New()

	uri := fmt.Sprintf("splicectl/v1/vault/rollbacksystemsettings?version=%d", ver)
	resp, resperr := restClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Token-Bearer", authClient.GetTokenBearer()).
		SetHeader("X-Token-Session", authClient.GetSessionID()).
		Post(fmt.Sprintf("%s/%s", apiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error rolling back System Settings")
		return "", resperr
	}

	return string(resp.Body()[:]), nil

}

func init() {
	rollbackCmd.AddCommand(rollbackSystemSettingsCmd)

	rollbackSystemSettingsCmd.Flags().Int("version", 0, "Specify the version to retrieve, default latest")
	rollbackSystemSettingsCmd.MarkFlagRequired("version")
}
