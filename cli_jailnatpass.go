package main

import (
	"github.com/gasiordev/go-cli"
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
	jail_natpass_create := c.AddCmd("jail_natpass_create", "Create NAT pass for jail", j.getCLIJailNATPassCreateHandler())
	jail_natpass_create.AddArg("jail", "JAIL", "", cli.TypeString|cli.Required)
	jail_natpass_create.AddArg("gw_if", "GATEWAY_INTERFACE", "", cli.TypeString|cli.Required)
	jail_natpass_create.AddFlag("quiet", "q", "", "Do not output anything", cli.TypeBool)
	jail_natpass_create.AddFlag("debug", "d", "", "Print more information", cli.TypeBool)

	jail_natpass_remove := c.AddCmd("jail_natpass_remove", "Remove NAT pass from jail", j.getCLIJailNATPassRemoveHandler())
	jail_natpass_remove.AddArg("jail", "JAIL", "", cli.TypeString|cli.Required)
	jail_natpass_remove.AddFlag("quiet", "q", "", "Do not output anything", cli.TypeBool)
	jail_natpass_remove.AddFlag("debug", "d", "", "Print more information", cli.TypeBool)

	jail_natpass_show := c.AddCmd("jail_natpass_show", "Show NAT pass gateway for jail", j.getCLIJailNATPassShowHandler())
	jail_natpass_show.AddArg("jail", "JAIL", "", cli.TypeString|cli.Required)
	jail_natpass_show.AddFlag("debug", "d", "", "Print more information", cli.TypeBool)
}
