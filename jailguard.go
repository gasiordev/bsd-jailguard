package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gasiordev/go-cli"
	//"io/ioutil"
	"log"
	"os"
	"os/exec"
	//"path/filepath"
	"regexp"
	"strings"
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

func (j *Jailguard) cmdOut(c string, a ...string) ([]byte, error) {
	cmd := exec.Command(c, a...)
	cmd.Stdin = os.Stdin
	return cmd.Output()
}

func (j *Jailguard) getStateFilePath() string {
	return PATHDATA + "/" + DIRSTATE + "/" + FILESTATE
}

func (j *Jailguard) getBaseDirPath(rls string) string {
	return PATHDATA + "/" + DIRBASES + "/" + rls
}

func (j *Jailguard) getJailDirPath(jl string) string {
	return PATHDATA + "/" + DIRJAILS + "/" + jl
}

func (j *Jailguard) getConfigFilePath(jl string) string {
	return PATHDATA + "/" + DIRCONFIGS + "/" + jl + ".jail"
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

func (j *Jailguard) getJail(cfg *JailConf) *Jail {
	jl := NewJail(cfg)
	jl.SetLogger(func(t int, s string) {
		j.Log(t, s)
	})
	jl.ConfigFilepath = cfg.Filepath
	jl.Dirpath = cfg.Config["path"]
	return jl
}

func (j *Jailguard) getJailConf(f string) (*JailConf, error) {
	cfg := NewJailConf()
	cfg.SetLogger(func(t int, s string) {
		j.Log(t, s)
	})
	err := cfg.ParseFile(f)
	if err != nil {
		return nil, err
	}
	err = cfg.Validate()
	if err != nil {
		return nil, err
	}

	if cfg.Config["host.hostname"] == "" {
		cfg.Config["host.hostname"] = cfg.Name
	}

	if cfg.Config["exec.start"] == "" {
		cfg.Config["exec.start"] = "/bin/sh /etc/rc"
	}
	if cfg.Config["exec.stop"] == "" {
		cfg.Config["exec.stop"] = "/bin/sh /etc/rc.shutdown"
	}

	return cfg, nil
}

func (j *Jailguard) getOSRelease() (string, error) {
	out, err := j.cmdOut("uname", "-m", "-r")
	if err != nil {
		return "", errors.New("Error getting OS release: " + err.Error())
	}
	a := strings.Split(string(out), " ")
	var re = regexp.MustCompile(`-p[0-9]+$`)
	return strings.TrimSpace(re.ReplaceAllString(a[0], "")), nil
}

func (j *Jailguard) jailExistsInOS(n string) (bool, error) {
	j.Log(LOGDBG, "Checking if jail "+n+" is running with jls")
	out, err := j.cmdOut("jls", "-Nn")
	if err != nil {
		return false, errors.New("Error running jls to check if jail is running: " + err.Error())
	}

	re := regexp.MustCompile("name=" + n + " ")
	if re.Match([]byte(string(out))) {
		return true, nil
	}

	return false, nil
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

func (j *Jailguard) ImportStateItem(t string, n string) error {
	st, err := j.getState()
	if err != nil {
		return err
	}

	if t == "base" {
		j.Log(LOGDBG, "Checking if base "+n+" already exists in the state file")
		bs, err := st.GetBase(n)
		if err != nil {
			return err
		}
		if bs != nil {
			j.Log(LOGERR, "Base "+n+" already exists. Remove it first before importing a new one")
			return errors.New("State item already exists")
		}

		bs = j.getBase(n)
		j.Log(LOGDBG, "Checking if base "+n+" can be imported into state file")
		err = bs.Import()
		if err != nil {
			return err
		}
		st.AddBase(n, bs)
	} else {
		return errors.New("Invalid item type")
	}

	j.Log(LOGINF, "Saving state file")
	err = st.Save()
	if err != nil {
		return err
	}
	return nil
}

func (j *Jailguard) ListStateItems(t string) error {
	st, err := j.getState()
	if err != nil {
		return err
	}

	err = st.PrintItems(j.cli.GetStdout(), t)
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

func (j *Jailguard) RemoveBase(rls string) error {
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
		return nil
	}
	bs.SetLogger(func(t int, s string) {
		j.Log(t, s)
	})

	j.Log(LOGINF, "Removing base "+rls)
	err = bs.Remove()
	if err != nil {
		return err
	}

	st.RemoveItem("base", rls)

	j.Log(LOGINF, "Saving state file")
	err = st.Save()
	if err != nil {
		return err
	}

	return nil
}

func (j *Jailguard) CreateJail(f string, rls string) error {
	st, err := j.getState()
	if err != nil {
		return err
	}

	cfg, err := j.getJailConf(f)
	if err != nil {
		return err
	}

	jl, err := st.GetJail(cfg.Name)
	if err != nil {
		return err
	}
	if jl != nil {
		return errors.New("Jail " + cfg.Name + " already exists in state file")
	}

	ex, err := j.jailExistsInOS(cfg.Name)
	if err != nil {
		return errors.New("Error checking if jail already exists in the system")
	}
	if ex {
		return errors.New("Jail " + cfg.Name + " already exists in the system")
	}

	if cfg.Config["path"] == "" {
		if rls == "" {
			j.Log(LOGDBG, "Getting OS release as base")
			rls, err = j.getOSRelease()
			if err != nil {
				return errors.New("Error getting OS release")
			}
		}

		j.Log(LOGDBG, "Checking if base "+rls+" exists")
		bs, err := st.GetBase(rls)
		if err != nil {
			return err
		}
		if bs == nil {
			return errors.New("Base " + rls + " not found in state file")
		}

		bs.SetLogger(func(t int, s string) {
			j.Log(t, s)
		})

		err = bs.CreateJailSource(j.getJailDirPath(cfg.Name))
		if err != nil {
			return errors.New("Error creating jail directory")
		}
		cfg.Config["path"] = j.getJailDirPath(cfg.Name)
	} else {
		if rls != "" {
			j.Log(LOGINF, "path is provided in the file so base flag will be ignored")
		}
	}

	j.Log(LOGDBG, "Writing jail config to a file")

	err = cfg.WriteToFile(j.getConfigFilePath(cfg.Name))
	if err != nil {
		return errors.New("Error creating config file")
	}

	jl = j.getJail(cfg)
	j.Log(LOGDBG, "Creating jail")
	err = jl.Create()
	if err != nil {
		return errors.New("Error creating jail")
	}

	st.AddJail(cfg.Name, jl)

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
