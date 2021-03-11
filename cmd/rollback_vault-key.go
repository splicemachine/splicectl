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

var rollbackVaultKeyCmd = &cobra.Command{
	Use:   "vault-key",
	Short: "Rollback a specified vault key to a specific vault version",
	Long: `EXAMPLES
	splicectl versions vault-key --keypath services/cloudmanager/config/default/ui
	splicectl rollback vault-key --keypath services/cloudmanager/config/default/ui --version 1
	`,
	Run: func(cmd *cobra.Command, args []string) {

		var sv semver.Version

		_, sv = versionDetail.RequirementMet("rollback_vault-key")

		keyPath, _ := cmd.Flags().GetString("keypath")
		if strings.HasPrefix(keyPath, "secrets/") {
			keyPath = strings.TrimPrefix(keyPath, "secrets/")
		}
		version, _ := cmd.Flags().GetInt("version")
		out, err := rollbackVaultKeyData(keyPath, version)
		if err != nil {
			logrus.WithError(err).Error("Error rolling back Vault Key")
		}

		if semverV1, err := semver.ParseRange(">=0.0.15 <0.0.17"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV1(sv) {
				displayRollbackVaultKeyV1(out)
			}
		}

		if semverV2, err := semver.ParseRange(">=0.0.17"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV2(sv) {
				displayRollbackVaultKeyV2(out)
			}
		}
	},
}

func displayRollbackVaultKeyV1(in string) {
	fmt.Println(in)
	os.Exit(0)
}

func displayRollbackVaultKeyV2(in string) {
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

func rollbackVaultKeyData(keypath string, ver int) (string, error) {
	restClient := resty.New()

	uri := fmt.Sprintf("splicectl/v1/vault/rollbackvaultkey?version=%d&keypath=%s", ver, keypath)
	resp, resperr := restClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Token-Bearer", authClient.GetTokenBearer()).
		SetHeader("X-Token-Session", authClient.GetSessionID()).
		Post(fmt.Sprintf("%s/%s", apiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error rolling back Vault Key")
		return "", resperr
	}

	return string(resp.Body()[:]), nil

}

func init() {
	rollbackCmd.AddCommand(rollbackVaultKeyCmd)

	rollbackVaultKeyCmd.Flags().String("keypath", "", "Specify the vault key path")
	rollbackVaultKeyCmd.Flags().String("output", "json", "Specify the output type")
	rollbackVaultKeyCmd.Flags().Int("version", 0, "Specify the version to retrieve, default latest")
	rollbackVaultKeyCmd.MarkFlagRequired("keypath")
	rollbackVaultKeyCmd.MarkFlagRequired("version")
}
