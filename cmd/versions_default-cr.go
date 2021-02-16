package cmd

import (
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/splicemachine/splicectl/common"
)

var versionsDefaultCRCmd = &cobra.Command{
	Use:   "default-cr",
	Short: "Retrieve the versions of the default CR in the cluster.",
	Long: `EXAMPLES
	splicectl versions default-cr
`,
	Run: func(cmd *cobra.Command, args []string) {
		out, err := getDefaultCRVersions()
		if err != nil {
			logrus.WithError(err).Error("Error getting Default CR Info")
		}

		crData, cerr := common.RestructureVersions(out)
		if cerr != nil {
			logrus.Fatal("Vault Version JSON conversion failed.")
		}

		if !formatOverridden {
			outputFormat = "text"
		}

		switch strings.ToLower(outputFormat) {
		case "json":
			crData.ToJSON()
		case "gron":
			crData.ToGRON()
		case "yaml":
			crData.ToYAML()
		case "text", "table":
			crData.ToTEXT(noHeaders)
		}

	},
}

func getDefaultCRVersions() (string, error) {
	restClient := resty.New()

	uri := "splicectl/v1/vault/defaultcrversions"
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
	versionsCmd.AddCommand(versionsDefaultCRCmd)

}
