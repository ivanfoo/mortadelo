package commands

import (
	"fmt"

	"gopkg.in/ini.v1"
)

type CmdConfigure struct {
	Alias string `short:"a" long:"alias" description:"role alias" default:"default"`
	Arn   string `short:"r" long:"role" description:"arn role" required:"true"`
}

func (c *CmdConfigure) Execute(args []string) error {
	cfg, err := ini.LooseLoad(aliasFile)

	if err != nil {
		fmt.Println("Creating new roles file...")
	}

	cfg.NewSection(c.Alias)
	cfg.Section(c.Alias).NewKey("arn", c.Arn)

	err = cfg.SaveTo(aliasFile)

	if err != nil {
		return err
	}

	return nil
}
