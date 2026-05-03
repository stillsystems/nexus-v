package nexusv

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidatePackageJSON(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "nexusv-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	t.Run("Valid PackageJSON", func(t *testing.T) {
		path := filepath.Join(tmpDir, "valid-package.json")
		content := `{
			"name": "test-pkg",
			"publisher": "stillsystems",
			"version": "1.0.0",
			"engines": { "vscode": "^1.90.0" },
			"categories": ["Programming Languages"],
			"contributes": {}
		}`
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		res, err := validatePackageJSON(path)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if res != nil {
			t.Errorf("expected nil result for valid file, got %v", res)
		}
	})

	t.Run("Missing Required Fields", func(t *testing.T) {
		path := filepath.Join(tmpDir, "missing-package.json")
		content := `{ "name": "test-pkg" }`
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		res, err := validatePackageJSON(path)
		if err != nil {
			t.Fatal(err)
		}
		if res == nil {
			t.Fatal("expected validation result for missing fields, got nil")
		}
		if len(res.Errors) < 5 {
			t.Errorf("expected at least 5 errors, got %d", len(res.Errors))
		}
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		path := filepath.Join(tmpDir, "invalid-package.json")
		content := `{ "name": "test-pkg", }` // Trailing comma
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		res, err := validatePackageJSON(path)
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Errors) == 0 {
			t.Error("expected error for invalid JSON")
		}
	})
}

func TestValidateTSConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "nexusv-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	t.Run("Valid TSConfig", func(t *testing.T) {
		path := filepath.Join(tmpDir, "valid-tsconfig.json")
		content := `{
			"compilerOptions": {
				"strict": true,
				"target": "ES2022"
			}
		}`
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		res, err := validateTSConfig(path)
		if err != nil {
			t.Fatal(err)
		}
		if res != nil {
			t.Errorf("expected nil result, got %v", res)
		}
	})

	t.Run("Non-strict Warning", func(t *testing.T) {
		path := filepath.Join(tmpDir, "loose-tsconfig.json")
		content := `{
			"compilerOptions": {
				"strict": false
			}
		}`
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		res, err := validateTSConfig(path)
		if err != nil {
			t.Fatal(err)
		}
		if len(res.Warnings) == 0 {
			t.Error("expected warning for non-strict mode")
		}
	})
}
