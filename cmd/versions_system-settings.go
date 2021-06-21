package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strings"

	"github.com/blang/semver/v4"
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

		var sv semver.Version

		_, sv = versionDetail.RequirementMet("versions_system-settings")

		out, err := getSystemSettingsVersions()
		if err != nil {
			logrus.WithError(err).Error("Error getting System Settings")
		}

		if semverV1, err := semver.ParseRange(">=0.0.15 <0.0.17"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV1(sv) {
				displayVersionsSystemSettingsV1(out)
			}
		}

		if semverV2, err := semver.ParseRange(">=0.0.17"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV2(sv) {
				displayVersionsSystemSettingsV2(out)
			}
		}
	},
}

func displayVersionsSystemSettingsV1(in string) {
	fmt.Println(in)
	os.Exit(0)
}

func displayVersionsSystemSettingsV2(in string) {
	if strings.ToLower(outputFormat) == "raw" {
		fmt.Println(in)
		os.Exit(0)
	}
	ssData, cerr := common.RestructureVersions(in)
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
}

func getSystemSettingsVersions() (string, error) {
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
