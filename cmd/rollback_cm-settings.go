package cmd

import (
	"crypto/tls"
	"crypto/x509"
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

		_, sv = versionDetail.RequirementMet("rollback_cm-settings")

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
	if strings.ToLower(outputFormat) == "raw" {
		fmt.Println(in)
		os.Exit(0)
	}
	var vvData objects.VaultVersion
	marshErr := json.Unmarshal([]byte(in), &vvData)
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
}

func rollbackCMSettings(comp string, ver int) (string, error) {
	restClient := resty.New()
	// Check if we've set a caBundle (via --ca-cert parameter)
	if len(caBundle) > 0 {
		roots := x509.NewCertPool()
		ok := roots.AppendCertsFromPEM([]byte(caBundle))
		if !ok {
			logrus.Info("Failed to parse CABundle")
		}
		restClient.SetTLSClientConfig(&tls.Config{RootCAs: roots})
	}

	uri := fmt.Sprintf("splicectl/v1/vault/rollbackcmsettings?component=%s&version=%d", comp, ver)
	resp, resperr := restClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Token-Bearer", authClient.GetTokenBearer()).
		SetHeader("X-Token-Session", authClient.GetSessionID()).
		Post(fmt.Sprintf("%s/%s", apiServer, uri))

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
