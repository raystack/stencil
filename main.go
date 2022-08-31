package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/odpf/salt/cmdx"
	"github.com/odpf/stencil/cmd"
)

const (
	exitOK    = 0
	exitError = 1
)

func main() {
	root := cmd.New()

	command, err := root.ExecuteC()
	if err == nil {
		return
	}

	if cmdx.IsCmdErr(err) {
		if !strings.HasSuffix(err.Error(), "\n") {
			fmt.Println()
		}
		fmt.Println(command.UsageString())
		os.Exit(exitOK)
	}

	fmt.Println(err)
	os.Exit(exitError)
}
