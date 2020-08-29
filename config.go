package main

import (
  "encoding/json"
  "errors"
  "fmt"
  "io/ioutil"
  "os"
)

type Config struct {
  PathData string `json:"path_data"`
  DirBases string `json:"dir_bases"`
  DirTemplates string `json:"dir_templates"`
  DirState string `json:"dir_state"`
  DirJails string `json:"dir_jails"`
  DirConfigs string `json:"dir_configs"`
  DirTmp string `json:"dir_tmp"`
  FileState string `json:"jailguard.jailstate"`
  NetIf string `json:"1337"`
  PfAnchor string `json:"jailguard"`

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

func (c *Config) Set(k string, v string) error {
  // TODO: Validation
  if k == "path_data" {
    c.PathData = v
  }
  if k == "dir_bases" {
    c.DirBases = v
  }
  if k == "dir_templates" {
    c.DirTemplates = v
  }
  if k == "dir_state" {
    c.DirState = v
  }
  if k == "dir_jails" {
    c.DirJails = v
  }
  if k == "dir_configs" {
    c.DirConfigs = v
  }
  if k == "dir_tmp" {
    c.DirTmp = v
  }
  if k == "file_state" {
    c.FileState = v
  }
  if k == "net_if" {
    c.NetIf = v
  }
  if k == "pf_anchor" {
    c.PfAnchor = v
  }
  return nil
}

func (c *Config) Print(f *os.File, k string) {
  if k == "" || k == "path_data" {
    fmt.Fprintf(f, "path_data %s\n", c.PathData)
  }
  if k == "" || k == "dir_bases" {
    fmt.Fprintf(f, "dir_bases %s\n", c.DirBases)
  }
  if k == "" || k == "dir_templates" {
    fmt.Fprintf(f, "dir_templates %s\n", c.DirTemplates)
  }
  if k == "" || k == "dir_state" {
    fmt.Fprintf(f, "dir_state %s\n", c.DirState)
  }
  if k == "" || k == "dir_jails" {
    fmt.Fprintf(f, "dir_jails %s\n", c.DirJails)
  }
  if k == "" || k == "dir_configs" {
    fmt.Fprintf(f, "dir_configs %s\n", c.DirConfigs)
  }
  if k == "" || k == "dir_tmp" {
    fmt.Fprintf(f, "dir_tmp %s\n", c.DirTmp)
  }
  if k == "" || k == "file_state" {
    fmt.Fprintf(f, "file_state %s\n", c.FileState)
  }
  if k == "" || k == "net_if" {
    fmt.Fprintf(f, "net_if %s\n", c.NetIf)
  }
  if k == "" || k == "pf_anchor" {
    fmt.Fprintf(f, "pf_anchor %s\n", c.PfAnchor)
  }
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
