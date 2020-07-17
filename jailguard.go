package main

import (
	"bytes"
	//"errors"
	"fmt"
	"github.com/gasiordev/go-cli"
	//"io/ioutil"
	"log"
	"os"
	//"os/exec"
	//"path/filepath"
	//"regexp"
	//"strings"
)

const PATHDATA = "/usr/local/jailguard"
const DIRBASES = "bases"
const DIRTEMPLATES = "templates"
const DIRSTATE = "state"
const DIRJAILS = "jails"
const DIRTMP = "tmp"
const FILESTATE = "jailguard.jailstate"
const NETIF = "1337"

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
	j.logger.Output(2, n+": "+s+"\n")

	if ((t == LOGINF || t == LOGERR) && !j.Quiet) || (t == LOGDBG && j.Debug && !j.Quiet) {
		fmt.Fprintf(j.cli.GetStdout(), "* "+s+"\n")
	}
}

func (j *Jailguard) getStateFilePath() string {
	return PATHDATA + "/" + DIRSTATE + "/" + FILESTATE
}

func (j *Jailguard) getBaseDirPath(rls string) string {
	return PATHDATA + "/" + DIRBASES + "/" + rls
}

func (j *Jailguard) getState() (*State, error) {
	st, err := NewState(j.getStateFilePath())
	if err != nil {
		return nil, err
	}
	st.SetLogger(func(t int, s string) {
		j.Log(t, s)
	})
	return st, nil
}

func (j *Jailguard) getBase(rls string) *Base {
	bs := NewBase(rls, j.getBaseDirPath(rls))
	bs.SetLogger(func(t int, s string) {
		j.Log(t, s)
	})
	return bs
}

///////////////////////////////////////////////////////////////////////////////
// Commands from CLI

func (j *Jailguard) RemoveStateItem(t string, n string) error {
	st, err := j.getState()
	if err != nil {
		return err
	}

	err = st.RemoveItem(t, n)
	if err != nil {
		return err
	}
	return nil
}

func (j *Jailguard) ListStateItems() error {
	st, err := j.getState()
	if err != nil {
		return err
	}

	err = st.PrintItems(j.cli.GetStdout())
	if err != nil {
		return err
	}
	return nil
}

func (j *Jailguard) DownloadBase(rls string, ow bool) error {
	st, err := j.getState()
	if err != nil {
		return err
	}

	j.Log(LOGDBG, "Checking for base "+rls+" in state file")
	bs, err := st.GetBase(rls)
	if err != nil {
		return err
	}
	if bs == nil {
		j.Log(LOGDBG, "Base "+rls+" not found in state file")
		bs = j.getBase(rls)
		err = bs.Download(ow)
		if err != nil {
			return err
		}
		st.AddBase(rls, bs)
	} else {
		bs.SetLogger(func(t int, s string) {
			j.Log(t, s)
		})

		if ow {
			j.Log(LOGINF, "Base "+rls+" already exists but downloading it again")
			err = bs.Download(ow)
			if err != nil {
				return err
			}
		} else {
			j.Log(LOGINF, "Base "+rls+" already exists. Use 'overwrite' flag to download it again")
			return nil
		}
	}

	j.Log(LOGINF, "Saving state file")
	err = st.Save()
	if err != nil {
		return err
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////

func NewJailguard() *Jailguard {
	j := &Jailguard{}
	return j
}
