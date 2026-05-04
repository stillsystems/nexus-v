package cli

import (
	"fmt"
	"os"
	"runtime"

	"nexus-v/internal/update"
	"nexus-v/internal/version"
)

func Run() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "init", "i":
		runInit(os.Args[2:])
	case "variants", "vars", "ls", "list":
		runList(os.Args[2:])
	case "version", "v", "-v", "--version":
		fmt.Printf("nexus-v %s (%s/%s)\n", version.Version, runtime.GOOS, runtime.GOARCH)
	case "update", "u":
		if err := update.CheckAndApply(); err != nil {
			Error("Update failed: " + err.Error())
			os.Exit(1)
		}
	case "doctor", "dr":
		runDoctor(os.Args[2:])
	case "validate", "lint":
		runValidate(os.Args[2:])
	case "config", "cfg":
		runConfig(os.Args[2:])
	case "search":
		runSearch(os.Args[2:])
	case "serve":
		runServe(os.Args[2:])
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Printf("%s🧱 Still Systems NEXUS-V%s\n", Bold, Reset)
	fmt.Printf("Modern developer tooling engineered for real-world conditions.\n\n")
	
	fmt.Printf("%sUsage:%s nexus-v <command> [options]\n", Bold, Reset)
	fmt.Println("\nCommands:")
	fmt.Println("  init, i    Scaffold a new project")
	fmt.Println("  list, ls   List template variants")
	fmt.Println("  doctor, dr Check environment health")
	fmt.Println("  validate   Lint and validate template structure")
	fmt.Println("  config     Manage global user preferences")
	fmt.Println("  search     Find templates in the Still Systems gallery")
	fmt.Println("  serve      Launch the Visual Scaffolder (Web UI)")
	fmt.Println("  update, u  Update to latest version")
	fmt.Println("  version, v Print version info")
	
	fmt.Println("\nFlags (init):")
	fmt.Println("  --out      Output directory")
	fmt.Println("  --variant  Template variant")
	fmt.Println("  --license  License type (MIT, Apache-2.0, etc.)")
	fmt.Println("  --dry-run  Preview without writing")
	fmt.Println("  --force    Overwrite existing files")
	
	fmt.Println("\nLearn more at: https://github.com/stillsystems/Nexus-V")
}
