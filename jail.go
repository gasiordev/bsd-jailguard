package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"time"
)

type Jail struct {
	Release        string    `json:"release"`
	SourceURL      string    `json:"source_url"`
	SourceTemplate string    `json:"source_template"`
	Name           string    `json:"name"`
	Created        string    `json:"created"`
	LastUpdated    string    `json:"last_updated"`
	Config         *JailConf `json:"config"`
	ConfigFilepath string    `json:"config_filepath"`
	Dirpath        string    `json:"dirpath"`
	Iteration      int
	logger         func(int, string)
}

func (jl *Jail) SetLogger(f func(int, string)) {
	jl.logger = f
}

func (jl *Jail) getCurrentDateTime() string {
	return time.Now().String()
}

func (jl *Jail) cmdOut(c string, a ...string) ([]byte, error) {
	cmd := exec.Command(c, a...)
	cmd.Stdin = os.Stdin
	return cmd.Output()
}

func (jl *Jail) cmdRun(c string, a ...string) error {
	cmd := exec.Command(c, a...)
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func (jl *Jail) existsInOS() (bool, error) {
	jl.logger(LOGDBG, "Checking if jail "+jl.Name+" is running with jls")
	out, err := jl.cmdOut("jls", "-Nn")
	if err != nil {
		return false, errors.New("Error running jls to check if jail is running: " + err.Error())
	}

	re := regexp.MustCompile("name=" + jl.Name + " ")
	if re.Match([]byte(string(out))) {
		return true, nil
	}

	return false, nil
}

func (jl *Jail) SetDefaultValues() {
	jl.Iteration = 1
}

func (jl *Jail) Create() error {
	jl.logger(LOGDBG, fmt.Sprintf("Running jail -c -f %s", jl.ConfigFilepath))

	err := jl.cmdRun("jail", "-c", "-f", jl.ConfigFilepath)
	if err != nil {
		return errors.New(fmt.Sprintf("Error executing jail command: %s", err.Error()))
	}
	return nil
}

func (jl *Jail) Destroy() error {
	ex, err := jl.existsInOS()
	if err != nil {
		return err
	}
	if ex {
		jl.logger(LOGDBG, fmt.Sprintf("Jail %s exists in the system", jl.Name))
		jl.logger(LOGDBG, fmt.Sprintf("Running jail -r %s", jl.Name))
		err := jl.cmdRun("jail", "-r", jl.Name)
		if err != nil {
			return errors.New(fmt.Sprintf("Error executing jail command to remove: %s", err.Error()))
		}
	}

	jl.logger(LOGDBG, "Checking if "+jl.Dirpath+" exists")
	_, err = os.Stat(jl.Dirpath)
	if err != nil {
		if os.IsNotExist(err) {
			jl.logger(LOGDBG, jl.Dirpath+"does not exist. Nothing to remove")
			return nil
		} else {
			jl.logger(LOGDBG, "Error with checking dir "+jl.Dirpath+" existance: "+err.Error())
			return err
		}
	}

	jl.logger(LOGDBG, fmt.Sprintf("Running chflags -R noschg on %s directory", jl.Dirpath))
	err = jl.cmdRun("chflags", "-R", "noschg", jl.Dirpath)
	if err != nil {
		return errors.New("Error running chflags -R noschg")
	}

	err = os.RemoveAll(jl.Dirpath)
	if err != nil {
		jl.logger(LOGERR, "Error removing dir "+jl.Dirpath+". Please remove the directory manually and remove the state")
		return errors.New("Error removing jail dir")
	}
	jl.logger(LOGDBG, jl.Dirpath+"has been removed")

	return nil
}

func NewJail(cfg *JailConf) *Jail {
	jl := &Jail{}
	jl.Config = cfg
	jl.Name = jl.Config.Name
	jl.SetDefaultValues()
	jl.Created = jl.getCurrentDateTime()
	return jl
}
