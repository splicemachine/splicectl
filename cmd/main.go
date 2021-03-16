package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/splicemachine/splicectl/auth"
	"github.com/splicemachine/splicectl/cmd/objects"
	"github.com/splicemachine/splicectl/common"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	semVer        string
	gitCommit     string
	buildDate     string
	gitRef        string
	versionDetail objects.Version
	versionJSON   string
	cfgFile       string
	serverURI     string
)

// var sessionID string
// var tokenBearer string
// var tokenValid bool
var apiServer string
var outputFormat string
var formatOverridden bool
var noHeaders bool
var authClient auth.Client

// rootCmd represents the base command when called without any subcommands
// splicectl doesn't have any functionality, other than to validate our auth
// token.  A VALID auth token is required to run ANY command other than the
// 'auth' command.
var rootCmd = &cobra.Command{
	Use:   "splicectl",
	Short: "Splice Machine control application for Kubernetes environments",
	Long: `splicectl is a CLI tool for making managment of Splice Machine
database clusters under Kubernetes easier to manage.`,
	Args: cobra.MinimumNArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		apiServer = getIngressDetail()
		if len(serverURI) > 0 {
			apiServer = serverURI
		}

		// Collect the version info, for use in determining valid commands based on SemVer
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
		marsherr := json.Unmarshal([]byte(versionJSON), &versionDetail)
		if marsherr != nil {
			logrus.WithError(marsherr).Error("Error decoding json for Version")
		}

		if os.Args[1] != "version" {
			environment := getEnvironmentName()
			authClient = auth.NewAuth(environment, common.SessionData{
				SessionID:  fmt.Sprintf("%s", viper.Get(fmt.Sprintf("%s-session_id", environment))),
				ValidUntil: fmt.Sprintf("%s", viper.Get(fmt.Sprintf("%s-valid_until", environment))),
			})
			isValid := authClient.CheckTokenValidity()
			if !isValid && os.Args[1] != "auth" {
				logrus.Info("Your session has expired, please run the 'auth' again.")
				os.Exit(1)
			}
		}

		// Validate global parameters here, BEFORE we start to waste time
		// and run any code.
		if outputFormat != "" {
			outputFormat = strings.ToLower(outputFormat)
			switch outputFormat {
			case "json":
			case "gron":
			case "yaml":
			case "text":
			case "table":
			case "raw":
			default:
				fmt.Println("Valid options for -o are [json|gron|[text|table]|yaml|raw]")
				os.Exit(1)
			}
			formatOverridden = true
		} else {
			formatOverridden = false
			outputFormat = "json"
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.splicectl/config.yml)")
	rootCmd.PersistentFlags().StringVar(&serverURI, "server-uri", "", "override the server uri for the API server http(s)://host.domain.name:overrideport")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "", "output types: json, text, yaml, gron")
	rootCmd.PersistentFlags().BoolVar(&noHeaders, "no-headers", false, "Suppress header output in Text output")
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
