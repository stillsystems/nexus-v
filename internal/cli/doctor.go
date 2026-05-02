package cli

import (
	"fmt"
	"os"

	"nexus-v/pkg/nexusv"
)

func runDoctor(_ []string) {
	cwd, _ := os.Getwd()
	fmt.Printf("%s🩺 NEXUS-V Doctor (%s)%s\n", Bold, nexusv.GetSystemInfo(), Reset)
	fmt.Println("-------------------------------------------")

	result := nexusv.RunFullDoctor(cwd)

	for _, issue := range result.Issues {
		icon := "✔"
		switch issue.Level {
		case "Error":
			icon = "✘"
		case "Warning":
			icon = "⚠️ "
		}

		fmt.Printf("%s %-12s ", icon, issue.Name)
		if issue.Level == "Success" {
			fmt.Println("(OK)")
		} else {
			fmt.Printf("(%s) - %s\n", issue.Level, issue.Message)
			if issue.Fix != "" {
				fmt.Printf("   💡 Fix: %s\n", issue.Fix)
			}
		}
	}

	fmt.Println("-------------------------------------------")
	if result.Healthy {
		fmt.Println("✨ Your environment looks great! You're ready to scaffold.")
	} else {
		fmt.Println("🚨 Critical issues found. Please address them to ensure stability.")
		os.Exit(1)
	}
}
