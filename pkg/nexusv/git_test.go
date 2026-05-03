package nexusv

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestGitAvailable(t *testing.T) {
	// This assumes git is installed in the test environment
	available := GitAvailable()
	if !available {
		t.Log("Git not available, skipping related tests")
	}
}

func TestGitWorkflow(t *testing.T) {
	if !GitAvailable() {
		t.Skip("Git not available")
	}

	tempDir, err := os.MkdirTemp("", "nexusv-git-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test GitInitRepo
	if err := GitInitRepo(tempDir); err != nil {
		t.Errorf("GitInitRepo failed: %v", err)
	}

	// Verify .git exists
	if _, err := os.Stat(filepath.Join(tempDir, ".git")); os.IsNotExist(err) {
		t.Error(".git directory was not created")
	}

	// Test GitAddAll
	dummyFile := filepath.Join(tempDir, "dummy.txt")
	if err := os.WriteFile(dummyFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to write dummy file: %v", err)
	}

	if err := GitAddAll(tempDir); err != nil {
		t.Errorf("GitAddAll failed: %v", err)
	}

	// Test GitFirstCommit
	// Note: This might fail if git user.name/email are not set.
	// We'll set them for the test if possible.
	setupGitUser(tempDir)

	if err := GitFirstCommit(tempDir); err != nil {
		t.Errorf("GitFirstCommit failed: %v", err)
	}
}

func setupGitUser(dir string) {
	// Set local config for the test repo
	cmd1 := execCommand("git", "config", "user.email", "test@example.com")
	cmd1.Dir = dir
	cmd1.Run()
	cmd2 := execCommand("git", "config", "user.name", "Test User")
	cmd2.Dir = dir
	cmd2.Run()
}

// execCommand is a helper to allow easier testing if we wanted to mock,
// but here we just use it to wrap the os/exec call for setupGitUser.
func execCommand(name string, arg ...string) *exec.Cmd {
	return exec.Command(name, arg...)
}
