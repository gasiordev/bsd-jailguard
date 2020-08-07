package main

import (
	"fmt"
	"regexp"
)

type PFAnchor struct {
	Name   string `json:"name"`
	logger func(int, string)
}

func (a *PFAnchor) SetLogger(f func(int, string)) {
	a.logger = f
}

func (a *PFAnchor) Exists() (bool, bool, bool, error) {
	a.logger(LOGDBG, fmt.Sprintf("Checking if anchors %s exist in the system...", a.Name))
	out, err := CmdOut(a.logger, "pfctl", "-s", "all")
	if err != nil {
		return false, false, false, err
	}

	var re *regexp.Regexp
	re = regexp.MustCompile(fmt.Sprintf("nat-anchor \"%s\\/\\*\"", a.Name))
	ex_nat := re.Match([]byte(string(out)))
	if !ex_nat {
		a.logger(LOGDBG, fmt.Sprintf("nat-anchor '%s/*' does not exist in pf", a.Name))
	} else {
		a.logger(LOGDBG, fmt.Sprintf("nat-anchor '%s/*' exists in pf", a.Name))

	}

	re = regexp.MustCompile(fmt.Sprintf("rdr-anchor \"%s\\/\\*\"", a.Name))
	ex_rdr := re.Match([]byte(string(out)))
	if !ex_rdr {
		a.logger(LOGDBG, fmt.Sprintf("rdr-anchor '%s/*' does not exist in pf", a.Name))
	} else {
		a.logger(LOGDBG, fmt.Sprintf("rdr-anchor '%s/*' exists in pf", a.Name))

	}

	re = regexp.MustCompile(fmt.Sprintf("anchor \"%s\\/\\*\"", a.Name))
	ex := re.Match([]byte(string(out)))
	if !ex {
		a.logger(LOGDBG, fmt.Sprintf("anchor '%s/*' does not exist in pf", a.Name))
	} else {
		a.logger(LOGDBG, fmt.Sprintf("anchor '%s/*' exists in pf", a.Name))

	}

	return ex_nat, ex_rdr, ex, nil
}

func (a *PFAnchor) PrintHelp() {
	a.logger(LOGINF, fmt.Sprintf(`To create an anchor, ensure you have the following lines in your pf.conf (or other path you use):

nat-anchor "%s/*"
rdr-anchor "%s/*"
anchor "%s/*"

Run 'pfctl -f /etc/pf.conf' to flush the PF rules.
Please note that it will remove all the rules and load ones from the file.
`, a.Name, a.Name, a.Name))
}

func NewPFAnchor(n string) *PFAnchor {
	a := &PFAnchor{Name: n}
	return a
}
