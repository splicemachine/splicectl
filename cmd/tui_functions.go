package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/sirupsen/logrus"
	"github.com/splicemachine/splicectl/cmd/objects"
)

type (
	cspAnswer struct {
		CSP string `survey:"cspselect"`
	}
	accountIDAnswer struct {
		AccountID string `survey:"accountid"`
	}
	databaseAnswer struct {
		DatabaseName string `survey:"databasename"` // or you can tag fields to match a specific name
	}
)

var (
	cspList = []string{
		"NONE",
		"OP",
		"AWS",
		"AZ",
		"GCP",
	}
	cspSurvey = []*survey.Question{
		{
			Name: "cspselect",
			Prompt: &survey.Select{
				Message: "Choose a Cloud Service Provider:",
				Options: cspList,
			},
		},
	}
	cspAnswers       = &cspAnswer{}
	accountIDAnswers = &accountIDAnswer{}
	databaseAnswers  = &databaseAnswer{}
)

// PromptForCSP - prompt the user on the command line for cloud service provider
func PromptForCSP() (string, error) {
	opts := survey.WithStdio(os.Stdin, os.Stderr, os.Stderr)
	if err := survey.Ask(cspSurvey, cspAnswers, opts); err != nil {
		return "", err
	}
	return cspAnswers.CSP, nil
}

// PromptForAccountID - prompt the user on the command line for account id
func PromptForAccountID() (string, error) {
	out, err := c.GetAccounts()
	if err != nil {
		logrus.WithError(err).Error("Error getting Default CR Info")
		return "", err
	}

	var accounts objects.AccountList
	if err := json.Unmarshal([]byte(out), &accounts); err != nil {
		logrus.Fatal("Could not unmarshall data", err)
	}

	acctArray := make([]string, 0)
	for _, v := range accounts.Accounts {
		accountID := fmt.Sprintf("%s (%s, %s <%s>)", v.AccountID, v.LastName, v.FirstName, v.EMail)
		acctArray = append(acctArray, accountID)
	}

	accountIDSurvey := []*survey.Question{
		{
			Name: "accountid",
			Prompt: &survey.Select{
				Message: "Choose an account:",
				Options: acctArray,
			},
		},
	}

	opts := survey.WithStdio(os.Stdin, os.Stderr, os.Stderr)

	// perform the questions
	if err = survey.Ask(accountIDSurvey, accountIDAnswers, opts); err != nil {
		return "", err
	}

	acctID := strings.TrimSpace(strings.Split(accountIDAnswers.AccountID, "(")[0])
	return acctID, nil
}

// PromptForDatabaseName - prompt the user on the command line for name
func PromptForDatabaseName() (string, error) {
	out, err := c.GetDatabaseList()
	if err != nil {
		logrus.WithError(err).Error("Error getting Database List")
		return "", err
	}
	var dbList objects.DatabaseList

	marshErr := json.Unmarshal([]byte(out), &dbList)
	if marshErr != nil {
		logrus.Fatal("Could not unmarshall data", marshErr)
	}
	var dbArray []string
	for _, v := range dbList.Clusters {
		dbArray = append(dbArray, v.DcosAppId)
	}

	// the questions to ask
	var databaseSurvey = []*survey.Question{
		{
			Name: "databasename",
			Prompt: &survey.Select{
				Message: "Choose a database:",
				Options: dbArray,
			},
		},
	}

	opts := survey.WithStdio(os.Stdin, os.Stderr, os.Stderr)

	// perform the questions
	if err = survey.Ask(databaseSurvey, databaseAnswers, opts); err != nil {
		logrus.Fatal("No databases on the list")
	}
	return databaseAnswers.DatabaseName, nil
}

func addTUIFunctionsToConfig() {
	c.PromptForCSP = PromptForCSP
	c.PromptForAccountID = PromptForAccountID
	c.PromptForDatabaseName = PromptForDatabaseName
}
