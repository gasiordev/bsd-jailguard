package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type State struct {
	Logger func(int, string)

	version     string `json:"version"`
	software    string `json:"software"`
	created     string `json:"created"`
	lastUpdated string `json:"last_updated"`
	iteration   int    `json:"iteration"`

	bases     map[string]*Base     `json:"bases"`
	templates map[string]*Template `json:"templates"`
	jails     map[string]*Jail     `json:"jails"`

	networkInterfaces map[string]*NetworkInterface `json:"network_interfaces"`
	pfRules           map[string]*PFRules          `json:"pf_rules"`

	history []*HistoryEntry `json:"history"`
}

func (s *State) SetDefaultValues() {
	s.version = "2"
	s.software = "jailguard " + VERSION
	s.iteration = 1
}

func (s *State) PrintItems(*os.File) error {
	return nil
}

func NewState(f string) (*State, error) {
	st := &State{}
	_, err := os.Stat(f)
	// TODO: Check if it's a file
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, errors.New("Error getting state file: " + err.Error())
		}
		st.SetDefaultValues()
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
