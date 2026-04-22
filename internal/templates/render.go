package templates

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

func renderTemplate(data []byte, name, outPath string, ctx Context) error {
	tmpl, err := template.New(name).Funcs(template.FuncMap{
		"currentYear": func() int {
			return time.Now().Year()
		},
	}).Parse(string(data))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, ctx); err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	if ctx.DryRun {
		return nil
	}

	return os.WriteFile(outPath, buf.Bytes(), 0o644)
}

func renderEmbeddedFile(srcPath, outPath string, ctx Context) error {
	data, err := templateFS.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("failed to read embedded template: %w", err)
	}
	return renderTemplate(data, filepath.Base(srcPath), outPath, ctx)
}

func renderLocalFile(srcPath, outPath string, ctx Context) error {
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("failed to read local template: %w", err)
	}
	return renderTemplate(data, filepath.Base(srcPath), outPath, ctx)
}
