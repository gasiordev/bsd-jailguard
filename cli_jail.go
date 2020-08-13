package main

import (
	"errors"
	"github.com/nicholasgasior/go-cli"
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
	create := c.AddCmd("jail_create", "Create jail source", j.getCLIJailCreateHandler())
	create.AddArg("file", "FILE.JAIL", "", cli.TypePathFile|cli.MustExist|cli.Required)
	create.AddFlag("base", "b", "", "Base to use", cli.TypeAlphanumeric|cli.AllowDots|cli.AllowUnderscore|cli.AllowHyphen)
	create.AddFlag("start", "s", "", "Start jail after creating", cli.TypeBool)

	remove := c.AddCmd("jail_remove", "Remove jail source", j.getCLIJailRemoveHandler())
	remove.AddArg("jail", "JAIL", "", cli.TypeString|cli.Required)
	remove.AddFlag("stop", "s", "", "Stop if running", cli.TypeBool)

	stop := c.AddCmd("jail_stop", "Stop jail", j.getCLIJailStopHandler())
	stop.AddArg("jail", "JAIL", "", cli.TypeString|cli.Required)

	start := c.AddCmd("jail_start", "Start jail", j.getCLIJailStartHandler())
	start.AddArg("jail", "JAIL", "", cli.TypeString|cli.Required)

	fn := func(c *cli.CLI) error {
		if !IsValidJailName(c.Arg("jail")) {
			return errors.New("Argument JAIL is not a valid jail name")
		}
		return nil
	}
	create.AddPostValidation(fn)
	remove.AddPostValidation(fn)
	stop.AddPostValidation(fn)
	start.AddPostValidation(fn)

	_ = c.AddCmd("jail_list", "List jails", j.getCLIJailListHandler())
}
