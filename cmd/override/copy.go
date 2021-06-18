package override

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var copyCmd = &cobra.Command{
	Use:   "copy",
	Args:  cobra.ExactArgs(2),
	Short: "TODO: Write description",
	Long: `EXAMPLES
	TODO: Make examples
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Check that version requirement is met
		// TODO: Add and an --all/-a flag that will do the whole thing
		out, err := copy(args[0], args[1], false)
		if err != nil {
			logrus.WithError(err).Error("could perform copy due to error")
		} else {
			logrus.Info(out)
		}
	},
}

func init() {
	overrideCmd.AddCommand(copyCmd)
}

// copy - performs the override setup by copying the default configuration into
// override location for the requested component and resource.
func copy(comp, rsrc string, all bool) (string, error) {
	_, _ = c.VersionDetail.RequirementMet("override_copy")
	component, err := GetComponent(comp)
	if err != nil {
		return "", err
	}
	resource, err := component.GetDefaultResource(rsrc)
	if err != nil {
		return "", err
	}
	// TODO: make changes to resource, update name and labels etc.
	_ = resource
	if err := component.PutOverrideResource(rsrc, resource); err != nil {
		return "", err
	}
	return "all clear", nil
}
