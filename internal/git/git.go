package git

import (
	"os/exec"
)

func Available() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

func InitRepo(dir string) error {
	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	return cmd.Run()
}

func AddAll(dir string) error {
	cmd := exec.Command("git", "add", ".")
	cmd.Dir = dir
	return cmd.Run()
}

func FirstCommit(dir string) error {
	cmd := exec.Command("git", "commit", "-m", "Initial commit")
	cmd.Dir = dir
	return cmd.Run()
}

func Clone(url, dest string) error {
	cmd := exec.Command("git", "clone", "--depth", "1", url, dest)
	return cmd.Run()
}

