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

var getSystemSettingsCmd = &cobra.Command{
	Use:   "system-settings",
	Short: "Get the default system settings for the cluster.",
	Long: `EXAMPLES
	splicectl get system-settings -o json > ~/tmp/system-settings.json
`,
	Run: func(cmd *cobra.Command, args []string) {

		var sv semver.Version

		_, sv = VersionDetail.RequirementMet("get_system-settings")

		version, _ := cmd.Flags().GetInt("version")
		decode, _ := cmd.Flags().GetBool("decode-values")

		out, err := getSystemSettings(version)
		if err != nil {
			logrus.WithError(err).Error("Error getting System Settings")
		}

		if semverV1, err := semver.ParseRange(">=0.0.14 <0.0.17"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV1(sv) {
				displayGetSystemSettingsV1(out)
			}
		}

		if semverV2, err := semver.ParseRange(">=0.0.17"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV2(sv) {
				displayGetSystemSettingsV2(out, decode)
			}
		}
	},
}

func displayGetSystemSettingsV1(in string) {
	fmt.Println(in)
	os.Exit(0)
}

func displayGetSystemSettingsV2(in string, dc bool) {
	if strings.ToLower(OutputFormat) == "raw" {
		fmt.Println(in)
		os.Exit(0)
	}
	var sessData objects.SystemSettings
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
		sessData.ToTEXT(NoHeaders, dc)
	}

}

func getSystemSettings(ver int) (string, error) {
	restClient := resty.New()

	uri := fmt.Sprintf("splicectl/v1/vault/systemsettings?version=%d", ver)
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
	getCmd.AddCommand(getSystemSettingsCmd)

	getSystemSettingsCmd.Flags().Int("version", 0, "Specify the version to retrieve, default latest")
	getSystemSettingsCmd.Flags().BoolP("decode-values", "d", false, "Decode Base64 Encoded Values")
}
