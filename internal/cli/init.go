package cli

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"

	"nexus-v/internal/prompts"
	"nexus-v/internal/telemetry"
	"nexus-v/pkg/nexusv"
)

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
	cfg, _ := nexusv.LoadConfig(cfgTarget)

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
		if err := nexusv.RunHooks(targetDir, cfg.Hooks.Pre); err != nil {
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
	meta, err := nexusv.GenerateProject(ctx, targetDir)
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

	results := nexusv.ValidateProject(targetDir)
	for _, res := range results {
		for _, e := range res.Errors {
			Warn(fmt.Sprintf("[%s] %s", res.File, e))
		}
		for _, w := range res.Warnings {
			Info(fmt.Sprintf("[%s] %s", res.File, w))
		}
	}

	if ctx.DryRun {
		Success("Dry run complete — no files were written")
		return
	}

	Success("Project created at " + filepath.Clean(targetDir))

	if useGit && nexusv.GitAvailable() {
		Info("Initializing Git repository...")
		if err := nexusv.GitInitRepo(targetDir); err == nil {
			nexusv.GitAddAll(targetDir)
			nexusv.GitFirstCommit(targetDir)
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

	if meta != nil && !*noHooks {
		if len(meta.Hooks.Pre) > 0 {
			Warn("Template defines pre-scaffold hooks, but generation is already underway. Consider moving these to post_scaffold.")
		}
		postHooks = append(postHooks, meta.Hooks.Post...)
	}

	if !*noHooks && len(postHooks) > 0 {
		Info("Running post-generation hooks...")
		if err := nexusv.RunHooks(targetDir, postHooks); err != nil {
			Warn("Some post-generation hooks failed: " + err.Error())
		} else {
			Success("All hooks completed")
		}
	}

	Info("Run `npm install` then press F5 to launch the extension")
}

func resolveContext(cfg nexusv.Config, name, id, desc, publisher, variant, templateDir, templateRef, license, targetDir string) (nexusv.Context, string, error) {
	// 1. Resolve from flags and config
	ctx, resolvedTarget := resolveFlags(cfg, name, id, desc, publisher, variant, templateDir, templateRef, license, targetDir)

	// 2. If essential info is missing, resolve interactively
	if ctx.Name == "" || ctx.Identifier == "" || ctx.Description == "" || ctx.Publisher == "" {
		if err := resolveInteractive(&ctx); err != nil {
			return nexusv.Context{}, "", err
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
		ctx.UserName = "developer"
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
			return nexusv.Context{}, "", fmt.Errorf("unsupported license: %s", ctx.License)
		}
	}

	return ctx, resolvedTarget, nil
}

func resolveFlags(cfg nexusv.Config, name, id, desc, publisher, variant, templateDir, templateRef, license, targetDir string) (nexusv.Context, string) {
	return nexusv.Context{
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

func resolveInteractive(ctx *nexusv.Context) error {
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
	if ctx.License == "" {
		ctx.License = answers.License
	}

	// Update CommandName if Identifier changed
	if ctx.CommandName == ".helloWorld" || ctx.CommandName == "" {
		ctx.CommandName = answers.CommandName
	}

	return nil
}

func isGitURL(path string) bool {
	return strings.HasPrefix(path, "http://") ||
		strings.HasPrefix(path, "https://") ||
		strings.HasPrefix(path, "git@") ||
		strings.HasPrefix(path, "file://")
}

func firstNonEmpty(a, b string) string {
	if a != "" {
		return a
	}
	return b
}
