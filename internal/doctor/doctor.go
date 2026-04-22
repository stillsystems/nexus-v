package doctor

import (
	"fmt"
	"os/exec"
	"runtime"
)

type Check struct {
	Name        string
	Command     string
	Args        []string
	Required    bool
	Description string
}

func RunChecks() {
	checks := []Check{
		{Name: "Node.js", Command: "node", Args: []string{"--version"}, Required: true, Description: "Required to run and develop VS Code extensions."},
		{Name: "npm", Command: "npm", Args: []string{"--version"}, Required: true, Description: "Standard package manager for extension dependencies."},
		{Name: "vsce", Command: "vsce", Args: []string{"--version"}, Required: false, Description: "Official tool for packaging and publishing (install with 'npm install -g @vscode/vsce')."},
		{Name: "git", Command: "git", Args: []string{"--version"}, Required: true, Description: "Used for version control and scaffolding initialization."},
	}

	fmt.Printf("NEXUS-V Doctor (%s/%s)\n", runtime.GOOS, runtime.GOARCH)
	fmt.Println("-------------------------------------------")

	allPassed := true
	for _, c := range checks {
		cmd := exec.Command(c.Command, c.Args...)
		err := cmd.Run()

		status := "✔"
		if err != nil {
			status = "✘"
			if c.Required {
				allPassed = false
			}
		}

		fmt.Printf("%s %-8s ", status, c.Name)
		if err != nil {
			fmt.Printf("(NOT FOUND) - %s\n", c.Description)
		} else {
			fmt.Println("(OK)")
		}
	}

	fmt.Println("-------------------------------------------")
	if allPassed {
		fmt.Println("✨ Your environment looks great! You're ready to scaffold.")
	} else {
		fmt.Println("⚠️  Some required tools are missing. Please install them to ensure NEXUS-V works correctly.")
	}
}
