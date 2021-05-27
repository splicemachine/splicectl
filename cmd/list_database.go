package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/splicemachine/splicectl/cmd/objects"

	"github.com/spf13/cobra"
)

var listDatabaseCmd = &cobra.Command{
	Use:     "database",
	Aliases: []string{"databases"},
	Short:   "Retrieve a list of splice databases in the cluster.",
	Long: `EXAMPLES
	splicectl list database
`,
	Run: func(cmd *cobra.Command, args []string) {

		var sv semver.Version

		_, sv = versionDetail.RequirementMet("list_database")

		// databaseName, _ := cmd.Flags().GetString("database-name")
		out, err := getDatabaseList()
		if err != nil {
			logrus.WithError(err).Error("Error getting Database CR Info")
		}

		if semverV1, err := semver.ParseRange(">=0.0.14 <0.0.17"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV1(sv) {
				displayListDatabaseV1(out)
			}
		}

		if semverV2, err := semver.ParseRange(">=0.0.17"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV2(sv) {
				displayListDatabaseV2(out)
			}
		}
	},
}

func displayListDatabaseV1(in string) {
	fmt.Println(in)
	os.Exit(0)
}

func displayListDatabaseV2(in string) {
	if strings.ToLower(outputFormat) == "raw" {
		fmt.Println(in)
		os.Exit(0)
	}
	var dbList objects.DatabaseList

	marshErr := json.Unmarshal([]byte(in), &dbList)
	if marshErr != nil {
		logrus.Fatal("Could not unmarshall data", marshErr)
	}

	if !formatOverridden {
		outputFormat = "table"
	}

	switch strings.ToLower(outputFormat) {
	case "json":
		dbList.ToJSON()
	case "gron":
		dbList.ToGRON()
	case "yaml":
		dbList.ToYAML()
	case "text", "table":
		dbList.ToTEXT(noHeaders)
	}

}
func getDatabaseList() (string, error) {
	restClient := resty.New()
	// Check if we've set a caBundle (via --ca-cert parameter)
	if len(caBundle) > 0 {
		roots := x509.NewCertPool()
		ok := roots.AppendCertsFromPEM([]byte(caBundle))
		if !ok {
			logrus.Info("Failed to parse CABundle")
		}
		restClient.SetTLSClientConfig(&tls.Config{RootCAs: roots})
	}

	uri := "splicectl/v1/splicedb/splicedatabase"
	resp, resperr := restClient.R().
		// SetHeader("Content-Type", "application/json").
		// SetHeader("Accept", "application/json").
		SetHeader("X-Token-Bearer", authClient.GetTokenBearer()).
		SetHeader("X-Token-Session", authClient.GetSessionID()).
		Execute("LIST", fmt.Sprintf("%s/%s", apiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error getting Database List")
		return "", resperr
	}

	return string(resp.Body()[:]), nil

}

func init() {
	listCmd.AddCommand(listDatabaseCmd)

	// getDatabaseCRCmd.Flags().String("database-name", "", "Specify the database name")

}
