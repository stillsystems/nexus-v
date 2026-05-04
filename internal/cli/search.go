package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"nexus-v/pkg/nexusv"
)

func runSearch(args []string) {
	searchCmd := flag.NewFlagSet("search", flag.ExitOnError)
	asJSON := searchCmd.Bool("json", false, "Output results as JSON")
	searchCmd.Parse(args)

	query := ""
	if searchCmd.NArg() > 0 {
		query = searchCmd.Arg(0)
	}

	// Official Still Systems registry index
	registryURL := "https://raw.githubusercontent.com/stillsystems/nexus-registry/main/templates.json"

	Info("Fetching template registry...")
	reg, err := nexusv.FetchRegistry(registryURL)
	if err != nil {
		Error("Failed to connect to registry: " + err.Error())
		os.Exit(1)
	}

	results := reg.Search(query)
	if *asJSON {
		data, _ := json.MarshalIndent(results, "", "  ")
		fmt.Println(string(data))
		return
	}

	if len(results) == 0 {
		Warn("No templates found matching '" + query + "'")
		return
	}

	Success(fmt.Sprintf("Found %d templates:", len(results)))
	fmt.Println()
	fmt.Printf("%-25s %-45s %-15s\n", "NAME", "DESCRIPTION", "TAGS")
	fmt.Println(strings.Repeat("-", 85))

	for _, t := range results {
		desc := t.Description
		if len(desc) > 42 {
			desc = desc[:39] + "..."
		}
		tags := strings.Join(t.Tags, ", ")
		fmt.Printf("%-25s %-45s %-15s\n", t.Name, desc, tags)
	}
	fmt.Println()
	Info("Use `nexus-v init --template <name>` to scaffold using one of these.")
}
