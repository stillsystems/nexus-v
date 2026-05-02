package nexusv

import (
	"fmt"
	"os"
	"path/filepath"
)

// LanguageSetup handles language-specific initialization (e.g., go mod init).
func LanguageSetup(dir string, ctx Context, meta *TemplateMetadata) ([]string, error) {
	lang := ""
	if meta != nil && meta.Language != "" {
		lang = meta.Language
	} else {
		lang = detectLanguage(dir)
	}

	var extraHooks []string

	switch lang {
	case "go":
		// If go.mod doesn't exist, queue initialization
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); os.IsNotExist(err) {
			moduleName := ctx.Identifier
			if ctx.Publisher != "" {
				moduleName = fmt.Sprintf("github.com/%s/%s", ctx.Publisher, ctx.Identifier)
			}
			extraHooks = append(extraHooks, fmt.Sprintf("go mod init %s", moduleName), "go mod tidy")
		}
	case "rust":
		// If Cargo.toml doesn't exist, queue initialization
		if _, err := os.Stat(filepath.Join(dir, "Cargo.toml")); os.IsNotExist(err) {
			extraHooks = append(extraHooks, "cargo init")
		}
	case "typescript", "node":
		// Standard Node project might need npm install, but that's already handled by flags
	}

	return extraHooks, nil
}

func detectLanguage(dir string) string {
	if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
		return "go"
	}
	if _, err := os.Stat(filepath.Join(dir, "Cargo.toml")); err == nil {
		return "rust"
	}
	if _, err := os.Stat(filepath.Join(dir, "package.json")); err == nil {
		return "typescript"
	}
	return ""
}
