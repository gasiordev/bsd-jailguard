package main

import (
	"bytes"
	"fmt"
	"github.com/gasiordev/go-cli"
	"log"
	"os"
)

const PATHDATA = "/usr/local/jailguard"
const DIRBASES = "bases"
const DIRTEMPLATES = "templates"
const DIRSTATE = "state"
const DIRJAILS = "jails"
const DIRCONFIGS = "configs"
const DIRTMP = "tmp"
const FILESTATE = "jailguard.jailstate"
const NETIF = "1337"
const PFANCHOR = "jailguard"

const LOGINF = 1
const LOGERR = -1
const LOGDBG = 2

type Jailguard struct {
	cli    *cli.CLI
	logBuf bytes.Buffer
	logger *log.Logger
	Quiet  bool
	Debug  bool
}

func (j *Jailguard) GetCLI() *cli.CLI {
	return j.cli
}

func (j *Jailguard) GetLogBuf() *bytes.Buffer {
	return &(j.logBuf)
}

func (j *Jailguard) initLogger() {
	j.logger = log.New(&(j.logBuf), "", 0)
}

func (j *Jailguard) Run() {
	j.initLogger()
	c := NewJailguardCLI(j)
	j.cli = c
	os.Exit(c.Run(os.Stdout, os.Stderr))
}

func (j *Jailguard) Log(t int, s string) {
	var n string
	if t == LOGINF {
		n = "INFO"
	} else if t == LOGDBG {
		n = "VERBOSE"
	} else if t == LOGERR {
		n = "ERROR"
	}
	j.logger.Output(2, fmt.Sprintf("%s: %s\n", n, s))

	if ((t == LOGINF || t == LOGERR) && !j.Quiet) || (t == LOGDBG && j.Debug && !j.Quiet) {
		fmt.Fprintf(j.cli.GetStdout(), fmt.Sprintf("* %s\n", s))
	}
}

func NewJailguard() *Jailguard {
	j := &Jailguard{}
	return j
}
