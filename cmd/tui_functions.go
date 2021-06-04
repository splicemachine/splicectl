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

func promptForCSP() (string, error) {

	cspList := []string{
		"NONE",
		"OP",
		"AWS",
		"AZ",
		"GCP",
	}

	// the questions to ask
	var qs = []*survey.Question{
		{
			Name: "cspselect",
			Prompt: &survey.Select{
				Message: "Choose a Cloud Service Provider:",
				Options: cspList,
			},
		},
	}
	// the answers will be written to this struct
	answers := struct {
		CSP string `survey:"cspselect"` // or you can tag fields to match a specific name
	}{}

	opts := survey.WithStdio(os.Stdin, os.Stderr, os.Stderr)

	// perform the questions
	err := survey.Ask(qs, &answers, opts)
	if err != nil {
		// logrus.Fatal("No accounts on the list")
		// // fmt.Println(err.Error())
		// // return
		return "", err
	}

	return answers.CSP, nil

}

func promptForAccountID() (string, error) {
	out, err := getAccounts()
	if err != nil {
		logrus.WithError(err).Error("Error getting Default CR Info")
		return "", err
	}

	var accounts objects.AccountList

	marshErr := json.Unmarshal([]byte(out), &accounts)
	if marshErr != nil {
		logrus.Fatal("Could not unmarshall data", marshErr)
	}

	var acctArray []string
	for _, v := range accounts.Accounts {
		acctArray = append(acctArray, fmt.Sprintf("%s (%s, %s <%s>)", v.AccountID, v.LastName, v.FirstName, v.EMail))
	}
	// the questions to ask
	var qs = []*survey.Question{
		{
			Name: "accountid",
			Prompt: &survey.Select{
				Message: "Choose an account:",
				Options: acctArray,
			},
		},
	}
	// the answers will be written to this struct
	answers := struct {
		AccountID string `survey:"accountid"` // or you can tag fields to match a specific name
	}{}

	opts := survey.WithStdio(os.Stdin, os.Stderr, os.Stderr)

	// perform the questions
	err = survey.Ask(qs, &answers, opts)
	if err != nil {
		// logrus.Fatal("No accounts on the list")
		// // fmt.Println(err.Error())
		// // return
		return "", err
	}

	acctID := strings.TrimSpace(strings.Split(answers.AccountID, "(")[0])
	return acctID, nil
}

func promptForDatabaseName() (string, error) {
	out, err := getDatabaseList()
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
	var qs = []*survey.Question{
		{
			Name: "databasename",
			Prompt: &survey.Select{
				Message: "Choose a database:",
				Options: dbArray,
			},
		},
	}
	// the answers will be written to this struct
	answers := struct {
		DatabaseName string `survey:"databasename"` // or you can tag fields to match a specific name
	}{}

	// opts := survey.AskOptions{
	// 	Stdio: terminal.Stdio{
	// 		In:  os.Stdin,
	// 		Out: os.Stderr,
	// 		Err: os.Stderr,
	// 	},
	// }
	opts := survey.WithStdio(os.Stdin, os.Stderr, os.Stderr)

	// perform the questions
	err = survey.Ask(qs, &answers, opts)
	if err != nil {
		logrus.Fatal("No databases on the list")
		// fmt.Println(err.Error())
		// return
	}
	return answers.DatabaseName, nil

}
