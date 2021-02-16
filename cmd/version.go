package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/splicemachine/splicectl/cmd/objects"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Express the 'version' of splicectl.",
	Aliases: []string{"v"},
	Run: func(cmd *cobra.Command, args []string) {
		versionData := objects.Version{}
		versionJSON := ""
		if apiServer != "" {
			version, err := getVersionInfo()
			if err != nil {
				logrus.WithError(err).Error("Error getting version info")
			}
			clientLine := fmt.Sprintf("\"Client\": {\"SemVer\": \"%s\", \"GitCommit\": \"%s\", \"BuildDate\": \"%s\"},", semVer, gitCommit, buildDate)
			serverLine := fmt.Sprintf("\"Server\": %s},", version)
			hostLine := fmt.Sprintf("\"Host\": \"%s\"", apiServer)
			versionJSON = fmt.Sprintf("{\"VersionInfo\": {\n%s\n%s\n%s\n}", clientLine, serverLine, hostLine)
		} else {
			clientLine := fmt.Sprintf("\"Client\": {\"SemVer\": \"%s\", \"GitCommit\": \"%s\", \"BuildDate\": \"%s\"}}", semVer, gitCommit, buildDate)
			versionJSON = fmt.Sprintf("{\"VersionInfo\": {%s}", clientLine)
		}
		marsherr := json.Unmarshal([]byte(versionJSON), &versionData)
		if marsherr != nil {
			logrus.WithError(marsherr).Error("Error decoding json for Version")
		}

		if !formatOverridden {
			outputFormat = "yaml"
		}

		switch strings.ToLower(outputFormat) {
		case "json":
			// We want to print the JSON in a condensed format
			fmt.Println(versionJSON)
		case "gron":
			versionData.ToGRON()
		case "yaml":
			versionData.ToYAML()
		case "text", "table":
			versionData.ToTEXT(noHeaders)
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
