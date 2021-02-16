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

// CMSettings - Structure to hold data for the cm-settings calls
type CMSettings struct {
	Data map[string]string `json:"data"`
}

// ToJSON - Write the output as JSON
func (settings *CMSettings) ToJSON() error {

	settingsJSON, enverr := json.MarshalIndent(settings, "", "  ")
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting json")
		return enverr
	} else {
		fmt.Println(string(settingsJSON[:]))
	}

	return nil

}

// ToGRON - Write the output as GRON
func (settings *CMSettings) ToGRON() error {
	settingsJSON, enverr := json.MarshalIndent(settings, "", "  ")
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting json")
		return enverr
	}

	subReader := strings.NewReader(string(settingsJSON[:]))
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
func (settings *CMSettings) ToYAML() error {

	settingsYAML, enverr := yaml.Marshal(settings)
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting yaml")
		return enverr
	} else {
		fmt.Println(string(settingsYAML[:]))
	}
	return nil

}

// ToTEXT - Write the output as TEXT
func (settings *CMSettings) ToTEXT(noHeaders bool) error {

	var row []string

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

	return nil

}
