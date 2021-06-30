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

// AccountList - Data to read from CM Postgres
type AccountList struct {
	Accounts []CMUserAccount
}

// CMUserAccount - Data fields from CM Accounts Postres Tables
type CMUserAccount struct {
	AccountID string `json:"accountId"`
	EMail     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// ToJSON - Write the output as JSON
func (accountList *AccountList) ToJSON() string {

	listJSON, enverr := json.MarshalIndent(accountList, "", "  ")
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting json")
		return ""
	}
	return string(listJSON[:])

}

// ToGRON - Write the output as GRON
func (accountList *AccountList) ToGRON() string {
	listJSON, enverr := json.MarshalIndent(accountList, "", "  ")
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting json")
		return ""
	}

	subReader := strings.NewReader(string(listJSON[:]))
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
func (accountList *AccountList) ToYAML() string {

	listYAML, enverr := yaml.Marshal(accountList)
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting yaml")
		return ""
	}
	return string(listYAML[:])
}

// ToText - Write the output as Text
func (accountList *AccountList) ToText(noHeaders bool) string {
	buf, row := new(bytes.Buffer), make([]string, 0)
	// ******************** TableWriter *******************************
	table := tablewriter.NewWriter(buf)
	if !noHeaders {
		table.SetHeader([]string{"ACCOUNTID", "EMAIL", "FIRSTNAME", "LASTNAME"})
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
	for _, v := range accountList.Accounts {
		row = []string{v.AccountID, v.EMail, v.FirstName, v.LastName}
		table.Append(row)
	}
	table.Render()

	return buf.String()
}
