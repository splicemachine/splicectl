package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Express the 'version' of splicectl.",
	Aliases: []string{"v"},
	Run: func(cmd *cobra.Command, args []string) {

		if !c.FormatOverridden {
			c.OutputFormat = "yaml"
		}

		switch strings.ToLower(c.OutputFormat) {
		case "raw":
			fmt.Println(c.VersionJSON)
		case "json":
			// We want to print the JSON in a condensed format
			fmt.Println(c.VersionJSON)
		case "gron":
			fmt.Println(c.VersionDetail.ToGRON())
		case "yaml":
			fmt.Println(c.VersionDetail.ToYAML())
		case "text", "table":
			fmt.Println(c.VersionDetail.ToTEXT(c.NoHeaders))
		}
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
