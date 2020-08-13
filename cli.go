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

	c.AddFlagToCmds("quiet", "q", "", "Do not output anything", cli.TypeBool)
	c.AddFlagToCmds("debug", "d", "", "Print more information", cli.TypeBool)

	// TODO
	// - get free ip address from interface begin-end range when ip address is not passed
	// - ip addr field validation
	// - jail_portfwd
	// - jail_natpass

	// whenever app is started, check for the configuration - if not then write about running
	// 'configure' or passing '--config'

	// whenever app is started, check if the configuration has correct values

	// config_list
	// config_get key            ------> /etc/jailguard.conf otherwise --config for all the commands
	// config_set key value
	// configure - before first use, tool has to be configured - all values can be passed interactively
	//   or as flags

	// when checking the jail conf: - if interface exists and if not then it has to be created
	//                              - if pf entries exist then pf has to be installed
	//                                                         and pf anchor created
	//                              - if networking then check the sysctl for allow_raw_sockets and if not present then display warning (which can be supressed by -q)

	// state_check for base and jail: checks if state is up-to-date

	_ = c.AddCmd("version", "Prints version", getCLIVersionHandler(j))

	if len(os.Args) == 2 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		os.Args = []string{"jailguard", "version"}
	}
	return c
}
