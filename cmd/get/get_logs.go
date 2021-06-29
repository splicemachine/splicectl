package get

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/splicemachine/splicectl/common"
	core "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var getLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Get logs of all pods in the namespace",
	Long: `EXAMPLES
	splicectl get logs --workspace splicedb
`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get kube client
		client, err := common.KubeClient()
		if err != nil {
			logrus.WithError(err).Error("could not get kube config to generate core urls")
			return
		}

		var dberr error
		databaseName := common.DatabaseName(cmd)
		if len(databaseName) == 0 {
			databaseName, dberr = c.PromptForDatabaseName()
			if dberr != nil {
				logrus.Fatal("Could not get a list of Databases", dberr)
			}
		}

		pi := client.CoreV1().Pods(databaseName)
		pods, err := pi.List(context.TODO(), v1.ListOptions{LabelSelector: ""})
		if err != nil {
			logrus.WithError(err).Error("could not list pods")
			return
		}
		for _, pod := range pods.Items {
			resp := pi.GetLogs(pod.Name, &core.PodLogOptions{}).Do(context.TODO())
			if err := resp.Error(); err != nil {
				// do something with the error
			}

		}
	},
}

func init() {
	getCmd.AddCommand(getLogsCmd)
}
