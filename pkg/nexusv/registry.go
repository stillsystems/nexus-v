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
	Name        string   `json:"name"`
	Description string   `json:"description"`
	URL         string   `json:"url"`
	Variant     string   `json:"variant"`
	Tags        []string `json:"tags"`
}

// Registry represents the collection of templates available in the Still Systems ecosystem.
type Registry struct {
	Templates []RegistryTemplate `json:"templates"`
}

// FetchRegistry retrieves the template index from a remote URL.
func FetchRegistry(url string) (*Registry, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("network error fetching registry: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("registry server returned error: %s", resp.Status)
	}

	var reg Registry
	if err := json.NewDecoder(resp.Body).Decode(&reg); err != nil {
		return nil, fmt.Errorf("failed to parse registry data: %w", err)
	}

	return &reg, nil
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
