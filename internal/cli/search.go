package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"nexus-v/pkg/nexusv"
)

func runSearch(args []string) {
	searchCmd := flag.NewFlagSet("search", flag.ExitOnError)
	searchCmd.Parse(args)

	query := ""
	if searchCmd.NArg() > 0 {
		query = searchCmd.Arg(0)
	}

	// Registry URL is hardcoded for now, pointing to the Still Systems official index.
	registryURL := "https://raw.githubusercontent.com/stillsystems/.github/main/nexus-v-registry.json"

	Info("Fetching template registry...")
	reg, err := nexusv.FetchRegistry(registryURL)
	if err != nil {
		Error("Failed to connect to registry: " + err.Error())
		os.Exit(1)
	}

	results := reg.Search(query)
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
