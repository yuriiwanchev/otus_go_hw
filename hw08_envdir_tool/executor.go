package main

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for key, value := range env {
		_, exists := os.LookupEnv(key)

		if exists {
			err := os.Unsetenv(key)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error unsetting environment variable: %v\n", err)
			}
		}

		if value.NeedRemove {
			continue
		}

		os.Setenv(key, value.Value)
	}

	command := cmd[0]
	args := cmd[1:]

	cmdExec := exec.Command(command, args...)
	cmdExec.Stdin = os.Stdin
	cmdExec.Stdout = os.Stdout
	cmdExec.Stderr = os.Stderr

	// slog.Info("running command", "cmd", cmdExec.String())

	err := cmdExec.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running command: %v\n", err)
	}

	returnCode = cmdExec.ProcessState.ExitCode()
	return returnCode
}
