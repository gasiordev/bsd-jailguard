package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"text/scanner"
	"regexp"
)

type JailConf struct {
	Name string
	Config map[string]string
}

func (jc *JailConf) ParseFile(f string) error {
	c, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}

	jc.Config = make(map[string]string)

	var s scanner.Scanner
	s.Init(strings.NewReader(string(c)))

	var validName = regexp.MustCompile(`^[a-z]+[a-z0-9_-]*$`)
	var validKey = regexp.MustCompile(`^[a-z]+[a-z0-9_]*$`)

	k := ""; v := ""; gotKey := false; prevToken := ""

	var i int
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		if i == 0 {
			if !validName.MatchString(s.TokenText()) {
				return errors.New("Invalid jail name")
			}
			jc.Name = s.TokenText()
		}
		if i == 1 && s.TokenText() != "{" {
			return errors.New("Missing opening curl bracket")
		}
		if i > 1 {
			if !gotKey {
				if s.TokenText() == "}" {
					if k == "" {
						break
					} else {
						return errors.New("Key not closed or its value missing")
					}
				}

				if !validKey.MatchString(s.TokenText()) && s.TokenText() != "." && s.TokenText() != "=" && s.TokenText() != ";" {
					return errors.New("Invalid key (1)")
				}
				if k == "" && (s.TokenText() == "." || s.TokenText() == "=" || s.TokenText() == ";") {
					return errors.New("Invalid key (2)")
				}
				if prevToken == "." && (s.TokenText() == "=" || s.TokenText() == ";") {
					return errors.New("Invalid key (3)")
				}
				if s.TokenText() == "=" {
					gotKey = true
				} else if s.TokenText() == ";" {
					gotKey = false
					v = "true"
					jc.Config[k] = v
					k = ""
					v = ""
				} else {
					k = k + s.TokenText()
				}
			} else {
				if s.TokenText() == ";" && v == "" {
					return errors.New("Invalid value")
				}
				if s.TokenText() == ";" {
					gotKey = false
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

	return nil
}

func NewJailConf() *JailConf {
	jc := &JailConf{}
	return jc
}
