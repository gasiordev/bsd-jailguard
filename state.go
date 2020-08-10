package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
)

type State struct {
	Version     string          `json:"version"`
	Software    string          `json:"software"`
	Created     string          `json:"created"`
	LastUpdated string          `json:"last_updated"`
	Iteration   int             `json:"iteration"`
	History     []*HistoryEntry `json:"history"`

	Bases  map[string]*Base  `json:"bases"`
	Jails  map[string]*Jail  `json:"jails"`
	Netifs map[string]*Netif `json:"network_interfaces"`

	JailPortFwds  map[string]*JailPortFwd `json:"jail_port_fwds"`
	JailNATPasses map[string]*JailNATPass `json:"jail_nat_passes"`

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
	st.History = append(st.History, he)
}

func (st *State) GetBase(rls string) (*Base, error) {
	st.logger(LOGDBG, fmt.Sprintf("Getting base %s from the state...", rls))
	if st.Bases[rls] == nil {
		st.logger(LOGDBG, fmt.Sprintf("Base %s has not been found in the state", rls))
		return nil, nil
	}
	st.logger(LOGDBG, fmt.Sprintf("Base %s has been found in the state", rls))
	return st.Bases[rls], nil
}

func (st *State) GetJail(jl string) (*Jail, error) {
	st.logger(LOGDBG, fmt.Sprintf("Getting jail %s from the state...", jl))
	if st.Jails[jl] == nil {
		st.logger(LOGDBG, fmt.Sprintf("Jail %s has not been found in the state", jl))
		return nil, nil
	}
	st.logger(LOGDBG, fmt.Sprintf("Jail %s has been found in the state", jl))
	return st.Jails[jl], nil
}

func (st *State) GetNetif(ni string) (*Netif, error) {
	st.logger(LOGDBG, fmt.Sprintf("Getting netif %s from the state...", ni))
	if st.Netifs[ni] == nil {
		st.logger(LOGDBG, fmt.Sprintf("Netif %s has not been found in the state", ni))
		return nil, nil
	}
	st.logger(LOGDBG, fmt.Sprintf("Netif %s has been found in the state", ni))
	return st.Netifs[ni], nil
}

func (st *State) GetJailPortFwd(n string) *JailPortFwd {
	st.logger(LOGDBG, fmt.Sprintf("Getting port fwd %s from the state...", n))
	if st.JailPortFwds[n] == nil {
		st.logger(LOGDBG, fmt.Sprintf("Jail port fwd %s has not been found in the state", n))
		return nil
	}
	st.logger(LOGDBG, fmt.Sprintf("Jail port fwd %s has been found in the state", n))
	return st.JailPortFwds[n]
}

func (st *State) GetJailNATPass(n string) *JailNATPass {
	st.logger(LOGDBG, fmt.Sprintf("Getting jail nat pass for jail %s from the state...", n))

	if st.JailNATPasses[n] == nil {
		st.logger(LOGDBG, fmt.Sprintf("Jail nat pass for jail %s has not been found in the state", n))
		return nil
	}
	st.logger(LOGDBG, fmt.Sprintf("NAT pass for jail %s has been found in the state", n))
	return st.JailNATPasses[n]
}

func (st *State) IsJailPortFwdPrefixExists(prfx string) bool {
	st.logger(LOGDBG, fmt.Sprintf("Checking if jail port fwd prefixed %s exists in the state file...", prfx))
	re := regexp.MustCompile("^" + prfx)
	for _, k := range reflect.ValueOf(st.JailPortFwds).MapKeys() {
		if re.Match([]byte(string(k.String()))) {
			st.logger(LOGDBG, fmt.Sprintf("Jail port fwd %s has been found in the state file", k.String()))
			return true
		}
	}
	st.logger(LOGDBG, fmt.Sprintf("No jail port fwds prefixed %s have been found in the state file...", prfx))
	return false
}

func (st *State) GetJailPortFwdsFilterJail(dst_jail string) map[string]*JailPortFwd {
	st.logger(LOGDBG, fmt.Sprintf("Getting jail port fwds with dst jail of %s...", dst_jail))
	m := make(map[string]*JailPortFwd)
	for k, v := range st.JailPortFwds {
		if v.DstJail == dst_jail {
			st.logger(LOGDBG, fmt.Sprintf("Found jail port fwd %s", k))
			m[k] = v
		}
	}
	return m
}

func (st *State) AddBase(rls string, bs *Base) {
	st.logger(LOGDBG, fmt.Sprintf("Adding base %s to the state...", rls))
	st.Bases[rls] = bs
	st.AddHistoryEntry(fmt.Sprintf("Add base %s", rls))
}

func (st *State) AddJail(n string, jl *Jail) {
	st.logger(LOGDBG, fmt.Sprintf("Adding jail %s to the state...", n))
	st.Jails[n] = jl
	st.AddHistoryEntry(fmt.Sprintf("Add jail %s", n))
}

func (st *State) AddNetif(n string, ni *Netif) {
	st.logger(LOGDBG, fmt.Sprintf("Adding netif %s to the state...", n))
	st.Netifs[n] = ni
	st.AddHistoryEntry(fmt.Sprintf("Add netif %s", n))
}

func (st *State) AddJailPortFwd(n string, fwd *JailPortFwd) {
	st.logger(LOGDBG, fmt.Sprintf("Adding jail port fwd from interface %s port %s to jail %s port %s to the state...", fwd.SrcIf, fwd.SrcPort, fwd.DstJail, fwd.DstPort))
	st.JailPortFwds[n] = fwd
	st.AddHistoryEntry(fmt.Sprintf("Add jailportfwd %s", n))
}

