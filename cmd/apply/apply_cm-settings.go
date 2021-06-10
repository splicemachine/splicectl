package apply

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	c "github.com/splicemachine/splicectl/cmd"
	"github.com/splicemachine/splicectl/cmd/objects"
	"github.com/splicemachine/splicectl/common"

	"github.com/spf13/cobra"
)

var applyCMSettingsCmd = &cobra.Command{
	Use:   "cm-settings",
	Short: "Submit new cm (cloud manager) settings to the cluster.",
	Long: `EXAMPLES
	splicectl get cm-settings --component ui -o json > ~/tmp/cm-ui.json
	#edit file
	splicectl apply cm-settings --component --file ~/tmp/cm-ui.json
`,
	Run: func(cmd *cobra.Command, args []string) {
		component, _ := cmd.Flags().GetString("component")

		var sv semver.Version

		_, sv = c.VersionDetail.RequirementMet("apply_cm-settings")

		component = strings.ToLower(component)
		if len(component) == 0 || !strings.Contains("ui api", component) {
			logrus.Fatal("--component needs to be 'ui' or 'api'")
		}
		filePath, _ := cmd.Flags().GetString("file")
		fileBytes, _ := ioutil.ReadFile(filePath)

		jsonBytes, cerr := common.WantJSON(fileBytes)
		if cerr != nil {
			logrus.Fatal("The input data MUST be in either JSON or YAML format")
		}

		out, err := setCMSettings(component, jsonBytes)
		if err != nil {
			logrus.WithError(err).Error("Error setting System Settings")
		}

		if semverV1, err := semver.ParseRange(">=0.1.6"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV1(sv) {
				displayApplyCmSettingsV1(out)
			}
		}
	},
}

func displayApplyCmSettingsV1(in string) {
	if strings.ToLower(c.OutputFormat) == "raw" {
		fmt.Println(in)
		os.Exit(0)
	}
	var vvData objects.VaultVersion
	marshErr := json.Unmarshal([]byte(in), &vvData)
	if marshErr != nil {
		logrus.Fatal("Could not unmarshall data", marshErr)
	}

	if !c.FormatOverridden {
		c.OutputFormat = "text"
	}

	switch strings.ToLower(c.OutputFormat) {
	case "json":
		vvData.ToJSON()
	case "gron":
		vvData.ToGRON()
	case "yaml":
		vvData.ToYAML()
	case "text", "table":
		vvData.ToTEXT(c.NoHeaders)
	}

}

func setCMSettings(comp string, in []byte) (string, error) {
	restClient := resty.New()
	uri := fmt.Sprintf("splicectl/v1/vault/cmsettings?component=%s", comp)
	resp, resperr := restClient.R().
		SetHeader("X-Token-Bearer", c.AuthClient.GetTokenBearer()).
		SetHeader("X-Token-Session", c.AuthClient.GetSessionID()).
		SetBody(in).
		SetResult(&c.AuthSuccess{}). // or SetResult(AuthSuccess{}).
		SetError(&c.AuthError{}).    // or SetError(AuthError{}).
		Post(fmt.Sprintf("%s/%s", c.ApiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error setting System Settings")
		return "", resperr
	}

	return string(resp.Body()[:]), nil

}

func init() {
	applyCmd.AddCommand(applyCMSettingsCmd)

	applyCMSettingsCmd.Flags().String("file", "", "Specify the input file")
	applyCMSettingsCmd.Flags().StringP("component", "c", "", "Specify the component, <ui|api>")
	applyCMSettingsCmd.MarkFlagRequired("file")

}
