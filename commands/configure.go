package commands

import "gopkg.in/ini.v1"

type CmdConfigure struct {
	Alias     string `short:"a" long:"alias" description:"role alias" default:"default"`
	Arn       string `short:"r" long:"role" description:"arn role" required:"true"`
	RolesFile string `short:"f" long:"file" description:"roles file" default:"~/.mortadelo/roles"`
}

func (c *CmdConfigure) Execute(args []string) error {
	cfg := ini.Empty()
	cfg.NewSection(c.Alias)
	cfg.Section(c.Alias).NewKey("arn", c.Arn)

	err := cfg.SaveTo(c.RolesFile)

	if err != nil {
		return err
	}

	return nil
}
