package templates

import (
	"embed"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"nexus-v/internal/git"
)

//go:embed files/**
var templateFS embed.FS

// Context holds the variables and flags used during the template rendering process.
type Context struct {
	Name              string
	Identifier        string
	Description       string
	Publisher         string
	CommandName       string
	Template          string
	TemplateRef       string
	CustomTemplateDir string
	License           string
	UserName          string
	NodeVersion       string
	Force             bool
	DryRun            bool
}

// ListTemplates returns all user-facing template variants in files/
func ListTemplates() ([]string, error) {
	templates, err := listFromDir("files", true)
	if err != nil {
		return nil, err
	}
	return filterInternal(templates), nil
}

// ListRemoteTemplates lists all template variants available in a remote Git repository.
func ListRemoteTemplates(url, ref string) ([]string, error) {
	// If it's a local directory, just list from there
	if !isGitURL(url) {
		return listFromDir(filepath.Join(url, "files"), false)
	}

	if !git.Available() {
		return nil, fmt.Errorf("git is not installed but is required for remote templates")
	}

	tmpDir, err := os.MkdirTemp("", "nexusv-ls-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	if err := git.CloneWithRef(url, ref, tmpDir); err != nil {
		return nil, fmt.Errorf("failed to clone remote template: %w", err)
	}

	// Remote templates are expected to have a "files/" directory for variants
	return listFromDir(filepath.Join(tmpDir, "files"), false)
}

func filterInternal(templates []string) []string {
	var filtered []string
	for _, t := range templates {
		if t == "default" || t == ".vscode" {
			continue
		}
		filtered = append(filtered, t)
	}
	return filtered
}

func listFromDir(dir string, embedded bool) ([]string, error) {
	var entries []os.DirEntry
	var err error

	if embedded {
		entries, err = templateFS.ReadDir(dir)
	} else {
		entries, err = os.ReadDir(dir)
	}

	if err != nil {
		if !embedded && os.IsNotExist(err) {
			return nil, fmt.Errorf("this repository does not follow the NEXUS-V plugin convention (missing 'files/' directory)")
		}
		return nil, err
	}

	var out []string
	for _, e := range entries {
		if e.IsDir() {
			out = append(out, e.Name())
		}
	}
	return out, nil
}

// GenerateProject scaffolds a new project by rendering templates from one or more sources
// into the target directory based on the provided Context.
func GenerateProject(ctx Context, targetDir string) error {
	if ctx.Template == "" && ctx.CustomTemplateDir == "" {
		ctx.Template = "default"
	}

	if !ctx.DryRun {
		if err := os.MkdirAll(targetDir, 0o755); err != nil {
			return fmt.Errorf("failed to create target directory: %w", err)
		}
	}

	seen := map[string]bool{}

	type source struct {
		path    string
		isLocal bool
	}

	// Handle remote templates (Plugins)
	if isGitURL(ctx.CustomTemplateDir) {
		if !git.Available() {
			return fmt.Errorf("git is not installed but is required for remote templates")
		}
		tmpDir, err := os.MkdirTemp("", "nexusv-tpl-*")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tmpDir)

		fmt.Printf("Cloning remote template: %s (ref: %s)...\n", ctx.CustomTemplateDir, ctx.TemplateRef)
		if err := git.CloneWithRef(ctx.CustomTemplateDir, ctx.TemplateRef, tmpDir); err != nil {
			return fmt.Errorf("failed to clone remote template: %w", err)
		}
		ctx.CustomTemplateDir = tmpDir
	}

	var sources []source
	if ctx.CustomTemplateDir != "" {
		sources = append(sources, source{path: ctx.CustomTemplateDir, isLocal: true})
	}
	if ctx.Template != "" && ctx.Template != "default" {
		sources = append(sources, source{path: path.Join("files", ctx.Template), isLocal: false})
	}
	sources = append(sources, source{path: path.Join("files", "default"), isLocal: false})

	for _, src := range sources {
		if src.path == "" {
			continue
		}

		var err error
		if src.isLocal {
			err = filepath.WalkDir(src.path, func(p string, d os.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if p == src.path {
					return nil
				}

				rel, _ := filepath.Rel(src.path, p)
				rel = filepath.ToSlash(rel)

				return processItem(rel, p, d.IsDir(), true, targetDir, ctx, seen)
			})
		} else {
			// Check if embedded path exists
			if _, err := templateFS.ReadDir(src.path); err != nil {
				if ctx.Template != "" && src.path == path.Join("files", ctx.Template) {
					return fmt.Errorf("unknown template variant: %s", ctx.Template)
				}
				continue
			}

			err = fsWalk(src.path, func(p string, isDir bool) error {
				rel := strings.TrimPrefix(p, src.path)
				rel = strings.TrimPrefix(rel, "/")

				return processItem(rel, p, isDir, false, targetDir, ctx, seen)
			})
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func processItem(rel, srcPath string, isDir bool, isLocal bool, targetDir string, ctx Context, seen map[string]bool) error {
	if seen[rel] {
		return nil
	}
	seen[rel] = true

	outPath := filepath.Join(targetDir, filepath.FromSlash(rel))

	// Support template rendering in filenames
	if strings.Contains(outPath, "{{") {
		t, err := template.New("path").Parse(outPath)
		if err != nil {
			return fmt.Errorf("failed to parse filename template %q: %w", outPath, err)
		}
		var buf strings.Builder
		if err := t.Execute(&buf, ctx); err != nil {
			return fmt.Errorf("failed to execute filename template %q: %w", outPath, err)
		}
		rendered := buf.String()
		if rendered == "" {
			return fmt.Errorf("filename template %q rendered to an empty string", outPath)
		}
		outPath = rendered
	}

	// Ensure outPath is within targetDir (Zip Slip protection)
	absTarget, err := filepath.Abs(targetDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute target path: %w", err)
	}
	absOut, err := filepath.Abs(outPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute output path: %w", err)
	}
	if !strings.HasPrefix(absOut, absTarget) {
		return fmt.Errorf("security violation: path %q attempts to write outside of target directory %q", outPath, targetDir)
	}

	if isDir {
		if ctx.DryRun {
			fmt.Println("[dir]  ", outPath)
			return nil
		}
		return os.MkdirAll(outPath, 0o755)
	}

	outPath = strings.TrimSuffix(outPath, ".tmpl")

	if !ctx.Force && !ctx.DryRun {
		if _, err := os.Stat(outPath); err == nil {
			return fmt.Errorf(
				"refusing to overwrite existing file: %s (use --force to override)",
				outPath,
			)
		}
	}

	if isLocal {
		return renderLocalFile(srcPath, outPath, ctx)
	}
	return renderEmbeddedFile(srcPath, outPath, ctx)
}

func isGitURL(path string) bool {
	return strings.HasPrefix(path, "http://") ||
		strings.HasPrefix(path, "https://") ||
		strings.HasPrefix(path, "git@")
}

