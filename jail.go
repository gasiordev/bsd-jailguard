package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
)

type Jail struct {
	Release        string          `json:"release"`
	SourceURL      string          `json:"source_url"`
	Name           string          `json:"name"`
	Created        string          `json:"created"`
	LastUpdated    string          `json:"last_updated"`
	Config         *JailConf       `json:"config"`
	ConfigFilepath string          `json:"config_filepath"`
	Dirpath        string          `json:"dirpath"`
	Iteration      int             `json:"iteration"`
	History        []*HistoryEntry `json:"history"`
	State          string          `json:"state"`
	logger         func(int, string)
}

func (jl *Jail) SetLogger(f func(int, string)) {
	jl.logger = f
}

func (jl *Jail) AddHistoryEntry(s string) {
	he := NewHistoryEntry(GetCurrentDateTime(), s)
	if jl.History == nil {
		jl.History = []*HistoryEntry{}
	}
	jl.History = append(jl.History, he)
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

func (jl *Jail) Start() error {
	jl.logger(LOGDBG, fmt.Sprintf("Running 'jail -c -f %s' command to start jail...", jl.ConfigFilepath))
	err := CmdRun("jail", "-c", "-f", jl.ConfigFilepath)
	if err != nil {
		jl.State = "error_starting"
		return errors.New(fmt.Sprintf("Error executing 'jail' command: %s", err.Error()))
	}
	jl.State = "started"

	jl.Iteration++
	jl.AddHistoryEntry("Start")

	return nil
}

func (jl *Jail) Stop() error {
	jl.logger(LOGDBG, fmt.Sprintf("Running 'jail -r %s' command to stop jail", jl.Name))
	err := CmdRun("jail", "-r", jl.Name)
	if err != nil {
		jl.State = "error_stopping"
		return errors.New(fmt.Sprintf("Error executing 'jail' command: %s", err.Error()))
	}
	jl.State = "stopped"

	jl.Iteration++
	jl.AddHistoryEntry("Stop")

	return nil
}

func (jl *Jail) Remove() error {
	_, _, errDirPath := StatWithLog(jl.Dirpath, jl.logger)
	_, _, errCfgPath := StatWithLog(jl.ConfigFilepath, jl.logger)

	var errRmDirPath1 error
	var errRmDirPath2 error
	if errDirPath == nil {
		jl.logger(LOGDBG, fmt.Sprintf("Running 'chflags -R noschg' on %s directory...", jl.Dirpath))
		errRmDirPath1 = CmdRun("chflags", "-R", "noschg", jl.Dirpath)
		errRmDirPath2 = RemoveAllWithLog(jl.Dirpath, jl.logger)
	}

	var errRmCfgPath error
	if errCfgPath == nil {
		errRmCfgPath = RemoveAllWithLog(jl.ConfigFilepath, jl.logger)
	}

	if (errDirPath != nil && !os.IsNotExist(errDirPath)) || (errCfgPath != nil && !os.IsNotExist(errCfgPath)) || errRmDirPath1 != nil || errRmDirPath2 != nil || errRmCfgPath != nil {
		return errors.New("Error has occurred while removing jail. Please remove the directories manually and remove the state")
	}

	return nil
}

func (jl *Jail) CleanAfterError() error {
	_ = jl.Remove()
	return nil
}

func (jl *Jail) Import() error {
	_, _, err := StatWithLog(jl.Dirpath, jl.logger)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("Jail directory has not been found")
		} else {
			return errors.New("Error has occurred when importing jail")
		}
	}

	_, _, err = StatWithLog(jl.ConfigFilepath, jl.logger)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("Jail config has not been found")
		} else {
			return errors.New("Error has occurred when importing jail")
		}
	}

	ex, err := JailExistsInOSWithLog(jl.Name, jl.logger)
	if err != nil {
		return errors.New("Errors has occurred when importing jail")
	}
	if ex {
		jl.State = "started"
	}

	jl.logger(LOGDBG, "Jail source directory and config exist and jail can be imported")
	return nil
}

func NewJailFromConfig(cfg *JailConf) *Jail {
	jl := &Jail{}
	jl.Config = cfg
	jl.Name = jl.Config.Name
	jl.SetDefaultValues()
	jl.Created = GetCurrentDateTime()
	jl.State = "created"
	return jl
}
