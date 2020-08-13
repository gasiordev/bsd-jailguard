package main

import (
	"github.com/nicholasgasior/go-cli"
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

	st_remove := c.AddCmd("state_remove", "Remove item from state", j.getCLIStateRemoveHandler())
	st_remove.AddArg("item_type", "TYPE", "", cli.TypeAlphanumeric|cli.AllowUnderscore|cli.AllowHyphen|cli.Required)
	st_remove.AddArg("item_name", "NAME", "", cli.TypeAlphanumeric|cli.AllowUnderscore|cli.AllowHyphen|cli.Required)

	st_import := c.AddCmd("state_import", "Import item to state", j.getCLIStateImportHandler())
	st_import.AddArg("item_type", "TYPE", "", cli.TypeAlphanumeric|cli.AllowUnderscore|cli.AllowHyphen|cli.Required)
	st_import.AddArg("item_name", "NAME", "", cli.TypeAlphanumeric|cli.AllowUnderscore|cli.AllowHyphen|cli.Required)
}
