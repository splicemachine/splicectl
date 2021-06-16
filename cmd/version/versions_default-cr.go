package version

import (
	"fmt"
	"os"
	"strings"

	"github.com/blang/semver/v4"
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
		_, sv := c.VersionDetail.RequirementMet("versions_default-cr")

		out, err := getDefaultCRVersions()
		if err != nil {
			logrus.WithError(err).Error("Error getting Default CR Info")
		}

		if semverV1, err := semver.ParseRange(">=0.0.15 <0.0.17"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV1(sv) {
				displayVersionsDefaultCRV1(out)
			}
		}

		if semverV2, err := semver.ParseRange(">=0.0.17"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV2(sv) {
				displayVersionsDefaultCRV2(out)
			}
		}
	},
}

func displayVersionsDefaultCRV1(in string) {
	fmt.Println(in)
	os.Exit(0)
}

func displayVersionsDefaultCRV2(in string) {
	if strings.ToLower(c.OutputFormat) == "raw" {
		fmt.Println(in)
		os.Exit(0)
	}
	crData, cerr := common.RestructureVersions(in)
	if cerr != nil {
		logrus.Fatal("Vault Version JSON conversion failed.")
	}

	c.OutputData(&crData)
}

func getDefaultCRVersions() (string, error) {
	uri := "splicectl/v1/vault/defaultcrversions"
	resp, resperr := c.RestyWithHeaders().
		Get(fmt.Sprintf("%s/%s", c.ApiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error getting Default CR Info")
		return "", resperr
	}

	return string(resp.Body()[:]), nil

}

func init() {
	versionsCmd.AddCommand(versionsDefaultCRCmd)

}
