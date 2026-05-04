package deps

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

type PackageManager string

const (
	Npm  PackageManager = "npm"
	Pnpm PackageManager = "pnpm"
	Yarn PackageManager = "yarn"
)

func DetectPackageManager() (PackageManager, error) {
	if exists("pnpm") {
		return Pnpm, nil
	}
	if exists("yarn") {
		return Yarn, nil
	}
	if exists("npm") {
		return Npm, nil
	}
	return "", errors.New("no supported package manager found (npm, pnpm, yarn)")
}

func exists(bin string) bool {
	_, err := exec.LookPath(bin)
	return err == nil
} // ← THIS BRACE WAS MISSING

func Install(pm PackageManager, cwd string) error {
	var cmd *exec.Cmd

	switch pm {
	case Pnpm:
		cmd = exec.Command("pnpm", "install")
	case Yarn:
		cmd = exec.Command("yarn", "install")
	case Npm:
		cmd = exec.Command("npm", "install")
	default:
		return fmt.Errorf("unsupported package manager: %s", pm)
	}

	cmd.Dir = cwd
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Installing dependencies using", pm)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("dependency installation failed: %w", err)
	}

	return nil
}

