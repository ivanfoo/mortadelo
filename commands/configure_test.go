package commands

import (
	"os"
	"testing"

	"gopkg.in/ini.v1"
)

const (
	testAlias     = "foobar"
	testArn       = "arn:aws:iam::777777777777:role/foobar"
	testAliasFile = "/tmp/alias_tests"
)

func TestCmdConfigure(t *testing.T) {
	cfg, _ := ini.LooseLoad(testAliasFile)
	cfg.NewSection(testAlias)
	cfg.Section(testAlias).NewKey("arn", testArn)
	cfg.SaveTo(testAliasFile)

	alias, _ := cfg.GetSection(testAlias)

	if alias.Name() != testAlias {
		t.Fatalf("Expected %s, got %s", testAlias, alias.Name)
	}

	arn, _ := alias.GetKey("arn")

	if arn.String() != testArn {
		t.Fatalf("Expected %s, got %s", testArn, arn.String())
	}

	os.Remove(testAliasFile)

}
