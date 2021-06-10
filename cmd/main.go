package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/splicemachine/splicectl/auth"
	"github.com/splicemachine/splicectl/cmd/objects"
	"github.com/splicemachine/splicectl/common"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	semVer    string
	gitCommit string
	buildDate string
	gitRef    string
	cfgFile   string
	serverURI string

	VersionDetail objects.Version
	VersionJSON   string

	ApiServer        string
	OutputFormat     string
	FormatOverridden bool
	NoHeaders        bool
	AuthClient       auth.Client
)

// RootCmd represents the base command when called without any subcommands
// splicectl doesn't have any functionality, other than to validate our auth
// token.  A VALID auth token is required to run ANY command other than the
// 'auth' command.
var RootCmd = &cobra.Command{
	Use:   "splicectl",
	Short: "Splice Machine control application for Kubernetes environments",
	Long: `splicectl is a CLI tool for making managment of Splice Machine
database clusters under Kubernetes easier to manage.`,
	Args: cobra.MinimumNArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		ApiServer = getIngressDetail()
		if len(serverURI) > 0 {
			ApiServer = serverURI
		}

		// Collect the version info, for use in determining valid commands based on SemVer
		if ApiServer != "" {
			version, err := GetVersionInfo()
			if err != nil {
				logrus.WithError(err).Error("Error getting version info")
			}
			clientLine := fmt.Sprintf("\"Client\": {\"SemVer\": \"%s\", \"GitCommit\": \"%s\", \"BuildDate\": \"%s\"},", semVer, gitCommit, buildDate)
			serverLine := fmt.Sprintf("\"Server\": %s},", version)
			hostLine := fmt.Sprintf("\"Host\": \"%s\"", ApiServer)
			VersionJSON = fmt.Sprintf("{\"VersionInfo\": {\n%s\n%s\n%s\n}", clientLine, serverLine, hostLine)
		} else {
			clientLine := fmt.Sprintf("\"Client\": {\"SemVer\": \"%s\", \"GitCommit\": \"%s\", \"BuildDate\": \"%s\"}}", semVer, gitCommit, buildDate)
			VersionJSON = fmt.Sprintf("{\"VersionInfo\": {%s}", clientLine)
		}

		if err := json.Unmarshal([]byte(VersionJSON), &VersionDetail); err != nil {
			logrus.WithError(err).Error("Error decoding json for Version")
		}

		if os.Args[1] != "version" {
			environment := getEnvironmentName()
			AuthClient = auth.NewAuth(environment, common.SessionData{
				SessionID:  fmt.Sprintf("%s", viper.Get(fmt.Sprintf("%s-session_id", environment))),
				ValidUntil: fmt.Sprintf("%s", viper.Get(fmt.Sprintf("%s-valid_until", environment))),
			})
			isValid := AuthClient.CheckTokenValidity()
			if !isValid && os.Args[1] != "auth" {
				logrus.Info("Your session has expired, please run the 'auth' again.")
				os.Exit(1)
			}
		}

		// Validate global parameters here, BEFORE we start to waste time
		// and run any code.
		if OutputFormat != "" {
			OutputFormat = strings.ToLower(OutputFormat)
			switch OutputFormat {
			case "json", "gron", "yaml", "text", "table", "raw":
				break
			default:
				fmt.Println("Valid options for -o are [json|gron|[text|table]|yaml|raw]")
				os.Exit(1)
			}
			FormatOverridden = true
		} else {
			FormatOverridden = false
			OutputFormat = "json"
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.splicectl/config.yml)")
	RootCmd.PersistentFlags().StringVar(&serverURI, "server-uri", "", "override the server uri for the API server http(s)://host.domain.name:overrideport")
	RootCmd.PersistentFlags().StringVarP(&OutputFormat, "output", "o", "", "output types: json, text, yaml, gron")
	RootCmd.PersistentFlags().BoolVar(&NoHeaders, "no-headers", false, "Suppress header output in Text output")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		if _, err := os.Stat(cfgFile); err != nil {
			if os.IsNotExist(err) {
				if os.Args[1] != "auth" {
					logrus.Info("Couldn't read the config file.  We require a session ID from the splicectl API.  Please run with 'auth'.")
					os.Exit(1)
				} else {
					createRestrictedConfigFile(cfgFile)
				}
			}
		}
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		directory := fmt.Sprintf("%s/%s", home, ".splicectl")
		if _, err := os.Stat(directory); err != nil {
			if os.IsNotExist(err) {
				os.Mkdir(directory, os.ModePerm)
			}
		}
		if stat, err := os.Stat(directory); err == nil && stat.IsDir() {
			configFile := fmt.Sprintf("%s/%s", home, ".splicectl/config.yml")
			createRestrictedConfigFile(configFile)
			viper.SetConfigFile(configFile)
		} else {
			logrus.Info("The ~/.splicectl path is a file and not a directory, please remove the .splicectl file.")
			os.Exit(1)
		}
	}

	// viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if os.Args[1] != "auth" {
			logrus.Info("Couldn't read the config file.  We require a session ID from the splicectl API.  Please run with 'auth'.")
			os.Exit(1)
		}
	}
}

func createRestrictedConfigFile(fileName string) {
	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			file, ferr := os.Create(fileName)
			if ferr != nil {
				logrus.Info("Unable to create the configfile.")
				os.Exit(1)
			}
			mode := int(0600)
			if cherr := file.Chmod(os.FileMode(mode)); cherr != nil {
				logrus.Info("Chmod for config file failed, please set the mode to 0600.")
			}
		}
	}
}

// GetDatabaseList - gets a list of databases
func GetDatabaseList() (string, error) {
	uri := "splicectl/v1/splicedb/splicedatabase"
	resp, resperr := RestyWithHeaders().
		Execute("LIST", fmt.Sprintf("%s/%s", ApiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error getting Database List")
		return "", resperr
	}

	return string(resp.Body()[:]), nil

}

// GetAccounts - get list of accounts
func GetAccounts() (string, error) {
	uri := "splicectl/v1/cm/accounts"
	resp, resperr := RestyWithHeaders().
		Get(fmt.Sprintf("%s/%s", ApiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error getting Account List Info")
		return "", resperr
	}

	return string(resp.Body()[:]), nil
}

// GetVersionInfo - gets version information
func GetVersionInfo() (string, error) {
	uri := "splicectl"
	resp, resperr := RestyWithHeaders().
		Get(fmt.Sprintf("%s/%s", ApiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error getting version info")
		return "", resperr
	}

	return strings.TrimSuffix(string(resp.Body()[:]), "\n"), nil
}

// RestyWithHeaders - new resty request with headers for auth and content-type.
func RestyWithHeaders() *resty.Request {
	return resty.
		New().
		R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Token-Bearer", AuthClient.GetTokenBearer()).
		SetHeader("X-Token-Session", AuthClient.GetSessionID())
}

// Outputable - defines ways that an object may need to present itself
type Outputable interface {
	ToJSON() string
	ToYAML() string
	ToGRON() string
	ToText(noHeaders bool) string
}

func outputData(data Outputable) string {
	switch strings.ToLower(OutputFormat) {
	case "json":
		return data.ToJSON()
	case "gron":
		return data.ToGRON()
	case "yaml":
		return data.ToYAML()
	case "text", "table":
		return data.ToText(NoHeaders)
	default:
		return ""
	}
}

// OutputData - outputs string representation of data in accordance with
// OutputFormat.
func OutputData(data Outputable) {
	if !FormatOverridden {
		OutputFormat = "text"
	}

	fmt.Println(outputData(data))
}
