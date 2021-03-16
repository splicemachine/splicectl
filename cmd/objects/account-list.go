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
func (accountList *AccountList) ToJSON() error {

	listJSON, enverr := json.MarshalIndent(accountList, "", "  ")
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting json")
		return enverr
	} else {
		fmt.Println(string(listJSON[:]))
	}

	return nil

}

// ToGRON - Write the output as GRON
func (accountList *AccountList) ToGRON() error {
	listJSON, enverr := json.MarshalIndent(accountList, "", "  ")
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting json")
		return enverr
	}

	subReader := strings.NewReader(string(listJSON[:]))
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
func (accountList *AccountList) ToYAML() error {

	listYAML, enverr := yaml.Marshal(accountList)
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting yaml")
		return enverr
	} else {
		fmt.Println(string(listYAML[:]))
	}
	return nil

}

// ToTEXT - Write the output as TEXT
func (accountList *AccountList) ToTEXT(noHeaders bool) error {

	var row []string

	// ******************** TableWriter *******************************
	table := tablewriter.NewWriter(os.Stdout)
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

	return nil

}
