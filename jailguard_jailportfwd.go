package main

import (
	"errors"
	"fmt"
)

func (j *Jailguard) getNewJailPortFwd(src_if string, src_port string, dst_jail string, dst_port string) *JailPortFwd {
	fwd := NewJailPortFwd(src_if, src_port, dst_jail, dst_port)
	fwd.SetLogger(func(t int, s string) {
		j.Log(t, s)
	})
	return fwd
}

func (j *Jailguard) AddJailPortFwd(src_if string, src_port string, dst_jail string, dst_port string) error {
	err := j.CheckPFAnchor(true)
	if err != nil {
		return err
	}

	st, jl, ex, err := j.getJailAndCheckIfExistsInOS(dst_jail, j.Log)
	if err != nil {
		return err
	}
	if jl == nil {
		return errors.New(fmt.Sprintf("Jail %s does not exist in state", dst_jail))
	}

	ex = st.IsJailPortFwdPrefixExists(fmt.Sprintf("%s__%s__", src_if, src_port))
	if ex {
		return errors.New(fmt.Sprintf("Interface %s port %s is already forwarded", src_if, src_port))
	}

	fwd := j.getNewJailPortFwd(src_if, src_port, dst_jail, dst_port)

	st.AddJailPortFwd(fmt.Sprintf("%s__%s__%s__%s", fwd.SrcIf, fwd.SrcPort, fwd.DstJail, fwd.DstPort), fwd)

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

func (j *Jailguard) DeleteJailPortFwd(src_if string, src_port string, dst_jail string, dst_port string) error {
	err := j.CheckPFAnchor(true)
	if err != nil {
		return err
	}

	st, jl, _, err := j.getJailAndCheckIfExistsInOS(dst_jail, j.Log)
	if err != nil {
		return err
	}
	if jl == nil {
		return errors.New(fmt.Sprintf("Jail %s does not exist in state", dst_jail))
	}

	fwd := st.GetJailPortFwd(fmt.Sprintf("%s__%s__%s__%s", src_if, src_port, dst_jail, dst_port))
	if fwd == nil {
		return nil
	}

	st.RemoveItem("jailportfwd", fmt.Sprintf("%s__%s__%s__%s", fwd.SrcIf, fwd.SrcPort, fwd.DstJail, fwd.DstPort))

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

func (j *Jailguard) DeleteJailAllPortFwds(dst_jail string) error {
	err := j.CheckPFAnchor(true)
	if err != nil {
		return err
	}

	st, jl, _, err := j.getJailAndCheckIfExistsInOS(dst_jail, j.Log)
	if err != nil {
		return err
	}
	if jl == nil {
		return errors.New(fmt.Sprintf("Jail %s does not exist in state", dst_jail))
	}

	err = nil
	fwds := st.GetJailPortFwdsFilterJail(dst_jail)
	for k, _ := range fwds {
		if err == nil {
			st.RemoveItem("jailportfwd", k)
		}
	}

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

func (j *Jailguard) ListJailPortFwds(n string) error {
	st, err := j.getState()
	if err != nil {
		return err
	}

	st.PrintJailPortFwds(j.cli.GetStdout(), n)

	return nil
}
