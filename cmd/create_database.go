package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/splicemachine/splicectl/cmd/objects"
	"github.com/splicemachine/splicectl/common"

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

var createDatabaseCmd = &cobra.Command{
	Use: "workspace",
	// splice-database used to be the command name, it is only being kept as an alias for back-compat as it does not match the naming system for other database/workspace commands
	Aliases: []string{"database", "splice-database"},
	Short:   "Create a Splice Machine Database",
	Long: `EXAMPLES
	splicectl get accounts

	While you can specify each of the parameters on the command line, it is much
	easier to create a SKEL yaml file and use that to create the database.

	splicectl create workspace --skel --account-id <accountid> --cloud-provider <aws|az|gcp|op|none> > ~/tmp/splicedb-create.yaml
	# edit the ~/tmp/splicedb-create.yaml
	splicectl create workspace --file ~/tmp/splicedb-create.yaml
	
	Note: --database-name and -d are the preferred way to supply the database name.
	However, --database and --workspace can also be used as well. In the event that
	more than one of them is supplied database-name and d are preferred over all
	and workspace is preferred over database. The most preferred option that is
	supplied will be used and a message will be displayed letting you know which
	option was chosen if more than one were supplied.`,
	Run: func(cmd *cobra.Command, args []string) {

		var sv semver.Version

		_, sv = VersionDetail.RequirementMet("create_database")

		// Look for --file first, load that into the structure, then read each
		// parameters and override the values loaded from the input file
		skel, _ := cmd.Flags().GetBool("skel")
		file, _ := cmd.Flags().GetString("file")
		fileProvided := false

		dbReq := objects.DatabaseRequest{}

		if len(file) > 0 {
			fileBytes, _ := ioutil.ReadFile(file)

			jsonBytes, cerr := common.WantJSON(fileBytes)
			if cerr != nil {
				logrus.Fatal("The input data MUST be in either JSON or YAML format")
			}
			if len(jsonBytes) > 0 {
				marshErr := json.Unmarshal(jsonBytes, &dbReq)
				if marshErr != nil {
					logrus.Fatal("Could not unmarshall data", marshErr)
				}
			}
			fileProvided = true
		}

		populateRequest(cmd, &dbReq, fileProvided)

		if skel {
			generateSkel(&dbReq)
			os.Exit(0)
		}

		out, err := createSpliceDatabase(&dbReq, false)
		if err != nil {
			logrus.WithError(err).Error("Error Generating Default CR Info")
		}

		if semverV1, err := semver.ParseRange(">=0.1.7"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV1(sv) {
				displayCreateSpliceDatabaseV1(out)
			}
		}
	},
}

func displayCreateSpliceDatabaseV1(in string) {
	fmt.Println(in)
	os.Exit(0)
}

