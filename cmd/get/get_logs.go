package get

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/splicemachine/splicectl/common"
	core "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	typed "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
)

type (
	requestBuffer chan core.Pod
	ticket        struct{}
)

const (
	// Common flag names
	selectorFlag  = "selector"
	directoryFlag = "directory"
	allFlag       = "all"
)

var (
	// Single value that ticket can take on
	ticketV ticket = struct{}{}
)

var getLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Get logs of all pods in the namespace",
	Long: `EXAMPLES
	splicectl get logs --workspace splicedb
`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			selector, _ = cmd.Flags().GetString(selectorFlag)
			dirName, _  = cmd.Flags().GetString(directoryFlag)
			all, _      = cmd.Flags().GetBool(allFlag)
		)

		dbNamespace, err := getDBNamespace(cmd)
		if err != nil {
			logrus.WithError(err).Fatal("could not stream logs due to database name error")
		}

		selector = allOrDefault(all, selector)

		client, err := common.KubeClient()
		if err != nil {
			logrus.WithError(err).Fatal("could not get kube config to generate core urls")
		}

		pi := client.CoreV1().Pods(dbNamespace)
		pods, err := pi.List(context.TODO(), v1.ListOptions{LabelSelector: selector})
		if err != nil {
			logrus.WithError(err).Fatal("could not list pods")
		}

		if err := makeDirectory(dirName); err != nil {
			logrus.WithError(err).Fatal("could not stream logs due to directory issue")
		}

		reqBuffer, wg := make(requestBuffer, 1), &sync.WaitGroup{}
		go streamLogs(reqBuffer, wg, dirName, pi)

		for _, pod := range pods.Items {
			wg.Add(1)
			reqBuffer <- pod
		}
		wg.Wait()
	},
}

// allOrDefault - return empty string if selector for all is requested, otherwise return supplied default
func allOrDefault(all bool, def string) string {
	if all {
		return ""
	}
	return def
}

// makeDirectory - makes the directory to output logs to, or clears a pre-existing one
func makeDirectory(dirName string) error {
	if _, err := os.Stat(dirName); err == nil {
		if err := os.RemoveAll(dirName); err != nil {
			return fmt.Errorf("%v; directory or file with name '%s', existed and could not be deleted", err, dirName)
		}
	}
	if err := os.Mkdir(dirName, 0755); err != nil {
		return fmt.Errorf("%v; could not make directory: %s", err, dirName)
	}
	return nil
}

// getDBNamespace - get the namespace of the desired db based on command line flags
func getDBNamespace(cmd *cobra.Command) (string, error) {
	var dberr error
	dbName := common.DatabaseName(cmd)
	if len(dbName) == 0 {
		dbName, dberr = c.PromptForDatabaseName()
		if dberr != nil {
			return "", fmt.Errorf("%v; could not get a list of databases", dberr)
		}
	}

	dbNamespace := ""
	list, err := c.GetDatabaseListStruct()
	if err != nil {
		return "", fmt.Errorf("%v; could not get list of databases", err)
	}
	for _, db := range list.Clusters {
		if db.DcosAppId == dbName {
			dbNamespace = db.Namespace
		}
	}
	if dbNamespace == "" {
		return "", fmt.Errorf("no database matched given name: '%s'", dbName)
	}

	return dbNamespace, nil
}

// podLogOptions - create PodLogOptions with just Container field being set
func podLogOptions(containerName string) *core.PodLogOptions {
	return &core.PodLogOptions{
		Container: containerName,
	}
}

// streamLogs - stream all logs that are requested from the names channel
func streamLogs(pods requestBuffer, waiter *sync.WaitGroup, dirName string, pi typed.PodInterface) {
	buf := make(chan ticket, cap(pods))
	for pod := range pods {
		buf <- ticketV
		go func(pod core.Pod) {
			defer waiter.Done()
			streamLog(pod, dirName, pi)
			<-buf
		}(pod)
	}
}

// streamLog - stream a single log from kubernetes into a file
func streamLog(pod core.Pod, dirName string, pi typed.PodInterface) {
	amount := int64(0)
	filName := dirName
	if len(pod.Spec.Containers) == 1 {
		filName = path.Join(filName, pod.Name+".log")
		amount = streamContainerLog(filName, pi.GetLogs(pod.Name, podLogOptions("")))
	} else {
		filName = path.Join(filName, pod.Name)
		if err := makeDirectory(filName); err != nil {
			logrus.WithError(err).Errorf("could not make directory: %s", filName)
			return
		}
		for _, container := range pod.Spec.Containers {
			innerFilName := path.Join(filName, container.Name+".log")
			amount = streamContainerLog(innerFilName, pi.GetLogs(pod.Name, podLogOptions(container.Name)))
		}
		for _, container := range pod.Spec.InitContainers {
			innerFilName := path.Join(filName, container.Name+".log")
			amount += streamContainerLog(innerFilName, pi.GetLogs(pod.Name, podLogOptions(container.Name)))
		}
	}
	logrus.Infof("wrote %10d bytes to %s", amount, filName)
}

// streamContainerLog - stream a log from a single container of a pod into kubernetes
func streamContainerLog(filName string, req *rest.Request) int64 {
	// Open file to write log data
	fil, err := os.OpenFile(filName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	if err != nil {
		logrus.WithError(err).Errorf("could not open file: %s", filName)
		return 0
	}
	defer fil.Close()

	// Start stream of log data
	stream, err := req.Stream(context.Background())
	if err != nil {
		logrus.WithError(err).Errorf("could not get log for %s", filName)
		return 0
	}
	defer stream.Close()

	// Stream data from http request to file
	amount, err := io.Copy(fil, stream)
	if err != nil {
		logrus.WithError(err).Errorf("could not write log to file: %s", filName)
		return 0
	}
	return amount
}

func init() {
	getLogsCmd.Flags().BoolP(allFlag, "a", false, "whether to get all logs from all pods with no selector")

	getLogsCmd.Flags().StringP(selectorFlag, "s", "app=hbase", "kubernetes selector expresssion to filter pods on")
	getLogsCmd.Flags().StringP(directoryFlag, "f", "", "name of folder to output logs to")

	// add database name and aliases
	getLogsCmd.Flags().StringP("database-name", "d", "", "Specify the database name")
	getLogsCmd.Flags().String("database", "", "Alias for database-name, prefer the use of -d and --database-name.")
	getLogsCmd.Flags().String("workspace", "", "Alias for database-name, prefer the use of -d and --database-name.")

	getLogsCmd.MarkFlagRequired(directoryFlag)
	getLogsCmd.MarkFlagDirname(directoryFlag)

	getCmd.AddCommand(getLogsCmd)
}
