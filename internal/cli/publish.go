package cli

import (
	"fmt"
	"os"

	"nexus-v/pkg/nexusv"
)

func runPublish(args []string) {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	fmt.Printf("%s\U0001f4e4 NEXUS-V Publish%s\n", Bold, Reset)
	fmt.Println("-------------------------------------------")

	// 1. Validate template before publishing
	fmt.Println("🔍 Auditing template structure...")
	results := nexusv.ValidateTemplate(dir)
	hasErrors := false
	for _, res := range results {
		if len(res.Errors) > 0 {
			hasErrors = true
			for _, err := range res.Errors {
				fmt.Printf("%s✘ Error:%s %s (%s)\n", Red, Reset, err, res.File)
			}
		}
	}

	if hasErrors {
		fmt.Printf("\n%sCannot publish template with validation errors.%s\n", Red, Reset)
		os.Exit(1)
	}

	// 2. Load metadata
	meta, err := nexusv.LoadTemplateMetadata(dir)
	if err != nil {
		fmt.Printf("%sError:%s Failed to read metadata: %v\n", Red, Reset, err)
		os.Exit(1)
	}

	fmt.Printf("✅ Validation passed for: %s (%s)\n", meta.Name, meta.Language)

	// 3. Publishing instructions (Foundation)
	// In a real scenario, this would authenticate and push to a registry API or GitHub.
	fmt.Printf("\n%sSubmission Instructions:%s\n", Cyan, Reset)
	fmt.Println("   To finalize your submission to the Nexus Hub:")
	fmt.Printf("   1. Ensure your template is pushed to a public GitHub repository.\n")
	fmt.Printf("   2. Visit: https://github.com/stillsystems/nexus-registry/issues/new\n")
	fmt.Printf("   3. Use the title: [SUBMISSION] %s\n", meta.Name)
	fmt.Printf("   4. Include the repository URL in the description.\n")

	fmt.Printf("\n%sAlternatively, if you have configured a PAT, run:%s\n", Bold, Reset)
	fmt.Println("   nexus-v publish --auto-pr")
	
	fmt.Println("\n-------------------------------------------")
	fmt.Println("🚀 Thank you for contributing to the Still Systems ecosystem!")
}
