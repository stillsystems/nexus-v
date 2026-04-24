package templates

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateProject(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "nexus-v-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	ctx := Context{
		Name:        "Test Extension",
		Identifier:  "test-extension",
		Description: "A test description",
		Publisher:   "test-publisher",
		CommandName: "test-extension.hello",
		Template:    "command",
		Force:       false,
		DryRun:      false,
	}

	target := filepath.Join(tempDir, "output")
	if err := GenerateProject(ctx, target); err != nil {
		t.Fatalf("GenerateProject failed: %v", err)
	}

	// Verify package.json
	pkgPath := filepath.Join(target, "package.json")
	data, err := os.ReadFile(pkgPath)
	if err != nil {
		t.Fatalf("failed to read package.json: %v", err)
	}

	var pkg map[string]interface{}
	if err := json.Unmarshal(data, &pkg); err != nil {
		t.Fatalf("failed to unmarshal package.json: %v", err)
	}

	if pkg["name"] != "test-extension" {
		t.Errorf("expected name 'test-extension', got %v", pkg["name"])
	}
	if pkg["publisher"] != "test-publisher" {
		t.Errorf("expected publisher 'test-publisher', got %v", pkg["publisher"])
	}

	// Verify src/extension.ts
	extPath := filepath.Join(target, "src", "extension.ts")
	if _, err := os.Stat(extPath); os.IsNotExist(err) {
		t.Errorf("src/extension.ts was not created")
	}
}

func TestGenerateVariants(t *testing.T) {
	variants := []string{"command", "webview", "language", "theme"}

	for _, v := range variants {
		t.Run(v, func(t *testing.T) {
			tempDir, _ := os.MkdirTemp("", "nexus-v-variant-test-*")
			defer os.RemoveAll(tempDir)

			ctx := Context{
				Name:       "Test " + v,
				Identifier: "test-" + v,
				Publisher:  "tester",
				Template:   v,
				DryRun:     false,
			}

			target := filepath.Join(tempDir, "output")
			if err := GenerateProject(ctx, target); err != nil {
				t.Fatalf("GenerateProject failed for variant %s: %v", v, err)
			}

			// Check for some variant specific files
			switch v {
			case "language":
				langConfig := filepath.Join(target, "language-configuration.json")
				if _, err := os.Stat(langConfig); os.IsNotExist(err) {
					t.Errorf("language-configuration.json missing for variant language")
				}
				grammar := filepath.Join(target, "syntaxes", "test-language.tmLanguage.json")
				if _, err := os.Stat(grammar); os.IsNotExist(err) {
					t.Errorf("grammar file missing for variant language")
				}
			case "theme":
				themeFile := filepath.Join(target, "themes", "test-theme-color-theme.json")
				if _, err := os.Stat(themeFile); os.IsNotExist(err) {
					t.Errorf("theme file missing for variant theme")
				}
			}
		})
	}
}
