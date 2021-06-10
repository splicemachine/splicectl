package apply

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/sirupsen/logrus"
	c "github.com/splicemachine/splicectl/cmd"
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
		_, sv := c.VersionDetail.RequirementMet("apply_system-settings")

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

		if semverV1, err := semver.ParseRange(">=0.0.14 <0.0.17"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV1(sv) {
				displayApplySystemSettingsV1(out)
			}
		}

		if semverV2, err := semver.ParseRange(">=0.0.17"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV2(sv) {
				displayApplySystemSettingsV2(out)
			}
		}

	},
}

func displayApplySystemSettingsV1(in string) {
	fmt.Println(in)
	os.Exit(0)
}

func displayApplySystemSettingsV2(in string) {
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

func setSystemSettings(in []byte) (string, error) {
	uri := "splicectl/v1/vault/systemsettings"
	resp, resperr := c.RestyWithHeaders().
		SetBody(in).
		Post(fmt.Sprintf("%s/%s", c.ApiServer, uri))

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
