package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <env_dir> <command> [args...]\n", os.Args[0])
		return
	}

	envDir := os.Args[1]
	cmdArgs := os.Args[2:]

	env, err := ReadDir(envDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading directory: %v\n", err)
		os.Exit(1)
	}

	code := RunCmd(cmdArgs, env)
	os.Exit(code)
}
