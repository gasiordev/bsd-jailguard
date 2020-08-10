package main

import (
	"errors"
	"fmt"
)

func (j *Jailguard) getNewPFAnchor(n string) *PFAnchor {
	a := NewPFAnchor(n)
	a.SetLogger(func(t int, s string) {
		j.Log(t, s)
	})
	return a
}

func (j *Jailguard) CheckPFAnchor(hide bool) error {
	a := j.getNewPFAnchor(PFANCHOR)
	a.SetLogger(func(t int, s string) {
		j.Log(t, s)
	})
	ex_nat, ex_rdr, ex, err := a.Exists()
	if err != nil {
		return errors.New("Error occurred while checking if pf anchors exist")
	}

	if !ex_nat {
		j.Log(LOGINF, fmt.Sprintf("nat-anchor \"%s/*\" is missing", PFANCHOR))
	}
	if !ex_rdr {
		j.Log(LOGINF, fmt.Sprintf("rdr-anchor \"%s/*\" is missing", PFANCHOR))
	}
	if !ex {
		j.Log(LOGINF, fmt.Sprintf("anchor \"%s/*\" is missing", PFANCHOR))
	}

	if !ex_nat || !ex_rdr || !ex {
		a.PrintHelp()
		return errors.New("Required anchors are missing")
	} else if !hide {
		j.Log(LOGINF, "All required anchors exist")
	}

	return nil
}
