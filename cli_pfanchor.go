package main

import (
	"github.com/gasiordev/go-cli"
)

func (j *Jailguard) getCLIPFAnchorCheckHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		if c.Flag("debug") == "true" {
			j.Debug = true
		}
		if c.Flag("quiet") == "true" {
			j.Quiet = true
		}

		err := j.CheckPFAnchor(false)
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) AddPFAnchorCmds(c *cli.CLI) {
	pf_anchor_check := c.AddCmd("pf_anchor_check", "Check if required PF anchors exist", j.getCLIPFAnchorCheckHandler())
	pf_anchor_check.AddFlag("quiet", "q", "", "Do not output anything", cli.TypeBool)
	pf_anchor_check.AddFlag("debug", "d", "", "Print more information", cli.TypeBool)
}
