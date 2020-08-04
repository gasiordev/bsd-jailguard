package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
)

type State struct {
	Version     string          `json:"version"`
	Software    string          `json:"software"`
	Created     string          `json:"created"`
	LastUpdated string          `json:"last_updated"`
	Iteration   int             `json:"iteration"`
	History     []*HistoryEntry `json:"history"`

	Bases             map[string]*Base             `json:"bases"`
	Jails             map[string]*Jail             `json:"jails"`
	NetworkInterfaces map[string]*NetworkInterface `json:"network_interfaces"`
	PFRules           map[string]*PFRules          `json:"pf_rules"`

	Filepath string `json:"filepath"`

	logger func(int, string)
}

func (st *State) SetLogger(f func(int, string)) {
	st.logger = f
}

func (st *State) SetDefaultValues() {
	st.Version = "2"
	st.Software = "jailguard " + VERSION
	st.Iteration++
}

func (st *State) AddHistoryEntry(s string) {
	he := NewHistoryEntry(GetCurrentDateTime(), s)
	if st.History == nil {
		st.History = []*HistoryEntry{}
	}
	st.History = append(st.History, he)
}

func (st *State) GetBase(rls string) (*Base, error) {
	st.logger(LOGDBG, fmt.Sprintf("Getting base %s from the state...", rls))
	if st.Bases == nil {
		st.logger(LOGDBG, "There are not bases in the state")
		return nil, nil
	}
	if st.Bases[rls] == nil {
		// If there's a "null" in the state file (so just empty base name key), we remove it
		for _, k := range reflect.ValueOf(st.Bases).MapKeys() {
			if k.String() == rls {
				st.RemoveItem("base", rls)
				st.logger(LOGDBG, fmt.Sprintf("Fixing base %s being null in the state", rls))
				_ = st.Save()
			}
		}

		st.logger(LOGDBG, fmt.Sprintf("Base %s has not been found in the state", rls))
		return nil, nil
	}
	st.logger(LOGDBG, fmt.Sprintf("Base %s has been found in the state", rls))
	return st.Bases[rls], nil
}

func (st *State) GetJail(jl string) (*Jail, error) {
	st.logger(LOGDBG, fmt.Sprintf("Getting jail %s from the state...", jl))
	if st.Jails == nil {
		st.logger(LOGDBG, "There no jails in the state")
		return nil, nil
	}
	if st.Jails[jl] == nil {
		// If there's a "null" in the state file (so just empty jail name key), we remove it
		for _, k := range reflect.ValueOf(st.Jails).MapKeys() {
			if k.String() == jl {
				st.RemoveItem("jail", jl)
				st.logger(LOGDBG, fmt.Sprintf("Fixing jail %s being null in the state", jl))
				_ = st.Save()
			}
		}

		st.logger(LOGDBG, fmt.Sprintf("Jail %s has not been found in the state", jl))
		return nil, nil
	}
	st.logger(LOGDBG, fmt.Sprintf("Jail %s has been found in the state", jl))
	return st.Jails[jl], nil
}

func (st *State) AddBase(rls string, bs *Base) {
	st.logger(LOGDBG, fmt.Sprintf("Adding base %s to the state...", rls))
	if st.Bases == nil {
		st.Bases = make(map[string]*Base)
	}
	st.Bases[rls] = bs

	st.AddHistoryEntry(fmt.Sprintf("Add base %s", rls))
}

func (st *State) AddJail(n string, jl *Jail) {
	st.logger(LOGDBG, fmt.Sprintf("Adding jail %s to the state...", n))
	if st.Jails == nil {
		st.Jails = make(map[string]*Jail)
	}
	st.Jails[n] = jl

	st.AddHistoryEntry(fmt.Sprintf("Add jail %s", n))
}

