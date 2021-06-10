package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/blang/semver/v4"
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

		var sv semver.Version

		_, sv = VersionDetail.RequirementMet("get_cm-settings")

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

		if semverV1, err := semver.ParseRange(">=0.1.6"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV1(sv) {
				displayGetCmSettingsV1(out)
			}
		}
	},
}

func displayGetCmSettingsV1(in string) {
	if strings.ToLower(OutputFormat) == "raw" {
		fmt.Println(in)
		os.Exit(0)
	}

	var sessData objects.CMSettings
	marshErr := json.Unmarshal([]byte(in), &sessData)
	if marshErr != nil {
		logrus.Fatal("Could not unmarshall data", marshErr)
	}

	if !FormatOverridden {
		OutputFormat = "yaml"
	}

	switch strings.ToLower(OutputFormat) {
	case "json":
		sessData.ToJSON()
	case "gron":
		sessData.ToGRON()
	case "yaml":
		sessData.ToYAML()
	case "text", "table":
		sessData.ToTEXT(NoHeaders)
	}

}

func getCMSettings(comp string, ver int) (string, error) {
	restClient := resty.New()

	uri := fmt.Sprintf("splicectl/v1/vault/cmsettings?component=%s&version=%d", comp, ver)
	resp, resperr := restClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Token-Bearer", AuthClient.GetTokenBearer()).
		SetHeader("X-Token-Session", AuthClient.GetSessionID()).
		Get(fmt.Sprintf("%s/%s", ApiServer, uri))

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
