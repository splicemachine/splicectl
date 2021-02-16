package objects

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/maahsome/gron"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type DatabaseCR struct {
	Data map[string]interface{} `json:"data"`
}

// DatabaseCRList - List of Databases
type DatabaseCRList struct {
	Data struct {
		Metadata struct {
			Name string `json:"name"`
		} `json:"metadata"`
		Spec struct {
			Condition struct {
				Haproxy struct {
					Enabled bool `json:"enabled"`
				} `json:"haproxy"`
				Hbase struct {
					Enabled bool `json:"enabled"`
				} `json:"hbase"`
				Hdfs struct {
					Enabled bool `json:"enabled"`
				} `json:"hdfs"`
				JupyterHub struct {
					Enabled bool `json:"enabled"`
				} `json:"jupyterhub"`
				JvmProfiler struct {
					Enabled bool `json:"enabled"`
				} `json:"jvmprofiler"`
				Kafka struct {
					Enabled bool `json:"enabled"`
				} `json:"kafka"`
				MlManager struct {
					Enabled bool `json:"enabled"`
				} `json:"mlmanager"`
				Rbac struct {
					Enabled bool `json:"enabled"`
				} `json:"rbac"`
				SpliceHTTP struct {
					Enabled bool `json:"enabled"`
				} `json:"splice-http"`
				Zookeeper struct {
					Enabled bool `json:"enabled"`
				} `json:"zookeeper"`
			} `json:"condition"`
			Global struct {
				Namespace     string `json:"dnsPrefix"`
				CloudProvider string `json:"cloudprovider"`
			} `json:"global"`
		} `json:"spec"`
	} `json:"data"`
}

// ToJSON - Write the output as JSON
func (cr *DatabaseCR) ToJSON(file string) error {

	crJSON, enverr := json.MarshalIndent(cr, "", "  ")
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting json")
		return enverr
	} else {
		if len(file) == 0 {
			fmt.Println(string(crJSON[:]))
		} else {
			WriteToFile(file, string(crJSON[:]))
		}
	}

	return nil

}

// ToGRON - Write output in GRON format
func (cr *DatabaseCR) ToGRON(file string) error {

	crJSON, enverr := json.MarshalIndent(cr, "", "  ")
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting json")
		return enverr
	}

	subReader := strings.NewReader(string(crJSON[:]))
	subValues := &bytes.Buffer{}
	ges := gron.NewGron(subReader, subValues)
	ges.SetMonochrome(false)
	serr := ges.ToGron()
	if serr != nil {
		logrus.Error("Problem generating gron syntax", serr)
		return serr
	}
	fmt.Println(string(subValues.Bytes()))

	return nil
}

// ToYAML - Write the output as YAML
func (cr *DatabaseCR) ToYAML(file string) error {

	crYAML, enverr := yaml.Marshal(cr)
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting yaml")
		return enverr
	} else {
		if len(file) == 0 {
			fmt.Println(string(crYAML[:]))
		} else {
			WriteToFile(file, string(crYAML[:]))
		}
	}
	return nil

}

// ToTEXT - Write the output as TEXT
func (cr *DatabaseCR) ToTEXT(noHeaders bool) error {

	var row []string

	out, _ := json.Marshal(&cr)

	var crList DatabaseCRList
	marshErr := json.Unmarshal([]byte(out), &crList)
	if marshErr != nil {
		logrus.Fatal("Could not unmarshall data", marshErr)
	}

	// ******************** TableWriter *******************************
	table := tablewriter.NewWriter(os.Stdout)
	if !noHeaders {
		table.SetHeader([]string{"DATABASE", "NAMESPACE", "HAPROXY", "HBASE", "HDFS", "JUPYTERHUB", "JVMPROFILER", "KAFKA", "MLMANAGER", "RBAC", "SPLICE-HTTP", "ZOOKEEPER"})
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	}
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)
	row = []string{crList.Data.Metadata.Name,
		crList.Data.Spec.Global.Namespace,
		fmt.Sprintf("%t", crList.Data.Spec.Condition.Haproxy.Enabled),
		fmt.Sprintf("%t", crList.Data.Spec.Condition.Hbase.Enabled),
		fmt.Sprintf("%t", crList.Data.Spec.Condition.Hdfs.Enabled),
		fmt.Sprintf("%t", crList.Data.Spec.Condition.JupyterHub.Enabled),
		fmt.Sprintf("%t", crList.Data.Spec.Condition.JvmProfiler.Enabled),
		fmt.Sprintf("%t", crList.Data.Spec.Condition.Kafka.Enabled),
		fmt.Sprintf("%t", crList.Data.Spec.Condition.MlManager.Enabled),
		fmt.Sprintf("%t", crList.Data.Spec.Condition.Rbac.Enabled),
		fmt.Sprintf("%t", crList.Data.Spec.Condition.SpliceHTTP.Enabled),
		fmt.Sprintf("%t", crList.Data.Spec.Condition.Zookeeper.Enabled),
	}
	table.Append(row)
	table.Render()

	return nil

}

// WriteToFile will print any string of text to a file safely by
// checking for errors and syncing at the end.
func WriteToFile(filename string, data string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}
	return file.Sync()
}
