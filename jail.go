package main

import (
	"errors"
	"fmt"
	"os"
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

func (jl *Jail) existsInOS() (bool, error) {
	jl.logger(LOGDBG, fmt.Sprintf("Running 'jls' to check if jail %s is running...", jl.Name))
	out, err := CmdOut("jls", "-Nn")
	if err != nil {
		return false, errors.New("Error has occurred when checking if jail is running: " + err.Error())
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
	jl.logger(LOGDBG, fmt.Sprintf("Running 'jail -c -f %s' command to create jail...", jl.ConfigFilepath))
	err := CmdRun("jail", "-c", "-f", jl.ConfigFilepath)
	if err != nil {
		return errors.New(fmt.Sprintf("Error executing 'jail' command: %s", err.Error()))
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
		jl.logger(LOGDBG, fmt.Sprintf("Running 'jail -r %s' command to remove jail", jl.Name))
		err := CmdRun("jail", "-r", jl.Name)
		if err != nil {
			return errors.New(fmt.Sprintf("Error executing 'jail' command: %s", err.Error()))
		}
	}

	_, _, err = StatWithLog(jl.Dirpath, jl.logger)
	if err != nil {
		if os.IsNotExist(err) {
			jl.logger(LOGDBG, "Jail directory does not exist. Nothing to remove")
			return nil
		} else {
			jl.logger(LOGDBG, "Error has occurred when checking jail directory")
			return err
		}
	}

	jl.logger(LOGDBG, fmt.Sprintf("Running 'chflags -R noschg' on %s directory...", jl.Dirpath))
	err = CmdRun("chflags", "-R", "noschg", jl.Dirpath)
	if err != nil {
		return errors.New("Error running chflags -R noschg")
	}

	err = RemoveAllWithLog(jl.Dirpath, jl.logger)
	if err != nil {
		return errors.New("Error removing jail directory. Please remove the directory manually and remove the state")
	}

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
