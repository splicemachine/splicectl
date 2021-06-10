package objects

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/maahsome/gron"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// VaultVersionList - array of versions
type VaultVersionList struct {
	Versions []VaultVersion
}

// VaultVersion - structure of the version output of Vault
type VaultVersion struct {
	Version      int    `json:"version"`
	CreatedTime  string `json:"created_time"`
	DeletionTime string `json:"deletion_time"`
	Destroyed    bool   `json:"destroyed"`
}

// ToJSON - Write the output as JSON
func (vv *VaultVersionList) ToJSON() string {
	vvJSON, enverr := json.MarshalIndent(vv, "", "  ")
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting json")
		return ""
	}
	return string(vvJSON)
}

// ToGRON - Write the output as GRON
func (vv *VaultVersionList) ToGRON() string {
	vvJSON, enverr := json.MarshalIndent(vv, "", "  ")
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting json")
		return ""
	}

	subReader := strings.NewReader(string(vvJSON[:]))
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
func (vv *VaultVersionList) ToYAML() string {
	vvYAML, enverr := yaml.Marshal(vv)
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting yaml")
		return ""
	}
	return string(vvYAML[:])
}

// ToText - Write the output as Text
func (vv *VaultVersionList) ToText(noHeaders bool) string {
	buf, row := new(bytes.Buffer), make([]string, 0)

	// ******************** TableWriter *******************************
	table := tablewriter.NewWriter(buf)
	if !noHeaders {
		table.SetHeader([]string{"VERSION", "CREATED_AT", "DELETED_AT", "DESTROYED"})
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
	for _, v := range vv.Versions {
		row = []string{fmt.Sprintf("%d", v.Version), v.CreatedTime, v.DeletionTime, fmt.Sprintf("%t", v.Destroyed)}
		table.Append(row)
	}
	table.Render()

	return buf.String()
}

// ToJSON - Write the output as JSON
func (vv *VaultVersion) ToJSON() string {
	vvJSON, enverr := json.MarshalIndent(vv, "", "  ")
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting json")
		return ""
	}
	return string(vvJSON)

}

// ToGRON - Write the output as GRON
func (vv *VaultVersion) ToGRON() string {
	vvJSON, enverr := json.MarshalIndent(vv, "", "  ")
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting json")
		return ""
	}

	subReader := strings.NewReader(string(vvJSON[:]))
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
func (vv *VaultVersion) ToYAML() string {
	vvYAML, enverr := yaml.Marshal(vv)
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting yaml")
		return ""
	}
	return string(vvYAML[:])
}

// ToText - Write the output as Text
func (vv *VaultVersion) ToText(noHeaders bool) string {
	buf, row := new(bytes.Buffer), make([]string, 0)

	// ******************** TableWriter *******************************
	table := tablewriter.NewWriter(buf)
	if !noHeaders {
		table.SetHeader([]string{"VERSION", "CREATED_AT", "DELETED_AT", "DESTROYED"})
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

	row = []string{fmt.Sprintf("%d", vv.Version), vv.CreatedTime, vv.DeletionTime, fmt.Sprintf("%t", vv.Destroyed)}
	table.Append(row)

	table.Render()

	return buf.String()
}
