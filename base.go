package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"
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

func (bs *Base) cmdOut(c string, a ...string) ([]byte, error) {
	cmd := exec.Command(c, a...)
	cmd.Stdin = os.Stdin
	return cmd.Output()
}

func (bs *Base) statWithLog(p string) (os.FileInfo, bool, error) {
	bs.logger(LOGDBG, fmt.Sprintf("Getting stat for path %s...", p))
	st, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			bs.logger(LOGDBG, fmt.Sprintf("Path %s does not exist", p))
		} else {
			bs.logger(LOGDBG, fmt.Sprintf("Error has occurred when getting stat for path %s: %s", p, err.Error()))
		}
		return st, false, err
	} else {
		bs.logger(LOGDBG, fmt.Sprintf("Found path %s", p))
		if st.IsDir() {
			bs.logger(LOGDBG, fmt.Sprintf("Path %s is a directory", p))
			return st, true, nil
		}
	}
	return st, false, nil
}

func (bs *Base) removeAllWithLog(p string) error {
	bs.logger(LOGDBG, fmt.Sprintf("Removing %s...", p))
	err := os.RemoveAll(p)
	if err != nil {
		bs.logger(LOGDBG, fmt.Sprintf("Error has occurred when removing %s: %s", p, err.Error()))
	}
	bs.logger(LOGDBG, fmt.Sprintf("Path %s has been removed", p))
	return err
}

func (bs *Base) createDirWithLog(p string) error {
	bs.logger(LOGDBG, fmt.Sprintf("Creating directory %s...", p))
	err := os.MkdirAll(p, os.ModePerm)
	if err != nil {
		bs.logger(LOGDBG, fmt.Sprintf("Error has occurred when creating %s: %s", p, err.Error()))
	}
	bs.logger(LOGDBG, fmt.Sprintf("Directory %s has been created", p))
	return err
}
func (bs *Base) getCurrentDateTime() string {
	return time.Now().String()
}

func (bs *Base) SetLogger(f func(int, string)) {
	bs.logger = f
}

func (bs *Base) SetDefaultValues() {
	bs.Iteration = 1
}

func (bs *Base) Download(ow bool) error {
	_, _, err := bs.statWithLog(bs.Dirpath)
	if err != nil {
		if os.IsNotExist(err) {
			bs.logger(LOGDBG, "Jail directory does not exist and it has to be created")
			err2 := bs.createDirWithLog(bs.Dirpath)
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
			err2 := bs.removeAllWithLog(bs.Dirpath)
			if err2 != nil {
				return errors.New("Error has occurred when removing existing base")
			}

			err2 = bs.createDirWithLog(bs.Dirpath)
			if err2 != nil {
				return errors.New("Error has occurred when creating new directory for base")
			}
		}
	}

	url := "http://ftp.freebsd.org/pub/FreeBSD/releases/amd64/" + bs.Release + "/base.txz"
	bs.logger(LOGDBG, "Trying to download base.txz from "+url+" using fetch...")
	_, err = bs.cmdOut("fetch", url, "-o", bs.Dirpath+"/base.txz")
	if err != nil {
		return errors.New(fmt.Sprintf("Error has occurred when downloading base.txz from %s: %s", url, err))
	}
	bs.logger(LOGDBG, fmt.Sprintf("File %s has been successfully saved in %s/base.txz", url, bs.Dirpath))

	bs.LastUpdated = bs.getCurrentDateTime()
	return nil
}

func (bs *Base) Import() error {
	_, _, err := bs.statWithLog(bs.Dirpath + "/base.txz")
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
	_, _, err := bs.statWithLog(bs.Dirpath)
	if err != nil {
		if os.IsNotExist(err) {
			bs.logger(LOGDBG, "Nothing to remove")
			return nil
		} else {
			return err
		}
	}

	return bs.removeAllWithLog(bs.Dirpath)
}

func (bs *Base) CreateJailSource(p string) error {
	_, _, err := bs.statWithLog(p)
	if err != nil && !os.IsNotExist(err) {
		return errors.New("Error has occurred when creating jail directory")
	}
	if err == nil {
		return errors.New("Jail directory already exists")
	}

	err = bs.createDirWithLog(p)
	if err != nil {
		return err
	}

	bs.logger(LOGDBG, fmt.Sprintf("Running tar to extract %s/base.txz to %s...", bs.Dirpath, p))
	_, err = bs.cmdOut("tar", "-xvf", bs.Dirpath+"/base.txz", "-C", p)
	if err != nil {
		return errors.New("Error has occurred when extracting base.txz")
	}

	bs.logger(LOGDBG, fmt.Sprintf("Jail source directory %s has been successfully created", p))
	return nil
}

func NewBase(rls string, dir string) *Base {
	bs := &Base{}
	bs.SetDefaultValues()
	bs.Release = rls
	bs.Dirpath = dir
	bs.Created = bs.getCurrentDateTime()
	return bs
}