func (st *State) RemoveItem(t string, n string) error {
	if n == "" {
		return errors.New("Invalid name")
	}

	st.logger(LOGDBG, fmt.Sprintf("Removing item %s %s from the state...", t, n))
	if t == "base" {
		m := make(map[string]*Base)
		for k, _ := range st.Bases {
			if k != n {
				m[k] = st.Bases[k]
			}
		}
		st.Bases = m
	} else if t == "jail" {
		m := make(map[string]*Jail)
		for k, _ := range st.Jails {
			if k != n {
				m[k] = st.Jails[k]
			}
		}
		st.Jails = m
	} else if t == "network_interface" {
		m := make(map[string]*NetworkInterface)
		for k, _ := range st.NetworkInterfaces {
			if k != n {
				m[k] = st.NetworkInterfaces[k]
			}
		}
		st.NetworkInterfaces = m
	} else if t == "pf_rule" {
		m := make(map[string]*PFRules)
		for k, _ := range st.PFRules {
			if k != n {
				m[k] = st.PFRules[k]
			}
		}
		st.PFRules = m
	} else {
		return errors.New("Invalid state item type")
	}

	err := st.Save()
	if err != nil {
		return errors.New("Cannot save state to a file")
	}
	st.logger(LOGDBG, fmt.Sprintf("Item %s %s has been removed from the state", t, n))

	st.AddHistoryEntry(fmt.Sprintf("Remove item %s %s", t, n))

	return nil
}

func (st *State) PrintItems(f *os.File, t string) error {
	st.logger(LOGDBG, "Printing out state items...")
	if t == "" || t == "bases" {
		for k, _ := range st.Bases {
			fmt.Fprintf(f, "base %s\n", k)
		}
	}
	if t == "" || t == "jails" {
		for k, jl := range st.Jails {
			fmt.Fprintf(f, "jail %s %s\n", k, jl.State)
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
	st.logger(LOGDBG, "Preparing the state to be saved into the file...")
	st.SetDefaultValues()
	st.LastUpdated = GetCurrentDateTime()
	if st.Created == "" {
		st.Created = st.LastUpdated
	}

	if st.Created == "" {
		st.Created = st.LastUpdated
	}
	if st.Bases == nil {
		st.Bases = make(map[string]*Base)
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

	st.logger(LOGDBG, "Generating state JSON...")
	o, err := json.Marshal(st)
	if err != nil {
		st.logger(LOGDBG, "Error has occurred while generating state JSON: "+err.Error())
		return err
	}

	d := filepath.Dir(st.Filepath)
	stat, _, err := StatWithLog(d, st.logger)
	if err != nil {
		if os.IsNotExist(err) {
			st.logger(LOGDBG, fmt.Sprintf("State directory %s does not exist, trying to create it...", d))
			err2 := CreateDirWithLog(d, st.logger)
			if err2 != nil {
				return errors.New("Error has occurred while saving the state")
			}
		} else {
			return errors.New("Error has occurred while saving the state")
		}
	} else if !stat.IsDir() {
		return errors.New("Path for state directory is not a directory")
	}

	st.Iteration++

	st.logger(LOGDBG, fmt.Sprintf("Writing the state to %s...", st.Filepath))
	err = ioutil.WriteFile(st.Filepath, o, 0644)
	if err != nil {
		return err
	}

	st.logger(LOGDBG, fmt.Sprintf("State has been successfully saved to %s", st.Filepath))
	return nil
}

func NewState(f string) (*State, error) {
	st := &State{}
	st.Filepath = f
	_, err := os.Stat(f)
	// TODO: Check if it's a file
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, errors.New("Error has occurred while getting state file: " + err.Error())
		}
		st.SetDefaultValues()
		st.Created = GetCurrentDateTime()
		return st, nil
	}

	b, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, errors.New("Error has occurred while reading state file: " + err.Error())
	}
	err = json.Unmarshal(b, st)
	if err != nil {
		return nil, errors.New("Error has occurred while parsing state file: " + err.Error())
	}

	return st, nil
}
