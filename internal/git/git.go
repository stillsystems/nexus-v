package git

import (
	"os/exec"
)

// Available checks if the git executable is present in the system's PATH.
func Available() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

// InitRepo initializes a new Git repository in the specified directory.
func InitRepo(dir string) error {
	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	return cmd.Run()
}

// AddAll stages all changes in the specified directory for the next commit.
func AddAll(dir string) error {
	cmd := exec.Command("git", "add", ".")
	cmd.Dir = dir
	return cmd.Run()
}

// FirstCommit creates the initial commit with the message "Initial commit".
func FirstCommit(dir string) error {
	cmd := exec.Command("git", "commit", "-m", "Initial commit")
	cmd.Dir = dir
	return cmd.Run()
}

// Clone clones a git repository from a URL to a destination path using default settings.
func Clone(url, dest string) error {
	return CloneWithRef(url, "", dest)
}

// CloneWithRef clones a git repository and checks out a specific reference (branch, tag, or SHA).
// It optimizes for speed by using shallow clones where possible.
func CloneWithRef(url, ref, dest string) error {
	args := []string{"clone"}

	// Use shallow clone for branches/tags, but full clone for SHAs.
	// We detect SHAs by checking for a length of 40 (SHA-1) or 64 (SHA-256).
	// This is a heuristic; technically a branch/tag could have these lengths,
	// but it is extremely rare in practice and a full clone remains safe.
	isSHA := len(ref) == 40 || len(ref) == 64
	if ref != "" && !isSHA {
		args = append(args, "--depth", "1", "--branch", ref)
	} else if ref == "" {
		args = append(args, "--depth", "1")
	}

	args = append(args, url, dest)
	cmd := exec.Command("git", args...)
	if err := cmd.Run(); err != nil {
		return err
	}

	// If it was a SHA, we need to manually checkout
	if isSHA {
		checkoutCmd := exec.Command("git", "checkout", ref)
		checkoutCmd.Dir = dest
		return checkoutCmd.Run()
	}

	return nil
}
