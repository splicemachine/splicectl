package cmd

import (
	"encoding/json"
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
		_, semVerNum := ClientSemVer()
		url := fmt.Sprintf("https://api.github.com/repos/splicemachine/splicectl/releases/tags/%s", semVerNum)
		client := resty.New()
		resp, err := client.R().Get(url)
		if err != nil {
			logrus.WithError(err).Fatal("request could not be completed due to error, changelog is retrieved through api call to Github that requires internet access")
		}
		jsonMap := make(map[string]interface{})
		if err := json.Unmarshal(resp.Body(), &jsonMap); err != nil {
			logrus.WithError(err).Fatal("could not get changelog, recieved error while parsing gh api response")
		}
		if _, ok := jsonMap["body"]; !ok {
			logrus.Fatal("an unexpected response was returned from the Github api, the changelog cannot be found")
		}
		// using fmt here instead of logrus/log in order to get rid of leading [INFO] or similar header
		fmt.Println(jsonMap["body"])
	},
}

func init() {
	rootCmd.AddCommand(changelogCmd)
}
