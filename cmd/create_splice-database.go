package cmd

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// AuthSuccess - Varibles to extract on a resty auth success
type AuthSuccess struct {
	/* variables */
}

// AuthError - Variables to extract on a resty auth error
type AuthError struct {
	/* variables */
}

var createSpliceDatabaseCmd = &cobra.Command{
	Use:   "splice-database",
	Short: "Create a SpliceDBCluster KIND",
	Long:  `Create a SpliceDBCluster KIND`,
	Run: func(cmd *cobra.Command, args []string) {

		dryRun, _ := cmd.Flags().GetBool("dry-run")
		databaseName, _ := cmd.Flags().GetString("database-name")
		dnsPrefix, _ := cmd.Flags().GetString("dns-prefix")

		if len(databaseName) == 0 {
			logrus.Fatal("--database-name is a required flag")
		}

		if len(dnsPrefix) == 0 {
			logrus.Fatal("--database-name is a required flag")
		}

		out, err := generateCR(databaseName, dnsPrefix, dryRun)
		if err != nil {
			logrus.WithError(err).Error("Error Generating Default CR Info")
		}
		fmt.Println(out)
	},
}

func generateCR(dbname string, dns string, outputonly bool) (string, error) {
	restClient := resty.New()

	uri := fmt.Sprintf("splicectl/v1/splicedb/splicedatabase?database-name=%s&dns-prefix=%s&user=%s", dbname, dns, os.Getenv("USER"))

	var resp *resty.Response
	var resperr error

	logrus.WithFields(logrus.Fields{
		"dbname":     dbname,
		"dns":        dns,
		"outputonly": outputonly,
	}).Info("Field Values")

	if outputonly {
		resp, resperr = restClient.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Accept", "application/json").
			SetHeader("X-Token-Bearer", authClient.GetTokenBearer()).
			SetHeader("X-Token-Session", authClient.GetSessionID()).
			Get(fmt.Sprintf("%s/%s", apiServer, uri))
	} else {
		resp, resperr = restClient.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Accept", "application/json").
			SetHeader("X-Token-Bearer", authClient.GetTokenBearer()).
			SetHeader("X-Token-Session", authClient.GetSessionID()).
			Post(fmt.Sprintf("%s/%s", apiServer, uri))
	}
	if resperr != nil {
		logrus.WithError(resperr).Error("Error getting Default CR Info")
		return "", resperr
	}

	return string(resp.Body()[:]), nil

}

func init() {
	createCmd.AddCommand(createSpliceDatabaseCmd)

	createSpliceDatabaseCmd.Flags().BoolP("dry-run", "", false, "Output data rather than create the K8s resource")
	createSpliceDatabaseCmd.Flags().StringP("database-name", "d", "test", "Specify the Splice Machine Database Cluster Name")
	createSpliceDatabaseCmd.Flags().StringP("dns-prefix", "p", "test", "Specify the Splice Machine Database DNS Prefix/K8s Namespace")

}
