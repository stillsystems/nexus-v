package cli

import (
	"fmt"
	"os"
	"strings"

	"nexus-v/pkg/nexusv"
)

func runGenerate(args []string) {
	if len(args) == 0 {
		fmt.Println("Error: No prompt provided.")
		fmt.Println("Usage: nexus-v generate <your prompt>")
		os.Exit(1)
	}

	prompt := strings.Join(args, " ")
	fmt.Printf("%s\U0001f916 Nexus AI Engine:%s Processing prompt...\n", Cyan, Reset)
	fmt.Printf("   > \"%s\"\n\n", prompt)

	meta, err := nexusv.GenerateFromPrompt(prompt, nil)
	if err != nil {
		fmt.Printf("%sError:%s %v\n", Red, Reset, err)
		os.Exit(1)
	}

	fmt.Printf("%s\U0001f3ab Blueprint Generated:%s\n", Green, Reset)
	fmt.Printf("   Name:        %s\n", meta.Name)
	fmt.Printf("   Description: %s\n", meta.Description)
	fmt.Printf("   Language:    %s\n", meta.Language)
	fmt.Printf("   Features:    ")
	if len(meta.Features) == 0 {
		fmt.Println("None")
	} else {
		var names []string
		for _, f := range meta.Features {
			names = append(names, f.Name)
		}
		fmt.Println(strings.Join(names, ", "))
	}

	fmt.Println("\n\U0001f680 Ready to scaffold. Run 'nexus-v init' to begin.")
}
