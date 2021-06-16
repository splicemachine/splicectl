package apply

import (
	"fmt"
	"os"

	"github.com/blang/semver/v4"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var applyImageTagCmd = &cobra.Command{
	Use:   "image-tag",
	Short: "Apply an image tag to a cluster resource",
	Long: `EXAMPLES
	splicectl apply image-tag --database-name cjdb --component-name zookeeper --tag master-0.0.4

	Supported component-name(s):
		allspark
		hbase
		hdfs
		kafka
		zookeeper
`,
	Run: func(cmd *cobra.Command, args []string) {
		var dberr error
		_, sv := c.VersionDetail.RequirementMet("apply_image-tag")

		componentName, _ := cmd.Flags().GetString("component-name")
		databaseName, _ := cmd.Flags().GetString("database-name")
		if len(databaseName) == 0 {
			databaseName, dberr = c.PromptForDatabaseName()
			if dberr != nil {
				logrus.Fatal("Could not get a list of Databases", dberr)
			}
		}

		tag, _ := cmd.Flags().GetString("tag")
		out, err := setDatabaseImageTag(componentName, databaseName, tag)
		if err != nil {
			logrus.WithError(err).Error("Error getting image tag for component")
		}

		if semverV1, err := semver.ParseRange(">=0.0.16"); err != nil {
			logrus.Fatal("Failed to parse SemVer")
		} else {
			if semverV1(sv) {
				displayApplyImageTagV1(out)
			}
		}
	},
}

func displayApplyImageTagV1(in string) {
	fmt.Println(in)
	os.Exit(0)
}

func setDatabaseImageTag(componentName string, databaseName string, imageTag string) (string, error) {
	uri := fmt.Sprintf("splicectl/v1/splicedb/imagetag?component-name=%s&database-name=%s&tag=%s",
		componentName, databaseName, imageTag)
	resp, resperr := c.RestyWithHeaders().
		Post(fmt.Sprintf("%s/%s", c.ApiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error setting TAG for database")
		return "", resperr
	}

	return string(resp.Body()[:]), nil
}

func init() {
	applyCmd.AddCommand(applyImageTagCmd)

	applyImageTagCmd.Flags().StringP("component-name", "c", "", "Specify the component")
	applyImageTagCmd.Flags().StringP("database-name", "d", "", "Specify the database name")
	applyImageTagCmd.Flags().StringP("tag", "t", "", "Specify the image tag, ie: master-246")

	applyImageTagCmd.MarkFlagRequired("component-name")
	applyImageTagCmd.MarkFlagRequired("tag")

}
