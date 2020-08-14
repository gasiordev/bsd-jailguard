package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Netif struct {
	Name        string   `json:"name"`
	SystemName  string   `json:"system_name"`
	IPAddrBegin string   `json:"ip_addr_begin"`
	IPAddrEnd   string   `json:"ip_addr_end"`
	Aliases     []string `json:"aliases"`
	Created     string   `json:"created"`
	LastUpdated string   `json:"last_updated"`

	Iteration int             `json:"iteration"`
	History   []*HistoryEntry `json:"history"`

	logger func(int, string)
}

func (ni *Netif) SetLogger(f func(int, string)) {
	ni.logger = f
}

func (ni *Netif) SetDefaultValues() {
	ni.Iteration = 1
}

func (ni *Netif) AddHistoryEntry(s string) {
	he := NewHistoryEntry(GetCurrentDateTime(), s)
	if ni.History == nil {
		ni.History = []*HistoryEntry{}
	}
	ni.History = append(ni.History, he)
}

func (ni *Netif) getFreeSystemName() (string, error) {
	ni.logger(LOGDBG, "Getting system name for new interface")
	f := 0
	for i := 1; i < 3000; i++ {
		_, err := CmdOut(ni.logger, "ifconfig", "lo"+strconv.Itoa(i))
		if err != nil {
			f = i
			break
		}
	}
	if f == 0 {
		return "", errors.New("Error has occurred whilst getting a name for a network interface")
	}
	s := "lo" + strconv.Itoa(f)
	ni.logger(LOGDBG, fmt.Sprintf("Found free network interface name %s", s))
	return s, nil
}

func (ni *Netif) isSystemNameValid() bool {
	re := regexp.MustCompile(`^lo[0-9]{1,3}$`)
	if re.Match([]byte(ni.SystemName)) {
		return true
	}
	return false
}

func (ni *Netif) isSystemNameExists() (bool, error) {
	ni.logger(LOGDBG, fmt.Sprintf("Checking if interface name %s already exists in the system", ni.SystemName))
	out, err := CmdOut(ni.logger, "ifconfig", ni.SystemName)
	if err != nil {
		return false, errors.New("Error has occurred whilst checking if interface exist")
	}
	re := regexp.MustCompile(ni.SystemName + ":")
	if re.Match([]byte(string(out))) {
		ni.logger(LOGDBG, fmt.Sprintf("Network interface %s already exists in the system", ni.SystemName))
		return true, nil
	}
	ni.logger(LOGDBG, fmt.Sprintf("Network interface %s does not exist in the system", ni.SystemName))
	return false, nil
}

func (ni *Netif) ifconfigUp() error {
	ni.logger(LOGDBG, fmt.Sprintf("Bringing interface %s up", ni.SystemName))
	err := CmdRun(ni.logger, "ifconfig", ni.SystemName, "up")
	if err != nil {
		return err
	}
	ni.logger(LOGDBG, fmt.Sprintf("Network interface %s (jailguard name: %s) is now up", ni.SystemName, ni.Name))
	return nil
}

func (ni *Netif) ifconfigCreate() error {
	ni.logger(LOGDBG, fmt.Sprintf("Creating network interface %s", ni.SystemName))
	err := CmdRun(ni.logger, "ifconfig", ni.SystemName, "create")
	if err != nil {
		return err
	}
	ni.logger(LOGDBG, fmt.Sprintf("Network interface %s (jailguard name: %s) is now created", ni.SystemName, ni.Name))
	return nil
}

func (ni *Netif) ifconfigDestroy() error {
	ni.logger(LOGDBG, fmt.Sprintf("Destroying network interface %s", ni.SystemName))
	err := CmdRun(ni.logger, "ifconfig", ni.SystemName, "destroy")
	if err != nil {
		return err
	}
	ni.logger(LOGDBG, fmt.Sprintf("Network interface %s (jailguard name: %s) is now destroyed", ni.SystemName, ni.Name))
	return nil
}

func (ni *Netif) isAliasExists(ip string) (bool, error) {
	ni.logger(LOGDBG, fmt.Sprintf("Checking if interface name %s has an alias of %s in the system", ni.SystemName, ip))
	out, err := CmdOut(ni.logger, "ifconfig", ni.SystemName)
	if err != nil {
		return false, errors.New("Error has occurred whilst checking if interface has an alias")
	}
	re := regexp.MustCompile("inet " + ip + " netmask")
	if re.Match([]byte(string(out))) {
		ni.logger(LOGDBG, fmt.Sprintf("Network interface %s already has an alias of %s in the system", ni.SystemName, ip))
		return true, nil
	}
	ni.logger(LOGDBG, fmt.Sprintf("Network interface %s does not have an alias of %s in the system", ni.SystemName, ip))
	return false, nil
}

