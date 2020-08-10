package main

import (
	"github.com/gasiordev/go-cli"
)

func (j *Jailguard) getCLINetifListHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		err := j.ListStateItems("netifs")
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) getCLINetifCreateHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		if c.Flag("debug") == "true" {
			j.Debug = true
		}
		if c.Flag("quiet") == "true" {
			j.Quiet = true
		}

		err := j.CreateNetif(c.Arg("name"), c.Arg("ip_addr_begin"), c.Arg("ip_addr_end"), c.Arg("if_name"))
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) getCLINetifDestroyHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		if c.Flag("debug") == "true" {
			j.Debug = true
		}
		if c.Flag("quiet") == "true" {
			j.Quiet = true
		}

		err := j.DestroyNetif(c.Arg("name"))
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) getCLINetifAliasAddHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		if c.Flag("debug") == "true" {
			j.Debug = true
		}
		if c.Flag("quiet") == "true" {
			j.Quiet = true
		}

		err := j.AddNetifAlias(c.Arg("name"), c.Arg("ip_addr"))
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) getCLINetifAliasDeleteHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		if c.Flag("debug") == "true" {
			j.Debug = true
		}
		if c.Flag("quiet") == "true" {
			j.Quiet = true
		}

		err := j.DeleteNetifAlias(c.Arg("name"), c.Arg("ip_addr"))
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) getCLINetifAliasListHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		if c.Flag("debug") == "true" {
			j.Debug = true
		}
		if c.Flag("quiet") == "true" {
			j.Quiet = true
		}

		err := j.ListNetifAliases(c.Arg("name"))
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) AddNetifCmds(c *cli.CLI) {
	netif_create := c.AddCmd("netif_create", "Create network interface", j.getCLINetifCreateHandler())
	netif_create.AddArg("name", "NAME", "", cli.TypeString|cli.Required)
	netif_create.AddArg("ip_addr_begin", "IP_ADDR_BEGIN", "", cli.TypeString|cli.Required)
	netif_create.AddArg("ip_addr_end", "IP_ADDR_END", "", cli.TypeString|cli.Required)
	netif_create.AddArg("if_name", "INTERFACE_NAME", "", cli.TypeString)
	netif_create.AddFlag("quiet", "q", "", "Do not output anything", cli.TypeBool)
	netif_create.AddFlag("debug", "d", "", "Print more information", cli.TypeBool)

	netif_destroy := c.AddCmd("netif_destroy", "Destroy network interface", j.getCLINetifDestroyHandler())
	netif_destroy.AddArg("name", "NAME", "", cli.TypeString|cli.Required)
	netif_destroy.AddFlag("quiet", "q", "", "Do not output anything", cli.TypeBool)
	netif_destroy.AddFlag("debug", "d", "", "Print more information", cli.TypeBool)

	_ = c.AddCmd("netif_list", "List network interfaces", j.getCLINetifListHandler())

	netif_alias_add := c.AddCmd("netif_alias_add", "Add alias IP address to a network interface", j.getCLINetifAliasAddHandler())
	netif_alias_add.AddArg("name", "NAME", "", cli.TypeString|cli.Required)
	netif_alias_add.AddArg("ip_addr", "IP_ADDR", "", cli.TypeString|cli.Required)
	netif_alias_add.AddFlag("quiet", "q", "", "Do not output anything", cli.TypeBool)
	netif_alias_add.AddFlag("debug", "d", "", "Print more information", cli.TypeBool)

	netif_alias_delete := c.AddCmd("netif_alias_delete", "Delete alias IP address from a network interface", j.getCLINetifAliasDeleteHandler())
	netif_alias_delete.AddArg("name", "NAME", "", cli.TypeString|cli.Required)
	netif_alias_delete.AddArg("ip_addr", "IP_ADDR", "", cli.TypeString|cli.Required)
	netif_alias_delete.AddFlag("quiet", "q", "", "Do not output anything", cli.TypeBool)
	netif_alias_delete.AddFlag("debug", "d", "", "Print more information", cli.TypeBool)

	netif_alias_list := c.AddCmd("netif_alias_list", "List network interface alias IP addresses", j.getCLINetifAliasListHandler())
	netif_alias_list.AddArg("name", "NAME", "", cli.TypeString|cli.Required)
}
