package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/splicemachine/splicectl/auth"
	"github.com/splicemachine/splicectl/cmd/apply"
	"github.com/splicemachine/splicectl/cmd/config"
	"github.com/splicemachine/splicectl/cmd/create"
	"github.com/splicemachine/splicectl/cmd/del"
	"github.com/splicemachine/splicectl/cmd/get"
	"github.com/splicemachine/splicectl/cmd/list"
	"github.com/splicemachine/splicectl/cmd/restart"
	"github.com/splicemachine/splicectl/cmd/rollback"
	"github.com/splicemachine/splicectl/cmd/version"
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

	// semVerReg - gets the semVer portion only, cutting off any other release details
	semVerReg = regexp.MustCompile(`(v[0-9]+\.[0-9]+\.[0-9]+).*`)

	c = &config.Config{}
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
		if len(c.CACert) > 0 {
			if _, err := os.Stat(c.CACert); err != nil {
				if os.IsNotExist(err) {
					logrus.Info("Couldn't read the ca-file, please check the path")
					os.Exit(1)
				}
			}
			fileBytes, _ := ioutil.ReadFile(c.CACert)
			c.CABundle = strings.TrimSpace(string(fileBytes[:]))
		} else {
			c.CACert = os.Getenv("SPLICECTL_CACERT")
			if len(c.CACert) > 0 {
				if _, err := os.Stat(c.CACert); err != nil {
					if os.IsNotExist(err) {
						logrus.Info("Couldn't read the ca-file, please check the path")
						os.Exit(1)
					}
				}
			}
			fileBytes, _ := ioutil.ReadFile(c.CACert)
			c.CABundle = strings.TrimSpace(string(fileBytes[:]))
		}

		c.ApiServer = getIngressDetail()
		if len(serverURI) > 0 {
			c.ApiServer = serverURI
		}

		// Collect the version info, for use in determining valid commands based on SemVer
		if c.ApiServer != "" {
			version, err := c.GetVersionInfo()
			if err != nil {
				logrus.WithError(err).Error("Error getting version info")
			}
			clientLine := fmt.Sprintf("\"Client\": {\"SemVer\": \"%s\", \"GitCommit\": \"%s\", \"BuildDate\": \"%s\"},", semVer, gitCommit, buildDate)
			serverLine := fmt.Sprintf("\"Server\": %s},", version)
			hostLine := fmt.Sprintf("\"Host\": \"%s\"", c.ApiServer)
			c.VersionJSON = fmt.Sprintf("{\"VersionInfo\": {\n%s\n%s\n%s\n}", clientLine, serverLine, hostLine)
		} else {
			clientLine := fmt.Sprintf("\"Client\": {\"SemVer\": \"%s\", \"GitCommit\": \"%s\", \"BuildDate\": \"%s\"}}", semVer, gitCommit, buildDate)
			c.VersionJSON = fmt.Sprintf("{\"VersionInfo\": {%s}", clientLine)
		}

		if err := json.Unmarshal([]byte(c.VersionJSON), &c.VersionDetail); err != nil {
			logrus.WithError(err).Error("Error decoding json for Version")
		}

		if os.Args[1] != "version" {
			environment := getEnvironmentName()
			c.AuthClient = auth.NewAuth(environment, common.SessionData{
				SessionID:  fmt.Sprintf("%s", viper.Get(fmt.Sprintf("%s-session_id", environment))),
				ValidUntil: fmt.Sprintf("%s", viper.Get(fmt.Sprintf("%s-valid_until", environment))),
			})
			isValid := c.AuthClient.CheckTokenValidity()
			if !isValid && os.Args[1] != "auth" {
				logrus.Info("Your session has expired, please run the 'auth' again.")
				os.Exit(1)
			}
		}

		// Validate global parameters here, BEFORE we start to waste time
		// and run any code.
		if c.OutputFormat != "" {
			c.OutputFormat = strings.ToLower(c.OutputFormat)
			switch c.OutputFormat {
			case "json", "gron", "yaml", "text", "table", "raw":
				break
			default:
				fmt.Println("Valid options for -o are [json|gron|[text|table]|yaml|raw]")
				os.Exit(1)
			}
			c.FormatOverridden = true
		} else {
			c.FormatOverridden = false
			c.OutputFormat = "json"
		}
	},
}

func buildRootCmd() *cobra.Command {
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.splicectl/config.yml)")
	RootCmd.PersistentFlags().StringVar(&serverURI, "server-uri", "", "override the server uri for the API server http(s)://host.domain.name:overrideport")
	RootCmd.PersistentFlags().StringVarP(&c.OutputFormat, "output", "o", "", "output types: json, text, yaml, gron")
	RootCmd.PersistentFlags().BoolVar(&c.NoHeaders, "no-headers", false, "Suppress header output in Text output")
	RootCmd.PersistentFlags().StringVar(&c.CACert, "cacert", "", "Specify a cacert file to use to authenticate the SSL certificate")

	return RootCmd
}

func addSubcommands() {
	RootCmd.AddCommand(
		apply.InitSubCommands(c),
		create.InitSubCommands(c),
		del.InitSubCommands(c),
		get.InitSubCommands(c),
		list.InitSubCommands(c),
		restart.InitSubCommands(c),
		rollback.InitSubCommands(c),
		version.InitSubCommands(c),
	)
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
	addTUIFunctionsToConfig()
	buildRootCmd()
	cobra.OnInitialize(initConfig)
	addSubcommands()
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

// ClientSemVer - returns the full semVer as the first string and the numerical
// portion as the second string, they may be identical. One example where they
// would not be is:
//         semVer: v0.1.1-cacert -> (v0.1.1-cacert, v0.1.1).
func ClientSemVer() (string, string) {
	submatches := semVerReg.FindStringSubmatch(semVer)
	if submatches == nil || len(submatches) < 2 {
		logrus.Fatalf("the semver in the current build is not valid: %s", semVer)
	}
	return submatches[0], submatches[1]
}
