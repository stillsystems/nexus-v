package cli

import (
	"flag"
	"fmt"
	"os"

	"nexus-v/pkg/nexusv"
)

func runList(args []string) {
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	templateDir := listCmd.String("template-dir", "", "Remote template directory to list")
	templateRef := listCmd.String("template-ref", "", "Git ref for remote template")
	listCmd.Parse(args)

	var templatesList []string
	var err error

	if *templateDir != "" {
		templatesList, err = nexusv.ListRemoteTemplates(*templateDir, *templateRef)
	} else {
		templatesList, err = nexusv.ListTemplates()
	}

	if err != nil {
		Error("Failed to list templates: " + err.Error())
		os.Exit(1)
	}
	Success("Available templates:")
	for _, t := range templatesList {
		fmt.Println(" -", t)
	}
}
