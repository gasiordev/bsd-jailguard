package main

import (
	"github.com/gasiordev/go-cli"
)

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

func (j *Jailguard) AddJailCmds(c *cli.CLI) {
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

	_ = c.AddCmd("jail_list", "List jails", j.getCLIJailListHandler())
}
