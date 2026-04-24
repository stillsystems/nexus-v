package cli

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

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
	exe := filepath.Base(os.Args[0])
	fmt.Printf("%s⚓ SailorOps NEXUS-V%s\n", Bold, Reset)
	fmt.Printf("Modern developer tooling engineered for real-world conditions.\n\n")
	
	fmt.Printf("%sUsage:%s %s <command> [options]\n", Bold, Reset, exe)
	fmt.Println("\nCommands:")
	fmt.Println("  init, i    Scaffold a new project")
	fmt.Println("  list, ls   List template variants")
	fmt.Println("  doctor, dr Check environment health")
	fmt.Println("  update, u  Update to latest version")
	fmt.Println("  version, v Print version info")
	
	fmt.Println("\nFlags (init):")
	fmt.Println("  --out      Output directory")
	fmt.Println("  --variant  Template variant")
	fmt.Println("  --license  License type (MIT, Apache-2.0, etc.)")
	fmt.Println("  --dry-run  Preview without writing")
	fmt.Println("  --force    Overwrite existing files")
	
	fmt.Println("\nLearn more at: https://github.com/SailorOps/Nexus-V")
}

func runList(args []string) {
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	templateDir := listCmd.String("template-dir", "", "Remote template directory to list")
	templateRef := listCmd.String("template-ref", "", "Git ref for remote template")
	listCmd.Parse(args)

	var templatesList []string
	var err error

	if *templateDir != "" {
		templatesList, err = templates.ListRemoteTemplates(*templateDir, *templateRef)
	} else {
		templatesList, err = templates.ListTemplates()
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

func runInit(args []string) {
	initCmd := flag.NewFlagSet("init", flag.ExitOnError)

	var (
		name        = initCmd.String("name", "", "Human-readable project name")
		identifier  = initCmd.String("id", "", "Extension identifier (e.g. my-extension)")
		description = initCmd.String("description", "", "Short description")
		publisher   = initCmd.String("publisher", "", "Publisher ID")
		variant     = initCmd.String("variant", "", "Template variant")
		templateDir = initCmd.String("template-dir", "", "Custom template directory (Git URL or local path)")
		templateRef = initCmd.String("template-ref", "", "Git ref (branch, tag, or SHA) for remote templates")
		license     = initCmd.String("license", "", "License type (e.g. MIT, Apache-2.0)")
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
	ctx, resolvedTarget, err := resolveContext(cfg, *name, *identifier, *description, *publisher, *variant, *templateDir, *templateRef, *license, targetDir)
	if err != nil {
		Error(err.Error())
		os.Exit(1)
	}

	if ctx.CustomTemplateDir != "" && isGitURL(ctx.CustomTemplateDir) {
		Warn("CAUTION: Scaffolding from a remote template. Only trust templates from sources you know.")
		if ctx.TemplateRef == "" {
			Info("TIP: Use --template-ref <branch/tag/sha> to pin this template for reproducible builds.")
		}
	}
	targetDir = resolvedTarget
	ctx.Force = *force
	ctx.DryRun = *dryRun

	useGit := cfg.Defaults.Git && !*noGit

	tel := telemetry.New(cfg.Telemetry.Enabled, cfg.Telemetry.Session, cfg.Telemetry.Local, cfg.Telemetry.Project)

	if !*noHooks && len(cfg.Hooks.Pre) > 0 {
		Info("Running pre-scaffold hooks...")
		if err := hooks.RunHooks(targetDir, cfg.Hooks.Pre); err != nil {
			Warn("Some pre-scaffold hooks failed: " + err.Error())
		}
	}

	spin := NewSpinner()
	
	// Handle termination signals to ensure spinner stops and cursor is restored
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		spin.Stop()
		os.Exit(1)
	}()

	spin.Start("Generating project...")
	err = templates.GenerateProject(ctx, targetDir)
	spin.Stop()
	signal.Stop(sigChan) // Stop listening for signals after generation

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

func resolveContext(cfg config.Config, name, id, desc, publisher, variant, templateDir, templateRef, license, targetDir string) (templates.Context, string, error) {
	// 1. Resolve from flags and config
	ctx, resolvedTarget := resolveFlags(cfg, name, id, desc, publisher, variant, templateDir, templateRef, license, targetDir)

	// 2. If essential info is missing, resolve interactively
	if ctx.Name == "" || ctx.Identifier == "" || ctx.Description == "" || ctx.Publisher == "" {
		if err := resolveInteractive(&ctx); err != nil {
			return templates.Context{}, "", err
		}
	}

	// 3. Post-resolution defaults
	if ctx.Template == "" && ctx.CustomTemplateDir == "" {
		ctx.Template = "command"
	}
	if resolvedTarget == "" {
		resolvedTarget = ctx.Identifier
	}

	// 4. System info
	u, _ := user.Current()
	if u != nil {
		ctx.UserName = u.Username
	} else {
		ctx.UserName = "unknown"
	}

	if out, err := exec.Command("node", "-v").Output(); err == nil {
		ctx.NodeVersion = strings.TrimSpace(string(out))
	} else {
		ctx.NodeVersion = "unknown"
	}

	// 5. Final validation (already partly done in config.LoadConfig, but good to be sure)
	if ctx.License != "" {
		supported := map[string]bool{"MIT": true, "Apache-2.0": true, "GPL-3.0": true, "BSD-3-Clause": true, "Unlicense": true, "None": true}
		if !supported[ctx.License] {
			return templates.Context{}, "", fmt.Errorf("unsupported license: %s", ctx.License)
		}
	}

	return ctx, resolvedTarget, nil
}

func resolveFlags(cfg config.Config, name, id, desc, publisher, variant, templateDir, templateRef, license, targetDir string) (templates.Context, string) {
	return templates.Context{
		Name:              name,
		Identifier:        id,
		Description:       desc,
		Publisher:         firstNonEmpty(publisher, cfg.Defaults.Publisher),
		CommandName:       id + ".helloWorld",
		Template:          firstNonEmpty(variant, cfg.Defaults.Variant),
		TemplateRef:       templateRef,
		CustomTemplateDir: templateDir,
		License:           firstNonEmpty(license, cfg.Defaults.License),
	}, targetDir
}

func resolveInteractive(ctx *templates.Context) error {
	answers, err := prompts.AskQuestions()
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	if ctx.Name == "" {
		ctx.Name = answers.Name
	}
	if ctx.Identifier == "" {
		ctx.Identifier = answers.Identifier
	}
	if ctx.Description == "" {
		ctx.Description = answers.Description
	}
	if ctx.Publisher == "" {
		ctx.Publisher = answers.Publisher
	}
	if ctx.Template == "" && ctx.CustomTemplateDir == "" && answers.Variant != "" {
		ctx.Template = answers.Variant
	}

	// Update CommandName if Identifier changed
	if ctx.CommandName == ".helloWorld" || ctx.CommandName == "" {
		ctx.CommandName = ctx.Identifier + ".helloWorld"
	}

	return nil
}

func isGitURL(path string) bool {
	return strings.HasPrefix(path, "http://") ||
		strings.HasPrefix(path, "https://") ||
		strings.HasPrefix(path, "git@") ||
		strings.HasPrefix(path, "file://")
}
