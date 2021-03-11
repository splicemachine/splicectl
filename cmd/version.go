package cmd

import (
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Express the 'version' of splicectl.",
	Aliases: []string{"v"},
	Run: func(cmd *cobra.Command, args []string) {

		if !formatOverridden {
			outputFormat = "yaml"
		}

		switch strings.ToLower(outputFormat) {
		case "raw":
			fmt.Println(versionJSON)
		case "json":
			// We want to print the JSON in a condensed format
			fmt.Println(versionJSON)
		case "gron":
			versionDetail.ToGRON()
		case "yaml":
			versionDetail.ToYAML()
		case "text", "table":
			versionDetail.ToTEXT(noHeaders)
		}
	},
}

func getVersionInfo() (string, error) {
	restClient := resty.New()

	uri := "splicectl"
	resp, resperr := restClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		// SetHeader("X-Token-Bearer", authClient.GetTokenBearer()).
		// SetHeader("X-Token-Session", authClient.GetSessionID()).
		Get(fmt.Sprintf("%s/%s", apiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error getting version info")
		return "", resperr
	}

	return strings.TrimSuffix(string(resp.Body()[:]), "\n"), nil

}

func init() {
	rootCmd.AddCommand(versionCmd)
}
