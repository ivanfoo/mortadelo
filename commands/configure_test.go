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

func TestConfigureCommand(t *testing.T) {
	cmd := &ConfigureCommand{
		Alias: testAlias,
		Arn:   testArn,
		File:  testAliasFile,
	}

	cmd.filePath = testAliasFile
	cmd.setAlias()

	cfg, _ := ini.LooseLoad(testAliasFile)
	alias, _ := cfg.GetSection(testAlias)

	if alias.Name() != testAlias {
		t.Fatalf("Expected %s, got %s", testAlias, alias.Name())
	}

	arn, _ := alias.GetKey("arn")

	if arn.String() != testArn {
		t.Fatalf("Expected %s, got %s", testArn, arn.String())
	}

	os.Remove(testAliasFile)

}
