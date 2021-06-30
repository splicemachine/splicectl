package create

import (
	"github.com/spf13/cobra"
	"github.com/splicemachine/splicectl/cmd/config"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Args:  cobra.MinimumNArgs(1),
	Short: "Create cluster resources",
	Long: `EXAMPLES
	splicectl create workspace --database-name splicedb --dnsprefix splicedb`,
	Run: func(cmd *cobra.Command, args []string) {},
}

var c *config.Config

func InitSubCommands(conf *config.Config) *cobra.Command {
	c = conf
	return createCmd
}
