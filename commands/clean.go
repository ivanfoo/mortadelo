package commands

import (
	"fmt"
	"io/ioutil"
	"os"
)

type CleanCommand struct {
	Force  bool   `long:"force" description:"do not ask for permission before cleaning"`
	File   string `short:"f" long:"file" description:"alias file" default:"~/.mortadelo/alias"`
	Backup string `short:"b" long:"backup" description:"create backup file" default:""`
}

func (c *CleanCommand) Execute(args []string) error {
	filePath := expandPath(c.File)

	var ok bool
	if c.Force {
		ok = true
	} else {
		ok = c.askUser(filePath)
	}

	if ok {
		if c.Backup != "" {
			err := c.backup(filePath)
			if err != nil {
				return err
			}
		}

		fmt.Printf("Removing %s...", filePath)
		err := os.Remove(filePath)
		if err != nil {
			return err
		}

		fmt.Println("Done")
	}

	return nil
}

func (c *CleanCommand) askUser(path string) bool {
	fmt.Printf("Delete %s (yes/no): ", path)

	var input string
	fmt.Scanf("%s", &input)

	if input == "yes" {
		return true
	} else {
		fmt.Println("Aborted")
		return false
	}
}

func (c *CleanCommand) backup(path string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.Backup, content, 0644)
	if err != nil {
		return err
	}

	return nil
}
