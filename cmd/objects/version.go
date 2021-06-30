package objects

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/maahsome/gron"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// ApplyCmSettings - lookup names
const ApplyCmSettings = "apply_cm-settings"

// CommandVersions - map used to determine compatibility of API server
var CommandVersions = map[string]string{
	"override_copy":            "0.2.0",
	"apply_cm-settings":        "0.1.6",
	"apply_database-cr":        "0.0.14",
	"apply_default-cr":         "0.0.14",
	"apply_image-tag":          "0.0.16",
	"apply_system-settings":    "0.0.14",
	"apply_vault-key":          "0.0.14",
	"create_database":          "0.1.7",
	"delete":                   "0.1.7",
	"get_accounts":             "0.1.7",
	"get_cm-settings":          "0.1.6",
	"get_database-cr":          "0.0.14",
	"get_database-status":      "0.1.6",
	"get_default-cr":           "0.0.14",
	"get_image-tag":            "0.0.16",
	"get_system-settings":      "0.0.14",
	"get_vault-key":            "0.0.14",
	"list_database":            "0.0.14",
	"pause":                    "0.1.7",
	"restart_database":         "0.1.6",
	"resume":                   "0.1.7",
	"rollback_cm-settings":     "0.1.6",
	"rollback_database-cr":     "0.0.15",
	"rollback_default-cr":      "0.0.15",
	"rollback_system-settings": "0.0.15",
	"rollback_vault-key":       "0.0.15",
	"versions_cm-settings":     "0.1.6",
	"versions_database-cr":     "0.0.15",
	"versions_default-cr":      "0.0.15",
	"versions_system-settings": "0.0.15",
	"versions_vault-key":       "0.0.15",
}

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

// RequirementMet - Check if the command is supported by the server version
func (v *Version) RequirementMet(command string) (semver.Version, semver.Version) {
	cv, err := semver.Parse(CommandVersions[command])
	if err != nil {
		logrus.Warn(fmt.Sprintf("Error parsing SemVer for %s", CommandVersions[command]))
	}

	sv, serr := semver.Parse(strings.Replace(v.VersionInfo.Server.SemVer, "v", "", 1))
	if serr != nil {
		logrus.Warn(fmt.Sprintf("Error parsing SemVer for %s", v.VersionInfo.Server.SemVer))
	}

	if sv.GTE(cv) {
		return cv, sv
	}
	logrus.Fatal(fmt.Sprintf("The API server, version %s, does not support this call, the version needs to be v%s or higher", v.VersionInfo.Server.SemVer, CommandVersions[command]))
	// Clearly we will never get here, though the compiler complains.
	return semver.Version{}, semver.Version{}
}

// ToJSON - Write the output as JSON
func (v *Version) ToJSON() string {
	versionJSON, enverr := json.MarshalIndent(v, "", "  ")
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting json")
		return ""
	}
	return string(versionJSON[:])
}

// ToGRON - Write the output as GRON
func (v *Version) ToGRON() string {
	verJSON, enverr := json.MarshalIndent(v, "", "  ")
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting json")
		return ""
	}

	subReader := strings.NewReader(string(verJSON[:]))
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
func (v *Version) ToYAML() string {
	versionYAML, enverr := yaml.Marshal(v)
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting yaml")
		return ""
	}
	return string(versionYAML[:])
}

// ToText - Write the output as Text
func (v *Version) ToTEXT(noHeaders bool) string {
	buf, row := new(bytes.Buffer), make([]string, 0)

	// ******************** TableWriter *******************************
	table := tablewriter.NewWriter(buf)
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

	return buf.String()
}
