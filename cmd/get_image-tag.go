package cmd

import (
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

var getImageTag = &cobra.Command{
	Use:   "image-tag",
	Short: "Get the image tag for a component of a database.",
	Long: `EXAMPLES
	splicectl get image-tag --component-name "hbase" --database-name "cjdb"
`,
	Run: func(cmd *cobra.Command, args []string) {
		var dberr error
		var sv semver.Version

		_, sv = versionDetail.RequirementMet("get_image-tag")

		componentName, _ := cmd.Flags().GetString("component-name")
		databaseName, _ := cmd.Flags().GetString("database-name")
		if len(databaseName) == 0 {
			databaseName, dberr = promptForDatabaseName()
			if dberr != nil {
				logrus.Fatal("Could not get a list of Databases", dberr)
			}
		}

		out, err := getImageTagData(componentName, databaseName)
		if err != nil {
			logrus.WithError(err).Error("Error getting image tag for component")
		}

		if semverV1, err := semver.ParseRange(">=0.0.16 <0.0.17"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV1(sv) {
				displayGetImageTagV1(out)
			}
		}

		if semverV2, err := semver.ParseRange(">=0.0.17"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV2(sv) {
				displayGetImageTagV2(out)
			}
		}
	},
}

func displayGetImageTagV1(in string) {
	fmt.Println(in)
	os.Exit(0)
}

func displayGetImageTagV2(in string) {
	if strings.ToLower(outputFormat) == "raw" {
		fmt.Println(in)
		os.Exit(0)
	}

	var tags []objects.ImageTag

	marshErr := json.Unmarshal([]byte(in), &tags)
	if marshErr != nil {
		logrus.Fatal("Could not unmarshall data", marshErr)
	}

	if !formatOverridden {
		outputFormat = "table"
	}

	tagList := objects.ImageTagList{
		ImageTags: tags,
	}

	switch strings.ToLower(outputFormat) {
	case "json":
		tagList.ToJSON()
	case "gron":
		tagList.ToGRON()
	case "yaml":
		tagList.ToYAML()
	case "text", "table":
		tagList.ToTEXT(noHeaders)
	}

}

func getImageTagData(componenetName string, databaseName string) (string, error) {

	restClient := resty.New()

	uri := fmt.Sprintf("splicectl/v1/splicedb/imagetag?component-name=%s&database-name=%s", componenetName, databaseName)
	resp, resperr := restClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Token-Bearer", authClient.GetTokenBearer()).
		SetHeader("X-Token-Session", authClient.GetSessionID()).
		Get(fmt.Sprintf("%s/%s", apiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error getting Default CR Info")
		return "", resperr
	}
	return string(resp.Body()[:]), nil
}

func init() {
	getCmd.AddCommand(getImageTag)

	getImageTag.Flags().StringP("component-name", "c", "", "Specify the component")
	getImageTag.Flags().StringP("database-name", "d", "", "Specify the database name")

	getImageTag.MarkFlagRequired("component-name")
	// getImageTag.MarkFlagRequired("database-name")

}
