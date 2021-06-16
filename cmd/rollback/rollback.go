package rollback

import (
	"github.com/spf13/cobra"
	"github.com/splicemachine/splicectl/cmd/config"
)

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Args:  cobra.MinimumNArgs(1),
	Short: "Rollback various resource vault records to a previous version.",
	Long: `EXAMPLES
	splicectl versions default-cr
	splicectl rollback default-cr --version 1
	splicectl list database
	splicectl versions database-cr --database-name splicedb
	splicectl rollback database-cr --database-name splicedb --version 2
	`,
	Run: func(cmd *cobra.Command, args []string) {},
}

var c *config.Config

func InitSubCommands(conf *config.Config) *cobra.Command {
	c = conf
	return rollbackCmd
}
