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

func (j *Jailguard) getCLIStateListHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		if c.Flag("debug") == "true" {
			j.Debug = true
		}

		err := j.ListStateItems()
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) getCLIStateRemoveHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		if c.Flag("debug") == "true" {
			j.Debug = true
		}

		err := j.RemoveStateItem(c.Arg("item_type"), c.Arg("item_name"))
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) getCLIStateImportHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		if c.Flag("debug") == "true" {
			j.Debug = true
		}

		j.Log(LOGERR, "Not implemented")
		return 10
	}
	return fn
}

func (j *Jailguard) getCLIBaseDownloadHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		if c.Flag("debug") == "true" {
			j.Debug = true
		}
		if c.Flag("quiet") == "true" {
			j.Quiet = true
		}

		ow := false
		if c.Flag("overwrite") == "true" {
			ow = true
		}

		err := j.DownloadBase(c.Arg("release"), ow)
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func NewJailguardCLI(j *Jailguard) *cli.CLI {
	c := cli.NewCLI("Jailguard", "Create and manage jails in FreeBSD", "Mikolaj Gasior")

	// state commands
	_ = c.AddCmd("state_list", "Lists saved state items", j.getCLIStateListHandler())

	state_remove := c.AddCmd("state_remove", "Remove item from state", j.getCLIStateRemoveHandler())
	state_remove.AddArg("item_type", "TYPE", "", cli.TypeString|cli.Required)
	state_remove.AddArg("item_name", "NAME", "", cli.TypeString|cli.Required)
	state_remove.AddFlag("quiet", "q", "", "Do not output anything", cli.TypeBool)
	state_remove.AddFlag("debug", "d", "", "Print more information", cli.TypeBool)

	state_import := c.AddCmd("state_import", "Import item to state", j.getCLIStateImportHandler())
	state_import.AddArg("item_type", "TYPE", "", cli.TypeString|cli.Required)
	state_import.AddArg("item_name", "NAME", "", cli.TypeAlphanumeric|cli.Required)
	state_import.AddFlag("quiet", "q", "", "Do not output anything", cli.TypeBool)
	state_import.AddFlag("debug", "d", "", "Print more information", cli.TypeBool)

	// base commands
	base_download := c.AddCmd("base_download", "Downloads FreeBSD base", j.getCLIBaseDownloadHandler())
	// TODO: Change 'release' flag to TypeAlphanumeric once AllowHyphen gets implemented in go-cli
	base_download.AddArg("release", "RELEASE", "", cli.TypeString|cli.Required)
	base_download.AddFlag("overwrite", "w", "", "Overwrite if exists", cli.TypeBool)
	base_download.AddFlag("quiet", "q", "", "Do not output anything", cli.TypeBool)
	base_download.AddFlag("debug", "d", "", "Print more information", cli.TypeBool)

	// state import type.name

	// apply -f file.jail
	// destroy [-n jail OR -f file.jail]

	// network show

	// template list
	// template remove template_name
	// template create jail_name template_name
	// template update template_name base_name

	// base list
	// base remove

	// TODO queue:
	// base_remove - with -y otherwise ask for removal (cli does not have Stdin, new issue?)
	// base_list -> needs to call state_list with 'base' filter on it
	// state_import -> just like download: check if exists in state, create new base obj, check if valid, add base, save

	_ = c.AddCmd("version", "Prints version", getCLIVersionHandler(j))

	if len(os.Args) == 2 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		os.Args = []string{"jailguard", "version"}
	}
	return c
}
