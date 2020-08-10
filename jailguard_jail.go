package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

func (j *Jailguard) getJailDirPath(jl string) string {
	return PATHDATA + "/" + DIRJAILS + "/" + jl
}

func (j *Jailguard) getConfigFilePath(jl string) string {
	return PATHDATA + "/" + DIRCONFIGS + "/" + jl + ".jail"
}

func (j *Jailguard) getNewJail(cfg *JailConf, dir *JailDir) *Jail {
	jl := NewJail(cfg, dir)
	jl.SetLogger(func(t int, s string) {
		j.Log(t, s)
	})
	return jl
}

func (j *Jailguard) getJailConf(f string) (*JailConf, error) {
	cfg := NewJailConf()
	cfg.SetLogger(func(t int, s string) {
		j.Log(t, s)
	})
	err := cfg.ParseFile(f)
	if err != nil {
		return nil, err
	}
	err = cfg.Validate()
	if err != nil {
		return nil, err
	}

	if cfg.Config["host.hostname"] == "" {
		cfg.Config["host.hostname"] = cfg.Name
	}

	if cfg.Config["exec.start"] == "" {
		cfg.Config["exec.start"] = "/bin/sh /etc/rc"
	}
	if cfg.Config["exec.stop"] == "" {
		cfg.Config["exec.stop"] = "/bin/sh /etc/rc.shutdown"
	}

	return cfg, nil
}

func (j *Jailguard) getJailDir(n string, d string) *JailDir {
	dir := NewJailDir(n, d)
	dir.SetLogger(func(t int, s string) {
		j.Log(t, s)
	})
	return dir
}

func (j *Jailguard) getOSRelease() (string, error) {
	out, err := CmdOut(j.Log, "uname", "-m", "-r")
	if err != nil {
		return "", errors.New("Error getting OS release: " + err.Error())
	}
	a := strings.Split(string(out), " ")
	var re = regexp.MustCompile(`-p[0-9]+$`)
	return strings.TrimSpace(re.ReplaceAllString(a[0], "")), nil
}

func (j *Jailguard) getJailAndCheckIfExistsInOS(n string, fn func(int, string)) (*State, *Jail, bool, error) {
	st, err := j.getState()
	if err != nil {
		return nil, nil, false, err
	}

	jl, err := st.GetJail(n)
	if err != nil {
		return st, nil, false, err
	}

	ex, err := JailExistsInOSWithLog(n, fn)
	if err != nil {
		return st, jl, false, err
	}

	if jl != nil {
		jl.SetLogger(func(t int, s string) {
			j.Log(t, s)
		})
		if jl.Config != nil {
			jl.Config.SetLogger(func(t int, s string) {
				j.Log(t, s)
			})
		}
		if jl.Dir != nil {
			jl.Dir.SetLogger(func(t int, s string) {
				j.Log(t, s)
			})
		}
	}
	return st, jl, ex, nil
}

func (j *Jailguard) StopJail(n string) error {
	st, jl, ex, err := j.getJailAndCheckIfExistsInOS(n, j.Log)
	if err != nil {
		return err
	}

	if !ex {
		return nil
	}

	if jl == nil {
		j.Log(LOGDBG, fmt.Sprintf("There is a jail %s running in the system", n))
		return errors.New("Jail does not exist in state file but there is a jail with same name running in the system. Stop it manually or import into the state")
	}

	err = jl.Stop()
	if err != nil {
		return errors.New("Error stopping jail")
	}

	st.AddHistoryEntry(fmt.Sprintf("Stop jail %s", n))

	err = st.Save()
	if err != nil {
		return err
	}

	return nil

}

func (j *Jailguard) StartJail(n string) error {
	st, jl, ex, err := j.getJailAndCheckIfExistsInOS(n, j.Log)
	if err != nil {
		return err
	}

	if ex {
		if jl == nil {
			j.Log(LOGDBG, fmt.Sprintf("There is a jail %s running in the system", n))
			return errors.New("Jail does not exist in state file but there is a jail with same name running in the system")
		}
		return nil
	}
	if jl == nil {
		return errors.New("Jail does not exist in state file")
	}

	err = jl.Start()
	if err != nil {
		return errors.New("Error starting jail")
	}

	st.AddHistoryEntry(fmt.Sprintf("Start jail %s", n))

	err = st.Save()
	if err != nil {
		return err
	}

	return nil

}

func (j *Jailguard) RemoveJail(n string, stop bool) error {
	st, jl, ex, err := j.getJailAndCheckIfExistsInOS(n, j.Log)
	if err != nil {
		return err
	}

	if ex {
		if jl == nil {
			j.Log(LOGDBG, fmt.Sprintf("Jail %s is running in the system", n))
			return errors.New("Jail does not exist in state file but there is a jail with same name running in the system. Remove it manually or import the state of it")
		} else if !stop {
			return errors.New("Please stop jail first or use --stop")
		}
	}
	if jl == nil {
		return nil
	}

	if ex && stop {
		err = jl.Stop()
		if err != nil {
			return errors.New("Error stopping jail")
		}
	}

	err = jl.Remove()
	if err != nil {
		return errors.New("Error removing jail")
	}

	st.RemoveItem("jail", n)

	err = st.Save()
	if err != nil {
		return err
	}

	return nil

}

func (j *Jailguard) CreateJail(f string, rls string, start bool) error {
	cfg, err := j.getJailConf(f)
	if err != nil {
		return err
	}

	dir := j.getJailDir(cfg.Name, j.getJailDirPath(cfg.Name))

	st, jl, ex, err := j.getJailAndCheckIfExistsInOS(cfg.Name, j.Log)
	if err != nil {
		return err
	}
	if jl != nil {
		return errors.New(fmt.Sprintf("Jail %s already exists in state file", cfg.Name))
	}
	if ex {
		return errors.New(fmt.Sprintf("Jail %s already exists in the system", cfg.Name))
	}

	var errCreateDir error
	var errWriteCfg error

	if cfg.Config["path"] == "" {
		if rls == "" {
			rls, err = j.getOSRelease()
			if err != nil {
				return errors.New("Error getting OS release")
			}
		}

		bs, err := st.GetBase(rls)
		if err != nil {
			return err
		}
		if bs == nil {
			return errors.New(fmt.Sprintf("Base %s not found in state file", rls))
		}

		bs.SetLogger(func(t int, s string) {
			j.Log(t, s)
		})

		errCreateDir = dir.CreateFromTarball(bs.GetBaseTarballPath())
		cfg.Config["path"] = j.getJailDirPath(cfg.Name)
	} else {
		if rls != "" {
			j.Log(LOGINF, "'path' is provided in the file so base flag will be ignored")
		}
	}

	j.Log(LOGDBG, "Writing jail config to a file...")
	errWriteCfg = cfg.Write(j.getConfigFilePath(cfg.Name))

	jl = j.getNewJail(cfg, dir)
	if errWriteCfg != nil || errCreateDir != nil {
		jl.CleanAfterError()
	}

	if errCreateDir != nil {
		return errors.New("Error creating jail source directory")
	}
	if errWriteCfg != nil {
		return errors.New("Error creating config file")
	}

	st.AddJail(cfg.Name, jl)

	err = st.Save()
	if err != nil {
		return errors.New("Jail has been created but there was an error with writing state. Try to import the state of the jail using state_import")
	}

	if start {
		j.Log(LOGDBG, "Starting jail")
		err1 := jl.Start()
		err2 := st.Save()
		if err2 != nil {
			return errors.New("Error has occurred while saving state")
		}
		if err1 != nil {
			return errors.New("Error has occurred while starting jail")
		}
	}

	return nil
}
