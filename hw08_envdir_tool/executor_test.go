package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunCmd(t *testing.T) {
	// Create a temporary script file for testing
	script := `#!/bin/sh
if [ "$1" = "success" ]; then
  exit 0
else
  exit 1
fi
`
	scriptFile := "/testscript.sh"
	err := os.WriteFile(scriptFile, []byte(script), 0755)
	if err != nil {
		t.Fatalf("Failed to create test script: %v", err)
	}

	defer os.Remove(scriptFile)

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
			for key, value := range tt.env {
				envValue, exists := os.LookupEnv(key)
				assert.True(t, exists)
				assert.Equal(t, value, envValue)
			}

			assert.Equal(t, tt.wantStatus, gotStatus)
		})
	}
}