func populateRequest(cmd *cobra.Command, req *objects.DatabaseRequest, fileData bool) {

	requiredList := []string{}

	databaseName := common.DatabaseName(cmd)
	password, _ := cmd.Flags().GetString("password")
	accountID, _ := cmd.Flags().GetString("account-id")
	authorizationCode, _ := cmd.Flags().GetString("authorization-code")
	backupFrequency, _ := cmd.Flags().GetString("backup-frequency")
	backupInterval, _ := cmd.Flags().GetInt("backup-interval")
	keepBackups, _ := cmd.Flags().GetInt("keep-backups")
	backupStartWindow, _ := cmd.Flags().GetString("backup-start-window")
	cloudProvider, _ := cmd.Flags().GetString("cloud-provider")
	sparkExecutors, _ := cmd.Flags().GetInt("spark-executors")
	regionServers, _ := cmd.Flags().GetInt("region-servers")
	dedicatedStorage, _ := cmd.Flags().GetBool("dedicated-storage")
	externalDatasetSize, _ := cmd.Flags().GetInt("external-dataset-size")
	internalDatasetSize, _ := cmd.Flags().GetInt("internal-dataset-size")
	enableMLManager, _ := cmd.Flags().GetBool("enable-mlmanager")
	notebookActiveUsers, _ := cmd.Flags().GetInt("notebook-active-users")
	notebookExecutors, _ := cmd.Flags().GetInt("notebook-executors")
	notebookTotalUsers, _ := cmd.Flags().GetInt("notebook-total-users")
	notebooksPerUser, _ := cmd.Flags().GetInt("notebooks-per-user")

	if cmd.Flags().Changed("database-name") || !fileData {
		req.Name = databaseName
	}
	if cmd.Flags().Changed("password") || !fileData {
		req.Password = password
	}
	if cmd.Flags().Changed("account-id") {
		req.AccountID = accountID
	} else {
		if len(req.AccountID) == 0 || !fileData {
			selectedAccountID, err := PromptForAccountID()
			if err != nil {
				requiredList = append(requiredList, "--account-id, obtain from 'splicectl get accounts', or select from list")
			}
			req.AccountID = selectedAccountID
		}
	}
	if cmd.Flags().Changed("authorization-code") || !fileData {
		req.AuthorizationCode = authorizationCode
	}
	if cmd.Flags().Changed("backup-frequency") || !fileData {
		req.BackupFrequency = backupFrequency
	}
	if cmd.Flags().Changed("backup-interval") || !fileData {
		req.BackupInterval = backupInterval
	}
	if cmd.Flags().Changed("keep-backups") || !fileData {
		req.BackupKeepCount = keepBackups
	}
	if cmd.Flags().Changed("backup-start-window") || !fileData {
		req.BackupStartWindow = backupStartWindow
	}
	if cmd.Flags().Changed("cloud-provider") {
		req.CloudProvider = cloudProvider
	} else {
		if len(req.CloudProvider) == 0 || !fileData {
			selectedCSP, err := PromptForCSP()
			if err != nil {
				requiredList = append(requiredList, "--cloud-provider, (aws|az|gcp|op|none)")
			}
			req.CloudProvider = selectedCSP
		}
	}
	if cmd.Flags().Changed("spark-executors") || !fileData {
		req.ClusterPowerOlap = sparkExecutors
	}
	if cmd.Flags().Changed("region-servers") || !fileData {
		req.ClusterPowerOltp = regionServers
	}
	if cmd.Flags().Changed("dedicated-storage") || !fileData {
		req.DedicatedStorage = dedicatedStorage
	}
	if cmd.Flags().Changed("external-dataset-size") || !fileData {
		req.ExternalDatasetSizeGb = externalDatasetSize
	}
	if cmd.Flags().Changed("internal-dataset-size") || !fileData {
		req.InternalDatasetSizeGb = internalDatasetSize
	}
	if cmd.Flags().Changed("enable-mlmanager") || !fileData {
		req.MlManager = enableMLManager
	}
	if cmd.Flags().Changed("notebook-active-users") || !fileData {
		req.NotebookActiveUsers = notebookActiveUsers
	}
	if cmd.Flags().Changed("notebook-executors") || !fileData {
		req.NotebookExecutorsPerNotebook = notebookExecutors
	}
	if cmd.Flags().Changed("notebook-total-users") || !fileData {
		req.NotebookTotalUsers = notebookTotalUsers
	}
	if cmd.Flags().Changed("notebooks-per-user") || !fileData {
		req.NotebooksPerUser = notebooksPerUser
	}

	req.CloudProvider = strings.ToUpper(req.CloudProvider)

	if len(requiredList) > 0 {
		for _, v := range requiredList {
			logrus.Warn(fmt.Sprintf("Required parameter not provided: %s", v))
		}
		os.Exit(1)
	}
}

