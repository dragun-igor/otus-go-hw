package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	for key, val := range env {
		os.Unsetenv(key)
		if !val.NeedRemove {
			os.Setenv(key, val.Value)
		}
	}

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	var ee *exec.ExitError
	if err := command.Run(); err != nil {
		if errors.As(err, &ee) {
			returnCode = ee.ExitCode()
			log.Println("exit code error:", returnCode)
		}
	}
	return
}
