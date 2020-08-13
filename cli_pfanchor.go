package main

import (
	"github.com/nicholasgasior/go-cli"
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
	_ = c.AddCmd("pf_anchor_check", "Check if required PF anchors exist", j.getCLIPFAnchorCheckHandler())
}
