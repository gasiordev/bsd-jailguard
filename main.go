package main

import (
	"fmt"
	"os"
)

func main() {
	if os.Geteuid() != 0 {
		fmt.Fprintf(os.Stderr, "You have to run jailguard as root.\n")
		os.Exit(1)
	}

	j := NewJailguard()
	j.Run()
}
