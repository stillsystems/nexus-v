package nexusv

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRenderTemplate(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "nexusv-render-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	ctx := Context{
		Name:        "Test Project",
		Identifier:  "test-project",
		Description: "A test description",
		Publisher:   "stillsystems",
	}

	t.Run("Basic Variable Injection", func(t *testing.T) {
		data := []byte("Project: {{.Name}}")
		outPath := filepath.Join(tmpDir, "basic.txt")
		if err := renderTemplate(data, "test", outPath, ctx); err != nil {
			t.Fatal(err)
		}

		content, err := os.ReadFile(outPath)
		if err != nil {
			t.Fatalf("failed to read test output: %v", err)
		}
		if string(content) != "Project: Test Project" {
			t.Errorf("expected 'Project: Test Project', got %q", string(content))
		}
	})

	t.Run("Custom Function: currentYear", func(t *testing.T) {
		data := []byte("Year: {{currentYear}}")
		outPath := filepath.Join(tmpDir, "year.txt")
		if err := renderTemplate(data, "test", outPath, ctx); err != nil {
			t.Fatal(err)
		}

		content, err := os.ReadFile(outPath)
		if err != nil {
			t.Fatalf("failed to read test output: %v", err)
		}
		expectedYear := fmt.Sprintf("%d", time.Now().Year())
		if !strings.Contains(string(content), expectedYear) {
			t.Errorf("expected current year %s, got %q", expectedYear, string(content))
		}
	})

	t.Run("Custom Function: licenseText", func(t *testing.T) {
		data := []byte("License: {{licenseText \"MIT\"}}")
		outPath := filepath.Join(tmpDir, "license.txt")
		if err := renderTemplate(data, "test", outPath, ctx); err != nil {
			t.Fatal(err)
		}

		content, err := os.ReadFile(outPath)
		if err != nil {
			t.Fatalf("failed to read test output: %v", err)
		}
		if !strings.Contains(string(content), "Permission is hereby granted") {
			t.Error("expected MIT License text to be present")
		}
	})

	t.Run("Security: Path Traversal Protection", func(t *testing.T) {
		data := []byte("test")
		outPath := filepath.Join(tmpDir, "safe.txt")
		if err := renderTemplate(data, "test", outPath, ctx); err != nil {
			t.Errorf("expected success for safe path, got %v", err)
		}
	})

	t.Run("Dry Run Mode", func(t *testing.T) {
		ctxDry := ctx
		ctxDry.DryRun = true
		data := []byte("Should not be written")
		outPath := filepath.Join(tmpDir, "dryrun.txt")
		if err := renderTemplate(data, "test", outPath, ctxDry); err != nil {
			t.Fatal(err)
		}

		if _, err := os.Stat(outPath); !os.IsNotExist(err) {
			t.Error("file should not exist in dry run mode")
		}
	})
}
