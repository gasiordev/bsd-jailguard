package main

import (
	// "errors"
	"fmt"
	"github.com/gasiordev/go-cli"
	"os"
)

func getCLIVersionHandler(j *Jailguard) func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		fmt.Fprintf(c.GetStdout(), VERSION+"\n")
		return 0
	}
	return fn
}

func NewJailguardCLI(j *Jailguard) *cli.CLI {
	c := cli.NewCLI("Jailguard", "Create and manage jails in FreeBSD", "Mikolaj Gasior")

	// state list
	// state remove type.name
	// state import type.name

	// apply -f file.jail
	// destroy [-n jail OR -f file.jail]

	// network show

	// template list
	// template remove template_name
	// template create jail_name template_name
	// template update template_name base_name

	// base list
	// base download base_name
	// base update base_name

	// TODO queue:
	// * implement simple state stored in a file in JSON; simple open() and save() - all based on a struct
	// * simple removing by [type][name]
	// * state_list cmd
	// * state_remove cmd
	//
	// * base management

	_ = c.AddCmd("version", "Prints version", getCLIVersionHandler(j))

	if len(os.Args) == 2 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		os.Args = []string{"jailguard", "version"}
	}
	return c
}
