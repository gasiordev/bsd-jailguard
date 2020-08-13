package main

import (
	"errors"
	"github.com/nicholasgasior/go-cli"
)

func (j *Jailguard) getCLIJailNATPassCreateHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		if c.Flag("debug") == "true" {
			j.Debug = true
		}
		if c.Flag("quiet") == "true" {
			j.Quiet = true
		}

		err := j.CreateJailNATPass(c.Arg("jail"), c.Arg("gw_if"))
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) getCLIJailNATPassRemoveHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		if c.Flag("debug") == "true" {
			j.Debug = true
		}
		if c.Flag("quiet") == "true" {
			j.Quiet = true
		}

		err := j.RemoveJailNATPass(c.Arg("jail"))
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) getCLIJailNATPassShowHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		if c.Flag("debug") == "true" {
			j.Debug = true
		}
		if c.Flag("quiet") == "true" {
			j.Quiet = true
		}

		err := j.ShowJailNATPass(c.Arg("jail"))
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) AddJailNATPassCmds(c *cli.CLI) {
	create := c.AddCmd("jail_natpass_create", "Create NAT pass for jail", j.getCLIJailNATPassCreateHandler())
	create.AddArg("jail", "JAIL", "", cli.TypeString|cli.Required)
	create.AddArg("gw_if", "GATEWAY_INTERFACE", "", cli.TypeString|cli.Required)

	remove := c.AddCmd("jail_natpass_remove", "Remove NAT pass from jail", j.getCLIJailNATPassRemoveHandler())
	remove.AddArg("jail", "JAIL", "", cli.TypeString|cli.Required)

	show := c.AddCmd("jail_natpass_show", "Show NAT pass gateway for jail", j.getCLIJailNATPassShowHandler())
	show.AddArg("jail", "JAIL", "", cli.TypeString|cli.Required)

	fn := func(c *cli.CLI) error {
		if !IsValidJailName(c.Arg("jail")) {
			return errors.New("Argument JAIL is not a valid jail name")
		}
		return nil
	}
	create.AddPostValidation(fn)
	remove.AddPostValidation(fn)
	show.AddPostValidation(fn)
}
