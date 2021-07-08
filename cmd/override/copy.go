package override

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
		// TODO: Add and an --all/-a flag that will do the whole thing
		err := copy(args[0], args[1], false)
		if err != nil {
			logrus.WithError(err).Error("could not perform copy due to error")
		} else {
			logrus.Infof("Copy completed successfully for component/resource: %s/%s.", args[0], args[1])
		}
	},
}

func init() {
	overrideCmd.AddCommand(copyCmd)
}

// copy - performs the override setup by copying the default configuration into
// override location for the requested component and resource.
func copy(comp, rsrc string, all bool) error {
	component, err := GetComponent(comp)
	if err != nil {
		return err
	}
	resource, err := component.GetDefaultResource(rsrc)
	if err != nil {
		return err
	}
	rm, ok := resource.(map[string]interface{})
	fmt.Println(rm)
	if !ok {
		return errors.New("resource was not expected type")
	}
	return component.PutOverrideResource(rsrc, rsrc)
}
