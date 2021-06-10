package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	c "github.com/splicemachine/splicectl/cmd"
	"github.com/splicemachine/splicectl/cmd/objects"
)

var rollbackDefaultCRCmd = &cobra.Command{
	Use:   "default-cr",
	Short: "Rollback the default CR for the cluster to a specific vault version.",
	Long: `EXAMPLES
	splicectl versions default-cr
	splicectl rollback default-cr --version 1
`,
	Run: func(cmd *cobra.Command, args []string) {
		_, sv := c.VersionDetail.RequirementMet("rollback_default-cr")

		version, _ := cmd.Flags().GetInt("version")
		out, err := rollbackDefaultCR(version)
		if err != nil {
			logrus.WithError(err).Error("Error getting Default CR Info")
		}

		if semverV1, err := semver.ParseRange(">=0.0.15 <0.0.17"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV1(sv) {
				displayRollbackDefaultCRV1(out)
			}
		}

		if semverV2, err := semver.ParseRange(">=0.0.17"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV2(sv) {
				displayRollbackDefaultCRV2(out)
			}
		}
	},
}

func displayRollbackDefaultCRV1(in string) {
	fmt.Println(in)
	os.Exit(0)
}

func displayRollbackDefaultCRV2(in string) {
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

func rollbackDefaultCR(ver int) (string, error) {
	uri := fmt.Sprintf("splicectl/v1/vault/rollbackdefaultcr?version=%d", ver)
	resp, resperr := c.RestyWithHeaders().
		Post(fmt.Sprintf("%s/%s", c.ApiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error getting Default CR Info")
		return "", resperr
	}

	return string(resp.Body()[:]), nil

}

func init() {
	rollbackCmd.AddCommand(rollbackDefaultCRCmd)

	rollbackDefaultCRCmd.Flags().String("output", "json", "Specify the output type")
	rollbackDefaultCRCmd.Flags().Int("version", 0, "Specify the version to retrieve, default latest")
	rollbackDefaultCRCmd.MarkFlagRequired("version")
}
