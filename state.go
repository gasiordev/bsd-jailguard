package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type State struct {
	Version     string          `json:"version"`
	Software    string          `json:"software"`
	Created     string          `json:"created"`
	LastUpdated string          `json:"last_updated"`
	Iteration   int             `json:"iteration"`
	History     []*HistoryEntry `json:"history"`

	Bases             map[string]*Base             `json:"bases"`
	Templates         map[string]*Template         `json:"templates"`
	Jails             map[string]*Jail             `json:"jails"`
	NetworkInterfaces map[string]*NetworkInterface `json:"network_interfaces"`
	PFRules           map[string]*PFRules          `json:"pf_rules"`

	Filepath string `json:"filepath"`

	logger func(int, string)
}

func (st *State) getCurrentDateTime() string {
	return time.Now().String()
}

func (st *State) SetLogger(f func(int, string)) {
	st.logger = f
}

func (st *State) SetDefaultValues() {
	st.Version = "2"
	st.Software = "jailguard " + VERSION
	st.Iteration = 1
}

func (st *State) RemoveItem(t string, n string) error {
	if n == "" {
		return errors.New("Invalid name")
	}

	if t == "base" {
		if st.Bases[n] != nil {
			st.Bases[n] = nil
		}
	} else if t == "template" {
		if st.Templates[n] != nil {
			st.Templates[n] = nil
		}
	} else if t == "jail" {
		if st.Jails[n] != nil {
			st.Jails[n] = nil
		}
	} else if t == "network_interface" {
		if st.NetworkInterfaces[n] != nil {
			st.NetworkInterfaces[n] = nil
		}
	} else if t == "pf_rule" {
		if st.PFRules[n] != nil {
			st.PFRules[n] = nil
		}
	} else {
		return errors.New("Invalid state item type")
	}

	err := st.Save()
	if err != nil {
		return errors.New("Cannot save state to a file")
	}

	return nil
}

func (st *State) PrintItems(f *os.File) error {
	for k, _ := range st.Bases {
		fmt.Fprintf(f, "base %s\n", k)
	}
	for k, _ := range st.Jails {
		fmt.Fprintf(f, "jail %s\n", k)
	}
	for k, _ := range st.Templates {
		fmt.Fprintf(f, "template %s\n", k)
	}
	for k, _ := range st.NetworkInterfaces {
		fmt.Fprintf(f, "netif %s\n", k)
	}
	for k, _ := range st.PFRules {
		fmt.Fprintf(f, "pfrule %s\n", k)
	}
	return nil
}

func (st *State) Save() error {
	st.SetDefaultValues()
	st.LastUpdated = st.getCurrentDateTime()
	if st.Created == "" {
		st.Created = st.LastUpdated
	}
	if st.Bases == nil {
		st.Bases = make(map[string]*Base)
	}
	if st.Templates == nil {
		st.Templates = make(map[string]*Template)
	}
	if st.Jails == nil {
		st.Jails = make(map[string]*Jail)
	}
	if st.NetworkInterfaces == nil {
		st.NetworkInterfaces = make(map[string]*NetworkInterface)
	}
	if st.PFRules == nil {
		st.PFRules = make(map[string]*PFRules)
	}

	st.logger(LOGDBG, "Generating state JSON")
	o, err := json.Marshal(st)
	if err != nil {
		st.logger(LOGDBG, "Error generating JSON: "+err.Error())
		return err
	}

	d := filepath.Dir(st.Filepath)
	st.logger(LOGDBG, "Checking if "+d+" exists and is dir")
	stat, err := os.Stat(d)
	if err != nil {
		if os.IsNotExist(err) {
			st.logger(LOGDBG, d+" does not exist, trying to create it")
			err2 := os.MkdirAll(d, os.ModePerm)
			if err2 != nil {
				st.logger(LOGDBG, "Error with creating "+d+" dir")
				return err2
			}
		} else {
			st.logger(LOGDBG, "Error with checking dir "+d+" existance: "+err.Error())
			return err
		}
	} else if !stat.IsDir() {
		st.logger(LOGDBG, "State dir "+d+" exists but it is not a dir")
		return errors.New("Path for state dir is not a dir")
	}

	st.logger(LOGDBG, "Writing state to "+st.Filepath)
	err = ioutil.WriteFile(st.Filepath, o, 0644)
	if err != nil {
		return err
	}

	return nil
}

func NewState(f string) (*State, error) {
	st := &State{}
	st.Filepath = f
	_, err := os.Stat(f)
	// TODO: Check if it's a file
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, errors.New("Error getting state file: " + err.Error())
		}
		st.SetDefaultValues()
		st.Created = st.getCurrentDateTime()
		return st, nil
	}

	b, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, errors.New("Error reading state file: " + err.Error())
	}
	err = json.Unmarshal(b, st)
	if err != nil {
		return nil, errors.New("Error parsing state file: " + err.Error())
	}

	return st, nil
}
