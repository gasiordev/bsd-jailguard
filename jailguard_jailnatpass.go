package main

import (
	"errors"
	"fmt"
)

func (j *Jailguard) getNewJailNATPass(jl string, gw_if string) *JailNATPass {
	np := NewJailNATPass(gw_if, jl)
	np.SetLogger(func(t int, s string) {
		j.Log(t, s)
	})
	return np
}

func (j *Jailguard) CreateJailNATPass(n string, if_gw string) error {
	err := j.CheckPFAnchor(true)
	if err != nil {
		return err
	}

	st, jl, _, err := j.getJailAndCheckIfExistsInOS(n, j.Log)
	if err != nil {
		return err
	}
	if jl == nil {
		return errors.New(fmt.Sprintf("Jail %s does not exist in state", n))
	}

	nat := st.GetJailNATPass(n)
	if nat != nil {
		return errors.New(fmt.Sprintf("NAT pass for jail %s already exists", n))
	}

	nat = j.getNewJailNATPass(n, if_gw)

	st.AddJailNATPass(n, nat)

	err = j.FlushJailPFRulesFromState(jl, st)
	if err != nil {
		return err
	}

	err = st.Save()
	if err != nil {
		return err
	}

	return nil
}

func (j *Jailguard) RemoveJailNATPass(n string) error {
	err := j.CheckPFAnchor(true)
	if err != nil {
		return err
	}

	st, jl, _, err := j.getJailAndCheckIfExistsInOS(n, j.Log)
	if err != nil {
		return err
	}
	if jl == nil {
		return errors.New(fmt.Sprintf("Jail %s does not exist in state", n))
	}

	nat := st.GetJailNATPass(n)
	if nat == nil {
		return nil
	}

	err = st.RemoveItem("jailnatpass", n)

	err = j.FlushJailPFRulesFromState(jl, st)
	if err != nil {
		return err
	}

	err = st.Save()
	if err != nil {
		return err
	}

	return nil
}

func (j *Jailguard) ShowJailNATPass(n string) error {
	err := j.CheckPFAnchor(true)
	if err != nil {
		return err
	}

	st, err := j.getState()
	if err != nil {
		return err
	}

	st.PrintJailNATPass(j.cli.GetStdout(), n)

	return nil
}
