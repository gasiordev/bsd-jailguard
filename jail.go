package main

import (
	"errors"
	"fmt"
	"os"
	// "path/filepath"
	"regexp"
)

type Jail struct {
	Release     string            `json:"release"`
	SourceURL   string            `json:"source_url"`
	Name        string            `json:"name"`
	Created     string            `json:"created"`
	LastUpdated string            `json:"last_updated"`
	Config      *JailConf         `json:"config"`
	Dir         *JailDir          `json:"dir"`
	Iteration   int               `json:"iteration"`
	History     []*HistoryEntry   `json:"history"`
	State       string            `json:"state"`
	PortFwds    map[string]string `json:"port_fwds"`
	NATPass     string            `json:"nat_pass"`
	logger      func(int, string)
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
	out, err := CmdOut(jl.logger, "jls", "-Nn")
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
	jl.logger(LOGDBG, fmt.Sprintf("Running 'jail -c -f %s' command to start jail...", jl.Config.Filepath))
	err := CmdRun(jl.logger, "jail", "-c", "-f", jl.Config.Filepath)
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
	err := CmdRun(jl.logger, "jail", "-r", jl.Name)
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
	err1 := jl.Dir.Remove()
	err2 := jl.Config.Remove()

	if err1 != nil || err2 != nil {
		return errors.New("Error has occurred while removing jail. Please remove the directories manually and remove the state")
	}

	return nil
}

func (jl *Jail) CleanAfterError() error {
	_ = jl.Remove()
	return nil
}

func (jl *Jail) Import() error {
	_, _, err := StatWithLog(jl.Dir.Dirpath, jl.logger)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("Jail directory has not been found")
		} else {
			return errors.New("Error has occurred when importing jail")
		}
	}

	_, _, err = StatWithLog(jl.Config.Filepath, jl.logger)
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

func (jl *Jail) AddPortFwd(host_port string, jail_port string) error {
	if jl.PortFwds == nil {
		jl.PortFwds = make(map[string]string)
	}
	for k, v := range jl.PortFwds {
		if k == host_port {
			if v == jail_port {
				return nil
			} else {
				return errors.New(fmt.Sprintf("Host port %s is already mapped", host_port))
			}
		}
	}
	jl.PortFwds[host_port] = jail_port
	return nil
}

func (jl *Jail) DeletePortFwd(host_port string, jail_port string) error {
	if jl.PortFwds == nil {
		jl.PortFwds = make(map[string]string)
	}
	m := make(map[string]string)
	for k, v := range jl.PortFwds {
		if k != host_port || v != jail_port {
			m[k] = v
		}
	}
	jl.PortFwds = m
	return nil
}

func (jl *Jail) CreateNATPass(if_gw string) error {
	jl.NATPass = if_gw
	return nil
}

func (jl *Jail) RemoveNATPass() error {
	jl.NATPass = ""
	return nil
}

func NewJail(cfg *JailConf, d *JailDir) *Jail {
	jl := &Jail{}
	jl.Config = cfg
	jl.Dir = d
	jl.Name = jl.Config.Name
	jl.SetDefaultValues()
	jl.Created = GetCurrentDateTime()
	jl.State = "created"
	return jl
}
