package version

import (
	"fmt"
	"os"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/sirupsen/logrus"
	"github.com/splicemachine/splicectl/common"

	"github.com/spf13/cobra"
)

var versionsVaultKeyCmd = &cobra.Command{
	Use:   "vault-key",
	Short: "Retrieve the versions of a specified vault key from the cluster.",
	Long: `EXAMPLES
	splicectl versions vault-key --keypath services/cloudmanager/config/default/ui
	`,
	Run: func(cmd *cobra.Command, args []string) {
		_, sv := c.VersionDetail.RequirementMet("versions_vault-key")

		keyPath, _ := cmd.Flags().GetString("keypath")
		if strings.HasPrefix(keyPath, "secrets/") {
			keyPath = strings.TrimPrefix(keyPath, "secrets/")
		}
		out, err := getVaultKeyVersionData(keyPath)
		if err != nil {
			logrus.WithError(err).Error("Error getting Default CR Info")
		}

		if semverV1, err := semver.ParseRange(">=0.0.15 <0.0.17"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV1(sv) {
				displayVersionsVaultKeyV1(out)
			}
		}

		if semverV2, err := semver.ParseRange(">=0.0.17"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV2(sv) {
				displayVersionsVaultKeyV2(out)
			}
		}
	},
}

func displayVersionsVaultKeyV1(in string) {
	fmt.Println(in)
	os.Exit(0)
}

func displayVersionsVaultKeyV2(in string) {
	if strings.ToLower(c.OutputFormat) == "raw" {
		fmt.Println(in)
		os.Exit(0)
	}
	vkData, cerr := common.RestructureVersions(in)
	if cerr != nil {
		logrus.Fatal("Vault Version JSON conversion failed.")
	}

	c.OutputData(&vkData)
}

func getVaultKeyVersionData(keypath string) (string, error) {
	uri := fmt.Sprintf("splicectl/v1/vault/vaultkeyversions?keypath=%s", keypath)
	resp, resperr := c.RestyWithHeaders().
		Get(fmt.Sprintf("%s/%s", c.ApiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error getting Default CR Info")
		return "", resperr
	}

	return string(resp.Body()[:]), nil

}

func init() {
	versionsCmd.AddCommand(versionsVaultKeyCmd)

	versionsVaultKeyCmd.Flags().String("keypath", "", "Specify the vault key path")
	versionsVaultKeyCmd.MarkFlagRequired("keypath")
}
