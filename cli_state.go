package main

import (
	"github.com/gasiordev/go-cli"
)

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

func (j *Jailguard) AddStateCmds(c *cli.CLI) {
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
}
