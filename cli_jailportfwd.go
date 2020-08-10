package main

import (
	"github.com/gasiordev/go-cli"
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
	jail_portfwd_add := c.AddCmd("jail_portfwd_add", "Add port forwarding from host to jail", j.getCLIJailPortFwdAddHandler())
	jail_portfwd_add.AddArg("src_if", "SOURCE_INTERFACE", "", cli.TypeString|cli.Required)
	jail_portfwd_add.AddArg("src_port", "SOURCE_PORT", "", cli.TypeString|cli.Required)
	jail_portfwd_add.AddArg("dst_jail", "DESTINATION_JAIL", "", cli.TypeString|cli.Required)
	jail_portfwd_add.AddArg("dst_port", "DESTINATION_PORT", "", cli.TypeString|cli.Required)
	jail_portfwd_add.AddFlag("quiet", "q", "", "Do not output anything", cli.TypeBool)
	jail_portfwd_add.AddFlag("debug", "d", "", "Print more information", cli.TypeBool)

	jail_portfwd_delete := c.AddCmd("jail_portfwd_delete", "Delete port forwarding", j.getCLIJailPortFwdDeleteHandler())
	jail_portfwd_delete.AddArg("src_if", "SOURCE_INTERFACE", "", cli.TypeString|cli.Required)
	jail_portfwd_delete.AddArg("src_port", "SOURCE_PORT", "", cli.TypeString|cli.Required)
	jail_portfwd_delete.AddArg("dst_jail", "DESTINATION_JAIL", "", cli.TypeString|cli.Required)
	jail_portfwd_delete.AddArg("dst_port", "DESTINATION_PORT", "", cli.TypeString|cli.Required)
	jail_portfwd_delete.AddFlag("quiet", "q", "", "Do not output anything", cli.TypeBool)
	jail_portfwd_delete.AddFlag("debug", "d", "", "Print more information", cli.TypeBool)

	jail_portfwd_delete_all := c.AddCmd("jail_portfwd_delete_all", "Delete all port forwarding for jail", j.getCLIJailPortFwdDeleteAllHandler())
	jail_portfwd_delete_all.AddArg("dst_jail", "DESTINATION_JAIL", "", cli.TypeString|cli.Required)
	jail_portfwd_delete_all.AddFlag("quiet", "q", "", "Do not output anything", cli.TypeBool)
	jail_portfwd_delete_all.AddFlag("debug", "d", "", "Print more information", cli.TypeBool)

	jail_portfwd_list := c.AddCmd("jail_portfwd_list", "List port forwarding from host to jail", j.getCLIJailPortFwdListHandler())
	jail_portfwd_list.AddArg("dst_jail", "DESTINATION_JAIL", "", cli.TypeString|cli.Required)
	jail_portfwd_list.AddFlag("quiet", "q", "", "Do not output anything", cli.TypeBool)
	jail_portfwd_list.AddFlag("debug", "d", "", "Print more information", cli.TypeBool)
}
