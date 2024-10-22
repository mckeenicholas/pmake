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
	argsString := strings.Join(a.args, " ")
	cmd := exec.Command("sh", "-c", argsString)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("execution failed with %v, %s", err, output)
	}
	return nil
}
