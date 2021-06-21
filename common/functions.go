package common

import (
	"encoding/json"
	"sort"
	"strconv"

	"github.com/sirupsen/logrus"
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
// provided for the database-name through its different aliases. A warning is
// logged at the warn level if multiple flags are populated.
//
// database name is resolved in accordance with the following note:
// Note: --database-name and -d are the preferred way to supply the database name.
// However, --database and --workspace can also be used as well. In the event that
// more than one of them is supplied database-name and d are preferred over all
// and workspace is preferred over database. The most preferred option that is
// supplied will be used and a message will be displayed letting you know which
// option was chosen if more than one were supplied.
func DatabaseName(cmd *cobra.Command) string {
	dbName, _ := cmd.Flags().GetString("database-name")
	workspace, _ := cmd.Flags().GetString("workspace")
	db, _ := cmd.Flags().GetString("database")
	prefName, numFlagsSupplied := "", 0
	for _, name := range []string{db, workspace, dbName} {
		if name != "" {
			numFlagsSupplied += 1
			prefName = name
		}
	}
	if numFlagsSupplied > 1 {
		logrus.Warn("multiple flags were supplied of [database-name|workspace|database], but this command may not use the expected name if multiple names are supplied")
		logrus.Warnf("this command will use name: %s", prefName)
	}
	return prefName
}
