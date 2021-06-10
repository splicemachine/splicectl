package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/go-resty/resty/v2"
	"github.com/maahsome/gron"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/spf13/cobra"
)

var getVaultKeyCmd = &cobra.Command{
	Use:   "vault-key",
	Short: "Get the data from a specific vault key",
	Long: `EXAMPLES
	splicectl get vault-key --keypath services/cloudmanager/config/default/ui -o json > ~/tmp/cm-ui.json
	`,
	Run: func(cmd *cobra.Command, args []string) {

		var sv semver.Version

		_, sv = VersionDetail.RequirementMet("get_vault-key")

		keyPath, _ := cmd.Flags().GetString("keypath")
		if strings.HasPrefix(keyPath, "secrets/") {
			keyPath = strings.TrimPrefix(keyPath, "secrets/")
		}
		version, _ := cmd.Flags().GetInt("version")
		out, err := getVaultKeyData(keyPath, version)
		if err != nil {
			logrus.WithError(err).Error("Error getting Default CR Info")
		}

		if semverV1, err := semver.ParseRange(">=0.0.14 <0.0.17"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV1(sv) {
				displayGetVaultKeyV1(out)
			}
		}

		if semverV2, err := semver.ParseRange(">=0.0.17"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV2(sv) {
				displayGetVaultKeyV2(out)
			}
		}
	},
}

func displayGetVaultKeyV1(in string) {
	fmt.Println(in)
	os.Exit(0)
}

func displayGetVaultKeyV2(in string) {
	if strings.ToLower(OutputFormat) == "raw" {
		fmt.Println(in)
		os.Exit(0)
	}
	var vaultKey map[string]interface{}

	marshErr := json.Unmarshal([]byte(in), &vaultKey)
	if marshErr != nil {
		logrus.Fatal("Could not unmarshall data", marshErr)
	}

	if !FormatOverridden {
		OutputFormat = "yaml"
	}

	switch strings.ToLower(OutputFormat) {
	case "json":
		fmt.Println(string(in[:]))
	case "gron":
		subReader := strings.NewReader(string(in[:]))
		subValues := &bytes.Buffer{}
		ges := gron.NewGron(subReader, subValues)
		ges.SetMonochrome(false)
		serr := ges.ToGron()
		if serr != nil {
			logrus.Error("Problem generating gron syntax", serr)
		} else {
			fmt.Println(string(subValues.Bytes()))
		}
	case "yaml", "text", "table":
		rawVaultKey, crerr := yaml.Marshal(vaultKey)
		if crerr != nil {
			logrus.WithError(crerr).Error("Failed to convert to YAML")
		}
		fmt.Println(string(rawVaultKey[:]))
	}

}

func getVaultKeyData(keypath string, ver int) (string, error) {
	restClient := resty.New()

	uri := fmt.Sprintf("splicectl/v1/vault/vaultkey?version=%d&keypath=%s", ver, keypath)
	resp, resperr := restClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Token-Bearer", AuthClient.GetTokenBearer()).
		SetHeader("X-Token-Session", AuthClient.GetSessionID()).
		Get(fmt.Sprintf("%s/%s", ApiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error getting Default CR Info")
		return "", resperr
	}

	return string(resp.Body()[:]), nil

}

func init() {
	getCmd.AddCommand(getVaultKeyCmd)

	getVaultKeyCmd.Flags().String("keypath", "", "Specify the vault key path")
	getVaultKeyCmd.Flags().Int("version", 0, "Specify the version to retrieve, default latest")
	getVaultKeyCmd.MarkFlagRequired("keypath")
}
