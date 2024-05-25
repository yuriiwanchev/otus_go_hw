package main

import (
	"fmt"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading directory: %v\n", err)
		return nil, err
	}

	env := make(Environment)

	for _, file := range files {
		if file.IsDir() || strings.HasPrefix(file.Name(), ".") {
			continue
		}

		fileName := file.Name()
		filePath := fmt.Sprintf("%s/%s", dir, fileName)

		data, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", fileName, err)
			return nil, err
		}

		env[fileName] = ProcessData(data)
	}

	return env, nil
}

func ProcessData(data []byte) EnvValue {
	if len(data) == 0 {
		return EnvValue{Value: "", NeedRemove: true}
	}

	lines := strings.Split(string(data), "\n")
	if len(lines) == 0 {
		return EnvValue{Value: "", NeedRemove: true}
	}

	value := lines[0]
	value = strings.TrimRight(value, " \t")
	value = strings.ReplaceAll(value, "\x00", "\n")
	return EnvValue{Value: value, NeedRemove: false}
}
