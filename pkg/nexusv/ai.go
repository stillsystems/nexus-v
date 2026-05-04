package nexusv

import (
	"fmt"
	"strings"
)

// AIEngine defines the interface for AI-powered scaffolding.
type AIEngine interface {
	GenerateTemplate(prompt string) (*TemplateMetadata, error)
}

// SimpleAIEngine is a heuristic-based implementation for initial Phase 7 rollout.
type SimpleAIEngine struct{}

func (e *SimpleAIEngine) GenerateTemplate(prompt string) (*TemplateMetadata, error) {
	prompt = strings.ToLower(prompt)
	
	meta := &TemplateMetadata{
		Name:        "AI Generated Template",
		Description: fmt.Sprintf("Template generated from prompt: %s", prompt),
		Language:    "Unknown",
		Version:     "0.1.0",
	}

	// Heuristics
	if strings.Contains(prompt, "go") || strings.Contains(prompt, "golang") {
		meta.Language = "Go"
		meta.Features = append(meta.Features, Feature{ID: "cli", Name: "CLI Support"})
	}
	
	if strings.Contains(prompt, "python") {
		meta.Language = "Python"
		meta.Features = append(meta.Features, Feature{ID: "lint", Name: "Ruff Linting"})
	}

	if strings.Contains(prompt, "tailwind") {
		meta.Features = append(meta.Features, Feature{ID: "tailwind", Name: "Tailwind CSS"})
	}

	return meta, nil
}

// GenerateFromPrompt is the entry point for AI scaffolding.
func GenerateFromPrompt(prompt string, engine AIEngine) (*TemplateMetadata, error) {
	if engine == nil {
		engine = &SimpleAIEngine{}
	}
	return engine.GenerateTemplate(prompt)
}
