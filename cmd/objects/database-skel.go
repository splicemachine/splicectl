package objects

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// DatabaseRequest - Cloud Manager Database Request
type DatabaseRequest struct {
	AccountID                    string `json:"accountId"`
	AuthorizationCode            string `json:"authorizationCode"`
	BackupFrequency              string `json:"backupFrequency"`
	BackupInterval               int    `json:"backupInterval"`
	BackupKeepCount              int    `json:"backupKeepCount"`
	BackupStartWindow            string `json:"backupStartWindow"`
	CloudProvider                string `json:"cloudProvider"`
	ClusterPowerOlap             int    `json:"clusterPowerOlap"`
	ClusterPowerOltp             int    `json:"clusterPowerOltp"`
	DedicatedStorage             bool   `json:"dedicatedStorage"`
	ExternalDatasetSizeGb        int    `json:"externalDatasetSizeGb"`
	InternalDatasetSizeGb        int    `json:"internalDatasetSizeGb"`
	MlManager                    bool   `json:"mlManager"`
	Name                         string `json:"name"`
	NotebookActiveUsers          int    `json:"notebookActiveUsers"`
	NotebookExecutorsPerNotebook int    `json:"notebookExecutorsPerNotebook"`
	NotebookTotalUsers           int    `json:"notebookTotalUsers"`
	NotebooksPerUser             int    `json:"notebooksPerUser"`
	Password                     string `json:"password"`
	// Region                       string `json:"region"`
}

// ToJSON - Write the output as JSON
func (r *DatabaseRequest) ToJSON() string {
	rJSON, enverr := json.MarshalIndent(r, "", "  ")
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting json")
		return ""
	}
	return string(rJSON[:])
}

// ToYAML - Write the output as YAML
func (r *DatabaseRequest) ToYAML() string {
	rYAML, enverr := yaml.Marshal(r)
	if enverr != nil {
		logrus.WithError(enverr).Error("Error extracting yaml")
		return ""
	}
	return string(rYAML[:])
}
