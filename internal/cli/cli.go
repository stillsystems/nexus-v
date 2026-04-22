package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"nexus-v/internal/config"
	"nexus-v/internal/doctor"
	"nexus-v/internal/git"
	"nexus-v/internal/hooks"
	"nexus-v/internal/prompts"
	"nexus-v/internal/telemetry"
	"nexus-v/internal/templates"
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
	case "variants", "vars", "ls":
		runVariants()
	case "version", "v", "-v", "--version":
		fmt.Printf("nexus-v %s (%s/%s)\n", version.Version, runtime.GOOS, runtime.GOARCH)
	case "update", "u":
		if err := update.CheckAndApply(); err != nil {
			Error("Update failed: " + err.Error())
			os.Exit(1)
		}
	case "doctor", "dr":
		doctor.RunChecks()
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: nexus-v <command> [options]")
	fmt.Println("\nCommands:")
	fmt.Println("  init      Scaffold a new VS Code extension project")
	fmt.Println("  variants  List available template variants")
	fmt.Println("  version   Print version information")
	fmt.Println("  update    Update nexus-v to the latest version")
	fmt.Println("  doctor    Check environment for required tools")
}

func runVariants() {
	templatesList, err := templates.ListTemplates()
	if err != nil {
		Error("Failed to list templates: " + err.Error())
		os.Exit(1)
	}
	Success("Available templates:")
	for _, t := range templatesList {
		fmt.Println(" -", t)
	}
}

func runInit(args []string) {
	initCmd := flag.NewFlagSet("init", flag.ExitOnError)

	var (
		name        = initCmd.String("name", "", "Human-readable project name")
		identifier  = initCmd.String("id", "", "Extension identifier (e.g. my-extension)")
		description = initCmd.String("description", "", "Short description")
		publisher   = initCmd.String("publisher", "", "Publisher ID")
		variant     = initCmd.String("variant", "", "Template variant")
		templateDir = initCmd.String("template-dir", "", "Custom template directory")
		force       = initCmd.Bool("force", false, "Overwrite existing files")
		dryRun      = initCmd.Bool("dry-run", false, "Preview files without writing them")
		out         = initCmd.String("out", "", "Output directory")
		noGit       = initCmd.Bool("no-git", false, "Do not initialize a Git repository")
		noHooks     = initCmd.Bool("no-hooks", false, "Disable pre/post-generation hooks")

		// Built-in hook shortcuts
		installDep = initCmd.Bool("install", false, "Run npm install after scaffold")
		openCode   = initCmd.Bool("open", false, "Open VS Code after scaffold")
		gitInit    = initCmd.Bool("git", false, "Run git init after scaffold")
	)

	initCmd.Parse(args)

	// Determine target directory early for config loading
	targetDir := *out
	if targetDir == "" && initCmd.NArg() > 0 {
		targetDir = initCmd.Arg(0)
	}

	cfgTarget := targetDir
	if cfgTarget == "" {
		cfgTarget = "."
	}
	cfg, _ := config.LoadConfig(cfgTarget)

	// Resolve the final context
	ctx, resolvedTarget, err := resolveContext(cfg, *name, *identifier, *description, *publisher, *variant, *templateDir, targetDir)
	if err != nil {
		Error(err.Error())
		os.Exit(1)
	}
	targetDir = resolvedTarget
	ctx.Force = *force
	ctx.DryRun = *dryRun

	useGit := cfg.Defaults.Git && !*noGit

	tel := telemetry.Telemetry{
		SessionEnabled: cfg.Telemetry.Enabled && cfg.Telemetry.Session,
		LocalEnabled:   cfg.Telemetry.Enabled && cfg.Telemetry.Local,
		ProjectEnabled: cfg.Telemetry.Enabled && cfg.Telemetry.Project,
		SessionSink:    &telemetry.SessionSink{},
		LocalSink:      &telemetry.LocalSink{},
		ProjectSink:    &telemetry.ProjectSink{},
	}

	if !*noHooks && len(cfg.Hooks.Pre) > 0 {
		Info("Running pre-scaffold hooks...")
		if err := hooks.RunHooks(targetDir, cfg.Hooks.Pre); err != nil {
			Warn("Some pre-scaffold hooks failed: " + err.Error())
		}
	}

	spin := NewSpinner()
	spin.Start("Generating project...")
	err = templates.GenerateProject(ctx, targetDir)
	spin.Stop()

	ev := telemetry.Event{
		Template:   ctx.Template,
		DryRun:     *dryRun,
		Force:      *force,
		ProjectDir: filepath.Base(targetDir),
	}
	tel.Record(ev)

	if err != nil {
		Error(err.Error())
		os.Exit(1)
	}

	if ctx.DryRun {
		Success("Dry run complete — no files were written")
		return
	}

	Success("Project created at " + filepath.Clean(targetDir))

	if useGit && git.Available() {
		Info("Initializing Git repository...")
		if err := git.InitRepo(targetDir); err == nil {
			git.AddAll(targetDir)
			git.FirstCommit(targetDir)
			Success("Git repository initialized")
		} else {
			Warn("Git is installed but initialization failed")
		}
	} else if useGit {
		Warn("Git not found — skipping repository initialization")
	}

	postHooks := append([]string{}, cfg.Hooks.Post...)
	if *installDep {
		postHooks = append(postHooks, "npm install")
	}
	if *gitInit {
		postHooks = append(postHooks, "git init", "git add -A")
	}
	if *openCode {
		postHooks = append(postHooks, "code .")
	}

	if !*noHooks && len(postHooks) > 0 {
		Info("Running post-generation hooks...")
		if err := hooks.RunHooks(targetDir, postHooks); err != nil {
			Warn("Some post-generation hooks failed: " + err.Error())
		} else {
			Success("All hooks completed")
		}
	}

	Info("Run `npm install` then press F5 to launch the extension")
}

func firstNonEmpty(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func resolveContext(cfg config.Config, name, id, desc, publisher, variant, templateDir, targetDir string) (templates.Context, string, error) {
	finalName := name
	finalID := id
	finalDesc := desc
	finalPublisher := firstNonEmpty(publisher, cfg.Defaults.Publisher)
	finalVariant := firstNonEmpty(variant, cfg.Defaults.Variant)
	finalTemplateDir := templateDir

	// If interactive mode is needed
	if finalName == "" || finalID == "" || finalDesc == "" || finalPublisher == "" {
		answers, err := prompts.AskQuestions()
		if err != nil {
			return templates.Context{}, "", fmt.Errorf("failed to read input: %w", err)
		}
		if finalName == "" {
			finalName = answers.Name
		}
		if finalID == "" {
			finalID = answers.Identifier
		}
		if finalDesc == "" {
			finalDesc = answers.Description
		}
		if finalPublisher == "" {
			finalPublisher = answers.Publisher
		}
		if finalVariant == "" && finalTemplateDir == "" && answers.Variant != "" {
			finalVariant = answers.Variant
		}
	}

	if finalVariant == "" && finalTemplateDir == "" {
		finalVariant = "command"
	}

	if targetDir == "" {
		targetDir = finalID
	}

	ctx := templates.Context{
		Name:              finalName,
		Identifier:        finalID,
		Description:       finalDesc,
		Publisher:         finalPublisher,
		CommandName:       finalID + ".helloWorld",
		Template:          finalVariant,
		CustomTemplateDir: finalTemplateDir,
	}

	return ctx, targetDir, nil
}
