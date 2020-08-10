package main

import (
	"github.com/gasiordev/go-cli"
)

func (j *Jailguard) getCLIBaseDownloadHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		if c.Flag("debug") == "true" {
			j.Debug = true
		}
		if c.Flag("quiet") == "true" {
			j.Quiet = true
		}

		ow := false
		if c.Flag("overwrite") == "true" {
			ow = true
		}

		err := j.DownloadBase(c.Arg("release"), ow)
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) getCLIBaseRemoveHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		if c.Flag("debug") == "true" {
			j.Debug = true
		}
		if c.Flag("quiet") == "true" {
			j.Quiet = true
		}

		err := j.RemoveBase(c.Arg("release"))
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) getCLIBaseListHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		err := j.ListStateItems("bases")
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) AddBaseCmds(c *cli.CLI) {
	base_download := c.AddCmd("base_download", "Downloads FreeBSD base", j.getCLIBaseDownloadHandler())
	// TODO: Change 'release' flag to TypeAlphanumeric once AllowHyphen gets implemented in go-cli
	base_download.AddArg("release", "RELEASE", "", cli.TypeString|cli.Required)
	base_download.AddFlag("overwrite", "w", "", "Overwrite if exists", cli.TypeBool)
	base_download.AddFlag("quiet", "q", "", "Do not output anything", cli.TypeBool)
	base_download.AddFlag("debug", "d", "", "Print more information", cli.TypeBool)

	_ = c.AddCmd("base_list", "List FreeBSD bases", j.getCLIBaseListHandler())

	base_remove := c.AddCmd("base_remove", "Removes FreeBSD base", j.getCLIBaseRemoveHandler())
	base_remove.AddArg("release", "RELEASE", "", cli.TypeString|cli.Required)
	base_remove.AddFlag("quiet", "q", "", "Do not output anything", cli.TypeBool)
	base_remove.AddFlag("debug", "d", "", "Print more information", cli.TypeBool)
}
