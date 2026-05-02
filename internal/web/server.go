package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"

	"nexus-v/pkg/nexusv"
)

//go:embed static
var staticFS embed.FS

// Start launches the Visual Scaffolder web server.
func Start(port int) error {
	// API Endpoints
	http.HandleFunc("/api/templates", handleTemplates)
	http.HandleFunc("/api/generate", handleGenerate)

	// Static Files (Frontend)
	public, err := fs.Sub(staticFS, "static")
	if err != nil {
		return err
	}
	http.Handle("/", http.FileServer(http.FS(public)))

	fmt.Printf("\n🧱 Still Systems Nexus-V\n")
	fmt.Printf("🚀 Visual Scaffolder active at http://localhost:%d\n\n", port)
	
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func handleTemplates(w http.ResponseWriter, r *http.Request) {
	templates, err := nexusv.GetAvailableTemplates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"templates": templates,
	})
}

type GenerateRequest struct {
	Name            string          `json:"name"`
	Identifier      string          `json:"identifier"`
	Publisher       string          `json:"publisher"`
	Description     string          `json:"description"`
	Template        string          `json:"template"`
	EnabledFeatures map[string]bool `json:"enabled_features"`
}

func handleGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Security: Ensure identifier is a simple name, not a path
	if strings.Contains(req.Identifier, "/") || strings.Contains(req.Identifier, "\\") || strings.Contains(req.Identifier, "..") {
		http.Error(w, "Invalid identifier: must be a simple folder name", http.StatusBadRequest)
		return
	}

	ctx := nexusv.Context{
		Name:            req.Name,
		Identifier:      req.Identifier,
		Description:     req.Description,
		Publisher:       req.Publisher,
		Template:        req.Template,
		EnabledFeatures: req.EnabledFeatures,
	}

	// We'll scaffold into a subfolder of the current directory for the web UI
	targetDir := req.Identifier
	_, err := nexusv.GenerateProject(ctx, targetDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Project generated at " + targetDir,
	})
}
