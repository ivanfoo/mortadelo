package commands

import (
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

type CmdConfigure struct {
	Alias string `short:"a" long:"alias" description:"role alias" required:"true"`
	Arn   string `short:"r" long:"role" description:"arn role" required:"true"`
	File  string `short:"f" long:"file" description:"alias file" default:"~/.mortadelo/alias"`

	filePath string
}

func (c *CmdConfigure) Execute(args []string) error {
	c.filePath = expandPath(c.File)

	err := c.createHome()
	if err != nil {
		return err
	}

	err = c.setAlias()
	if err != nil {
		return err
	}

	return nil
}

func (c *CmdConfigure) createHome() error {
	path := filepath.Dir(c.filePath)
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		err = os.Mkdir(path, 0755)
		if err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}

func (c *CmdConfigure) setAlias() error {
	cfg, err := ini.LooseLoad(c.filePath)

	if err != nil {
		return err
	}

	_, err = cfg.NewSection(c.Alias)
	if err != nil {
		return err
	}

	cfg.Section(c.Alias).NewKey("arn", c.Arn)
	if err != nil {
		return err
	}

	err = cfg.SaveTo(c.filePath)
	if err != nil {
		return err
	}

	return nil
}
