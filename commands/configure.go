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
		fmt.Println("creating new alias file...")
	}

	_, err = cfg.NewSection(c.Alias)

	if err != nil {
		return fmt.Errorf("creating alias failed")
	}

	cfg.Section(c.Alias).NewKey("arn", c.Arn)

	if err != nil {
		return fmt.Errorf("setting arn value failed")
	}

	err = cfg.SaveTo(aliasFile)

	if err != nil {
		return fmt.Errorf("saving alias file failed")
	}

	return nil
}
