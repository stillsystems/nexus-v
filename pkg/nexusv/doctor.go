package nexusv

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Issue represents a diagnostic finding.
type Issue struct {
	Name    string
	Level   string // "Success", "Warning", "Error"
	Message string
	Fix     string
	FixID   string
}

// DoctorResult contains the summary of a diagnostic run.
type DoctorResult struct {
	Issues  []Issue
	Healthy bool
}

// RunFullDoctor performs both environment and project-specific checks.
func RunFullDoctor(projectDir string) *DoctorResult {
	result := &DoctorResult{Healthy: true}

	// 1. Environment Checks (Tools)
	toolChecks := []struct {
		name string
		cmd  string
		req  bool
		desc string
	}{
		{"Node.js", "node", true, "Required to run and develop VS Code extensions."},
		{"npm", "npm", true, "Standard package manager for extension dependencies."},
		{"git", "git", true, "Used for version control and scaffolding initialization."},
		{"vsce", "vsce", false, "Official tool for packaging and publishing."},
	}

	for _, c := range toolChecks {
		_, err := exec.LookPath(c.cmd)
		if err != nil {
			level := "Warning"
			if c.req {
				level = "Error"
				result.Healthy = false
			}
			result.Issues = append(result.Issues, Issue{
				Name:    c.name,
				Level:   level,
				Message: fmt.Sprintf("%s not found. %s", c.name, c.desc),
				Fix:     fmt.Sprintf("Install %s to proceed.", c.name),
				FixID:   "tool_missing_" + c.cmd,
			})
		} else {
			result.Issues = append(result.Issues, Issue{
				Name:  c.name,
				Level: "Success",
			})
		}
	}

	// 2. Project Checks (if inside a project)
	if projectDir != "" {
		checkProjectIntegrity(projectDir, result)
	}

	return result
}

func checkProjectIntegrity(dir string, result *DoctorResult) {
	pkgPath := filepath.Join(dir, "package.json")
	if _, err := os.Stat(pkgPath); os.IsNotExist(err) {
		return // Not a project directory, skip project checks
	}

	data, err := os.ReadFile(pkgPath)
	if err != nil {
		result.Issues = append(result.Issues, Issue{
			Name:    "Manifest",
			Level:   "Error",
			Message: "Could not read package.json",
		})
		result.Healthy = false
		return
	}

	var pkg map[string]interface{}
	if err := json.Unmarshal(data, &pkg); err != nil {
		result.Issues = append(result.Issues, Issue{
			Name:    "Manifest",
			Level:   "Error",
			Message: "package.json is invalid JSON",
		})
		result.Healthy = false
		return
	}

	// Check for vital fields
	vital := []string{"name", "publisher", "engines", "activationEvents"}
	for _, f := range vital {
		if _, ok := pkg[f]; !ok {
			result.Issues = append(result.Issues, Issue{
				Name:    "Manifest",
				Level:   "Warning",
				Message: fmt.Sprintf("Missing vital field: %s", f),
				Fix:     fmt.Sprintf("Add '%s' to package.json", f),
				FixID:   "manifest_missing_" + f,
			})
		}
	}
}

// FixIssue attempts to automatically repair a diagnosed issue.
func FixIssue(dir string, issue Issue) error {
	if strings.HasPrefix(issue.FixID, "manifest_missing_") {
		field := strings.TrimPrefix(issue.FixID, "manifest_missing_")
		return fixMissingManifestField(dir, field)
	}
	if strings.HasPrefix(issue.FixID, "tool_missing_") {
		tool := strings.TrimPrefix(issue.FixID, "tool_missing_")
		return fixMissingTool(tool)
	}
	return fmt.Errorf("no automated fix available for %q", issue.FixID)
}

func fixMissingTool(tool string) error {
	var cmd *exec.Cmd
	switch tool {
	case "vsce":
		fmt.Printf(" (Running 'npm install -g @vscode/vsce')...")
		cmd = exec.Command("npm", "install", "-g", "@vscode/vsce")
	case "ovsx":
		fmt.Printf(" (Running 'npm install -g ovsx')...")
		cmd = exec.Command("npm", "install", "-g", "ovsx")
	default:
		return fmt.Errorf("don't know how to install %s automatically", tool)
	}

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("install failed: %s", string(output))
	}
	return nil
}
func fixMissingManifestField(dir, field string) error {
	pkgPath := filepath.Join(dir, "package.json")
	data, err := os.ReadFile(pkgPath)
	if err != nil {
		return err
	}

	var pkg map[string]interface{}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return err
	}

	// Default values for missing vital fields
	switch field {
	case "name":
		pkg["name"] = filepath.Base(dir)
	case "publisher":
		pkg["publisher"] = "stillsystems"
	case "engines":
		pkg["engines"] = map[string]string{"vscode": "^1.90.0"}
	case "activationEvents":
		pkg["activationEvents"] = []string{}
	default:
		return fmt.Errorf("don't know how to fix missing field %q", field)
	}

	newData, err := json.MarshalIndent(pkg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(pkgPath, newData, 0o644)
}

// GetSystemInfo returns OS/Arch info.
func GetSystemInfo() string {
	return fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
}
