package objects

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/maahsome/gron"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// ActionStatus - Status of various actions
type ActionStatus struct {
	Process  string `json:"Process"`
	Success  bool   `json:"Success"`
	Database string `json:"database"`
	Error    string `json:"error"`
}

// ToJSON - Write the output as JSON
func (as *ActionStatus) ToJSON() string {

	asJSON, enverr := json.MarshalIndent(as, "", "  ")
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting json")
		return ""
	}
	return string(asJSON[:])
}

// ToGRON - Write the output as GRON
func (as *ActionStatus) ToGRON() string {
	asJSON, enverr := json.MarshalIndent(as, "", "  ")
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting json")
		return ""
	}

	subReader := strings.NewReader(string(asJSON[:]))
	subValues := &bytes.Buffer{}
	ges := gron.NewGron(subReader, subValues)
	ges.SetMonochrome(false)
	serr := ges.ToGron()
	if serr != nil {
		logrus.Error("Problem generating gron syntax", serr)
		return ""
	}
	return string(subValues.Bytes())
}

// ToYAML - Write the output as YAML
func (as *ActionStatus) ToYAML() string {
	asYAML, enverr := yaml.Marshal(as)
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting yaml")
		return ""
	}
	return string(asYAML[:])
}

// ToText - Write the output as Text
func (as *ActionStatus) ToText(noHeaders bool) string {
	buf, row := new(bytes.Buffer), make([]string, 0)
	// ******************** TableWriter *******************************
	table := tablewriter.NewWriter(os.Stdout)
	if !noHeaders {
		table.SetHeader([]string{"PROCESS", "SUCCESS", "DATABASE", "ERROR"})
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
	row = []string{as.Process, fmt.Sprintf("%t", as.Success), as.Database, as.Error}
	table.Append(row)
	table.Render()

	return buf.String()
}
