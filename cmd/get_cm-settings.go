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

var getCMSettingsCmd = &cobra.Command{
	Use:   "cm-settings",
	Short: "Get the cm (cloud manager) settings for the cluster.",
	Long: `EXAMPLES
	splicectl get cm-settings --component ui -o json > ~/tmp/cm-ui.json
`,
	Run: func(cmd *cobra.Command, args []string) {
		version, _ := cmd.Flags().GetInt("version")
		component, _ := cmd.Flags().GetString("component")

		component = strings.ToLower(component)
		if len(component) == 0 || !strings.Contains("ui api", component) {
			logrus.Fatal("--component needs to be 'ui' or 'api'")
		}
		out, err := getCMSettings(component, version)
		if err != nil {
			logrus.WithError(err).Error("Error getting CM Settings")
		}

		var sessData objects.CMSettings
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
			sessData.ToTEXT(noHeaders)
		}
	},
}

func getCMSettings(comp string, ver int) (string, error) {
	restClient := resty.New()

	uri := fmt.Sprintf("splicectl/v1/vault/cmsettings?component=%s&version=%d", comp, ver)
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
	getCmd.AddCommand(getCMSettingsCmd)

	getCMSettingsCmd.Flags().Int("version", 0, "Specify the version to retrieve, default latest")
	getCMSettingsCmd.Flags().StringP("component", "c", "", "Specify the component, <ui|api>")
}
