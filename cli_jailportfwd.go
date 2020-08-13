package main

import (
	"errors"
	"github.com/nicholasgasior/go-cli"
)

func (j *Jailguard) getCLIJailPortFwdAddHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		if c.Flag("debug") == "true" {
			j.Debug = true
		}
		if c.Flag("quiet") == "true" {
			j.Quiet = true
		}

		err := j.AddJailPortFwd(c.Arg("src_if"), c.Arg("src_port"), c.Arg("dst_jail"), c.Arg("dst_port"))
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) getCLIJailPortFwdDeleteHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		if c.Flag("debug") == "true" {
			j.Debug = true
		}
		if c.Flag("quiet") == "true" {
			j.Quiet = true
		}

		err := j.DeleteJailPortFwd(c.Arg("src_if"), c.Arg("src_port"), c.Arg("dst_jail"), c.Arg("dst_port"))
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) getCLIJailPortFwdDeleteAllHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		if c.Flag("debug") == "true" {
			j.Debug = true
		}
		if c.Flag("quiet") == "true" {
			j.Quiet = true
		}

		err := j.DeleteJailAllPortFwds(c.Arg("dst_jail"))
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) getCLIJailPortFwdListHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		if c.Flag("debug") == "true" {
			j.Debug = true
		}
		if c.Flag("quiet") == "true" {
			j.Quiet = true
		}

		err := j.ListJailPortFwds(c.Arg("dst_jail"))
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) AddJailPortFwdCmds(c *cli.CLI) {
	add := c.AddCmd("jail_portfwd_add", "Add port forwarding from host to jail", j.getCLIJailPortFwdAddHandler())
	add.AddArg("src_if", "SOURCE_INTERFACE", "", cli.TypeString|cli.Required)
	add.AddArg("src_port", "SOURCE_PORT", "", cli.TypeString|cli.Required)
	add.AddArg("dst_jail", "DESTINATION_JAIL", "", cli.TypeString|cli.Required)
	add.AddArg("dst_port", "DESTINATION_PORT", "", cli.TypeString|cli.Required)

	delete := c.AddCmd("jail_portfwd_delete", "Delete port forwarding", j.getCLIJailPortFwdDeleteHandler())
	delete.AddArg("src_if", "SOURCE_INTERFACE", "", cli.TypeString|cli.Required)
	delete.AddArg("src_port", "SOURCE_PORT", "", cli.TypeString|cli.Required)
	delete.AddArg("dst_jail", "DESTINATION_JAIL", "", cli.TypeString|cli.Required)
	delete.AddArg("dst_port", "DESTINATION_PORT", "", cli.TypeString|cli.Required)

	delete_all := c.AddCmd("jail_portfwd_delete_all", "Delete all port forwarding for jail", j.getCLIJailPortFwdDeleteAllHandler())
	delete_all.AddArg("dst_jail", "DESTINATION_JAIL", "", cli.TypeString|cli.Required)

	list := c.AddCmd("jail_portfwd_list", "List port forwarding from host to jail", j.getCLIJailPortFwdListHandler())
	list.AddArg("dst_jail", "DESTINATION_JAIL", "", cli.TypeString|cli.Required)

	fn := func(c *cli.CLI) error {
		if !IsValidJailName(c.Arg("dst_jail")) {
			return errors.New("Argument DESTINATION_JAIL is not a valid jail name")
		}
		return nil
	}
	add.AddPostValidation(fn)
	delete.AddPostValidation(fn)
	delete_all.AddPostValidation(fn)
	list.AddPostValidation(fn)
}
