package main

import (
	"errors"
	"fmt"
	"os"
)

type Base struct {
	Release     string `json:"release"`
	SourceURL   string `json:"source_url"`
	Created     string `json:"created"`
	LastUpdated string `json:"last_updated"`
	Iteration   int    `json:"iteration"`

	Dirpath string          `json:"dirpath"`
	History []*HistoryEntry `json:"history"`

	logger func(int, string)
}

func (bs *Base) SetLogger(f func(int, string)) {
	bs.logger = f
}

func (bs *Base) SetDefaultValues() {
	bs.Iteration = 1
}

func (bs *Base) AddHistoryEntry(s string) {
	he := NewHistoryEntry(GetCurrentDateTime(), s)
	if bs.History == nil {
		bs.History = []*HistoryEntry{}
	}
	bs.History = append(bs.History, he)
}

func (bs *Base) Download(ow bool) error {
	_, _, err := StatWithLog(bs.Dirpath, bs.logger)

	if err != nil {
		if os.IsNotExist(err) {
			bs.logger(LOGDBG, "Jail directory does not exist and it has to be created")
			err2 := CreateDirWithLog(bs.Dirpath, bs.logger)
			if err2 != nil {
				return err2
			}
		} else {
			return errors.New("Error has occurred when downloading base")
		}
	} else {
		if !ow {
			return errors.New(fmt.Sprintf("Base %s already exists. Use 'overwrite' flag to remove it and download again", bs.Release))
		} else {
			bs.logger(LOGDBG, fmt.Sprintf("Base %s already exists but 'overwrite' flag was provided so it will be re-created", bs.Release))
			bs.Iteration++

			err2 := RemoveAllWithLog(bs.Dirpath, bs.logger)
			if err2 != nil {
				return errors.New("Error has occurred when removing existing base")
			}

			err2 = CreateDirWithLog(bs.Dirpath, bs.logger)
			if err2 != nil {
				return errors.New("Error has occurred when creating new directory for base")
			}
		}
	}

	if err == nil && ow {
		bs.AddHistoryEntry("Download again (overwrite)")
	}

	url := fmt.Sprintf("http://ftp.freebsd.org/pub/FreeBSD/releases/amd64/%s/base.txz", bs.Release)
	err = CmdFetchWithLog(url, bs.Dirpath+"/base.txz", bs.logger)
	if err != nil {
		return errors.New("Error has occurred when downloading base. Please try again or fix base manually")
	}
	bs.LastUpdated = GetCurrentDateTime()
	bs.SourceURL = url

	return nil
}

func (bs *Base) Import() error {
	_, _, err := StatWithLog(bs.Dirpath+"/base.txz", bs.logger)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("Base file (base.txz) has not been found")
		} else {
			return errors.New("Error has occurred when importing base")
		}
	}
	bs.logger(LOGDBG, "Base source exists and it can be imported")
	return nil
}

func (bs *Base) Remove() error {
	_, _, err := StatWithLog(bs.Dirpath, bs.logger)
	if err != nil {
		if os.IsNotExist(err) {
			bs.logger(LOGDBG, "Nothing to remove")
			return nil
		} else {
			return err
		}
	}

	return RemoveAllWithLog(bs.Dirpath, bs.logger)
}

func (bs *Base) CreateJailSource(p string) error {
	_, _, err := StatWithLog(p, bs.logger)
	if err != nil && !os.IsNotExist(err) {
		return errors.New("Error has occurred when creating jail directory")
	}
	if err == nil {
		return errors.New("Jail directory already exists")
	}

	err = CreateDirWithLog(p, bs.logger)
	if err != nil {
		return err
	}

	err = CmdTarExtractWithLog(bs.Dirpath+"/base.txz", p, bs.logger)
	if err != nil {
		return errors.New("Error has occurred when extracting base")
	}
	bs.logger(LOGDBG, fmt.Sprintf("Jail source directory %s has been successfully created", p))

	bs.AddHistoryEntry(fmt.Sprintf("Create jail source directory %s", p))

	return nil
}

func (bs *Base) GetBaseTarballPath() string {
	return bs.Dirpath + "/base.txz"
}

func NewBase(rls string, dir string) *Base {
	bs := &Base{}
	bs.SetDefaultValues()
	bs.Release = rls
	bs.Dirpath = dir
	bs.Created = GetCurrentDateTime()
	return bs
}
