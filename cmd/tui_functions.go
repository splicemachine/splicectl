package cmd

import (
	"encoding/json"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/sirupsen/logrus"
	"github.com/splicemachine/splicectl/cmd/objects"
)

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
