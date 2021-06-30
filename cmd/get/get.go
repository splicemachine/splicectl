package get

import (
	"github.com/spf13/cobra"
	"github.com/splicemachine/splicectl/cmd/config"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Args:  cobra.MinimumNArgs(1),
	Short: "Get various resources of the Splice Machine Database Cluster",
	Long: `EXAMPLES
	splicectl get default-cr > ~/tmp/default-cr.json
	splicectl get system-settings > ~/tmp/system-settings.json
	splicectl get database-status --database-name "test" `,
	Run: func(cmd *cobra.Command, args []string) {},
}

var c *config.Config

func InitSubCommands(conf *config.Config) *cobra.Command {
	c = conf
	return getCmd
}
