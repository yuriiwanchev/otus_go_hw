package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func buildTestCommand(t *testing.T) string {
	t.Helper()
	cmdName := "testcommand"
	if runtime.GOOS == "windows" {
		cmdName += ".exe"
	}

	cmdPath := filepath.Join(os.TempDir(), cmdName)
	cmd := exec.Command("go", "build", "-o", cmdPath, "./testdata/testcommand.go")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to build test command: %v", err)
	}
	return cmdPath
}

func TestRunCmd(t *testing.T) {
	scriptFile := buildTestCommand(t)

	tests := []struct {
		name       string
		cmd        []string
		env        Environment
		wantStatus int
	}{
		{
			name:       "successful command",
			cmd:        []string{scriptFile, "success"},
			env:        Environment{},
			wantStatus: 0,
		},
		{
			name:       "failing command",
			cmd:        []string{scriptFile, "fail"},
			env:        Environment{},
			wantStatus: 1,
		},
		{
			name:       "environment variable set",
			cmd:        []string{scriptFile, "success"},
			env:        Environment{"TEST_ENV": {Value: "value", NeedRemove: false}},
			wantStatus: 0,
		},
		{
			name:       "environment variable removed",
			cmd:        []string{scriptFile, "success"},
			env:        Environment{"TO_DELETE_ENV": {Value: "", NeedRemove: true}},
			wantStatus: 0,
		},
		{
			name:       "environment variable set to delete and set again",
			cmd:        []string{scriptFile, "success"},
			env:        Environment{"SWAP_ENV": {Value: "other_value", NeedRemove: false}},
			wantStatus: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("TO_DELETE_ENV", "value")
			os.Setenv("SWAP_ENV", "value")

			gotStatus := RunCmd(tt.cmd, tt.env)
			for key, env := range tt.env {
				envValue, exists := os.LookupEnv(key)
				if env.NeedRemove {
					assert.False(t, exists, "Environment variable %s should not exist", key)
					continue
				}
				assert.True(t, exists, "Environment variable %s should exist", key)
				assert.Equal(t, env.Value, envValue, "Environment variable %s should have value %s", key, env.Value)
			}

			assert.Equal(t, tt.wantStatus, gotStatus)
		})
	}
}
