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

func (st *State) GetBase(rls string) (*Base, error) {
	if st.Bases == nil {
		return nil, nil
	}
	if st.Bases[rls] == nil {
		return nil, nil
	}
	return st.Bases[rls], nil
}

func (st *State) AddBase(rls string, bs *Base) {
	if st.Bases == nil {
		st.Bases = make(map[string]*Base)
	}
	st.Bases[rls] = bs
}

func (st *State) RemoveItem(t string, n string) error {
	if n == "" {
		return errors.New("Invalid name")
	}

	if t == "base" {
		m := make(map[string]*Base)
		for k, _ := range st.Bases {
			if k != n {
				m[k] = st.Bases[k]
			}
		}
		st.Bases = m
	} else if t == "template" {
		m := make(map[string]*Template)
		for k, _ := range st.Templates {
			if k != n {
				m[k] = st.Templates[k]
			}
		}
	} else if t == "jail" {
		m := make(map[string]*Jail)
		for k, _ := range st.Jails {
			if k != n {
				m[k] = st.Jails[k]
			}
		}
	} else if t == "network_interface" {
		m := make(map[string]*NetworkInterface)
		for k, _ := range st.NetworkInterfaces {
			if k != n {
				m[k] = st.NetworkInterfaces[k]
			}
		}
	} else if t == "pf_rule" {
		m := make(map[string]*PFRules)
		for k, _ := range st.PFRules {
			if k != n {
				m[k] = st.PFRules[k]
			}
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

func (st *State) PrintItems(f *os.File, t string) error {
	if t == "" || t == "bases" {
		for k, _ := range st.Bases {
			fmt.Fprintf(f, "base %s\n", k)
		}
	}
	if t == "" || t == "jails" {
		for k, _ := range st.Jails {
			fmt.Fprintf(f, "jail %s\n", k)
		}
	}
	if t == "" || t == "templates" {
		for k, _ := range st.Templates {
			fmt.Fprintf(f, "template %s\n", k)
		}
	}
	if t == "" || t == "networkinterfaces" {
		for k, _ := range st.NetworkInterfaces {
			fmt.Fprintf(f, "netif %s\n", k)
		}
	}
	if t == "" || t == "pfrules" {
		for k, _ := range st.PFRules {
			fmt.Fprintf(f, "pfrule %s\n", k)
		}
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
