package main

import (
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return 1
	}

	ex := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	envStrings := []string{}

	for k, v := range env {
		if v.NeedRemove {
			os.Unsetenv(k)

			continue
		}

		envStrings = append(envStrings,
			k+"="+v.Value,
		)
	}

	ex.Env = os.Environ()
	ex.Env = append(ex.Env, envStrings...)

	ex.Stdout = os.Stdout
	if err := ex.Run(); err != nil {
		log.Fatal(err)
	}

	return ex.ProcessState.ExitCode()
}
