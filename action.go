package main

import (
	"fmt"
	"os/exec"
	"strings"
)

type Action struct {
	args []string
}

func (a *Action) Execute() error {
	if len(a.args) == 0 {
		return fmt.Errorf("no action to execute")
	}
	argsString := strings.Join(a.args, " ")
	cmd := exec.Command("sh", "-c")
	cmd.Args = append(cmd.Args, argsString)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("execution failed: %v, output: %s", err, output)
	}
	return nil
}
