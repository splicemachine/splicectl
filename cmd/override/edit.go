package override

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/splicemachine/splicectl/common"
	"k8s.io/kubectl/pkg/cmd/util/editor"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Args:  cobra.ExactArgs(2),
	Short: "Downloads the overriden resource from the cluster and opens it in the user's local editor.",
	Long: `EXAMPLES
	# sets up the override of splicedb-hbase-config with default values for fairscheduler.xml
	splicectl override copy splicedb-hbase-config fairscheduler.xml
	# downloads the newly created override resource, opens it in local editor, and when changes are saved pushes it to the cluster
	splicectl override edit splicedb-hbase-config fairscheduler.xml
	`,
	Run: func(cmd *cobra.Command, args []string) {
		_, _ = c.VersionDetail.RequirementMet("override_edit")
		var err error
		dbName := common.DatabaseName(cmd)
		if len(dbName) == 0 {
			dbName, err = c.PromptForDatabaseName()
			if err != nil {
				logrus.Fatal("Could not get a list of workspaces", err)
			}
		}
		err = edit(dbName, args[0], args[1])
		if err != nil {
			logrus.WithError(err).Error("could not perform edit due to error")
		} else {
			fmt.Printf("Edit completed successfully for component/resource: %s/%s\n", args[0], args[1])
		}
	},
}

func init() {
	// add database name and aliases
	editCmd.Flags().StringP("database-name", "d", "", "Specify the database name")
	editCmd.Flags().String("database", "", "Alias for database-name, prefer the use of -d and --database-name.")
	editCmd.Flags().String("workspace", "", "Alias for database-name, prefer the use of -d and --database-name.")

	overrideCmd.AddCommand(editCmd)
}

// edit - edits an existing override by downloading the override resource from
// the cluster, opening it in a local terminal, then pushing any changes back
// to the cluster.
func edit(dbName, comp, rsrc string) error {
	// Get resource
	component, err := GetComponent(comp)
	if err != nil {
		return err
	}
	resource, err := component.GetOverrideResource(dbName, rsrc)
	if err != nil {
		return err
	}

	// Open in editor
	in, e := bytes.NewBufferString(resource), editor.NewDefaultEditor([]string{})
	data, path, err := e.LaunchTempFile("override-", ".yaml", in)
	defer os.Remove(path)
	if err != nil {
		return err
	}

	// Put resource
	strData := string(data)
	if strData == resource {
		return errors.New("no changes were made, the resource was not rewritten")
	}
	return component.PutOverrideResource(dbName, rsrc, strData)
}
