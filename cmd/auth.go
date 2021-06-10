package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/viper"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// SessionData - Session Authorization Info
type SessionData struct {
	SessionID  string `json:"session_id"`
	ValidUntil string `json:"valid_until"`
}

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Request an auth session",
	Long: `EXAMPLES
	splicectl auth`,
	Aliases: []string{"login"},
	Run: func(cmd *cobra.Command, args []string) {

		// var authClient auth.Client
		// environment := getEnvironmentName()
		// authClient = auth.NewAuth(environment, common.SessionData{
		// 	SessionID:  fmt.Sprintf("%s", viper.Get(fmt.Sprintf("%s-session_id", environment))),
		// 	ValidUntil: fmt.Sprintf("%s", viper.Get(fmt.Sprintf("%s-valid_until", environment))),
		// })
		if pass := AuthClient.CheckTokenValidity(); pass {
			sessData, err := json.Marshal(AuthClient.GetSession())
			if err != nil {
				logrus.WithError(err).Error("Error converting session data to JSON")
				return
			}
			fmt.Println(string(sessData[:]))
		} else {
			out, err := performAuth()
			if err != nil {
				logrus.WithError(err).Error("Error getting AUTH Info")
			}
			var response SessionData
			marsherr := json.Unmarshal([]byte(out), &response)
			if marsherr != nil {
				logrus.WithError(marsherr).Error("Error decoding json")
			}
			environment := getEnvironmentName()
			viper.Set(fmt.Sprintf("%s-session_id", environment), response.SessionID)
			viper.Set(fmt.Sprintf("%s-valid_until", environment), response.ValidUntil)
			verr := viper.WriteConfig()
			if verr != nil {
				logrus.WithError(verr).Info("Failed to write config")
			}
			fmt.Println(out)
		}

	},
}

func performAuth() (string, error) {
	restClient := resty.New()

	uri := "splicectl/v1/auth"
	resp, resperr := restClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		Get(fmt.Sprintf("%s/%s", ApiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error getting Default CR Info")
		return "", resperr
	}

	return string(resp.Body()[:]), nil

}

func init() {
	RootCmd.AddCommand(authCmd)
}
