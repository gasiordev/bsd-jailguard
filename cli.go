package main

import (
	// "errors"
	"fmt"
	"github.com/nicholasgasior/go-cli"
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
	c := cli.NewCLI("Jailguard", "Create and manage jails in FreeBSD", "Nicholas Gasior")

	j.AddStateCmds(c)
	j.AddBaseCmds(c)
	j.AddJailCmds(c)
	j.AddNetifCmds(c)
	j.AddPFAnchorCmds(c)
	j.AddJailPortFwdCmds(c)
	j.AddJailNATPassCmds(c)
	j.AddConfigCmds(c)

	c.AddFlagToCmds("quiet", "q", "", "Do not output anything", cli.TypeBool)
	c.AddFlagToCmds("debug", "d", "", "Print more information", cli.TypeBool)

	_ = c.AddCmd("version", "Prints version", getCLIVersionHandler(j))

	if len(os.Args) == 2 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		os.Args = []string{"jailguard", "version"}
	}
	return c
}
