package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(command []string, env Environment) (returnCode int) {
	for name, val := range env {
		switch val.NeedRemove {
		case true:
			if err := os.Unsetenv(name); err != nil {
				fmt.Printf("error while removing env %s: %s", name, err)
				return 1
			}
		default:
		}

		if val.Value != "" {
			err := os.Setenv(name, val.Value)
			if err != nil {
				fmt.Printf("error while setting env %s: %s", name, err)
				return 1
			}
		}
	}

	//nolint: gosec
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("error while running comand %v: %s", command, err)
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}

		return 1
	}

	return 0
}
