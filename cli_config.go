package main

import (
	"github.com/nicholasgasior/go-cli"
)

func (j *Jailguard) getCLIConfigGetHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		if c.Flag("debug") == "true" {
			j.Debug = true
		}
		if c.Flag("quiet") == "true" {
			j.Quiet = true
		}

		err := j.ShowConfigValue(c.Arg("key"))
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) getCLIConfigSetHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		if c.Flag("debug") == "true" {
			j.Debug = true
		}
		if c.Flag("quiet") == "true" {
			j.Quiet = true
		}

		err := j.SetConfigValue(c.Arg("key"), c.Arg("value"))
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) getCLIConfigListHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		err := j.ListConfig()
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) AddConfigCmds(c *cli.CLI) {
	_ = c.AddCmd("config_list", "List configuration", j.getCLIConfigListHandler())

	cfg_get := c.AddCmd("config_get", "Get specific configuration value", j.getCLIConfigGetHandler())
	cfg_get.AddArg("key", "KEY", "", cli.TypeAlphanumeric|cli.AllowUnderscore|cli.Required)

	cfg_set := c.AddCmd("config_set", "Set specific configuration value", j.getCLIConfigSetHandler())
	cfg_set.AddArg("key", "KEY", "", cli.TypeAlphanumeric|cli.AllowUnderscore|cli.Required)
	cfg_set.AddArg("value", "VALUE", "", cli.TypeString|cli.Required)
}
