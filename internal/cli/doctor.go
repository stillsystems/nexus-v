package cli

import (
	"fmt"
	"os"

	"github.com/stillsystems/nexus-v/pkg/nexusv"
)

func runDoctor(args []string) {
	fix := false
	for _, arg := range args {
		if arg == "--fix" {
			fix = true
		}
	}

	cwd, _ := os.Getwd()
	fmt.Printf("%s\U0001fa7a NEXUS-V Doctor (%s)%s\n", Bold, nexusv.GetSystemInfo(), Reset)
	fmt.Println("-------------------------------------------")

	result := nexusv.RunFullDoctor(cwd)

	fixedCount := 0
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
			if fix && issue.Fix != "" {
				fmt.Printf("   🛠️ Attempting auto-fix...")
				if err := nexusv.FixIssue(cwd, issue); err == nil {
					fmt.Printf("%s FIXED%s\n", Green, Reset)
					fixedCount++
				} else {
					fmt.Printf("%s FAILED: %v%s\n", Red, err, Reset)
				}
			} else if issue.Fix != "" {
				fmt.Printf("   💡 Fix: %s\n", issue.Fix)
			}
		}
	}

	fmt.Println("-------------------------------------------")
	if result.Healthy || (fix && fixedCount > 0) {
		if fixedCount > 0 {
			fmt.Printf("✨ Applied %d automated fixes. Environment is now stable.\n", fixedCount)
		} else {
			fmt.Println("✨ Your environment looks great! You're ready to scaffold.")
		}
	} else {
		fmt.Println("🚨 Critical issues found. Please address them to ensure stability.")
		if !fix {
			fmt.Println("   (Try 'nexus-v doctor --fix' for automated repair)")
		}
		os.Exit(1)
	}
}

