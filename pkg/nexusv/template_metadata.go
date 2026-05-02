package nexusv

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// TemplateMetadata holds information defined by the template author.
type TemplateMetadata struct {
	Name        string     `yaml:"name"`
	Description string     `yaml:"description"`
	Language    string     `yaml:"language"`
	Features    []Feature  `yaml:"features"`
	Hooks       HookConfig `yaml:"hooks"`
}

type Feature struct {
	ID          string   `yaml:"id"`
	Prompt      string   `yaml:"prompt"`
	Default     bool     `yaml:"default"`
	Files       []string `yaml:"files"`
}

// LoadTemplateMetadata attempts to load a nexus-template.yaml from the given directory.
func LoadTemplateMetadata(dir string) (*TemplateMetadata, error) {
	path := filepath.Join(dir, "nexus-template.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // No metadata is acceptable
		}
		return nil, err
	}

	var meta TemplateMetadata
	if err := yaml.Unmarshal(data, &meta); err != nil {
		return nil, err
	}

	return &meta, nil
}
