package main

import (
	"bytes"
	"fmt"
	"github.com/nicholasgasior/go-cli"
	"log"
	"os"
)

const DIRCONFIG = "/usr/local/etc/jailguard.conf.json"

const LOGINF = 1
const LOGERR = -1
const LOGDBG = 2

type Jailguard struct {
	cli    *cli.CLI
	config *Config
	logBuf bytes.Buffer
	logger *log.Logger
	Quiet  bool
	Debug  bool
}

func (j *Jailguard) GetCLI() *cli.CLI {
	return j.cli
}

func (j *Jailguard) GetConfig() *Config {
	return j.config
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
	cfg, err := NewConfig(DIRCONFIG)
	if err != nil {
		fmt.Fprintf(j.cli.GetStderr(), err.Error())
		os.Exit(1)
	}
	j.config = cfg
	// TODO: Validate config
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
