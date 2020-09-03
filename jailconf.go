package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const CHAR_END_KEYVAL = ";"
const CHAR_BEGIN_VAL = "="
const CHAR_BEGIN_BLK = "{"
const CHAR_END_BLK = "}"

type JailConf struct {
	Name      string            `json:"name"`
	Config    map[string]string `json:"config"`
	Filepath  string            `json:"filepath"`
	Iteration int               `json:"iteration"`
	History   []*HistoryEntry   `json:"history"`

	logger func(int, string)
}

type JailConfJSON struct {
	Version string            `json:"version"`
	Jail    map[string]string `json:"jail"`
}

func (jc *JailConf) SetLogger(f func(int, string)) {
	jc.logger = f
}

func (jc *JailConf) AddHistoryEntry(s string) {
	he := NewHistoryEntry(GetCurrentDateTime(), s)
	if jc.History == nil {
		jc.History = []*HistoryEntry{}
	}
	jc.History = append(jc.History, he)
}

func (jc *JailConf) SetDefaultValues() {
	jc.Iteration = 1
}

func (jc *JailConf) isValidKey(s string) bool {
	var r = regexp.MustCompile(`^[a-z]+[a-z0-9_]*$`)
	return r.MatchString(s)
}

func (jc *JailConf) isKeyValValid(k string, v string) error {
	return nil
}

func (jc *JailConf) ParseFile(f string) error {
	jc.logger(LOGDBG, fmt.Sprintf("Opening %s to parse...", f))
	c, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}

	v := &JailConfJSON{}
	err = json.Unmarshal(c, &v)
	if err != nil {
		return errors.New(fmt.Sprintf("Error has occurred while unmarshaling file: %s", err.Error()))
	}

	if v.Jail["name"] == "" {
		return errors.New("Jail name is missing from the jail file")
	}
	jc.Name = v.Jail["name"]
	jc.Config = v.Jail

	jc.logger(LOGDBG, fmt.Sprintf("File %s has been successfully parsed", f))
	return nil
}

func (jc *JailConf) Validate() error {
	jc.logger(LOGDBG, "Checking if key-value pairs in config are valid...")
	for k, v := range jc.Config {
		err := jc.isKeyValValid(k, v)
		if err != nil {
			return err
		}
	}
	jc.logger(LOGDBG, "Checking for required values in config...")

	// If path is not empty then exec.start and exec.stop are necessary
	if jc.Config["path"] != "" {
		jc.logger(LOGDBG, "Checking if 'path' leads to an existing directory...")
		stat, _, err := StatWithLog(jc.Config["path"], jc.logger)
		if err != nil {
			if os.IsNotExist(err) {
				return errors.New("'path' is invalid")
			} else {
				return errors.New("Error has occurred when checking for 'path' existance")
			}
		} else if !stat.IsDir() {
			return errors.New("'path' is not a directory")
		}

		for _, k := range []string{"exec.start", "exec.stop"} {
			if jc.Config[k] == "" {
				return errors.New(fmt.Sprintf("'%s' is missing!", k))
			}
		}
	}

	return nil
}

func (jc *JailConf) Write(p string) error {
	jc.Filepath = p
	d := filepath.Dir(jc.Filepath)

	stat, _, err := StatWithLog(d, jc.logger)
	if err != nil {
		if os.IsNotExist(err) {
			jc.logger(LOGDBG, fmt.Sprintf("%s does not exist and it has to be created", d))
			err2 := CreateDirWithLog(d, jc.logger)
			if err2 != nil {
				return err2
			}
		} else {
			jc.logger(LOGDBG, "Error has occurred when writing config to a file")
			return err
		}
	} else if !stat.IsDir() {
		return errors.New("Error has occurred when writing config to a file")
	}

	jc.Iteration++
	o := jc.Name + " {\n"
	for k, v := range jc.Config {
		if v == "true" {
			o = o + fmt.Sprintf("  %s;\n", k)
		} else {
			if strings.Contains(v, " ") {
				o = o + fmt.Sprintf("  %s = \"%s\";\n", k, v)
			} else {
				o = o + fmt.Sprintf("  %s = %s;\n", k, v)
			}
		}
	}
	o = o + "}\n"

	jc.logger(LOGDBG, fmt.Sprintf("Writing jail config..."))
	err = ioutil.WriteFile(jc.Filepath, []byte(o), 0644)
	if err != nil {
		return err
	}
	jc.logger(LOGDBG, "Config has been written to a file")

	return nil
}

func (jc *JailConf) Remove() error {
	if jc.Filepath == "" {
		return nil
	}

	_, _, err := StatWithLog(jc.Filepath, jc.logger)
	if err != nil {
		if os.IsNotExist(err) {
			jc.logger(LOGDBG, "Nothing to remove")
			return nil
		} else {
			return err
		}
	}

	err = RemoveAllWithLog(jc.Filepath, jc.logger)
	if err != nil {
		return errors.New("Error has occurred while removing jail config. Please remove the file manually")
	}

	return nil
}

func NewJailConf() *JailConf {
	jc := &JailConf{}
	jc.Config = make(map[string]string)
	return jc
}
