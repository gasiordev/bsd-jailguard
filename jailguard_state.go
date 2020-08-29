package main

import (
	"errors"
	"fmt"
)

func (j *Jailguard) getStateFilePath() string {
	c := j.GetConfig()
	return c.PathData + "/" + c.DirState + "/" + c.FileState
}

func (j *Jailguard) getState() (*State, error) {
	st, err := NewState(j.getStateFilePath())
	if err != nil {
		return nil, err
	}
	st.SetLogger(func(t int, s string) {
		j.Log(t, s)
	})
	return st, nil
}

func (j *Jailguard) RemoveStateItem(t string, n string) error {
	st, err := j.getState()
	if err != nil {
		return err
	}

	err = st.RemoveItem(t, n)
	if err != nil {
		return err
	}
	return nil
}

func (j *Jailguard) ImportStateItem(t string, n string) error {
	st, err := j.getState()
	if err != nil {
		return err
	}

	if t == "base" {
		bs, err := st.GetBase(n)
		if err != nil {
			return err
		}
		if bs != nil {
			j.Log(LOGERR, fmt.Sprintf("Base %s already exists. Remove it first before importing a new one"))
			return errors.New("State item already exists")
		}

		bs = j.getNewBase(n)
		j.Log(LOGDBG, fmt.Sprintf("Checking if base %s can be imported into state file...", n))
		err = bs.Import()
		if err != nil {
			return err
		}
		st.AddBase(n, bs)
	} else if t == "jail" {
		jl, err := st.GetJail(n)
		if err != nil {
			return err
		}
		if jl != nil {
			j.Log(LOGERR, fmt.Sprintf("Jail %s already exists. Remove it first before import a new one"))
			return errors.New("State item already exists")
		}

		cfg, err := j.getJailConf(j.getConfigFilePath(n))
		if err != nil {
			return err
		}
		cfg.Filepath = j.getConfigFilePath(n)

		dir := j.getJailDir(n, j.getJailDirPath(n))

		jl = j.getNewJail(cfg, dir)
		j.Log(LOGDBG, fmt.Sprintf("Checking if base %s can be imported into state file...", n))
		err = jl.Import()
		if err != nil {
			return err
		}
		st.AddJail(n, jl)
		// netif
		// pf anchor
	} else {
		return errors.New("Invalid item type")
	}

	err = st.Save()
	if err != nil {
		return err
	}
	return nil
}

func (j *Jailguard) ListStateItems(t string) error {
	st, err := j.getState()
	if err != nil {
		return err
	}

	err = st.PrintItems(j.cli.GetStdout(), t)
	if err != nil {
		return err
	}
	return nil
}
