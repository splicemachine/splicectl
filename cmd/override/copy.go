package override

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/splicemachine/splicectl/common"
)

var copyCmd = &cobra.Command{
	Use:   "copy",
	Args:  cobra.ExactArgs(2),
	Short: "Sets up an override of a resource by copying the default resource to its override location.",
	Long: `EXAMPLES
	# sets up the override of splicedb-hbase-config with default values for fairscheduler.xml
	splicectl override copy splicedb-hbase-config fairscheduler.xml
	`,
	Run: func(cmd *cobra.Command, args []string) {
		_, _ = c.VersionDetail.RequirementMet("override_copy")
		var err error
		dbName := common.DatabaseName(cmd)
		if len(dbName) == 0 {
			dbName, err = c.PromptForDatabaseName()
			if err != nil {
				logrus.Fatal("Could not get a list of workspaces", err)
			}
		}
		err = copy(dbName, args[0], args[1])
		if err != nil {
			logrus.WithError(err).Error("could not perform copy due to error")
		} else {
			fmt.Printf("Copy completed successfully for component/resource: %s/%s\n", args[0], args[1])
		}
	},
}

func init() {
	// add database name and aliases
	copyCmd.Flags().StringP("database-name", "d", "", "Specify the database name")
	copyCmd.Flags().String("database", "", "Alias for database-name, prefer the use of -d and --database-name.")
	copyCmd.Flags().String("workspace", "", "Alias for database-name, prefer the use of -d and --database-name.")

	overrideCmd.AddCommand(copyCmd)
}

// copy - performs the override setup by copying the default configuration into
// override location for the requested component and resource.
func copy(dbName, comp, rsrc string) error {
	component, err := GetComponent(comp)
	if err != nil {
		return err
	}
	resource, err := component.GetDefaultResource(dbName, rsrc)
	if err != nil {
		return err
	}
	return component.PutOverrideResource(dbName, rsrc, resource)
}
