package main

import (
	"log"

	"github.com/ivanfoo/joe/commands"

	"github.com/jessevdk/go-flags"
)

func main() {
	parser := flags.NewNamedParser("joe", flags.Default)
	parser.AddCommand("assume", "assume role", "", commands.NewCmdAssume())

	_, err := parser.Parse()
	if err != nil {
		log.Fatal(err)
	}
}
