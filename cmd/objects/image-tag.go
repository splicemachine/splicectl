package objects

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/maahsome/gron"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// ImageTagList - An array of image tags
type ImageTagList struct {
	ImageTags []ImageTag
}

// ImageTag - Structure for Version Info
type ImageTag struct {
	Component       string `json:"Component"`
	DatabaseCRImage string `json:"DatabaseCRImage"`
	ActiveImage     string `json:"ActiveImage"`
}

// ToJSON - Write the output as JSON
func (i *ImageTagList) ToJSON() string {
	imageTagJSON, enverr := json.MarshalIndent(i, "", "  ")
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting json")
		return ""
	}
	return string(imageTagJSON[:])
}

// ToGRON - Write the output as GRON
func (i *ImageTagList) ToGRON() string {
	tagJSON, enverr := json.MarshalIndent(i, "", "  ")
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting json")
		return ""
	}

	subReader := strings.NewReader(string(tagJSON[:]))
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
func (i *ImageTagList) ToYAML() string {
	imageTagYAML, enverr := yaml.Marshal(i)
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting yaml")
		return ""
	}
	return string(imageTagYAML[:])
}

// ToText - Write the output as Text
func (i *ImageTagList) ToText(noHeaders bool) string {
	buf, row := new(bytes.Buffer), make([]string, 0)

	// ******************** TableWriter *******************************
	table := tablewriter.NewWriter(buf)
	if !noHeaders {
		table.SetHeader([]string{"COMPONENT", "DB_CR_IMAGE", "ACTIVE_IMAGE"})
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
	for _, v := range i.ImageTags {
		row = []string{v.Component, v.DatabaseCRImage, v.ActiveImage}
		table.Append(row)
	}
	table.Render()

	return buf.String()
}
