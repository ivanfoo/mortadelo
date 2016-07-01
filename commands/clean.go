package commands

import (
	"fmt"
	"os"
)

type CmdClean struct {
	Force bool `short:"f" long:"force" description:"do not ask for permission before cleaning"`
}

func (c *CmdClean) Execute(args []string) error {
	if c.Force {
		os.Remove(awsCredentialsFile)
		fmt.Println("deleted " + awsCredentialsFile)

		os.Remove(aliasFile)
		fmt.Println("deleted " + aliasFile)

		return nil
	}

	if userSayingYesToMessage("About to delete " + aliasFile + ". You sure? [y/n] ") {
		os.Remove(aliasFile)
		fmt.Println("done")
	}

	if userSayingYesToMessage("About to delete " + awsCredentialsFile + ". You sure? [y/n] ") {
		os.Remove(awsCredentialsFile)
		fmt.Println("done")
	}

	return nil
}

func userSayingYesToMessage(message string) bool {
	fmt.Print(message)
	var userInput string
	fmt.Scanf("%s", &userInput)

	if userInput == "y" {
		return true
	}

	return false
}
