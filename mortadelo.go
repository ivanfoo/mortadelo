package main

import (
	"log"

	"github.com/ivanfoo/mortadelo/commands"

	"github.com/jessevdk/go-flags"
)

func main() {
	parser := flags.NewNamedParser("mortadelo", flags.Default)
	parser.AddCommand("assume", "assume role", "", &commands.CmdAssume{})
	parser.AddCommand("clean", "clean generated files", "", &commands.CmdClean{})

	_, err := parser.Parse()
	if err != nil {
		log.Fatal(err)
	}
}
