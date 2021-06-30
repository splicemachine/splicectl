package objects

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
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
func (settings *SystemSettings) ToJSON() string {

	settingsJSON, enverr := json.MarshalIndent(settings, "", "  ")
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting json")
		return ""
	}
	return string(settingsJSON[:])
}

// ToGRON - Write the output as GRON
func (settings *SystemSettings) ToGRON() string {
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
func (settings *SystemSettings) ToYAML() string {
	settingsYAML, enverr := yaml.Marshal(settings)
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting yaml")
		return ""
	}
	return string(settingsYAML[:])
}

// ToText - Write the output as Text
func (settings *SystemSettings) ToText(noHeaders bool, decode bool) string {
	buf, row := new(bytes.Buffer), make([]string, 0)
	// ******************** TableWriter *******************************
	table := tablewriter.NewWriter(buf)
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

	return buf.String()
}
