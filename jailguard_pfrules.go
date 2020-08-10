package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

func (j *Jailguard) getJailPFRulesFilePath(jl string) string {
	return PATHDATA + "/" + DIRCONFIGS + "/" + jl + ".pf"
}

func (j *Jailguard) FlushJailPFRulesFromState(jl *Jail, st *State) error {
	err := j.RecreateJailPFRulesFromState(jl, st)
	if err != nil {
		return err
	}

	err = j.FlushJailPFRulesFromFile(jl.Name)
	if err != nil {
		return err
	}

	return nil
}

func (j *Jailguard) RecreateJailPFRulesFromState(jl *Jail, st *State) error {
	// This function is performed before state is saved so everything needs to be
	// checked if it is not nil
	c := ""
	if jl.Config.Config["ip4.addr"] != "" {
		nat := st.GetJailNATPass(jl.Name)
		fwds := st.GetJailPortFwdsFilterJail(jl.Name)

		if nat != nil && nat.GwIf != "" {
			jl.logger(LOGDBG, fmt.Sprintf("Checking if gateway network interface %s exists...", nat.GwIf))
			err := CmdRun(jl.logger, "ifconfig", nat.GwIf)
			if err != nil {
				return errors.New(fmt.Sprintf("Gateway network interface %s does not exist", nat.GwIf))
			}
			c += fmt.Sprintf("nat pass on %s from %s/32 to any -> (%s:0)\n", nat.GwIf, jl.Config.Config["ip4.addr"], nat.GwIf)
		}

		for _, v := range fwds {
			if v != nil {
				c += fmt.Sprintf("rdr pass on %s inet proto tcp from any to (%s:0) port %s -> %s port %s\n", v.SrcIf, v.SrcIf, v.SrcPort, jl.Config.Config["ip4.addr"], v.DstPort)
			}
		}
	}

	j.Log(LOGDBG, fmt.Sprintf("Writing jail pf rules file to %s...", j.getJailPFRulesFilePath(jl.Name)))
	err := ioutil.WriteFile(j.getJailPFRulesFilePath(jl.Name), []byte(c), 0644)
	if err != nil {
		return err
	}
	j.Log(LOGDBG, fmt.Sprintf("Jail pf rules have been written to %s", j.getJailPFRulesFilePath(jl.Name)))

	return nil
}

func (j *Jailguard) FlushJailPFRulesFromFile(n string) error {
	_, isdir, err := StatWithLog(j.getJailPFRulesFilePath(n), j.Log)
	if err != nil {
		if os.IsNotExist(err) {
			j.Log(LOGINF, fmt.Sprintf("Jail %s does not have any pf rules", n))
			return nil
		} else {
			return errors.New("Error has occurred while flushing jail pf rules from file")
		}
	}
	if isdir {
		return errors.New("Error has occurred while flusing jail pf rules from file. It seems file is actually directory. Fix it manually please")
	}

	err = CmdRun(j.Log, "pfctl", "-a", PFANCHOR+"/"+n, "-f", j.getJailPFRulesFilePath(n))
	if err != nil {
		return errors.New("Error has occurred while flushing jail pf rules from file")
	}

	return nil
}
