package cmd

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var changelogCmd = &cobra.Command{
	Use:   "changelog",
	Short: "List the most recent changes on the changelog for splicectl.",
	Long: `EXAMPLES
	splicectl changelog
`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: output the changelog var...
		logrus.Info(versionDetail.VersionInfo.Client.SemVer)
		url := fmt.Sprintf("https://api.github.com/repos/splicemachine/splicectl/releases/tags/%s", versionDetail.VersionInfo.Client.SemVer)
		client := resty.New()
		resp, err := client.R().Get(url)
		if err != nil {
			logrus.WithError(err).Fatal("request could not be completed due to error, changelog is retrieved through api call to Github that requires internet access")
		}
		logrus.Info(string(resp.Body()))
	},
}

func init() {
	rootCmd.AddCommand(changelogCmd)
}
