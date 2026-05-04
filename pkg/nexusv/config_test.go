package nexusv

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	t.Run("Default Values", func(t *testing.T) {
		cfg, err := LoadConfig("")
		if err != nil {
			t.Fatal(err)
		}
		if cfg.Defaults.License != "MIT" {
			t.Errorf("expected default license MIT, got %s", cfg.Defaults.License)
		}
		if cfg.Telemetry.Enabled != false {
			t.Error("expected telemetry to be disabled by default")
		}
	})

	t.Run("Environment Overrides", func(t *testing.T) {
		os.Setenv("NEXUSV_PUBLISHER", "test-pub")
		os.Setenv("NEXUSV_TELEMETRY", "on")
		defer os.Unsetenv("NEXUSV_PUBLISHER")
		defer os.Unsetenv("NEXUSV_TELEMETRY")

		cfg, err := LoadConfig("")
		if err != nil {
			t.Fatal(err)
		}
		if cfg.Defaults.Publisher != "test-pub" {
			t.Errorf("expected publisher test-pub, got %s", cfg.Defaults.Publisher)
		}
		if cfg.Telemetry.Enabled != true {
			t.Error("expected telemetry to be enabled via env")
		}
	})

	t.Run("DO_NOT_TRACK Support", func(t *testing.T) {
		os.Setenv("DO_NOT_TRACK", "1")
		os.Setenv("NEXUSV_TELEMETRY", "on")
		defer os.Unsetenv("DO_NOT_TRACK")
		defer os.Unsetenv("NEXUSV_TELEMETRY")

		cfg, err := LoadConfig("")
		if err != nil {
			t.Fatal(err)
		}
		if cfg.Telemetry.Enabled != false {
			t.Error("expected DO_NOT_TRACK to override telemetry setting")
		}
	})

	t.Run("Project Level Config", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "nexusv-proj-*")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		projCfg := filepath.Join(tmpDir, ".nexusvrc.yaml")
		content := "defaults:\n  license: Apache-2.0\n  variant: multi-command"
		if err := os.WriteFile(projCfg, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test config: %v", err)
		}

		cfg, err := LoadConfig(tmpDir)
		if err != nil {
			t.Fatal(err)
		}
		if cfg.Defaults.License != "Apache-2.0" {
			t.Errorf("expected Apache-2.0 from project config, got %s", cfg.Defaults.License)
		}
		if cfg.Defaults.Variant != "multi-command" {
			t.Errorf("expected multi-command, got %s", cfg.Defaults.Variant)
		}
	})
}

func TestConfigValidation(t *testing.T) {
	cfg := Config{}
	cfg.Defaults.License = "MIT"
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected MIT to be valid, got %v", err)
	}

	cfg.Defaults.License = "INVALID"
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for invalid license, got nil")
	}

	cfg.Defaults.License = ""
	cfg.Defaults.Variant = ""
	if err := cfg.Validate(); err != nil {
		t.Fatal(err)
	}
	if cfg.Defaults.Variant != "command" {
		t.Error("expected default variant to be set during validation")
	}
}

