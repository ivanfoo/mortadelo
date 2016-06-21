package commands

import "os"

type CmdClean struct {
	//Interactive bool `short:"i" long:"interactive" description:"ask for permission before cleaning" default:"false"`
}

func (c *CmdClean) Execute(args []string) error {
	os.Remove(rolesFile)
	os.Remove(awsCredentialsFile)
	return nil
}
