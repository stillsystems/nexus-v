package cli

import (
	"fmt"
	"os"
	"github.com/stillsystems/nexus-v/pkg/nexusv"
)

func runConfig(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: nexus-v config <get|set|list> [key] [value]")
		return
	}

	cfg, err := nexusv.LoadConfig("")
	if err != nil {
		Error("Failed to load config: " + err.Error())
		os.Exit(1)
	}

	action := args[0]
	switch action {
	case "get":
		if len(args) < 2 {
			Error("Missing key")
			return
		}
		key := args[1]
		switch key {
		case "publisher":
			fmt.Println(cfg.Defaults.Publisher)
		case "license":
			fmt.Println(cfg.Defaults.License)
		case "variant":
			fmt.Println(cfg.Defaults.Variant)
		case "git":
			fmt.Println(cfg.Defaults.Git)
		default:
			Error("Unknown key: " + key)
		}
	case "set":
		if len(args) < 3 {
			Error("Missing key or value")
			return
		}
		key := args[1]
		val := args[2]
		switch key {
		case "publisher":
			cfg.Defaults.Publisher = val
		case "license":
			cfg.Defaults.License = val
		case "variant":
			cfg.Defaults.Variant = val
		case "git":
			cfg.Defaults.Git = (val == "true" || val == "1")
		default:
			Error("Unknown key: " + key)
			return
		}
		
		if err := cfg.Validate(); err != nil {
			Error("Invalid value: " + err.Error())
			return
		}

		if err := cfg.Save(); err != nil {
			Error("Failed to save config: " + err.Error())
			os.Exit(1)
		}
		Success(fmt.Sprintf("Set %s to %s", key, val))
	case "list":
		fmt.Printf("publisher: %s\n", cfg.Defaults.Publisher)
		fmt.Printf("license: %s\n", cfg.Defaults.License)
		fmt.Printf("variant: %s\n", cfg.Defaults.Variant)
		fmt.Printf("git: %v\n", cfg.Defaults.Git)
	default:
		Error("Unknown action: " + action)
	}
}

