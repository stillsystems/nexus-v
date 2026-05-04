package nexusv

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func RunHooks(dir string, commands []string) error {
	for _, cmdStr := range commands {
		parts := strings.Fields(cmdStr)
		if len(parts) == 0 {
			continue
		}

		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Dir = dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		fmt.Println("→ Running:", cmdStr)

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("hook failed (%s): %w", cmdStr, err)
		}
	}

	return nil
}

