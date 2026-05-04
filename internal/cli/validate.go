package cli

import (
	"fmt"
	"os"

	"github.com/stillsystems/nexus-v/pkg/nexusv"
)

func runValidate(args []string) {
	dir, _ := os.Getwd()
	if len(args) > 0 {
		dir = args[0]
	}

	fmt.Printf("%s🔍 NEXUS-V Linter: Validating %s...%s\n", Bold, dir, Reset)
	fmt.Println("-------------------------------------------")

	results := nexusv.ValidateTemplate(dir)
	
	if len(results) == 0 {
		fmt.Println("✨ Template structure is valid!")
		return
	}

	hasErrors := false
	for _, res := range results {
		fmt.Printf("File: %s\n", res.File)
		for _, err := range res.Errors {
			fmt.Printf("  %s✘ Error: %s%s\n", Red, err, Reset)
			hasErrors = true
		}
		for _, warn := range res.Warnings {
			fmt.Printf("  %s⚠️ Warning: %s%s\n", Yellow, warn, Reset)
		}
	}

	fmt.Println("-------------------------------------------")
	if hasErrors {
		fmt.Println("🚨 Validation failed. Please address the errors above.")
		os.Exit(1)
	} else {
		fmt.Println("✅ Validation passed with warnings.")
	}
}

