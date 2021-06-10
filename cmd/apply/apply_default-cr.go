package apply

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	c "github.com/splicemachine/splicectl/cmd"
	"github.com/splicemachine/splicectl/cmd/objects"
	"github.com/splicemachine/splicectl/common"
)

var applyDefaultCRCmd = &cobra.Command{
	Use:   "default-cr",
	Short: "Submit a new default-cr to the cluster",
	Long: `EXAMPLES
	splicectl get default-cr -o json > ~/tmp/default-cr.json
	# edit file
	splicectl apply default-cr --file ~/tmp/default-cr.json
`,
	Run: func(cmd *cobra.Command, args []string) {
		_, sv := c.VersionDetail.RequirementMet("apply_default-cr")

		filePath, _ := cmd.Flags().GetString("file")
		fileBytes, _ := ioutil.ReadFile(filePath)

		jsonBytes, cerr := common.WantJSON(fileBytes)
		if cerr != nil {
			logrus.Fatal("The input data MUST be in either JSON or YAML format")
		}

		out, err := setDefaultCR(jsonBytes)
		if err != nil {
			logrus.WithError(err).Error("Error setting Default CR Info")
		}

		if semverV1, err := semver.ParseRange(">=0.0.14 <0.0.17"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV1(sv) {
				displayApplyDefaultCRV1(out)
			}
		}

		if semverV2, err := semver.ParseRange(">=0.0.17"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV2(sv) {
				displayApplyDefaultCRV2(out)
			}
		}
	},
}

func displayApplyDefaultCRV1(in string) {
	fmt.Println(in)
	os.Exit(0)
}

func displayApplyDefaultCRV2(in string) {
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

func setDefaultCR(in []byte) (string, error) {
	uri := "splicectl/v1/vault/defaultcr"
	resp, resperr := c.RestyWithHeaders().
		SetBody(in).
		Post(fmt.Sprintf("%s/%s", c.ApiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error setting Default CR Info")
		return "", resperr
	}

	return string(resp.Body()[:]), nil

}

func init() {
	applyCmd.AddCommand(applyDefaultCRCmd)

	applyDefaultCRCmd.Flags().String("file", "", "Specify the input file")
	applyDefaultCRCmd.MarkFlagRequired("file")
}
