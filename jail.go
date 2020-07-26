package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
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

func (jl *Jail) cmdRun(c string, a ...string) error {
	cmd := exec.Command(c, a...)
	cmd.Stdin = os.Stdin
	return cmd.Run()
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

func NewJail(cfg *JailConf) *Jail {
	jl := &Jail{}
	jl.Config = cfg
	jl.Name = jl.Config.Name
	jl.SetDefaultValues()
	jl.Created = jl.getCurrentDateTime()
	return jl
}
