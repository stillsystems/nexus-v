package nexusv

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ValidationResult captures linting issues for a specific file.
type ValidationResult struct {
	File     string
	Errors   []string
	Warnings []string
}

// ValidateProject performs deep-linting on generated project files.
func ValidateProject(dir string) []ValidationResult {
	var results []ValidationResult

	// Validate package.json
	if res, err := validatePackageJSON(filepath.Join(dir, "package.json")); err == nil && res != nil {
		results = append(results, *res)
	}

	// Validate tsconfig.json
	if res, err := validateTSConfig(filepath.Join(dir, "tsconfig.json")); err == nil && res != nil {
		results = append(results, *res)
	}

	return results
}

func validatePackageJSON(path string) (*ValidationResult, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var pkg map[string]interface{}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return &ValidationResult{
			File:   "package.json",
			Errors: []string{fmt.Sprintf("Invalid JSON: %v", err)},
		}, nil
	}

	res := &ValidationResult{File: "package.json"}

	// Required fields
	required := []string{"name", "publisher", "version", "engines", "categories", "contributes"}
	for _, field := range required {
		if _, ok := pkg[field]; !ok {
			res.Errors = append(res.Errors, fmt.Sprintf("Missing required field: %q", field))
		}
	}

	// Still Systems Best Practices
	if engines, ok := pkg["engines"].(map[string]interface{}); ok {
		if vscode, ok := engines["vscode"].(string); ok {
			if !strings.HasPrefix(vscode, "^1.") {
				res.Warnings = append(res.Warnings, "Engines.vscode should usually target a recent version (e.g., ^1.90.0)")
			}
		}
	}

	if len(res.Errors) == 0 && len(res.Warnings) == 0 {
		return nil, nil
	}
	return res, nil
}

func validateTSConfig(path string) (*ValidationResult, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	// tsconfig often allows comments, but standard json/unmarshal doesn't.
	// For simplicity in this core logic, we assume standard JSON or just basic parsing.
	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return &ValidationResult{
			File:   "tsconfig.json",
			Errors: []string{fmt.Sprintf("Invalid JSON (ensure no comments): %v", err)},
		}, nil
	}

	res := &ValidationResult{File: "tsconfig.json"}

	if opts, ok := config["compilerOptions"].(map[string]interface{}); ok {
		if strict, ok := opts["strict"].(bool); !ok || !strict {
			res.Warnings = append(res.Warnings, "compilerOptions.strict should be set to true for Still Systems projects")
		}
		if target, ok := opts["target"].(string); ok {
			if target == "ES5" || target == "ES6" {
				res.Warnings = append(res.Warnings, fmt.Sprintf("compilerOptions.target is %s; consider ES2022 or newer", target))
			}
		}
	} else {
		res.Errors = append(res.Errors, "Missing compilerOptions")
	}

	if len(res.Errors) == 0 && len(res.Warnings) == 0 {
		return nil, nil
	}
	return res, nil
}

// ValidateTemplate performs deep-linting on a Nexus-V template repository.
func ValidateTemplate(dir string) []ValidationResult {
	var results []ValidationResult

	// Validate nexus-template.yaml
	meta, err := LoadTemplateMetadata(dir)
	if err != nil {
		results = append(results, ValidationResult{
			File:   "nexus-template.yaml",
			Errors: []string{fmt.Sprintf("Failed to load metadata: %v", err)},
		})
		return results
	}

	if meta == nil {
		results = append(results, ValidationResult{
			File:   "nexus-template.yaml",
			Errors: []string{"Missing required file: nexus-template.yaml"},
		})
		return results
	}

	res := &ValidationResult{File: "nexus-template.yaml"}

	if meta.Name == "" {
		res.Errors = append(res.Errors, "Missing required field: \"name\"")
	}
	if meta.Description == "" {
		res.Warnings = append(res.Warnings, "Recommended field missing: \"description\"")
	}
	if meta.Language == "" {
		res.Errors = append(res.Errors, "Missing required field: \"language\"")
	}

	for i, feature := range meta.Features {
		if feature.ID == "" {
			res.Errors = append(res.Errors, fmt.Sprintf("Feature at index %d is missing \"id\"", i))
		}
		if len(feature.Files) == 0 {
			res.Warnings = append(res.Warnings, fmt.Sprintf("Feature %q has no associated files", feature.ID))
		}
	}

	if len(res.Errors) > 0 || len(res.Warnings) > 0 {
		results = append(results, *res)
	}

	return results
}

