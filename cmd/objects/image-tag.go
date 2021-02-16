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
func (i *ImageTagList) ToJSON() error {

	imageTagJSON, enverr := json.MarshalIndent(i, "", "  ")
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting json")
		return enverr
	}
	fmt.Println(string(imageTagJSON[:]))

	return nil

}

// ToGRON - Write the output as GRON
func (i *ImageTagList) ToGRON() error {
	tagJSON, enverr := json.MarshalIndent(i, "", "  ")
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting json")
		return enverr
	}

	subReader := strings.NewReader(string(tagJSON[:]))
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
func (i *ImageTagList) ToYAML() error {

	imageTagYAML, enverr := yaml.Marshal(i)
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting yaml")
		return enverr
	}
	fmt.Println(string(imageTagYAML[:]))

	return nil

}

// ToTEXT - Write the output as TEXT
func (i *ImageTagList) ToTEXT(noHeaders bool) error {

	var row []string

	// ******************** TableWriter *******************************
	table := tablewriter.NewWriter(os.Stdout)
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

	return nil

}
