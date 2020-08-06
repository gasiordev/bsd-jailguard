package main

import (
	"errors"
	"fmt"
	"os"
)

type JailDir struct {
	Name        string `json:"name"`
	Created     string `json:"created"`
	LastUpdated string `json:"last_updated"`
	Dirpath     string `json:"dirpath"`

	Iteration int             `json:"iteration"`
	History   []*HistoryEntry `json:"history"`

	logger func(int, string)
}

func (jd *JailDir) SetLogger(f func(int, string)) {
	jd.logger = f
}

func (jd *JailDir) SetDefaultValues() {
	jd.Iteration = 1
}

func (jd *JailDir) AddHistoryEntry(s string) {
	he := NewHistoryEntry(GetCurrentDateTime(), s)
	if jd.History == nil {
		jd.History = []*HistoryEntry{}
	}
	jd.History = append(jd.History, he)
}

func (jd *JailDir) Remove() error {
	_, _, err := StatWithLog(jd.Dirpath, jd.logger)
	if err != nil {
		if os.IsNotExist(err) {
			jd.logger(LOGDBG, "Nothing to remove")
			return nil
		} else {
			return err
		}
	}

	err1 := CmdRun(jd.logger, "chflags", "-R", "noschg", jd.Dirpath)
	err2 := RemoveAllWithLog(jd.Dirpath, jd.logger)
	if err1 != nil || err2 != nil {
		return errors.New("Error has occurred while removing jail directory. Please remove the directories manually and remove the state")
	}

	return nil
}

func (jd *JailDir) CreateFromTarball(t string) error {
	_, _, err := StatWithLog(jd.Dirpath, jd.logger)
	if err != nil && !os.IsNotExist(err) {
		return errors.New("Error has occurred when creating jail directory")
	}
	if err == nil {
		return errors.New("Jail directory already exists")
	}

	err = CreateDirWithLog(jd.Dirpath, jd.logger)
	if err != nil {
		return err
	}

	err = CmdTarExtractWithLog(t, jd.Dirpath, jd.logger)
	if err != nil {
		return errors.New("Error has occurred when extracting tarball")
	}
	jd.logger(LOGDBG, fmt.Sprintf("Jail source directory %s has been successfully created", jd.Dirpath))

	jd.AddHistoryEntry(fmt.Sprintf("Create jail source directory %s", jd.Dirpath))

	return nil
}

func NewJailDir(n string, dir string) *JailDir {
	jd := &JailDir{}
	jd.SetDefaultValues()
	jd.Name = n
	jd.Dirpath = dir
	jd.Created = GetCurrentDateTime()
	return jd
}
