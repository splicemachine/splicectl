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

// Version - Structure for Version Info
type Version struct {
	VersionInfo struct {
		Client BaseVersion `json:"Client"`
		Server BaseVersion `json:"Server"`
	} `json:"VersionInfo"`
	Host string `json:"Host"`
}

// BaseVersion - structure returned by the server
type BaseVersion struct {
	SemVer    string `json:"SemVer"`
	GitCommit string `json:"GitCommit"`
	BuildDate string `json:"BuildDate"`
}

// ToJSON - Write the output as JSON
func (v *Version) ToJSON() error {

	versionJSON, enverr := json.MarshalIndent(v, "", "  ")
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting json")
		return enverr
	}
	fmt.Println(string(versionJSON[:]))

	return nil

}

// ToGRON - Write the output as GRON
func (v *Version) ToGRON() error {
	verJSON, enverr := json.MarshalIndent(v, "", "  ")
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting json")
		return enverr
	}

	subReader := strings.NewReader(string(verJSON[:]))
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
func (v *Version) ToYAML() error {

	versionYAML, enverr := yaml.Marshal(v)
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting yaml")
		return enverr
	}
	fmt.Println(string(versionYAML[:]))

	return nil

}

// ToTEXT - Write the output as TEXT
func (v *Version) ToTEXT(noHeaders bool) error {

	var row []string

	// ******************** TableWriter *******************************
	table := tablewriter.NewWriter(os.Stdout)
	if !noHeaders {
		table.SetHeader([]string{"COMPONENT", "VERSION"})
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
	row = []string{"Client", v.VersionInfo.Client.SemVer}
	table.Append(row)
	row = []string{"Server", v.VersionInfo.Server.SemVer}
	table.Append(row)

	table.Render()

	return nil

}
