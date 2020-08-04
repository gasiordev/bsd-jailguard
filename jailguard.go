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
	"time"
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
	j.logger.Output(2, fmt.Sprintf("%s: %s\n", n, s))

	if ((t == LOGINF || t == LOGERR) && !j.Quiet) || (t == LOGDBG && j.Debug && !j.Quiet) {
		fmt.Fprintf(j.cli.GetStdout(), fmt.Sprintf("* %s\n", s))
	}
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

func (j *Jailguard) getNewBase(rls string) *Base {
	bs := NewBase(rls, j.getBaseDirPath(rls))
	bs.SetLogger(func(t int, s string) {
		j.Log(t, s)
	})
	return bs
}

func (j *Jailguard) getNewJail(cfg *JailConf, dir *JailDir) *Jail {
	jl := NewJail(cfg, dir)
	jl.SetLogger(func(t int, s string) {
		j.Log(t, s)
	})
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

func (j *Jailguard) getJailDir(n string, d string) *JailDir {
	dir := NewJailDir(n, d)
	dir.SetLogger(func(t int, s string) {
		j.Log(t, s)
	})
	return dir
}

func (j *Jailguard) getOSRelease() (string, error) {
	out, err := CmdOut("uname", "-m", "-r")
	if err != nil {
		return "", errors.New("Error getting OS release: " + err.Error())
	}
	a := strings.Split(string(out), " ")
	var re = regexp.MustCompile(`-p[0-9]+$`)
	return strings.TrimSpace(re.ReplaceAllString(a[0], "")), nil
}

func (j *Jailguard) getJailAndCheckIfExistsInOS(n string, fn func(int, string)) (*State, *Jail, bool, error) {
	st, err := j.getState()
	if err != nil {
		return nil, nil, false, err
	}

	jl, err := st.GetJail(n)
	if err != nil {
		return st, nil, false, err
	}

	ex, err := JailExistsInOSWithLog(n, fn)
	if err != nil {
		return st, jl, false, err
	}

	if jl != nil {
		jl.SetLogger(func(t int, s string) {
			j.Log(t, s)
		})
		if jl.Config != nil {
			jl.Config.SetLogger(func(t int, s string) {
				j.Log(t, s)
			})
		}
		if jl.Dir != nil {
			jl.Dir.SetLogger(func(t int, s string) {
				j.Log(t, s)
			})
		}
	}
	return st, jl, ex, nil
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
		bs, err := st.GetBase(n)
		if err != nil {
			return err
		}
		if bs != nil {
			j.Log(LOGERR, fmt.Sprintf("Base %s already exists. Remove it first before importing a new one"))
			return errors.New("State item already exists")
		}

		bs = j.getNewBase(n)
		j.Log(LOGDBG, fmt.Sprintf("Checking if base %s can be imported into state file...", n))
		err = bs.Import()
		if err != nil {
			return err
		}
		st.AddBase(n, bs)
	} else if t == "jail" {
		jl, err := st.GetJail(n)
		if err != nil {
			return err
		}
		if jl != nil {
			j.Log(LOGERR, fmt.Sprintf("Jail %s already exists. Remove it first before import a new one"))
			return errors.New("State item already exists")
		}

		cfg, err := j.getJailConf(j.getConfigFilePath(n))
		if err != nil {
			return err
		}
		cfg.Filepath = j.getConfigFilePath(n)

		dir := j.getJailDir(n, j.getJailDirPath(n))

		jl = j.getNewJail(cfg, dir)
		j.Log(LOGDBG, fmt.Sprintf("Checking if base %s can be imported into state file...", n))
		err = jl.Import()
		if err != nil {
			return err
		}
		st.AddJail(n, jl)
	} else {
		return errors.New("Invalid item type")
	}

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

	bs, err := st.GetBase(rls)
	if err != nil {
		return err
	}
	if bs == nil {
		bs = j.getNewBase(rls)
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
			j.Log(LOGINF, fmt.Sprintf("Base %s already exists but downloading it again...", rls))
			err = bs.Download(ow)
			if err != nil {
				return err
			}
		} else {
			j.Log(LOGINF, fmt.Sprintf("Base %s already exists. Use 'overwrite' flag to download it again", rls))
			return nil
		}
	}

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

	bs, err := st.GetBase(rls)
	if err != nil {
		return err
	}
	if bs == nil {
		return nil
	}
	bs.SetLogger(func(t int, s string) {
		j.Log(t, s)
	})

	err = bs.Remove()
	if err != nil {
		return err
	}

	st.RemoveItem("base", rls)

	err = st.Save()
	if err != nil {
		return err
	}

	return nil
}

