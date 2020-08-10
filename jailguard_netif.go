package main

import (
	"errors"
	"fmt"
)

func (j *Jailguard) getNewNetif(n string, ip_begin string, ip_end string, if_name string) *Netif {
	ni := NewNetif(n, ip_begin, ip_end, if_name)
	ni.SetLogger(func(t int, s string) {
		j.Log(t, s)
	})
	return ni
}

func (j *Jailguard) CreateNetif(n string, ip_begin string, ip_end string, if_n string) error {
	st, err := j.getState()
	if err != nil {
		return err
	}

	ni, err := st.GetNetif(n)
	if err != nil {
		return err
	}
	if ni != nil {
		return errors.New(fmt.Sprintf("Network interface %s already exists in the state", n))
	}

	ni = j.getNewNetif(n, ip_begin, ip_end, if_n)
	err = ni.Create()
	if err != nil {
		return err
	}

	st.AddNetif(n, ni)

	err = st.Save()
	if err != nil {
		return err
	}

	return nil
}

func (j *Jailguard) DestroyNetif(n string) error {
	st, err := j.getState()
	if err != nil {
		return err
	}

	ni, err := st.GetNetif(n)
	if err != nil {
		return err
	}
	if ni == nil {
		return nil
	}
	ni.SetLogger(func(t int, s string) {
		j.Log(t, s)
	})

	err = ni.Destroy()
	if err != nil {
		return err
	}

	st.RemoveItem("netif", n)

	err = st.Save()
	if err != nil {
		return err
	}

	return nil
}

func (j *Jailguard) AddNetifAlias(n string, ip string) error {
	st, err := j.getState()
	if err != nil {
		return err
	}

	ni, err := st.GetNetif(n)
	if err != nil {
		return err
	}
	if ni == nil {
		return errors.New(fmt.Sprintf("Network interface %s does not exist", n))
	}
	ni.SetLogger(func(t int, s string) {
		j.Log(t, s)
	})

	_, err = ni.AddAlias(ip)
	if err != nil {
		return err
	}

	err = st.Save()
	if err != nil {
		return err
	}

	return nil
}

func (j *Jailguard) DeleteNetifAlias(n string, ip string) error {
	st, err := j.getState()
	if err != nil {
		return err
	}

	ni, err := st.GetNetif(n)
	if err != nil {
		return err
	}
	if ni == nil {
		return errors.New(fmt.Sprintf("Network interface %s does not exist", n))
	}
	ni.SetLogger(func(t int, s string) {
		j.Log(t, s)
	})

	err = ni.DeleteAlias(ip)
	if err != nil {
		return err
	}

	err = st.Save()
	if err != nil {
		return err
	}

	return nil
}

func (j *Jailguard) ListNetifAliases(n string) error {
	st, err := j.getState()
	if err != nil {
		return err
	}

	err = st.PrintItemItems(j.cli.GetStdout(), "netif", n, "aliases")
	if err != nil {
		return err
	}
	return nil
}
