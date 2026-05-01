package prompts

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"nexus-v/internal/tui"
	"nexus-v/pkg/nexusv"
)

type Answers struct {
	Name        string
	Identifier  string
	Description string
	Publisher   string
	Variant     string
	CommandName string
	License     string
}

func AskQuestions() (*Answers, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("🧱 Still Systems NEXUS-V")
	fmt.Println("-------------------------------------------")

	name := ask(reader, "Extension name", "my-extension")
	identifier := ask(reader, "Extension identifier", sanitizeIdentifier(name))
	description := ask(reader, "Description", "A helpful VS Code extension")
	publisher := ask(reader, "Publisher", "your-publisher-id")

	// Use TUI for variant selection
	vars, _ := nexusv.ListTemplates()
	variant, err := tui.SelectVariant(vars)
	if err != nil {
		variant = "command" // fallback
	}

	commandName := ask(reader, "Command name", identifier+".helloWorld")
	license := ask(reader, "License", "MIT")

	return &Answers{
		Name:        name,
		Identifier:  identifier,
		Description: description,
		Publisher:   publisher,
		Variant:     variant,
		CommandName: commandName,
		License:     license,
	}, nil
}

func ask(reader *bufio.Reader, label, defaultValue string) string {
	fmt.Printf("? %s (%s): ", label, defaultValue)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return defaultValue
	}
	return input
}

func sanitizeIdentifier(name string) string {
	id := strings.ToLower(name)
	var b strings.Builder
	for _, r := range id {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			b.WriteRune(r)
		} else if r == ' ' || r == '_' {
			b.WriteRune('-')
		}
	}
	// Remove leading/trailing dashes
	return strings.Trim(b.String(), "-")
}
