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

var rollbackCMSettingsCmd = &cobra.Command{
	Use:   "cm-settings",
	Short: "Rollback the cm (cloud manager) settings to a specific vault version",
	Long: `EXAMPLES
	splicectl versions cm-settings --component ui
	splicectl rollback cm-settings --component ui --version 2
`,
	Run: func(cmd *cobra.Command, args []string) {

		var sv semver.Version

		_, sv = VersionDetail.RequirementMet("rollback_cm-settings")

		component, _ := cmd.Flags().GetString("component")

		component = strings.ToLower(component)
		if len(component) == 0 || !strings.Contains("ui api", component) {
			logrus.Fatal("--component needs to be 'ui' or 'api'")
		}
		version, _ := cmd.Flags().GetInt("version")
		out, err := rollbackCMSettings(component, version)
		if err != nil {
			logrus.WithError(err).Error("Error rolling back CM Settings")
		}

		if semverV1, err := semver.ParseRange(">=0.1.6"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV1(sv) {
				displayRollbackCmSettingsV1(out)
			}
		}
	},
}

func displayRollbackCmSettingsV1(in string) {
	if strings.ToLower(OutputFormat) == "raw" {
		fmt.Println(in)
		os.Exit(0)
	}
	var vvData objects.VaultVersion
	marshErr := json.Unmarshal([]byte(in), &vvData)
	if marshErr != nil {
		logrus.Fatal("Could not unmarshall data", marshErr)
	}

	if !FormatOverridden {
		OutputFormat = "text"
	}

	switch strings.ToLower(OutputFormat) {
	case "json":
		vvData.ToJSON()
	case "gron":
		vvData.ToGRON()
	case "yaml":
		vvData.ToYAML()
	case "text", "table":
		vvData.ToTEXT(NoHeaders)
	}
}

func rollbackCMSettings(comp string, ver int) (string, error) {
	restClient := resty.New()

	uri := fmt.Sprintf("splicectl/v1/vault/rollbackcmsettings?component=%s&version=%d", comp, ver)
	resp, resperr := restClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Token-Bearer", AuthClient.GetTokenBearer()).
		SetHeader("X-Token-Session", AuthClient.GetSessionID()).
		Post(fmt.Sprintf("%s/%s", ApiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error rolling back CM Settings")
		return "", resperr
	}

	return string(resp.Body()[:]), nil

}

func init() {
	rollbackCmd.AddCommand(rollbackCMSettingsCmd)

	rollbackCMSettingsCmd.Flags().StringP("component", "c", "", "Specify the component, <ui|api>")
	rollbackCMSettingsCmd.Flags().Int("version", 0, "Specify the version to retrieve, default latest")
	rollbackCMSettingsCmd.MarkFlagRequired("version")
}
