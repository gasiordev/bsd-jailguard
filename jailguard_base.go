package main

import (
	"fmt"
)

func (j *Jailguard) getBaseDirPath(rls string) string {
	return PATHDATA + "/" + DIRBASES + "/" + rls
}

func (j *Jailguard) getNewBase(rls string) *Base {
	bs := NewBase(rls, j.getBaseDirPath(rls))
	bs.SetLogger(func(t int, s string) {
		j.Log(t, s)
	})
	return bs
}

func (j *Jailguard) DownloadBase(rls string, ow bool) error {
	st, err := j.getState()
	if err != nil {
		return err
	}

	bs, err := st.GetBase(rls)
	if err != nil {
		return err
	}
	if bs == nil {
		bs = j.getNewBase(rls)
		err = bs.Download(ow)
		if err != nil {
			return err
		}
		st.AddBase(rls, bs)
	} else {
		bs.SetLogger(func(t int, s string) {
			j.Log(t, s)
		})

		if ow {
			j.Log(LOGINF, fmt.Sprintf("Base %s already exists but downloading it again...", rls))
			err = bs.Download(ow)
			if err != nil {
				return err
			}
		} else {
			j.Log(LOGINF, fmt.Sprintf("Base %s already exists. Use 'overwrite' flag to download it again", rls))
			return nil
		}
	}

	err = st.Save()
	if err != nil {
		return err
	}

	return nil
}

func (j *Jailguard) RemoveBase(rls string) error {
	st, err := j.getState()
	if err != nil {
		return err
	}

	bs, err := st.GetBase(rls)
	if err != nil {
		return err
	}
	if bs == nil {
		return nil
	}
	bs.SetLogger(func(t int, s string) {
		j.Log(t, s)
	})

	err = bs.Remove()
	if err != nil {
		return err
	}

	st.RemoveItem("base", rls)

	err = st.Save()
	if err != nil {
		return err
	}

	return nil
}
