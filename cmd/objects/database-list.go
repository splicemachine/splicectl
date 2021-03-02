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

// DatabaseList - Data to read from the CM API
type DatabaseList struct {
	Clusters []CMClusterInfo `json:"clusters"`
}

// CMClusterInfo - Cluster Info
type CMClusterInfo struct {
	CreatedAt             string                   `json:"createdAt"`
	UpdatedAt             string                   `json:"updatedAt"`
	DeletedAt             string                   `json:"deletedAt"`
	ClusterId             string                   `json:"clusterId"`
	DcosAppId             string                   `json:"dcosAppId"`
	Name                  string                   `json:"name"`
	Namespace             string                   `json:"namespace"`
	Status                string                   `json:"status"`
	ClusterConfigurations []CMClusterConfiguration `json:"clusterConfigurations"`
	Account               CMAccount                `json:"account"`
	User                  CMUser                   `json:"user"`
}

// CMClusterConfiguration - CM Cluster Config Info
type CMClusterConfiguration struct {
	CreatedAt          string `json:"createdAt"`
	UpdatedAt          string `json:"updatedAt"`
	EffectiveStartDate string `json:"effectiveStartDt"`
	EffectiveEndDate   string `json:"effectiveEndDt"`
	FreeTier           bool   `json:"freeTier"`
}

// CMAccount - Account Info
type CMAccount struct {
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	AccountName string `json:"accountName"`
	AccountId   string `json:"accountId"`
}

// CMUser - User Info
type CMUser struct {
	LastLogin string `json:"lastLogin"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Email     string `json:"email"`
}

// ToJSON - Write the output as JSON
func (databaseList *DatabaseList) ToJSON() error {

	dblJSON, enverr := json.MarshalIndent(databaseList, "", "  ")
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting json")
		return enverr
	} else {
		fmt.Println(string(dblJSON[:]))
	}

	return nil

}

// ToGRON - Write the output as GRON
func (databaseList *DatabaseList) ToGRON() error {
	dbJSON, enverr := json.MarshalIndent(databaseList, "", "  ")
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting json")
		return enverr
	}

	subReader := strings.NewReader(string(dbJSON[:]))
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
func (databaseList *DatabaseList) ToYAML() error {

	dblYAML, enverr := yaml.Marshal(databaseList)
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting yaml")
		return enverr
	} else {
		fmt.Println(string(dblYAML[:]))
	}
	return nil

}

// ToTEXT - Write the output as TEXT
func (databaseList *DatabaseList) ToTEXT(noHeaders bool) error {

	var row []string

	// ******************** TableWriter *******************************
	table := tablewriter.NewWriter(os.Stdout)
	if !noHeaders {
		table.SetHeader([]string{"DATABASE", "NAMESPACE", "STATUS", "CLUSTER_ID"})
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
	for _, v := range databaseList.Clusters {
		row = []string{v.DcosAppId, v.Namespace, v.Status, v.ClusterId}
		table.Append(row)
	}
	table.Render()

	return nil

}
