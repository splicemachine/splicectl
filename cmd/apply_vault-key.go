package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/splicemachine/splicectl/cmd/objects"
	"github.com/splicemachine/splicectl/common"

	"github.com/spf13/cobra"
)

var applyVaultKeyCmd = &cobra.Command{
	Use:   "vault-key",
	Short: "Submit new data to a specific key path in vault.",
	Long: `EXAMPLES
	splicectl get vault-key --keypath services/cloudmanager/config/default/ui -o json > ~/tmp/cm-ui.json
	# edit file
	splicectl apply vault-key --keypath services/cloudmanager/config/default/ui --file ~/tmp/cm-ui.json
`,
	Run: func(cmd *cobra.Command, args []string) {

		var sv semver.Version

		_, sv = versionDetail.RequirementMet("apply_vault-key")

		keyPath, _ := cmd.Flags().GetString("keypath")
		if strings.HasPrefix(keyPath, "secrets/") {
			keyPath = strings.TrimPrefix(keyPath, "secrets/")
		}
		filePath, _ := cmd.Flags().GetString("file")
		fileBytes, _ := ioutil.ReadFile(filePath)

		jsonBytes, cerr := common.WantJSON(fileBytes)
		if cerr != nil {
			logrus.Fatal("The input data MUST be in either JSON or YAML format")
		}

		out, err := setVaultKeyData(keyPath, jsonBytes)
		if err != nil {
			logrus.WithError(err).Error("Error setting Vault-Key Data")
		}

		if semverV1, err := semver.ParseRange(">=0.0.14 <0.0.17"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV1(sv) {
				displayApplyVaultKeyV1(out)
			}
		}

		if semverV2, err := semver.ParseRange(">=0.0.17"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV2(sv) {
				displayApplyVaultKeyV2(out)
			}
		}
	},
}

func displayApplyVaultKeyV1(in string) {
	fmt.Println(in)
	os.Exit(0)
}

func displayApplyVaultKeyV2(in string) {
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

func setVaultKeyData(keypath string, in []byte) (string, error) {
	restClient := resty.New()
	uri := fmt.Sprintf("splicectl/v1/vault/vaultkey?keypath=%s", keypath)
	resp, resperr := restClient.R().
		SetHeader("X-Token-Bearer", authClient.GetTokenBearer()).
		SetHeader("X-Token-Session", authClient.GetSessionID()).
		SetBody(in).
		SetResult(&AuthSuccess{}). // or SetResult(AuthSuccess{}).
		SetError(&AuthError{}).    // or SetError(AuthError{}).
		Post(fmt.Sprintf("%s/%s", apiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error setting Default CR Info")
		return "", resperr
	}

	return string(resp.Body()[:]), nil

}

func init() {
	applyCmd.AddCommand(applyVaultKeyCmd)

	applyVaultKeyCmd.Flags().String("keypath", "", "Specify the vault key path")
	applyVaultKeyCmd.Flags().String("file", "", "Specify the input file")
	applyVaultKeyCmd.MarkFlagRequired("keypath")
	applyVaultKeyCmd.MarkFlagRequired("file")
}
