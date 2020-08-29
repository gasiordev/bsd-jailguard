package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	PathData     string `json:"path_data"`
	DirBases     string `json:"dir_bases"`
	DirTemplates string `json:"dir_templates"`
	DirState     string `json:"dir_state"`
	DirJails     string `json:"dir_jails"`
	DirConfigs   string `json:"dir_configs"`
	DirTmp       string `json:"dir_tmp"`
	FileState    string `json:"jailguard.jailstate"`
	NetIf        string `json:"1337"`
	PfAnchor     string `json:"jailguard"`

	Filepath string `json:"filepath"`

	logger func(int, string)
}

func (c *Config) SetLogger(f func(int, string)) {
	c.logger = f
}

func (c *Config) SetDefaultValues() {
	c.PathData = "/usr/local/jailguard"
	c.DirBases = "bases"
	c.DirTemplates = "templates"
	c.DirState = "state"
	c.DirJails = "jails"
	c.DirConfigs = "configs"
	c.DirTmp = "tmp"
	c.FileState = "jailguard.jailstate"
	c.NetIf = "1337"
	c.PfAnchor = "jailguard"
}

func (c *Config) Save() error {
	c.logger(LOGDBG, "Generating config JSON...")
	o, err := json.Marshal(c)
	if err != nil {
		c.logger(LOGDBG, fmt.Sprintf("Error has occurred while generating config JSON: %s", err.Error()))
		return err
	}
	c.logger(LOGDBG, fmt.Sprintf("Writing the config to %s...", c.Filepath))
	err = ioutil.WriteFile(c.Filepath, o, 0644)
	if err != nil {
		return err
	}
	c.logger(LOGDBG, fmt.Sprintf("Config has been successfully saved to %s", c.Filepath))
	return nil
}

func NewConfig(f string) (*Config, error) {
	c := &Config{}
	c.Filepath = f
	_, err := os.Stat(f)
	// TODO: Check if it's a file
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, errors.New(fmt.Sprintf("Error has occurred while getting config file: %s", err.Error()))
		}
		c.SetDefaultValues()
		return c, nil
	}

	b, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error has occurred while getting config file: %s", err.Error()))
	}
	err = json.Unmarshal(b, c)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error has occurred while getting config file: %s", err.Error()))
	}

	return c, nil
}
