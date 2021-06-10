package version

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

var versionsSystemSettingsCmd = &cobra.Command{
	Use:   "system-settings",
	Short: "Retrieve the versions of the system settings in the cluster.",
	Long: `EXAMPLES
	splicectl versions system-settings
`,
	Run: func(cmd *cobra.Command, args []string) {
		_, sv := c.VersionDetail.RequirementMet("versions_system-settings")

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

func getSystemSettingsVersions() (string, error) {
	uri := "splicectl/v1/vault/systemsettingsversions"
	resp, resperr := c.RestyWithHeaders().
		Get(fmt.Sprintf("%s/%s", c.ApiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error getting System Settings")
		return "", resperr
	}

	return string(resp.Body()[:]), nil

}

func init() {
	versionsCmd.AddCommand(versionsSystemSettingsCmd)

}
