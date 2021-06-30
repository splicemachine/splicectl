package objects

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"

	"github.com/maahsome/gron"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// CMSettings - Structure to hold data for the cm-settings calls
type CMSettings struct {
	Data map[string]string `json:"data"`
}

// ToJSON - Write the output as JSON
func (settings *CMSettings) ToJSON() string {
	settingsJSON, enverr := json.MarshalIndent(settings, "", "  ")
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting json")
		return ""
	}
	return string(settingsJSON[:])
}

// ToGRON - Write the output as GRON
func (settings *CMSettings) ToGRON() string {
	settingsJSON, enverr := json.MarshalIndent(settings, "", "  ")
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting json")
		return ""
	}

	subReader := strings.NewReader(string(settingsJSON[:]))
	subValues := &bytes.Buffer{}
	ges := gron.NewGron(subReader, subValues)
	ges.SetMonochrome(false)
	if serr := ges.ToGron(); serr != nil {
		logrus.Error("Problem generating gron syntax", serr)
		return ""
	}
	return string(subValues.Bytes())
}

// ToYAML - Write the output as YAML
func (settings *CMSettings) ToYAML() string {
	settingsYAML, enverr := yaml.Marshal(settings)
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting yaml")
		return ""
	}
	return string(settingsYAML[:])
}

// ToText - Write the output as Text
func (settings *CMSettings) ToText(noHeaders bool) string {
	buf, row := new(bytes.Buffer), make([]string, 0)
	// ******************** TableWriter *******************************
	table := tablewriter.NewWriter(os.Stdout)
	if !noHeaders {
		table.SetHeader([]string{"KEY", "VALUE"})
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
	for k, v := range settings.Data {
		row = []string{k, v}
		table.Append(row)

	}
	table.Render()

	return buf.String()
}
