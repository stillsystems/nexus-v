package nexusv

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// RegistryTemplate defines the structure of a template entry in the remote registry.
type RegistryTemplate struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	URL         string   `json:"url"`
	Language    string   `json:"language"`
	Tags        []string `json:"tags"`
}

// Registry represents the collection of templates available in the Still Systems ecosystem.
type Registry struct {
	Templates []RegistryTemplate
}

// FetchRegistry retrieves the template index from a remote URL.
func FetchRegistry(url string) (*Registry, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add Still Systems telemetry headers
	req.Header.Set("User-Agent", "github.com/stillsystems/nexus-v/0.2.8 (Still Systems)")
	req.Header.Set("X-Nexus-Version", "0.2.8")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("network error fetching registry: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("registry server returned error: %s", resp.Status)
	}

	var templates []RegistryTemplate
	if err := json.NewDecoder(resp.Body).Decode(&templates); err != nil {
		return nil, fmt.Errorf("failed to parse registry data: %w", err)
	}

	return &Registry{Templates: templates}, nil
}

// Search filters the registry templates based on a keyword match in name or description.
func (r *Registry) Search(query string) []RegistryTemplate {
	if query == "" {
		return r.Templates
	}

	var results []RegistryTemplate
	query = strings.ToLower(query)

	for _, t := range r.Templates {
		if strings.Contains(strings.ToLower(t.Name), query) ||
			strings.Contains(strings.ToLower(t.Description), query) ||
			matchTags(t.Tags, query) {
			results = append(results, t)
		}
	}

	return results
}

func matchTags(tags []string, query string) bool {
	for _, tag := range tags {
		if strings.Contains(strings.ToLower(tag), query) {
			return true
		}
	}
	return false
}

