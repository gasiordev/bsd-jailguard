package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/scanner"
)

const CHAR_END_KEYVAL = ";"
const CHAR_BEGIN_VAL = "="
const CHAR_BEGIN_BLK = "{"
const CHAR_END_BLK = "}"

type JailConf struct {
	Name     string            `json:"name"`
	Config   map[string]string `json:"config"`
	Filepath string            `json:"filepath"`
	logger   func(int, string)
}

func (jc *JailConf) SetLogger(f func(int, string)) {
	jc.logger = f
}

func (jc *JailConf) isValidName(s string) bool {
	var r = regexp.MustCompile(`^[a-z]+[a-z0-9_-]*$`)
	return r.MatchString(s)
}

func (jc *JailConf) isValidKey(s string) bool {
	var r = regexp.MustCompile(`^[a-z]+[a-z0-9_]*$`)
	return r.MatchString(s)
}

func (jc *JailConf) isKeyValValid(k string, v string) error {
	return nil
}

func (jc *JailConf) getScanner(f string) (*scanner.Scanner, error) {
	c, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	var s scanner.Scanner
	s.Init(strings.NewReader(string(c)))
	return &s, nil
}

func (jc *JailConf) ParseFile(f string) error {
	jc.logger(LOGDBG, fmt.Sprintf("Opening %s to parse...", f))
	s, err := jc.getScanner(f)
	if err != nil {
		return err
	}

	k := ""
	v := ""
	gotKey := false
	prevToken := ""

	var i int
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		if i == 0 {
			if !jc.isValidName(s.TokenText()) {
				return errors.New("Invalid jail name")
			}
			jc.Name = s.TokenText()
		}
		if i == 1 && s.TokenText() != CHAR_BEGIN_BLK {
			return errors.New("Missing opening curl bracket")
		}
		if i > 1 {
			if !gotKey {
				if s.TokenText() == CHAR_END_BLK {
					if k == "" {
						break
					} else {
						return errors.New("Key not closed or its value missing")
					}
				}

				if !jc.isValidKey(s.TokenText()) && s.TokenText() != "." && s.TokenText() != CHAR_BEGIN_VAL && s.TokenText() != CHAR_END_KEYVAL {
					return errors.New("Invalid key (1)")
				}
				if k == "" && (s.TokenText() == "." || s.TokenText() == CHAR_BEGIN_VAL || s.TokenText() == CHAR_END_KEYVAL) {
					return errors.New("Invalid key (2)")
				}
				if prevToken == "." && (s.TokenText() == CHAR_BEGIN_VAL || s.TokenText() == CHAR_END_KEYVAL) {
					return errors.New("Invalid key (3)")
				}
				if s.TokenText() == CHAR_BEGIN_VAL {
					gotKey = true
				} else if s.TokenText() == CHAR_END_KEYVAL {
					gotKey = false
					v = "true"
					jc.Config[k] = v
					k = ""
					v = ""
				} else {
					k = k + s.TokenText()
				}
			} else {
				if s.TokenText() == CHAR_END_KEYVAL && v == "" {
					return errors.New("Invalid value")
				}
				if s.TokenText() == CHAR_END_KEYVAL {
					gotKey = false
					// TODO: Replacing double quote needs a better implementation
					v = strings.Trim(v, "\"")
					jc.Config[k] = v
					k = ""
					v = ""
				} else {
					v = v + s.TokenText()
				}
			}
		}
		prevToken = s.TokenText()
		i++
	}
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

func (jc *JailConf) WriteToFile(p string) error {
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

func NewJailConf() *JailConf {
	jc := &JailConf{}
	jc.Config = make(map[string]string)
	return jc
}