func (j *Jailguard) StopJail(n string) error {
	st, jl, ex, err := j.getJailAndCheckIfExistsInOS(n, j.Log)
	if err != nil {
		return err
	}

	if !ex {
		return nil
	}

	if jl == nil {
		j.Log(LOGDBG, fmt.Sprintf("There is a jail %s running in the system", n))
		return errors.New("Jail does not exist in state file but there is a jail with same name running in the system. Stop it manually or import into the state")
	}

	err = jl.Stop()
	if err != nil {
		return errors.New("Error stopping jail")
	}

	st.AddHistoryEntry(fmt.Sprintf("Stop jail %s", n))

	err = st.Save()
	if err != nil {
		return err
	}

	return nil

}

func (j *Jailguard) StartJail(n string) error {
	st, jl, ex, err := j.getJailAndCheckIfExistsInOS(n, j.Log)
	if err != nil {
		return err
	}

	if ex {
		if jl == nil {
			j.Log(LOGDBG, fmt.Sprintf("There is a jail %s running in the system", n))
			return errors.New("Jail does not exist in state file but there is a jail with same name running in the system")
		}
		return nil
	}
	if jl == nil {
		return errors.New("Jail does not exist in state file")
	}

	err = jl.Start()
	if err != nil {
		return errors.New("Error starting jail")
	}

	st.AddHistoryEntry(fmt.Sprintf("Start jail %s", n))

	err = st.Save()
	if err != nil {
		return err
	}

	return nil

}

func (j *Jailguard) RemoveJail(n string, stop bool) error {
	st, jl, ex, err := j.getJailAndCheckIfExistsInOS(n, j.Log)
	if err != nil {
		return err
	}

	if ex {
		if jl == nil {
			j.Log(LOGDBG, fmt.Sprintf("Jail %s is running in the system", n))
			return errors.New("Jail does not exist in state file but there is a jail with same name running in the system. Remove it manually or import the state of it")
		} else if !stop {
			return errors.New("Please stop jail first or use --stop")
		}
	}
	if jl == nil {
		return nil
	}

	if ex && stop {
		err = jl.Stop()
		if err != nil {
			return errors.New("Error stopping jail")
		}
	}

	err = jl.Remove()
	if err != nil {
		return errors.New("Error removing jail")
	}

	st.RemoveItem("jail", n)

	err = st.Save()
	if err != nil {
		return err
	}

	return nil

}

func (j *Jailguard) CreateJail(f string, rls string, start bool) error {
	cfg, err := j.getJailConf(f)
	if err != nil {
		return err
	}

	dir := j.getJailDir(cfg.Name, j.getJailDirPath(cfg.Name))

	st, jl, ex, err := j.getJailAndCheckIfExistsInOS(cfg.Name, j.Log)
	if err != nil {
		return err
	}
	if jl != nil {
		return errors.New(fmt.Sprintf("Jail %s already exists in state file", cfg.Name))
	}
	if ex {
		return errors.New(fmt.Sprintf("Jail %s already exists in the system", cfg.Name))
	}

	var errCreateDir error
	var errWriteCfg error

	if cfg.Config["path"] == "" {
		if rls == "" {
			rls, err = j.getOSRelease()
			if err != nil {
				return errors.New("Error getting OS release")
			}
		}

		bs, err := st.GetBase(rls)
		if err != nil {
			return err
		}
		if bs == nil {
			return errors.New(fmt.Sprintf("Base %s not found in state file", rls))
		}

		bs.SetLogger(func(t int, s string) {
			j.Log(t, s)
		})

		errCreateDir = dir.CreateFromTarball(bs.GetBaseTarballPath())
		cfg.Config["path"] = j.getJailDirPath(cfg.Name)
	} else {
		if rls != "" {
			j.Log(LOGINF, "'path' is provided in the file so base flag will be ignored")
		}
	}

	j.Log(LOGDBG, "Writing jail config to a file...")
	errWriteCfg = cfg.Write(j.getConfigFilePath(cfg.Name))

	jl = j.getNewJail(cfg, dir)
	if errWriteCfg != nil || errCreateDir != nil {
		jl.CleanAfterError()
	}

	if errCreateDir != nil {
		return errors.New("Error creating jail source directory")
	}
	if errWriteCfg != nil {
		return errors.New("Error creating config file")
	}

	st.AddJail(cfg.Name, jl)

	err = st.Save()
	if err != nil {
		return errors.New("Jail has been created but there was an error with writing state. Try to import the state of the jail using state_import")
	}

	if start {
		j.Log(LOGDBG, "Starting jail")
		err1 := jl.Start()
		err2 := st.Save()
		if err2 != nil {
			return errors.New("Error has occurred while saving state")
		}
		if err1 != nil {
			return errors.New("Error has occurred while starting jail")
		}
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////

func NewJailguard() *Jailguard {
	j := &Jailguard{}
	return j
}

func StatWithLog(p string, fn func(int, string)) (os.FileInfo, bool, error) {
	fn(LOGDBG, fmt.Sprintf("Getting stat for path %s...", p))
	st, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			fn(LOGDBG, fmt.Sprintf("Path %s does not exist", p))
		} else {
			fn(LOGDBG, fmt.Sprintf("Error has occurred when getting stat for path %s: %s", p, err.Error()))
		}
		return st, false, err
	} else {
		fn(LOGDBG, fmt.Sprintf("Found path %s", p))
		if st.IsDir() {
			fn(LOGDBG, fmt.Sprintf("Path %s is a directory", p))
			return st, true, nil
		}
	}
	return st, false, nil
}

