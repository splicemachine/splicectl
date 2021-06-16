package restart

import (
	"github.com/spf13/cobra"
	"github.com/splicemachine/splicectl/cmd/config"
)

var restartCmd = &cobra.Command{
	Use:   "restart",
	Args:  cobra.MinimumNArgs(1),
	Short: "Restart components of the Splice DB Cluster",
	Long: `EXAMPLES
	splicectl restart workspace --database-name <dbname>`,
	Run: func(cmd *cobra.Command, args []string) {},
}

var c *config.Config

func InitSubCommands(conf *config.Config) *cobra.Command {
	c = conf
	return restartCmd
}
