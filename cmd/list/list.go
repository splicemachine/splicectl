package list

import (
	"github.com/spf13/cobra"
	"github.com/splicemachine/splicectl/cmd/config"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Args:  cobra.MinimumNArgs(1),
	Short: "List resources in the Splice Machine workspace Cluster",
	Long: `EXAMPLES
	splicectl list workspace`,
	Run: func(cmd *cobra.Command, args []string) {},
}

var c *config.Config

func InitSubCommands(conf *config.Config) *cobra.Command {
	c = conf
	return listCmd
}