func RemoveAllWithLog(p string, fn func(int, string)) error {
	fn(LOGDBG, fmt.Sprintf("Removing %s...", p))
	err := os.RemoveAll(p)
	if err != nil {
		fn(LOGDBG, fmt.Sprintf("Error has occurred when removing %s: %s", p, err.Error()))
	}
	fn(LOGDBG, fmt.Sprintf("Path %s has been removed", p))
	return err
}

func CreateDirWithLog(p string, fn func(int, string)) error {
	fn(LOGDBG, fmt.Sprintf("Creating directory %s...", p))
	err := os.MkdirAll(p, os.ModePerm)
	if err != nil {
		fn(LOGDBG, fmt.Sprintf("Error has occurred when creating %s: %s", p, err.Error()))
	}
	fn(LOGDBG, fmt.Sprintf("Directory %s has been created", p))
	return err
}

func CmdFetchWithLog(url string, o string, fn func(int, string)) error {
	fn(LOGDBG, fmt.Sprintf("Running 'fetch' to download %s to %s...", url, o))
	_, err := CmdOut("fetch", url, "-o", o)
	if err != nil {
		fn(LOGDBG, fmt.Sprintf("Error has occurred when downloading %s to %s", url, o))
		return err
	}
	fn(LOGDBG, fmt.Sprintf("File %s has been successfully saved in %s", url, o))
	return nil
}

func CmdTarExtractWithLog(f string, d string, fn func(int, string)) error {
	fn(LOGDBG, fmt.Sprintf("Running 'tar' to extract %s to %s directory...", f, d))
	_, err := CmdOut("tar", "-xvf", f, "-C", d)
	if err != nil {
		fn(LOGDBG, fmt.Sprintf("Error has occurred when extracting %s to %s", f, d))
		return err
	}
	fn(LOGDBG, fmt.Sprintf("File %s has been successfully extracted to %s", f, d))
	return nil
}

func CmdOut(c string, a ...string) ([]byte, error) {
	cmd := exec.Command(c, a...)
	cmd.Stdin = os.Stdin
	return cmd.Output()
}

func CmdRun(c string, a ...string) error {
	cmd := exec.Command(c, a...)
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func JailExistsInOSWithLog(n string, fn func(int, string)) (bool, error) {
	fn(LOGDBG, fmt.Sprintf("Running 'jls' to check if jail %s is running...", n))
	out, err := CmdOut("jls", "-Nn")
	if err != nil {
		return false, errors.New("Error running 'jls' to check if jail is running: " + err.Error())
	}

	re := regexp.MustCompile("name=" + n + " ")
	if re.Match([]byte(string(out))) {
		fn(LOGDBG, fmt.Sprintf("Jail %s is running (it was found in 'jls' output)", n))
		return true, nil
	}

	fn(LOGDBG, fmt.Sprintf("Jail %s does not seem to be running", n))
	return false, nil
}

func GetCurrentDateTime() string {
	return time.Now().String()
}
