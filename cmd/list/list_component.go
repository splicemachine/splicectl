package list

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/splicemachine/splicectl/cmd/override"
)

var listComponentCmd = &cobra.Command{
	Use:     "component",
	Args:    cobra.MaximumNArgs(1),
	Aliases: []string{"components"},
	Short:   "Retrieve a list of overridable splice components.",
	Long: `EXAMPLES
	# lists all components that can be overridden
	splicectl list component
	or
	# lists all resources that a given component has that can be overriden
	splicectl list component splicedb-hbase-config
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 { // list components
			fmt.Printf("Overridable Components: \n%s", override.PrettyListComponents())
		} else { // list resources
			component, err := override.GetComponent(args[0])
			if err != nil {
				logrus.WithError(err).Fatal("could not get component")
			}
			fmt.Printf("Overridable Resources: \n%s", component.PrettyListResources())
		}
	},
}

func init() {
	listCmd.AddCommand(listComponentCmd)
}
