package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/sirupsen/logrus"
	"github.com/splicemachine/splicectl/common"

	"github.com/spf13/cobra"
	c "github.com/splicemachine/splicectl/cmd"
)

var versionsCMSettingsCmd = &cobra.Command{
	Use:   "cm-settings",
	Short: "Retrieve the versions of the cm (cloud manager) settings in the cluster.",
	Long: `EXAMPLES
	splicectl versions cm-settings --component ui
	splicectl versions cm-settings --component api
`,
	Run: func(cmd *cobra.Command, args []string) {
		_, sv := c.VersionDetail.RequirementMet("versions_cm-settings")

		component, _ := cmd.Flags().GetString("component")

		component = strings.ToLower(component)
		if len(component) == 0 || !strings.Contains("ui api", component) {
			logrus.Fatal("--component needs to be 'ui' or 'api'")
		}
		out, err := getCMSettingsVersions(component)
		if err != nil {
			logrus.WithError(err).Error("Error getting CM Settings")
		}

		if semverV1, err := semver.ParseRange(">=0.1.6"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV1(sv) {
				displayVersionsCmSettingsV1(out)
			}
		}
	},
}

func displayVersionsCmSettingsV1(in string) {
	if strings.ToLower(c.OutputFormat) == "raw" {
		fmt.Println(in)
		os.Exit(0)
	}
	ssData, cerr := common.RestructureVersions(in)
	if cerr != nil {
		logrus.Fatal("Vault Version JSON conversion failed.")
	}

	c.OutputData(&ssData)
}

func getCMSettingsVersions(comp string) (string, error) {
	uri := fmt.Sprintf("splicectl/v1/vault/cmsettingsversions?component=%s", comp)
	resp, resperr := c.RestyWithHeaders().
		Get(fmt.Sprintf("%s/%s", c.ApiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error getting System Settings")
		return "", resperr
	}

	return string(resp.Body()[:]), nil

}

func init() {
	versionsCmd.AddCommand(versionsCMSettingsCmd)
	versionsCMSettingsCmd.Flags().StringP("component", "c", "", "Specify the component, <ui|api>")

}
