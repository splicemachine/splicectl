package apply

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/splicemachine/splicectl/cmd/objects"
	"github.com/splicemachine/splicectl/common"
)

const dataKey = "data"

var (
	noTopLevelDataError   = errors.New("Default CR did not contain top level 'data' element that is required")
	dataIsWrongTypeError  = errors.New("data element in Default CR is not an object, but should be")
	doubleNestedDataError = errors.New("Default CR appears to contain a second level 'data' element, your Default CR appears to be double nested")
)

// validateDefaultCR - validate that the data representing default-cr contains a top
// level field named 'data'.
func validateDefaultCR(defaultCR []byte) (interface{}, error) {
	// get map representation of Default CR
	crMap := make(map[string]interface{})
	if err := json.Unmarshal(defaultCR, &crMap); err != nil {
		return nil, err
	}

	// get the data element of the Default CR
	crData, ok := crMap[dataKey]
	if !ok {
		return crMap, noTopLevelDataError
	}

	// verify that the data element is a map
	crDataMap, ok := crData.(map[string]interface{})
	if !ok {
		return crMap, dataIsWrongTypeError
	}

	// verify that there is not a data element in the top level data element,
	// would imply double nesting of Default CR
	if _, ok := crDataMap[dataKey]; ok {
		return crData, doubleNestedDataError
	}

	return crMap, nil
}

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
		if _, err := validateDefaultCR(jsonBytes); err != nil {
			logrus.WithError(err).Fatal("Error validating Default CR")
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
