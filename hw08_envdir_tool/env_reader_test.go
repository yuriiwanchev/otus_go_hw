package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadDir(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	normalTests := []struct {
		filename    string
		content     string
		expectedVal string
		needRemove  bool
	}{
		{"NORMAL", "value", "value", false},
		{"EMPTY", "", "", true},
		{"EMPTYFIRSTLINE", "\nvalue", "", false},
		{"NEWLINE", "value1\nvalue2", "value1", false},
		{"NULLS", "value1\x00value2", "value1\nvalue2", false},
		{"TRIM", "value  \t\n", "value", false},
	}

	testsForNoIncludeFiles := []string{"=HIDDEN", "HID=DEN", ".HIDDEN"}

	for _, tt := range normalTests {
		err := os.WriteFile(filepath.Join(tempDir, tt.filename), []byte(tt.content), 0o644)
		if err != nil {
			t.Fatalf("Failed to write file: %v", err)
		}
	}
	for _, tt := range testsForNoIncludeFiles {
		err := os.WriteFile(filepath.Join(tempDir, tt), []byte(tt), 0o644)
		if err != nil {
			t.Fatalf("Failed to write file: %v", err)
		}
	}

	// Act
	env, err := ReadDir(tempDir)
	if err != nil {
		t.Fatalf("ReadDir returned an error: %v", err)
	}

	// Assert
	for _, tt := range normalTests {
		val, exists := env[tt.filename]

		assert.True(t, exists, "Expected entry for file %s", tt.filename)
		assert.Equal(t, tt.expectedVal, val.Value,
			"Expected value %s for file %s, got %s", tt.expectedVal, tt.filename, val.Value)
		assert.Equal(t, tt.needRemove, val.NeedRemove,
			"Expected value %s for file %s, got %s", tt.needRemove, tt.filename, val.NeedRemove)
	}

	for _, tt := range testsForNoIncludeFiles {
		_, exists := env[tt]
		assert.False(t, exists, "Expected no entry for file %s", tt)
	}
}
