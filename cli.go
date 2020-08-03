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

		err := j.ListStateItems("")
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

		err := j.ImportStateItem(c.Arg("item_type"), c.Arg("item_name"))
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}

		return 0
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

func (j *Jailguard) getCLIBaseRemoveHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		if c.Flag("debug") == "true" {
			j.Debug = true
		}
		if c.Flag("quiet") == "true" {
			j.Quiet = true
		}

		err := j.RemoveBase(c.Arg("release"))
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) getCLIBaseListHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		err := j.ListStateItems("bases")
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) getCLIJailListHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		err := j.ListStateItems("jails")
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) getCLIJailCreateHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		if c.Flag("debug") == "true" {
			j.Debug = true
		}
		if c.Flag("quiet") == "true" {
			j.Quiet = true
		}

		start := false
		if c.Flag("start") == "true" {
			start = true
		}
		err := j.CreateJail(c.Arg("file"), c.Flag("base"), start)
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) getCLIJailRemoveHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		if c.Flag("debug") == "true" {
			j.Debug = true
		}
		if c.Flag("quiet") == "true" {
			j.Quiet = true
		}

		stop := false
		if c.Flag("stop") == "true" {
			stop = true
		}
		err := j.RemoveJail(c.Arg("jail"), stop)
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) getCLIJailStopHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		if c.Flag("debug") == "true" {
			j.Debug = true
		}
		if c.Flag("quiet") == "true" {
			j.Quiet = true
		}

		err := j.StopJail(c.Arg("jail"))
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) getCLIJailStartHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		if c.Flag("debug") == "true" {
			j.Debug = true
		}
		if c.Flag("quiet") == "true" {
			j.Quiet = true
		}

		err := j.StartJail(c.Arg("jail"))
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

	base_download := c.AddCmd("base_download", "Downloads FreeBSD base", j.getCLIBaseDownloadHandler())
	// TODO: Change 'release' flag to TypeAlphanumeric once AllowHyphen gets implemented in go-cli
	base_download.AddArg("release", "RELEASE", "", cli.TypeString|cli.Required)
	base_download.AddFlag("overwrite", "w", "", "Overwrite if exists", cli.TypeBool)
	base_download.AddFlag("quiet", "q", "", "Do not output anything", cli.TypeBool)
	base_download.AddFlag("debug", "d", "", "Print more information", cli.TypeBool)

	_ = c.AddCmd("base_list", "List FreeBSD bases", j.getCLIBaseListHandler())

	base_remove := c.AddCmd("base_remove", "Removes FreeBSD base", j.getCLIBaseRemoveHandler())
	base_remove.AddArg("release", "RELEASE", "", cli.TypeString|cli.Required)
	base_remove.AddFlag("quiet", "q", "", "Do not output anything", cli.TypeBool)
	base_remove.AddFlag("debug", "d", "", "Print more information", cli.TypeBool)

	jail_create := c.AddCmd("jail_create", "Create jail source", j.getCLIJailCreateHandler())
	jail_create.AddArg("file", "FILE.JAIL", "", cli.TypePathFile|cli.MustExist|cli.Required)
	jail_create.AddFlag("quiet", "q", "", "Do not output anything", cli.TypeBool)
	jail_create.AddFlag("debug", "d", "", "Print more information", cli.TypeBool)
	jail_create.AddFlag("base", "b", "", "Base to use", cli.TypeString)
	jail_create.AddFlag("start", "s", "", "Start jail after creating", cli.TypeBool)

	jail_remove := c.AddCmd("jail_remove", "Remove jail source", j.getCLIJailRemoveHandler())
	jail_remove.AddArg("jail", "JAIL", "", cli.TypeString|cli.Required)
	jail_remove.AddFlag("quiet", "q", "", "Do not output anything", cli.TypeBool)
	jail_remove.AddFlag("debug", "d", "", "Print more information", cli.TypeBool)
	jail_remove.AddFlag("stop", "s", "", "Stop if running", cli.TypeBool)

	jail_stop := c.AddCmd("jail_stop", "Stop jail", j.getCLIJailStopHandler())
	jail_stop.AddArg("jail", "JAIL", "", cli.TypeString|cli.Required)
	jail_stop.AddFlag("quiet", "q", "", "Do not output anything", cli.TypeBool)
	jail_stop.AddFlag("debug", "d", "", "Print more information", cli.TypeBool)

	jail_start := c.AddCmd("jail_start", "Start jail", j.getCLIJailStartHandler())
	jail_start.AddArg("jail", "JAIL", "", cli.TypeString|cli.Required)
	jail_start.AddFlag("quiet", "q", "", "Do not output anything", cli.TypeBool)
	jail_start.AddFlag("debug", "d", "", "Print more information", cli.TypeBool)

	jail_list := c.AddCmd("jail_list", "List jails", j.getCLIJailListHandler())
	jail_list.AddFlag("all", "a", "", "Show all", cli.TypeBool)

	// state_check for base and jail: checks if state is up-to-date

	// networking

	_ = c.AddCmd("version", "Prints version", getCLIVersionHandler(j))

	if len(os.Args) == 2 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		os.Args = []string{"jailguard", "version"}
	}
	return c
}
