package nexusv

import (
	"os"
	"testing"
)

func TestRunHooks(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "nexusv-hooks-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test with simple portable command
	commands := []string{"go version"}
	if err := RunHooks(tempDir, commands); err != nil {
		t.Errorf("RunHooks failed: %v", err)
	}

	// Test with empty command list
	if err := RunHooks(tempDir, []string{}); err != nil {
		t.Errorf("RunHooks failed with empty list: %v", err)
	}

	// Test with failing command
	failingCommands := []string{"nonexistentcommand"}
	if err := RunHooks(tempDir, failingCommands); err == nil {
		t.Error("RunHooks should have failed for nonexistent command")
	}
}
