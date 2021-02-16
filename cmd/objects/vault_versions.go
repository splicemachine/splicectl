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
func (vv *VaultVersionList) ToJSON() error {

	vvJSON, enverr := json.MarshalIndent(vv, "", "  ")
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting json")
		return enverr
	}
	fmt.Println(string(vvJSON[:]))

	return nil

}

// ToGRON - Write the output as GRON
func (vv *VaultVersionList) ToGRON() error {
	vvJSON, enverr := json.MarshalIndent(vv, "", "  ")
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting json")
		return enverr
	}

	subReader := strings.NewReader(string(vvJSON[:]))
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
func (vv *VaultVersionList) ToYAML() error {

	vvYAML, enverr := yaml.Marshal(vv)
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting yaml")
		return enverr
	}
	fmt.Println(string(vvYAML[:]))

	return nil

}

// ToTEXT - Write the output as TEXT
func (vv *VaultVersionList) ToTEXT(noHeaders bool) error {

	var row []string

	// ******************** TableWriter *******************************
	table := tablewriter.NewWriter(os.Stdout)
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

	return nil

}

// ToJSON - Write the output as JSON
func (vv *VaultVersion) ToJSON() error {

	vvJSON, enverr := json.MarshalIndent(vv, "", "  ")
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting json")
		return enverr
	}
	fmt.Println(string(vvJSON[:]))

	return nil

}

// ToGRON - Write the output as GRON
func (vv *VaultVersion) ToGRON() error {
	vvJSON, enverr := json.MarshalIndent(vv, "", "  ")
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting json")
		return enverr
	}

	subReader := strings.NewReader(string(vvJSON[:]))
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
func (vv *VaultVersion) ToYAML() error {

	vvYAML, enverr := yaml.Marshal(vv)
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting yaml")
		return enverr
	}
	fmt.Println(string(vvYAML[:]))

	return nil

}

// ToTEXT - Write the output as TEXT
func (vv *VaultVersion) ToTEXT(noHeaders bool) error {

	var row []string

	// ******************** TableWriter *******************************
	table := tablewriter.NewWriter(os.Stdout)
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

	return nil

}
