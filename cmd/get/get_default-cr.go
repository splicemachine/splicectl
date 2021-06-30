package get

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/sirupsen/logrus"

	"github.com/maahsome/gron"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var getDefaultCRCmd = &cobra.Command{
	Use:   "default-cr",
	Short: "Get the default CR for the cluster.",
	Long: `EXAMPLES
	splicectl get default-cr -o json > ~/tmp/default-cr.json
`,
	Run: func(cmd *cobra.Command, args []string) {
		_, sv := c.VersionDetail.RequirementMet("get_default-cr")

		version, _ := cmd.Flags().GetInt("version")
		out, err := getDefaultCR(version)
		if err != nil {
			logrus.WithError(err).Error("Error getting Default CR Info")
		}

		if semverV1, err := semver.ParseRange(">=0.0.14 <0.0.17"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV1(sv) {
				displayGetDefaultCRV1(out)
			}
		}

		if semverV2, err := semver.ParseRange(">=0.0.17"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV2(sv) {
				displayGetDefaultCRV2(out)
			}
		}

	},
}

func displayGetDefaultCRV1(in string) {
	fmt.Println(in)
	os.Exit(0)
}

func displayGetDefaultCRV2(in string) {
	if strings.ToLower(c.OutputFormat) == "raw" {
		fmt.Println(in)
		os.Exit(0)
	}

	var defaultCr map[string]interface{}

	marshErr := json.Unmarshal([]byte(in), &defaultCr)
	if marshErr != nil {
		logrus.Fatal("Could not unmarshall data", marshErr)
	}

	if !c.FormatOverridden {
		c.OutputFormat = "yaml"
	}

	switch strings.ToLower(c.OutputFormat) {
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
		rawDefaultCr, crerr := yaml.Marshal(defaultCr)
		if crerr != nil {
			logrus.WithError(crerr).Error("Failed to convert to YAML")
		}
		fmt.Println(string(rawDefaultCr[:]))
	}

}

func getDefaultCR(ver int) (string, error) {
	uri := fmt.Sprintf("splicectl/v1/vault/defaultcr?version=%d", ver)
	resp, resperr := c.RestyWithHeaders().
		Get(fmt.Sprintf("%s/%s", c.ApiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error getting Default CR Info")
		return "", resperr
	}

	return string(resp.Body()[:]), nil

}

func init() {
	getCmd.AddCommand(getDefaultCRCmd)

	getDefaultCRCmd.Flags().Int("version", 0, "Specify the version to retrieve, default latest")
}
