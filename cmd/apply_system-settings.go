package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/splicemachine/splicectl/cmd/objects"
	"github.com/splicemachine/splicectl/common"

	"github.com/spf13/cobra"
)

var applySystemSettingsCmd = &cobra.Command{
	Use:   "system-settings",
	Short: "Submit new system settings to the cluster.",
	Long: `EXAMPLES
	splicectl get system-settings -o json > ~/tmp/system-settings.json
	#edit file
	splicectl apply system-settings --file ~/tmp/system-settings.json
`,
	Run: func(cmd *cobra.Command, args []string) {
		filePath, _ := cmd.Flags().GetString("file")

		fileBytes, _ := ioutil.ReadFile(filePath)

		jsonBytes, cerr := common.WantJSON(fileBytes)
		if cerr != nil {
			logrus.Fatal("The input data MUST be in either JSON or YAML format")
		}

		out, err := setSystemSettings(jsonBytes)
		if err != nil {
			logrus.WithError(err).Error("Error setting System Settings")
		}
		var vvData objects.VaultVersion
		marshErr := json.Unmarshal([]byte(out), &vvData)
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

	},
}

func setSystemSettings(in []byte) (string, error) {
	restClient := resty.New()
	uri := "splicectl/v1/vault/systemsettings"
	resp, resperr := restClient.R().
		SetHeader("X-Token-Bearer", authClient.GetTokenBearer()).
		SetHeader("X-Token-Session", authClient.GetSessionID()).
		SetBody(in).
		SetResult(&AuthSuccess{}). // or SetResult(AuthSuccess{}).
		SetError(&AuthError{}).    // or SetError(AuthError{}).
		Post(fmt.Sprintf("%s/%s", apiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error setting System Settings")
		return "", resperr
	}

	return string(resp.Body()[:]), nil

}

func init() {
	applyCmd.AddCommand(applySystemSettingsCmd)

	applySystemSettingsCmd.Flags().String("file", "", "Specify the input file")
	applySystemSettingsCmd.MarkFlagRequired("file")
}
