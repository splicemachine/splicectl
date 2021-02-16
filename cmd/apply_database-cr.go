package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/splicemachine/splicectl/cmd/objects"
	"github.com/splicemachine/splicectl/common"

	"github.com/spf13/cobra"
)

var applyDatabaseCRCmd = &cobra.Command{
	Use:   "database-cr",
	Short: "Submit a new CR for a specified database in the cluster.",
	Long: `EXAMPLES
	splicectl list database
    splicectl get database-cr --database-name splicedb -o json > ~/tmp/splicedb.json
    # edit the file
    splicectl apply database-cr --database-name splicedb --file ~/tmp/splicedb.json
`,
	Run: func(cmd *cobra.Command, args []string) {
		var dberr error
		databaseName, _ := cmd.Flags().GetString("database-name")
		if len(databaseName) == 0 {
			databaseName, dberr = promptForDatabaseName()
			if dberr != nil {
				logrus.Fatal("Could not get a list of Databases", dberr)
			}
		}
		filePath, _ := cmd.Flags().GetString("file")
		fileBytes, _ := ioutil.ReadFile(filePath)

		jsonBytes, cerr := common.WantJSON(fileBytes)
		if cerr != nil {
			logrus.Fatal("The input data MUST be in either JSON or YAML format")
		}

		out, err := setDatabaseCR(databaseName, jsonBytes)
		if err != nil {
			logrus.WithError(err).Error("Error setting Database CR Info")
		}
		var vvData objects.VaultVersion
		marshErr := json.Unmarshal([]byte(out), &vvData)
		if marshErr != nil {
			logrus.Fatal("Could not unmarshall data", marshErr)
		}

		if !formatOverridden {
			outputFormat = "text"
		}

		switch strings.ToLower(outputFormat) {
		case "json":
			vvData.ToJSON()
		case "gron":
			vvData.ToGRON()
		case "yaml":
			vvData.ToYAML()
		case "text", "table":
			vvData.ToTEXT(noHeaders)
		}

	},
}

func setDatabaseCR(dbname string, in []byte) (string, error) {
	restClient := resty.New()
	uri := fmt.Sprintf("splicectl/v1/vault/databasecr?database-name=%s", dbname)
	resp, resperr := restClient.R().
		SetHeader("X-Token-Bearer", authClient.GetTokenBearer()).
		SetHeader("X-Token-Session", authClient.GetSessionID()).
		SetBody(in).
		SetResult(&AuthSuccess{}). // or SetResult(AuthSuccess{}).
		SetError(&AuthError{}).    // or SetError(AuthError{}).
		Post(fmt.Sprintf("%s/%s", apiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error setting Database CR Info")
		return "", resperr
	}

	return string(resp.Body()[:]), nil

}

func init() {
	applyCmd.AddCommand(applyDatabaseCRCmd)

	applyDatabaseCRCmd.Flags().StringP("database-name", "d", "", "Specify the database name")
	applyDatabaseCRCmd.Flags().StringP("file", "f", "", "Specify the input file")
	// applyDatabaseCRCmd.MarkFlagRequired("database-name")
	applyDatabaseCRCmd.MarkFlagRequired("file")
}
