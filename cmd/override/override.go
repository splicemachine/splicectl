package override

import (
	"github.com/spf13/cobra"
	"github.com/splicemachine/splicectl/cmd/config"
)

var overrideCmd = &cobra.Command{
	Use:   "override",
	Args:  cobra.MinimumNArgs(1),
	Short: "Override configurations to various resources of the Splice Machine Database Cluster",
	Long: `EXAMPLES
	TODO: Make examples
	`,
	Run: func(cmd *cobra.Command, args []string) {},
}

var c *config.Config

func InitSubCommands(conf *config.Config) *cobra.Command {
	c = conf
	return overrideCmd
}
