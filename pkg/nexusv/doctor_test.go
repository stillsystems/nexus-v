package nexusv

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunFullDoctor(t *testing.T) {
	// Test environment checks (at least some tools should be found in a dev env)
	result := RunFullDoctor("")
	if len(result.Issues) == 0 {
		t.Error("expected at least some tool checks in doctor result")
	}

	// Test project checks
	tempDir, err := os.MkdirTemp("", "nexusv-doctor-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Case 1: Invalid package.json
	pkgPath := filepath.Join(tempDir, "package.json")
	if err := os.WriteFile(pkgPath, []byte("{invalid json"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	result = RunFullDoctor(tempDir)
	foundManifestError := false
	for _, issue := range result.Issues {
		if issue.Name == "Manifest" && issue.Level == "Error" {
			foundManifestError = true
			break
		}
	}
	if !foundManifestError {
		t.Error("expected manifest error for invalid JSON")
	}

	// Case 2: Missing vital fields
	if err := os.WriteFile(pkgPath, []byte(`{"name": "test"}`), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	result = RunFullDoctor(tempDir)
	foundWarning := false
	for _, issue := range result.Issues {
		if issue.Name == "Manifest" && issue.Level == "Warning" {
			foundWarning = true
		}
	}
	if !foundWarning {
		t.Error("expected manifest warning for missing fields")
	}
}

func TestGetSystemInfo(t *testing.T) {
	info := GetSystemInfo()
	if info == "" {
		t.Error("GetSystemInfo returned empty string")
	}
}

