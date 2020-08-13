package main

import (
	"errors"
	"github.com/nicholasgasior/go-cli"
	"regexp"
	"strconv"
	"strings"
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
	create := c.AddCmd("netif_create", "Create network interface", j.getCLINetifCreateHandler())
	create.AddArg("name", "NAME", "", cli.TypeAlphanumeric|cli.AllowUnderscore|cli.AllowHyphen|cli.Required)
	create.AddArg("ip_addr_begin", "FIRST_IPv4_ADDRESS", "", cli.TypeString|cli.Required)
	create.AddArg("ip_addr_end", "LAST_IPv4_ADDRESS", "", cli.TypeString|cli.Required)
	create.AddArg("if_name", "INTERFACE_NAME", "", cli.TypeAlphanumeric|cli.AllowUnderscore|cli.AllowHyphen)

	fn := func(c *cli.CLI) error {
		if !IsValidIPAddress(c.Arg("ip_addr_begin")) {
			return errors.New("Argument FIRST_IPv4_ADDRESS has invalid value")
		}
		if !IsValidIPAddress(c.Arg("ip_addr_end")) {
			return errors.New("Argument LAST_IPv4_ADDRESS has invalid value")
		}

		ip_b := strings.Split(c.Arg("ip_addr_begin"), ".")
		last, _ := strconv.Atoi(ip_b[3])

		ip_e := strings.Split(c.Arg("ip_addr_end"), ".")
		for i := 0; i < 4; i++ {
			if i < 3 && ip_e[i] != ip_b[i] {
				return errors.New("Argument LAST_IPv4_ADDRESS has invalid value. It should be the same IP range with last digit higher")
			}
			v, _ := strconv.Atoi(ip_e[i])
			if i == 3 && v <= last {
				return errors.New("Argument LAST_IPv4_ADDRESS has invalid value. It should be the same IP range with last digit higher")
			}
		}

		if c.Arg("if_name") != "" {
			re := regexp.MustCompile(`^[a-z][a-z0-9]{1,20}$`)
			m := re.Match([]byte(c.Arg("if_name")))
			if !m {
				return errors.New("Argument INTERFACE_NAME has invalid value")
			}
		}
		return nil
	}
	create.AddPostValidation(fn)

	destroy := c.AddCmd("netif_destroy", "Destroy network interface", j.getCLINetifDestroyHandler())
	destroy.AddArg("name", "NAME", "", cli.TypeAlphanumeric|cli.AllowUnderscore|cli.AllowHyphen|cli.Required)

	_ = c.AddCmd("netif_list", "List network interfaces", j.getCLINetifListHandler())

	alias_add := c.AddCmd("netif_alias_add", "Add alias IP address to a network interface", j.getCLINetifAliasAddHandler())
	alias_add.AddArg("name", "NAME", "", cli.TypeAlphanumeric|cli.AllowUnderscore|cli.AllowHyphen|cli.Required)
	alias_add.AddArg("ip_addr", "IPv4_ADDR", "", cli.TypeString|cli.Required)

	alias_delete := c.AddCmd("netif_alias_delete", "Delete alias IP address from a network interface", j.getCLINetifAliasDeleteHandler())
	alias_delete.AddArg("name", "NAME", "", cli.TypeAlphanumeric|cli.AllowUnderscore|cli.AllowHyphen|cli.Required)
	alias_delete.AddArg("ip_addr", "IPv4_ADDR", "", cli.TypeString|cli.Required)

	alias_list := c.AddCmd("netif_alias_list", "List network interface alias IP addresses", j.getCLINetifAliasListHandler())
	alias_list.AddArg("name", "NAME", "", cli.TypeAlphanumeric|cli.AllowUnderscore|cli.AllowHyphen|cli.Required)

	fn2 := func(c *cli.CLI) error {
		if !IsValidIPAddress(c.Arg("ip_addr")) {
			return errors.New("Argument IPv4_ADDR has invalid value")
		}
		return nil
	}
	alias_add.AddPostValidation(fn2)
	alias_delete.AddPostValidation(fn2)
}
