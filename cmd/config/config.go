package config

import (
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/splicemachine/splicectl/auth"
	"github.com/splicemachine/splicectl/cmd/objects"
)

type (
	Config struct {
		VersionDetail objects.Version
		VersionJSON   string

		ApiServer        string
		OutputFormat     string
		FormatOverridden bool
		NoHeaders        bool
		AuthClient       auth.Client

		// tui functions
		PromptForCSP          func() (string, error)
		PromptForAccountID    func() (string, error)
		PromptForDatabaseName func() (string, error)
	}
	// Outputable - defines ways that an object may need to present itself
	Outputable interface {
		ToJSON() string
		ToYAML() string
		ToGRON() string
		ToText(noHeaders bool) string
	}
)

// GetDatabaseList - gets a list of databases
func (c *Config) GetDatabaseList() (string, error) {
	uri := "splicectl/v1/splicedb/splicedatabase"
	resp, resperr := c.
		RestyWithHeaders().
		Execute("LIST", fmt.Sprintf("%s/%s", c.ApiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error getting Database List")
		return "", resperr
	}

	return string(resp.Body()[:]), nil

}

// GetAccounts - get list of accounts
func (c *Config) GetAccounts() (string, error) {
	uri := "splicectl/v1/cm/accounts"
	resp, resperr := c.
		RestyWithHeaders().
		Get(fmt.Sprintf("%s/%s", c.ApiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error getting Account List Info")
		return "", resperr
	}

	return string(resp.Body()[:]), nil
}

// GetVersionInfo - gets version information
func (c *Config) GetVersionInfo() (string, error) {
	uri := "splicectl"
	resp, resperr := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		Get(fmt.Sprintf("%s/%s", c.ApiServer, uri))

	if resperr != nil {
		logrus.WithError(resperr).Error("Error getting version info")
		return "", resperr
	}

	return strings.TrimSuffix(string(resp.Body()[:]), "\n"), nil
}

// RestyWithHeaders - new resty request with headers for auth and content-type.
func (c *Config) RestyWithHeaders() *resty.Request {
	return resty.
		New().
		R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("X-Token-Bearer", c.AuthClient.GetTokenBearer()).
		SetHeader("X-Token-Session", c.AuthClient.GetSessionID())
}

func (c *Config) outputData(data Outputable) string {
	switch strings.ToLower(c.OutputFormat) {
	case "json":
		return data.ToJSON()
	case "gron":
		return data.ToGRON()
	case "yaml":
		return data.ToYAML()
	case "text", "table":
		return data.ToText(c.NoHeaders)
	default:
		return ""
	}
}

// OutputData - outputs string representation of data in accordance with
// OutputFormat.
func (c *Config) OutputData(data Outputable) {
	if !c.FormatOverridden {
		c.OutputFormat = "text"
	}

	fmt.Println(c.outputData(data))
}
