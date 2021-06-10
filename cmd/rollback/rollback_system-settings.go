package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/sirupsen/logrus"
	c "github.com/splicemachine/splicectl/cmd"
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
		_, sv := c.VersionDetail.RequirementMet("rollback_system-settings")

		version, _ := cmd.Flags().GetInt("version")
		out, err := rollbackSystemSettings(version)
		if err != nil {
			logrus.WithError(err).Error("Error rolling back System Settings")
		}

		if semverV1, err := semver.ParseRange(">=0.0.15 <0.0.17"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV1(sv) {
				displayRollbackSystemSettingsV1(out)
			}
		}

		if semverV2, err := semver.ParseRange(">=0.0.17"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV2(sv) {
				displayRollbackSystemSettingsV2(out)
			}
		}
	},
}

func displayRollbackSystemSettingsV1(in string) {
	fmt.Println(in)
	os.Exit(0)
}

func displayRollbackSystemSettingsV2(in string) {
	if strings.ToLower(c.OutputFormat) == "raw" {
		fmt.Println(in)
		os.Exit(0)
	}
	var vvData objects.VaultVersion
	marshErr := json.Unmarshal([]byte(in), &vvData)
	if marshErr != nil {
		logrus.Fatal("Could not unmarshall data", marshErr)
	}

	c.OutputData(&vvData)
}

func rollbackSystemSettings(ver int) (string, error) {
	uri := fmt.Sprintf("splicectl/v1/vault/rollbacksystemsettings?version=%d", ver)
	resp, resperr := c.RestyWithHeaders().
		Post(fmt.Sprintf("%s/%s", c.ApiServer, uri))

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
