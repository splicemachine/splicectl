package get

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/sirupsen/logrus"
	c "github.com/splicemachine/splicectl/cmd"
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
		_, sv := c.VersionDetail.RequirementMet("get_image-tag")

		componentName, _ := cmd.Flags().GetString("component-name")
		databaseName, _ := cmd.Flags().GetString("database-name")
		if len(databaseName) == 0 {
			databaseName, dberr = c.PromptForDatabaseName()
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
	if strings.ToLower(c.OutputFormat) == "raw" {
		fmt.Println(in)
		os.Exit(0)
	}

	var tags []objects.ImageTag

	marshErr := json.Unmarshal([]byte(in), &tags)
	if marshErr != nil {
		logrus.Fatal("Could not unmarshall data", marshErr)
	}

	tagList := objects.ImageTagList{
		ImageTags: tags,
	}

	c.OutputData(&tagList)

}

func getImageTagData(componenetName string, databaseName string) (string, error) {
	uri := fmt.Sprintf("splicectl/v1/splicedb/imagetag?component-name=%s&database-name=%s", componenetName, databaseName)
	resp, resperr := c.RestyWithHeaders().
		Get(fmt.Sprintf("%s/%s", c.ApiServer, uri))

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
