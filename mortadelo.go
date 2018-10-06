package main

import (
	"os"

	"github.com/ivanfoo/mortadelo/commands"

	"github.com/jessevdk/go-flags"
)

func main() {
	parser := flags.NewNamedParser("mortadelo", flags.Default)
	parser.AddCommand("assume", "assume role", "", &commands.AssumeCommand{})
	parser.AddCommand("clean", "clean generated files", "", &commands.CleanCommand{})
	parser.AddCommand("configure", "configure roles alias file", "", &commands.ConfigureCommand{})

	_, err := parser.Parse()
	if err != nil {
		os.Exit(1)
	}
}