func (ni *Netif) ifconfigAliasAdd(ip string) error {
	ni.logger(LOGDBG, fmt.Sprintf("Adding alias %s to network interface %s...", ip, ni.SystemName))
	err := CmdRun(ni.logger, "ifconfig", ni.SystemName, "inet", ip+"/24", "alias")
	if err != nil {
		return err
	}
	ni.logger(LOGDBG, fmt.Sprintf("Alias %s is now added to network interface %s", ip, ni.SystemName))
	return nil
}

func (ni *Netif) ifconfigAliasDelete(ip string) error {
	ni.logger(LOGDBG, fmt.Sprintf("Deleting alias %s from network interface %s...", ip, ni.SystemName))
	err := CmdRun(ni.logger, "ifconfig", ni.SystemName, "inet", ip+"/24", "-alias")
	if err != nil {
		return err
	}
	ni.logger(LOGDBG, fmt.Sprintf("Alias %s is now deleted from network interface %s", ip, ni.SystemName))
	return nil
}

// when ip is empty then pick next available address
func (ni *Netif) AddAlias(ip string) (string, error) {
	if ni.Aliases == nil {
		ni.Aliases = []string{}
	}

	if ip == "" {
		// Find free IP address within a range defined in the network interface
		ip_b := strings.Split(ni.IPAddrBegin, ".")
		ip_e := strings.Split(ni.IPAddrEnd, ".")
		num_b, _ := strconv.Atoi(ip_b[3])
		num_e, _ := strconv.Atoi(ip_e[3])
		found := 0
		for i := num_b; i <= num_e; i++ {
			taken := false
			for _, v := range ni.Aliases {
				i_ip := fmt.Sprintf("%s.%s.%s.%s", ip_b[0], ip_b[1], ip_b[2], strconv.Itoa(i))
				if v == i_ip {
					taken = true
					break
				}
			}
			if !taken {
				found = i
				break
			}
		}
		if found == 0 {
			return "", errors.New("Error has occurred whilst finding a free IP address")
		}
		ip = fmt.Sprintf("%s.%s.%s.%s", ip_b[0], ip_b[1], ip_b[2], strconv.Itoa(found))
	} else {
		for _, v := range ni.Aliases {
			if v == ip {
				return ip, nil
			}
		}
	}

	ex, err := ni.isAliasExists(ip)
	if ex {
		return "", errors.New(fmt.Sprintf("Network interface %s already has an alias of %s", ni.SystemName, ip))
	}

	err = ni.ifconfigAliasAdd(ip)
	if err != nil {
		return "", errors.New("Error has occurred whilst adding alias")
	}

	ni.Aliases = append(ni.Aliases, ip)
	return ip, nil
}

func (ni *Netif) DeleteAlias(ip string) error {
	if ni.Aliases == nil {
		return nil
	}

	ex, err := ni.isAliasExists(ip)
	if ex {
		err = ni.ifconfigAliasDelete(ip)
		if err != nil {
			return errors.New("Error has occurred whilst deleting alias")
		}
	}

	as := []string{}
	for _, v := range ni.Aliases {
		if v != ip {
			as = append(as, v)
		}
	}
	return nil
}

func (ni *Netif) Import() error {
	return nil
}

func (ni *Netif) Create() error {
	ni.logger(LOGDBG, "Interface name was passed so checking if it exists and it's valid")
	if ni.SystemName != "" {
		ex, err := ni.isSystemNameExists()
		if err != nil {
			return err
		}
		if ex {
			return errors.New(fmt.Sprintf("Network interface already exists in the system", ni.SystemName))
		}

		v := ni.isSystemNameValid()
		if !v {
			return errors.New(fmt.Sprintf("Network interface should be 'lo' suffixed with number", ni.SystemName))
		}
	} else {
		nu, err := ni.getFreeSystemName()
		if err != nil {
			return errors.New("Error occurred whilst getting a new network interface name")
		}
		ni.SystemName = nu
	}

	err := ni.ifconfigCreate()
	if err != nil {
		return errors.New("Error creating network interface")
	}

	err = ni.ifconfigUp()
	if err != nil {
		return errors.New("Error bringing network interface up")
	}

	return nil
}

func (ni *Netif) Destroy() error {
	if ni.SystemName != "" {
		ex, err := ni.isSystemNameExists()
		if err != nil {
			return err
		}
		if !ex {
			return nil
		}

		err = ni.ifconfigDestroy()
		if err != nil {
			return errors.New("Error destroying network interface")
		}

	}
	return nil
}

func NewNetif(n string, ip_begin string, ip_end string, if_name string) *Netif {
	ni := &Netif{Name: n, IPAddrBegin: ip_begin, IPAddrEnd: ip_end, SystemName: if_name}
	ni.SetDefaultValues()
	ni.Created = GetCurrentDateTime()
	ni.LastUpdated = GetCurrentDateTime()
	return ni
}
