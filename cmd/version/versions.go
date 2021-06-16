package version

import (
	"github.com/spf13/cobra"
	"github.com/splicemachine/splicectl/cmd/config"
)

var versionsCmd = &cobra.Command{
	Use:   "versions",
	Args:  cobra.MinimumNArgs(1),
	Short: "Get the vault versions for a specific resource in the cluster.",
	Long: `EXAMPLES
	splicectl versions system-settings
	splicectl versions default-cr
	splicectl list workspace
	splicectl versions database-cr --database-name splicedb
	`,
	Run: func(cmd *cobra.Command, args []string) {},
}

var c *config.Config

func InitSubCommands(conf *config.Config) *cobra.Command {
	c = conf
	return versionsCmd
}
