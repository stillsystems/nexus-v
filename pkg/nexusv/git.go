package nexusv

import (
	"os/exec"
)

// GitAvailable checks if the git executable is present in the system's PATH.
func GitAvailable() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

// GitInitRepo initializes a new Git repository in the specified directory.
func GitInitRepo(dir string) error {
	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	return cmd.Run()
}

// GitAddAll stages all changes in the specified directory for the next commit.
func GitAddAll(dir string) error {
	cmd := exec.Command("git", "add", ".")
	cmd.Dir = dir
	return cmd.Run()
}

// GitFirstCommit creates the initial commit with the message "Initial commit".
func GitFirstCommit(dir string) error {
	cmd := exec.Command("git", "commit", "-m", "Initial commit")
	cmd.Dir = dir
	return cmd.Run()
}

// GitClone clones a git repository from a URL to a destination path using default settings.
func GitClone(url, dest string) error {
	return GitCloneWithRef(url, "", dest)
}

// GitCloneWithRef clones a git repository and checks out a specific reference (branch, tag, or SHA).
// It optimizes for speed by using shallow clones where possible.
func GitCloneWithRef(url, ref, dest string) error {
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

