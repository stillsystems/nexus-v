package nexusv

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type DefaultsConfig struct {
	Publisher string `yaml:"publisher"`
	Variant   string `yaml:"variant"`
	Git       bool   `yaml:"git"`
	License   string `yaml:"license"`
}

type TelemetryConfig struct {
	Enabled bool `yaml:"enabled"`
	Session bool `yaml:"session"`
	Local   bool `yaml:"local"`
	Project bool `yaml:"project"`
}

type HookConfig struct {
	Pre  []string `yaml:"pre_scaffold"`
	Post []string `yaml:"post_scaffold"`
}

type Config struct {
	Defaults  DefaultsConfig  `yaml:"defaults"`
	Telemetry TelemetryConfig `yaml:"telemetry"`
	Hooks     HookConfig      `yaml:"hooks"`
}

func (c *Config) Validate() error {
	supportedLicenses := map[string]bool{
		"MIT":          true,
		"Apache-2.0":   true,
		"GPL-3.0":      true,
		"BSD-3-Clause": true,
		"Unlicense":    true,
		"None":         true,
	}

	if c.Defaults.License != "" && !supportedLicenses[c.Defaults.License] {
		return fmt.Errorf("unsupported license in config: %s", c.Defaults.License)
	}

	// Basic validation for variant - though custom templates might use anything,
	// we can at least check if it's empty if not using customTemplateDir.
	if c.Defaults.Variant == "" {
		c.Defaults.Variant = "command"
	}

	return nil
}

func UserConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".nexusvrc.yaml"), nil
}

func (c *Config) Save() error {
	path, err := UserConfigPath()
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}

func LoadConfig(targetDir string) (Config, error) {
	var cfg Config

	// Default values
	cfg.Defaults.Git = true
	cfg.Defaults.License = "MIT"
	cfg.Defaults.Variant = "command"
	cfg.Telemetry.Enabled = false // Off by default per README
	cfg.Telemetry.Session = true  // Internal default if enabled

	// User-level config: ~/.nexusvrc.yaml
	if userCfg, err := UserConfigPath(); err == nil {
		if data, err := os.ReadFile(userCfg); err == nil {
			_ = yaml.Unmarshal(data, &cfg)
		}
	}

	// Project-level config: <targetDir>/.nexusvrc.yaml
	if targetDir != "" {
		projectCfg := filepath.Join(targetDir, ".nexusvrc.yaml")
		if data, err := os.ReadFile(projectCfg); err == nil {
			_ = yaml.Unmarshal(data, &cfg)
		}
	}

	// Environment variables override
	if pub := os.Getenv("NEXUSV_PUBLISHER"); pub != "" {
		cfg.Defaults.Publisher = pub
	}
	if variant := os.Getenv("NEXUSV_DEFAULT_VARIANT"); variant != "" {
		cfg.Defaults.Variant = variant
	}

	telemetryEnv := strings.ToLower(os.Getenv("NEXUSV_TELEMETRY"))
	doNotTrack := os.Getenv("DO_NOT_TRACK")
	if doNotTrack == "1" || doNotTrack == "true" || telemetryEnv == "off" || telemetryEnv == "false" || telemetryEnv == "0" {
		cfg.Telemetry.Enabled = false
	} else if telemetryEnv == "on" || telemetryEnv == "true" || telemetryEnv == "1" {
		cfg.Telemetry.Enabled = true
	}

	return cfg, cfg.Validate()
}