func (st *State) AddJailNATPass(n string, np *JailNATPass) {
	st.logger(LOGDBG, fmt.Sprintf("Adding nat pass for jail %s to the state...", n))
	st.JailNATPasses[n] = np
	st.AddHistoryEntry(fmt.Sprintf("Add nat pass for jail %s", n))
}

func (st *State) RemoveItem(t string, n string) error {
	if n == "" {
		return errors.New("Invalid name")
	}

	st.logger(LOGDBG, fmt.Sprintf("Removing item %s %s from the state...", t, n))
	if t == "base" {
		st.Bases[n] = nil
	} else if t == "jail" {
		st.Jails[n] = nil
	} else if t == "netif" {
		st.Netifs[n] = nil
	} else if t == "jailportfwd" {
		st.JailPortFwds[n] = nil
	} else if t == "jailnatpass" {
		st.JailNATPasses[n] = nil
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
	if t == "" || t == "netifs" {
		for k, _ := range st.Netifs {
			fmt.Fprintf(f, "netif %s\n", k)
		}
	}
	if t == "" || t == "jailportfwds" {
		for _, v := range st.JailPortFwds {
			fmt.Fprintf(f, "jailportfwd srcif %s srcport %s dstjail %s dstport %s\n", v.SrcIf, v.SrcPort, v.DstJail, v.DstPort)
		}
	}
	if t == "" || t == "jailnatpasses" {
		for _, v := range st.JailNATPasses {
			fmt.Fprintf(f, "jailnatpass jail %s gwif %s\n", v.JailName, v.GwIf)
		}
	}
	return nil
}

func (st *State) PrintJailPortFwds(f *os.File, n string) {
	for _, v := range st.JailPortFwds {
		if v.DstJail == n {
			fmt.Fprintf(f, "jailportfwd srcif %s srcport %s dstjail %s dstport %s\n", v.SrcIf, v.SrcPort, v.DstJail, v.DstPort)
		}
	}
}

func (st *State) PrintJailNATPass(f *os.File, n string) {
	if st.JailNATPasses[n] != nil {
		fmt.Fprintf(f, "jailnatpass jail %s gwif %s\n", st.JailNATPasses[n].JailName, st.JailNATPasses[n].GwIf)
	}
}

func (st *State) PrintItemItems(f *os.File, t string, n string, k string) error {
	st.logger(LOGDBG, "Printing out state items...")
	if t == "netif" && n != "" && k == "aliases" {
		if st.Netifs[n] == nil {
			return nil
		}
		if st.Netifs[n].Aliases != nil {
			for _, v := range st.Netifs[n].Aliases {
				fmt.Fprintf(f, "netif %s alias %s\n", n, v)
			}
		}
	}
	return nil
}

func (st *State) removeNilFields() {
	if st.History == nil {
		st.History = []*HistoryEntry{}
	}
	if st.Bases == nil {
		st.Bases = make(map[string]*Base)
	}
	if st.Jails == nil {
		st.Jails = make(map[string]*Jail)
	}
	if st.Netifs == nil {
		st.Netifs = make(map[string]*Netif)
	}
	if st.JailPortFwds == nil {
		st.JailPortFwds = make(map[string]*JailPortFwd)
	}
	if st.JailNATPasses == nil {
		st.JailNATPasses = make(map[string]*JailNATPass)
	}
}

func (st *State) removeNilItemsBase() {
	m := make(map[string]*Base)
	for _, k := range reflect.ValueOf(st.Bases).MapKeys() {
		if st.Bases[k.String()] != nil {
			m[k.String()] = st.Bases[k.String()]
		}
	}
	st.Bases = m
}
func (st *State) removeNilItemsJail() {
	m := make(map[string]*Jail)
	for _, k := range reflect.ValueOf(st.Jails).MapKeys() {
		if st.Jails[k.String()] != nil {
			m[k.String()] = st.Jails[k.String()]
		}
	}
	st.Jails = m
}
func (st *State) removeNilItemsNetif() {
	m := make(map[string]*Netif)
	for _, k := range reflect.ValueOf(st.Netifs).MapKeys() {
		if st.Netifs[k.String()] != nil {
			m[k.String()] = st.Netifs[k.String()]
		}
	}
	st.Netifs = m
}
func (st *State) removeNilItemsJailPortFwd() {
	m := make(map[string]*JailPortFwd)
	for _, k := range reflect.ValueOf(st.JailPortFwds).MapKeys() {
		if st.JailPortFwds[k.String()] != nil {
			m[k.String()] = st.JailPortFwds[k.String()]
		}
	}
	st.JailPortFwds = m
}
func (st *State) removeNilItemsJailNATPass() {
	m := make(map[string]*JailNATPass)
	for _, k := range reflect.ValueOf(st.JailNATPasses).MapKeys() {
		if st.JailNATPasses[k.String()] != nil {
			m[k.String()] = st.JailNATPasses[k.String()]
		}
	}
	st.JailNATPasses = m
}

func (st *State) removeNilItems() {
	st.removeNilItemsBase()
	st.removeNilItemsJail()
	st.removeNilItemsNetif()
	st.removeNilItemsJailPortFwd()
	st.removeNilItemsJailNATPass()
}

func (st *State) Save() error {
	st.logger(LOGDBG, "Preparing the state to be saved into the file...")
	st.SetDefaultValues()
	st.LastUpdated = GetCurrentDateTime()
	if st.Created == "" {
		st.Created = st.LastUpdated
	}

	st.removeNilFields()
	st.removeNilItems()

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

	st.removeNilFields()
	st.removeNilItems()

	return st, nil
}
