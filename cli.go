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

	//        cmdUp := c.AddCmd("up", "Creates and starts jails", getCLIUpHandler(b))
	//        cmdUp.AddFlag("jailguardfile", "f", "file", "Jailguardfile", cli.TypePathFile|cli.MustExist)
	//        cmdUp.AddFlag("quiet", "q", "", "Do not output anything", cli.TypeBool)
	//        cmdUp.AddFlag("debug", "d", "", "Print more information", cli.TypeBool)
	//
	//        cmdStop := c.AddCmd("stop", "Stop jails", getCLIStopHandler(b))
	//        cmdStop.AddFlag("jailguardfile", "f", "file", "Jailguardfile", cli.TypePathFile|cli.MustExist)
	//        cmdStop.AddFlag("quiet", "q", "", "Do not output anything", cli.TypeBool)
	//        cmdStop.AddFlag("debug", "d", "", "Print more information", cli.TypeBool)
	//
	//        cmdShow := c.AddCmd("show", "Shows jails", getCLIShowHandler(b))
	//        cmdShow.AddFlag("jailguardfile", "f", "file", "Jailguardfile", cli.TypePathFile|cli.MustExist)
	//        cmdShow.AddFlag("quiet", "q", "", "Do not output anything", cli.TypeBool)
	//        cmdShow.AddFlag("debug", "d", "", "Print more information", cli.TypeBool)
	//
	//        cmdDestroy := c.AddCmd("destroy", "Destroys jails", getCLIDestroyHandler(b))
	//        cmdDestroy.AddFlag("jailguardfile", "f", "file", "Jailguardfile", cli.TypePathFile|cli.MustExist)
	//        cmdDestroy.AddFlag("quiet", "q", "", "Do not output anything", cli.TypeBool)
	//        cmdDestroy.AddFlag("debug", "d", "", "Print more information", cli.TypeBool)
	//
	//        cmdExec := c.AddCmd("exec", "Exec commands on jail", getCLIExecHandler(b))
	//        cmdExec.AddFlag("jailguardfile", "f", "file", "Jailguardfile", cli.TypePathFile|cli.MustExist)
	//        cmdExec.AddArg("guest", "JAIL", "", cli.TypeString|cli.Required)
	//        cmdExec.AddArg("command", "COMMAND", "", cli.TypeString|cli.Required)

	_ = c.AddCmd("version", "Prints version", getCLIVersionHandler(j))

	if len(os.Args) == 2 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		os.Args = []string{"jailguard", "version"}
	}
	return c
}
