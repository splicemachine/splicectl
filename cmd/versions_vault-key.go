package cmd

import (
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
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
		keyPath, _ := cmd.Flags().GetString("keypath")
		if strings.HasPrefix(keyPath, "secrets/") {
			keyPath = strings.TrimPrefix(keyPath, "secrets/")
		}
		out, err := getVaultKeyVersionData(keyPath)
		if err != nil {
			logrus.WithError(err).Error("Error getting Default CR Info")
		}
		vkData, cerr := common.RestructureVersions(out)
		if cerr != nil {
			logrus.Fatal("Vault Version JSON conversion failed.")
		}

		if !formatOverridden {
			outputFormat = "text"
		}

		switch strings.ToLower(outputFormat) {
		case "json":
			vkData.ToJSON()
		case "gron":
			vkData.ToGRON()
		case "yaml":
			vkData.ToYAML()
		case "text", "table":
			vkData.ToTEXT(noHeaders)
		}
	},
}

func getVaultKeyVersionData(keypath string) (string, error) {
	restClient := resty.New()

	uri := fmt.Sprintf("splicectl/v1/vault/vaultkeyversions?keypath=%s", keypath)
	resp, resperr := restClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Token-Bearer", authClient.GetTokenBearer()).
		SetHeader("X-Token-Session", authClient.GetSessionID()).
		Get(fmt.Sprintf("%s/%s", apiServer, uri))

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