func generateSkel(dbReq *objects.DatabaseRequest) {

	if !FormatOverridden {
		OutputFormat = "yaml"
	}

	switch strings.ToLower(OutputFormat) {
	case "json", "gron":
		dbReq.ToJSON()
	case "yaml", "text", "table":
		dbReq.ToYAML()
	}

}
func createSpliceDatabase(dbReq *objects.DatabaseRequest, outputonly bool) (string, error) {

	restClient := resty.New()

	uri := "splicectl/v1/splicedb/splicedatabase"

	var resp *resty.Response
	var resperr error

	if outputonly {
		resp, resperr = restClient.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Accept", "application/json").
			SetHeader("X-Token-Bearer", AuthClient.GetTokenBearer()).
			SetHeader("X-Token-Session", AuthClient.GetSessionID()).
			Get(fmt.Sprintf("%s/%s", ApiServer, uri))
	} else {

		reqJSON, enverr := json.MarshalIndent(dbReq, "", "  ")
		if enverr != nil {
			logrus.WithError(enverr).Error("Error extracting json")
			return "", enverr
		}
		resp, resperr = restClient.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Accept", "application/json").
			SetHeader("X-Token-Bearer", AuthClient.GetTokenBearer()).
			SetHeader("X-Token-Session", AuthClient.GetSessionID()).
			SetBody(reqJSON).
			Post(fmt.Sprintf("%s/%s", ApiServer, uri))
	}
	if resperr != nil {
		logrus.WithError(resperr).Error("Error getting Default CR Info")
		return "", resperr
	}

	return string(resp.Body()[:]), nil

}

func init() {
	createCmd.AddCommand(createDatabaseCmd)

	createDatabaseCmd.Flags().BoolP("skel", "s", false, "Generate a skeleton values file for submission")
	createDatabaseCmd.Flags().StringP("file", "f", "", "Specify the input file")

	// add database name and aliases
	createDatabaseCmd.Flags().StringP("database-name", "d", "", "Specify the database name")
	createDatabaseCmd.Flags().String("database", "", "Alias for database-name, prefer the use of -d and --database-name.")
	createDatabaseCmd.Flags().String("workspace", "", "Alias for database-name, prefer the use of -d and --database-name.")

	createDatabaseCmd.Flags().String("password", "admin", "Specify the Splice Machine Database Password (default=admin)")
	createDatabaseCmd.Flags().String("account-id", "", "Specify the Cloud Manager Account ID to associate the database to")
	createDatabaseCmd.Flags().String("authorization-code", "", "Specify the Authorization Code")
	createDatabaseCmd.Flags().String("backup-frequency", "daily", "Specify the Backup Frequency (default=daily)")
	createDatabaseCmd.Flags().Int("backup-interval", 1, "Specify the Backup Interval (default=1)")
	createDatabaseCmd.Flags().Int("keep-backups", 1, "Specify the Backup Keep Count (default=1)")
	createDatabaseCmd.Flags().String("backup-start-window", "02:30", "Specify the Backup Start Window (default=02:30)")
	createDatabaseCmd.Flags().String("cloud-provider", "", "Specify the Cloud Provider (az|aws|gcp|op|none)")
	createDatabaseCmd.Flags().Int("spark-executors", 4, "Specify the number of Spark Executors/OLAP (default=4)")
	createDatabaseCmd.Flags().Int("region-servers", 4, "Specify the number of Region Servers/OLTP (default=4)")
	createDatabaseCmd.Flags().Bool("dedicated-storage", false, "Specify if dedicated storage should be used")
	createDatabaseCmd.Flags().Int("external-dataset-size", 0, "Specify the size (GB) of the external storage (default=0)")
	createDatabaseCmd.Flags().Int("internal-dataset-size", 1, "Specify the size (GB) of the internal storage (default=1)")
	createDatabaseCmd.Flags().Bool("enable-mlmanager", false, "Enable the ML Manager features of the database (default=false)")
	createDatabaseCmd.Flags().Int("notebook-active-users", 4, "Specify the max number of active Jupyter notebook sessions (default=4)")
	createDatabaseCmd.Flags().Int("notebook-executors", 2, "Specify the max number Spark Executors per notebook (default=2)")
	createDatabaseCmd.Flags().Int("notebook-total-users", 10, "Specify the max notebook users (default=10)")
	createDatabaseCmd.Flags().Int("notebooks-per-user", 2, "Specify the max number of notebooks per user (default=2)")

}
