package objects

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/maahsome/gron"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// SystemSettings - Structure to hold data for the system-settings calls
type SystemSettings struct {
	Data map[string]string `json:"data"`
}

// ToJSON - Write the output as JSON
func (settings *SystemSettings) ToJSON() error {

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
func (settings *SystemSettings) ToGRON() error {
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
func (settings *SystemSettings) ToYAML() error {

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
func (settings *SystemSettings) ToTEXT(noHeaders bool, decode bool) error {

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
	if decode {
		table.SetFooter([]string{"* denotes a field where values were base64 decoded", ""})
	}
	for k, v := range settings.Data {
		if decode {
			switch k {
			case "POSTGRES_BACKUP_AZURE_ACCOUNT_NAME", "POSTGRES_BACKUP_AZURE_ACCOUNT_KEY", "POSTGRES_BACKUP_AWS_SECRET_ACCESS_KEY", "POSTGRES_PASSWORD", "POSTGRES_USER":
				data, wasEncoded := base64.StdEncoding.DecodeString(v)
				if wasEncoded == nil {
					row = []string{fmt.Sprintf("%s *", k), string(data)}
				} else {
					row = []string{k, v}
				}
			default:
				row = []string{k, v}
			}
		} else {
			row = []string{k, v}
		}
		table.Append(row)

	}
	table.Render()

	return nil

}
