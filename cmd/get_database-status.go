package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/blang/semver/v4"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var getDatabaseStatus = &cobra.Command{
	Use:   "database-status",
	Short: "Get the status of database.",
	Long: `EXAMPLES
	splicectl get database-status --database-name "test"
`,
	Run: func(cmd *cobra.Command, args []string) {
		var dberr error
		var sv semver.Version

		_, sv = versionDetail.RequirementMet("get_database-status")

		databaseName, _ := cmd.Flags().GetString("database-name")
		if len(databaseName) == 0 {
			databaseName, dberr = promptForDatabaseName()
			if dberr != nil {
				logrus.Fatal("Could not get name of Database", dberr)
			}
		}

		out, err := getDatabaseStatusData(databaseName)
		if err != nil {
			logrus.WithError(err).Error("Error getting status of database ")
		}

		if semverV1, err := semver.ParseRange(">=0.1.6"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV1(sv) {
				displayGetDatabaseStatusV1(out)
			}
		}
	},
}

func displayGetDatabaseStatusV1(in string) {
	fmt.Println(in)
	os.Exit(0)
}

func getDatabaseStatusData(databaseName string) (string, error) {

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

	uri := fmt.Sprintf("splicectl/v1/splicedb/splicedatabasestatus?database-name=%s", databaseName)
	resp, resperr := restClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Token-Bearer", authClient.GetTokenBearer()).
		SetHeader("X-Token-Session", authClient.GetSessionID()).
		Get(fmt.Sprintf("%s/%s", apiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error getting Database status info")
		return "", resperr
	}
	return string(resp.Body()[:]), nil
}

func init() {
	getCmd.AddCommand(getDatabaseStatus)
	getDatabaseStatus.Flags().StringP("database-name", "d", "", "Specify the database name")
}
