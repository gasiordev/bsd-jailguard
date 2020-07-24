package main

import (
	"errors"
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
	bs.logger(LOGDBG, "Checking if "+bs.Dirpath+" exists")
	_, err := os.Stat(bs.Dirpath)
	if err != nil {
		if os.IsNotExist(err) {
			bs.logger(LOGDBG, bs.Dirpath+" does not exist, trying to create it")
			err2 := os.MkdirAll(bs.Dirpath, os.ModePerm)
			if err2 != nil {
				bs.logger(LOGERR, "Error with creating "+bs.Dirpath+" dir")
				return err2
			}
		} else {
			bs.logger(LOGDBG, "Error with checking dir "+bs.Dirpath+" existance: "+err.Error())
			return err
		}
	} else {
		if !ow {
			bs.logger(LOGERR, "Base "+bs.Release+" exists already. Use 'overwrite' flag to remove it and download again")
			return errors.New("Base " + bs.Release + " already exists. Use 'overwrite' flag to remove it and download again")
		} else {
			bs.logger(LOGDBG, "Base "+bs.Release+" already exists but 'overwrite' flag was provided so trying to remove the directory")
			err2 := os.RemoveAll(bs.Dirpath)
			if err2 != nil {
				bs.logger(LOGERR, "Error removing dir "+bs.Dirpath+". Please remove the directory manually and remove the state")
				return errors.New("Error removing base dir")
			}

			bs.logger(LOGDBG, "Creating directory "+bs.Dirpath)
			err2 = os.MkdirAll(bs.Dirpath, os.ModePerm)
			if err2 != nil {
				bs.logger(LOGERR, "Error with creating "+bs.Dirpath+" dir")
				return err2
			}

		}
	}

	url := "http://ftp.freebsd.org/pub/FreeBSD/releases/amd64/" + bs.Release + "/base.txz"
	bs.logger(LOGDBG, "Trying to download base.txz from "+url+" using fetch")
	_, err = bs.cmdOut("fetch", url, "-o", bs.Dirpath+"/base.txz")
	if err != nil {
		return errors.New("Error downloading base.txz from " + url)
	}

	bs.LastUpdated = bs.getCurrentDateTime()
	return nil
}

func (bs *Base) Import() error {
	bs.logger(LOGDBG, "Checking if "+bs.Dirpath+"/base.txz exists")
	_, err := os.Stat(bs.Dirpath + "/base.txz")
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("Base file " + bs.Dirpath + "/base.txz does not exist")
		} else {
			bs.logger(LOGDBG, "Error with checking "+bs.Dirpath+"/base.txz existance: "+err.Error())
		}
	}
	bs.logger(LOGDBG, "Base source exist so it can be imported")
	return nil
}

func (bs *Base) Remove() error {
	bs.logger(LOGDBG, "Checking if "+bs.Dirpath+" exists")
	_, err := os.Stat(bs.Dirpath)
	if err != nil {
		if os.IsNotExist(err) {
			bs.logger(LOGDBG, bs.Dirpath+"does not exist. Nothing to remove")
			return nil
		} else {
			bs.logger(LOGDBG, "Error with checking dir "+bs.Dirpath+" existance: "+err.Error())
			return err
		}
	}

	err = os.RemoveAll(bs.Dirpath)
	if err != nil {
		bs.logger(LOGERR, "Error removing dir "+bs.Dirpath+". Please remove the directory manually and remove the state")
		return errors.New("Error removing base dir")
	}
	bs.logger(LOGDBG, bs.Dirpath+"has been removed")

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
