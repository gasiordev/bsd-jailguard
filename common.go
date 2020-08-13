package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func StatWithLog(p string, fn func(int, string)) (os.FileInfo, bool, error) {
	fn(LOGDBG, fmt.Sprintf("Getting stat for path %s...", p))
	st, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			fn(LOGDBG, fmt.Sprintf("Path %s does not exist", p))
		} else {
			fn(LOGDBG, fmt.Sprintf("Error has occurred when getting stat for path %s: %s", p, err.Error()))
		}
		return st, false, err
	} else {
		fn(LOGDBG, fmt.Sprintf("Found path %s", p))
		if st.IsDir() {
			fn(LOGDBG, fmt.Sprintf("Path %s is a directory", p))
			return st, true, nil
		}
	}
	return st, false, nil
}

func RemoveAllWithLog(p string, fn func(int, string)) error {
	fn(LOGDBG, fmt.Sprintf("Removing %s...", p))
	err := os.RemoveAll(p)
	if err != nil {
		fn(LOGDBG, fmt.Sprintf("Error has occurred when removing %s: %s", p, err.Error()))
	}
	fn(LOGDBG, fmt.Sprintf("Path %s has been removed", p))
	return err
}

func CreateDirWithLog(p string, fn func(int, string)) error {
	fn(LOGDBG, fmt.Sprintf("Creating directory %s...", p))
	err := os.MkdirAll(p, os.ModePerm)
	if err != nil {
		fn(LOGDBG, fmt.Sprintf("Error has occurred when creating %s: %s", p, err.Error()))
	}
	fn(LOGDBG, fmt.Sprintf("Directory %s has been created", p))
	return err
}

func CmdFetchWithLog(url string, o string, fn func(int, string)) error {
	fn(LOGDBG, fmt.Sprintf("Running 'fetch' to download %s to %s...", url, o))
	_, err := CmdOut(fn, "fetch", url, "-o", o)
	if err != nil {
		fn(LOGDBG, fmt.Sprintf("Error has occurred when downloading %s to %s", url, o))
		return err
	}
	fn(LOGDBG, fmt.Sprintf("File %s has been successfully saved in %s", url, o))
	return nil
}

func CmdTarExtractWithLog(f string, d string, fn func(int, string)) error {
	fn(LOGDBG, fmt.Sprintf("Running 'tar' to extract %s to %s directory...", f, d))
	_, err := CmdOut(fn, "tar", "-xvf", f, "-C", d)
	if err != nil {
		fn(LOGDBG, fmt.Sprintf("Error has occurred when extracting %s to %s", f, d))
		return err
	}
	fn(LOGDBG, fmt.Sprintf("File %s has been successfully extracted to %s", f, d))
	return nil
}

func CmdOut(fn func(int, string), c string, a ...string) ([]byte, error) {
	fn(LOGDBG, fmt.Sprintf("Running command '%s %s'...", c, strings.Join(a, "")))
	cmd := exec.Command(c, a...)
	cmd.Stdin = os.Stdin
	return cmd.Output()
}

func CmdRun(fn func(int, string), c string, a ...string) error {
	fn(LOGDBG, fmt.Sprintf("Running command '%s %s'...", c, strings.Join(a, "")))
	cmd := exec.Command(c, a...)
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func JailExistsInOSWithLog(n string, fn func(int, string)) (bool, error) {
	fn(LOGDBG, fmt.Sprintf("Running 'jls' to check if jail %s is running...", n))
	out, err := CmdOut(fn, "jls", "-Nn")
	if err != nil {
		return false, errors.New("Error running 'jls' to check if jail is running: " + err.Error())
	}

	re := regexp.MustCompile("name=" + n + " ")
	if re.Match([]byte(string(out))) {
		fn(LOGDBG, fmt.Sprintf("Jail %s is running (it was found in 'jls' output)", n))
		return true, nil
	}

	fn(LOGDBG, fmt.Sprintf("Jail %s does not seem to be running", n))
	return false, nil
}

func IsValidJailName(n string) bool {
	// TODO: 'jail' man page doesn't say too much about name restrictions.
	// Needs checking elsewhere
	re := regexp.MustCompile(`^[a-z][a-z0-9_\-]{1,31}$`)
	if !re.Match([]byte(n)) {
		return false
	}
	if strings.Contains(n, "--") {
		return false
	}
	return true
}

func IsValidIPAddress(ip string) bool {
	re := regexp.MustCompile(`^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$`)
	if !re.Match([]byte(ip)) {
		return false
	}

	ip_b := strings.Split(ip, ".")
	for i := 0; i < 4; i++ {
		v, _ := strconv.Atoi(ip_b[i])
		if (i == 0 && (v < 1 || v > 255)) || (i > 0 && (v < 0 || v > 255)) {
			return false
		}
	}
	return true
}

func GetCurrentDateTime() string {
	return time.Now().String()
}
