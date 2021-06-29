package get

import (
	"context"
	"io"
	"os"
	"path"

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
			logrus.WithError(err).Fatal("could not get kube config to generate core urls")
		}

		// Get database name
		var dberr error
		databaseName := common.DatabaseName(cmd)
		if len(databaseName) == 0 {
			databaseName, dberr = c.PromptForDatabaseName()
			if dberr != nil {
				logrus.Fatal("Could not get a list of Databases", dberr)
			}
		}

		// Get list of pods that match selector
		pi := client.CoreV1().Pods(databaseName)
		pods, err := pi.List(context.TODO(), v1.ListOptions{LabelSelector: ""})
		if err != nil {
			logrus.WithError(err).Fatal("could not list pods")
		}

		// Make directory for log files
		dirName := ""
		if _, err := os.Stat(dirName); err == nil {
			if err := os.RemoveAll(dirName); err != nil {
				logrus.WithError(err).Fatalf("directory or file existed and could not be deleted: %s", dirName)
			}
		}
		if err := os.Mkdir(dirName, 0755); err != nil {
			logrus.WithError(err).Fatalf("could not make directory: %s", dirName)
		}

		// Get log from each pod
		for _, pod := range pods.Items {
			// Get log data from api
			stream, err := pi.GetLogs(pod.Name, &core.PodLogOptions{}).Stream(context.TODO())
			if err != nil {
				logrus.WithError(err).Errorf("could not get log for %s", pod.Name)
				continue
			}

			// Open file to write log data
			filName := path.Join(dirName, pod.Name)
			fil, err := os.OpenFile(filName, os.O_CREATE, 0777)
			if err != nil {
				logrus.WithError(err).Errorf("could not open file for pod: %s", pod.Name)
			}

			// Stream data from http request to file
			if _, err := io.Copy(fil, stream); err != nil {
				logrus.WithError(err).Errorf("could not write log to file from pod: %s", pod.Name)
			}
		}
	},
}

func init() {
	// TODO: wire flags into exec func, stream concurrently, test
	getLogsCmd.Flags().BoolP("all", "a", false, "whether to get all logs from all pods with no selector")
	getLogsCmd.Flags().Bool("sequential", false, "whether to get logs one at a time, slower but safer")

	getLogsCmd.Flags().StringP("selector", "s", "app=hbase", "kubernetes selector expresssion to filter pods on")
	getLogsCmd.Flags().StringP("directory", "f", "", "name of folder to output logs to")

	// add database name and aliases
	getUrlsCmd.Flags().StringP("database-name", "d", "", "Specify the database name")
	getUrlsCmd.Flags().String("database", "", "Alias for database-name, prefer the use of -d and --database-name.")
	getUrlsCmd.Flags().String("workspace", "", "Alias for database-name, prefer the use of -d and --database-name.")

	getLogsCmd.MarkFlagRequired("directory")

	getCmd.AddCommand(getLogsCmd)
}
