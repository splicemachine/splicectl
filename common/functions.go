package common

import (
	"encoding/json"
	"sort"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/splicemachine/splicectl/cmd/objects"
	"sigs.k8s.io/yaml"
)

// WantJSON - This function takes in JSON or YAML and will return JSON, or err
func WantJSON(raw []byte) ([]byte, error) {

	var jsonStructure interface{}
	var yamlStructure interface{}
	if err := json.Unmarshal(raw, &jsonStructure); err != nil {
		// The data isn't JSON, try YAML
		if err := yaml.Unmarshal(raw, &yamlStructure); err != nil {
			return []byte(""), err
		}
		jsonRaw, cerr := json.MarshalIndent(yamlStructure, "", "  ")
		if cerr != nil {
			return []byte(""), cerr
		}
		return jsonRaw, nil
	}
	return raw, nil

}

// RestructureVersions - Vault Version JSON is not well, needs some help.
func RestructureVersions(in string) (objects.VaultVersionList, error) {
	// The raw data out of Hashicorp Vault for versions uses JSON keys that
	// are numeric (Version), rather than "version": <versionnum>, so we
	// strip out the first level map[string] and re-build the struct with
	// a Version field and populate it with the retured key.  Messy and only
	// needed when we want to display in other formats.
	rawData := map[string]interface{}{}
	marshErr := json.Unmarshal([]byte(in), &rawData)
	if marshErr != nil {
		return objects.VaultVersionList{}, marshErr
	}
	var versionList []objects.VaultVersion
	for k, v := range rawData {
		var crvData objects.VaultVersion
		crJSON, enverr := json.MarshalIndent(v, "", "  ")
		if enverr != nil {
			return objects.VaultVersionList{}, enverr
		}
		marshErr := json.Unmarshal([]byte(crJSON), &crvData)
		if marshErr != nil {
			return objects.VaultVersionList{}, marshErr
		}
		i, cnverr := strconv.ParseInt(k, 10, 0)
		if cnverr != nil {
			return objects.VaultVersionList{}, marshErr
		}
		crvData.Version = int(i)
		versionList = append(versionList, crvData)
	}
	sort.Slice(versionList, func(i, j int) bool {
		return versionList[i].Version < versionList[j].Version
	})

	crData := objects.VaultVersionList{
		Versions: versionList,
	}

	return crData, nil
}

// DatabaseName - gets the most preferred database name from the command flags.
// This is meant to be used to pick the best option when multiple flags are
// provided for the database-name through its different aliases.
func DatabaseName(cmd *cobra.Command) (string, string) {
	return "", ""
}
