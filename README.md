# NEXUS-V

**A modern, zero-install VS Code extension scaffolder — built in Go, outputs TypeScript.**

NEXUS-V replaces the legacy Yeoman `yo code` generator with a single static binary that produces clean, dependency-light VS Code extension projects. No global installs, no ecosystem rot, no hidden state — just a portable tool that does one thing well.

---

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Template Variants](#template-variants)
- [Configuration](#configuration)
- [Hooks](#hooks)
- [Telemetry](#telemetry)
- [Project Structure](#project-structure)
- [Next Steps & Roadmap](#next-steps--roadmap)
- [Contributing](#contributing)
- [License](#license)

---

## Overview

The VS Code extension ecosystem depends heavily on Yeoman — a scaffolding tool with deep Node.js dependency trees, global installs, and a maintenance surface area that drifts over time. NEXUS-V takes a different approach:

- **Written in Go** — compiles to a single static binary with zero runtime dependencies.
- **Outputs TypeScript** — every scaffolded project is a clean, modern TypeScript VS Code extension ready for development.
- **Modular template engine** — variant-based template system that renders projects from embedded Go templates, not fragile file-copy heuristics.
- **Portable by design** — no PATH pollution, no global packages, no persistent state. Drop it in a directory and run it.

NEXUS-V is built for developers who want a scaffolding tool that stays out of the way and never breaks between uses.

---

## Features

| Feature | Description |
|---|---|
| **Single binary** | One executable, zero runtime dependencies — works on Windows, macOS, and Linux |
| **Embedded templates** | All project templates are compiled into the binary via Go's `embed` package |
| **Template variants** | Modular variant system for different extension types (commands, webviews, language support, etc.) |
| **Dynamic Filenames** | Support for template variables in output filenames (e.g., `{{.Identifier}}.tmLanguage.json`) |
| **Template renderer** | Dedicated `render` module that processes Go `text/template` files with project-specific data |
| **Interactive prompts** | CLI prompts collect project metadata (name, publisher, description) at scaffold time |
| **Custom Templates** | Support for loading external template directories via `--template-dir` |
| **Self-Update** | Built-in `update` command to fetch the latest version from GitHub |
| **Zero-install** | No `npm install -g`, no Yeoman, no generator packages — download and run |
| **Hook system** | Pre- and post-scaffold hooks for custom automation |
| **Opt-in telemetry** | Anonymous, minimal usage telemetry — off by default |

---

## Installation

### Download the binary

Grab the latest release for your platform from the [Releases](https://github.com/geriatric-sailor/nexus-v/releases) page.

```bash
# Linux / macOS
curl -L https://github.com/geriatric-sailor/nexus-v/releases/latest/download/nexus-v-linux-amd64 -o nexus-v
chmod +x nexus-v
```

```powershell
# Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/geriatric-sailor/nexus-v/releases/latest/download/nexus-v-windows-amd64.exe" -OutFile "nexus-v.exe"
```

### Build from source

Requires Go 1.22+.

```bash
git clone https://github.com/geriatric-sailor/nexus-v.git
cd nexus-v
go build -o nexus-v ./cmd/nexus-v.go
```

### Verify installation

```bash
nexus-v version  # or: nexus-v v
```

---

## Usage

### Interactive mode (default)

```bash
nexus-v init  # or: nexus-v i
```

NEXUS-V walks you through a series of prompts:

```
? Extension name (my-extension): My Extension
? Extension identifier (my-extension): 
? Description (A helpful VS Code extension): 
? Publisher (your-publisher-id): 
? Template variant (command): 
```

### Non-interactive mode

Pass all values as flags for scripting and CI use:

```bash
nexus-v init \
  --name "My Extension" \
  --id my-extension \
  --description "A helpful VS Code extension" \
  --publisher your-publisher-id \
  --variant command
```

### Custom Templates

Use a local directory as a template source:

```bash
nexus-v init --template-dir ./my-custom-templates
```

### List available template variants

```bash
nexus-v variants  # or: nexus-v ls / nexus-v vars
```

### Update NEXUS-V

```bash
nexus-v update  # or: nexus-v u
```

### Check Environment

```bash
nexus-v doctor  # or: nexus-v dr
```

### Dry run (preview without writing files)

```bash
nexus-v init --dry-run
```

---

## Template Variants

| Variant | Description |
|---|---|
| `command` | Basic extension with a registered command and activation event |
| `webview` | Extension with a webview panel boilerplate |
| `language` | Language support with syntax highlighting (tmLanguage.json) and config |
| `theme` | Color theme extension with a base theme JSON |

---

## Configuration

NEXUS-V can be configured at three levels, in order of precedence:

### 1. CLI flags (highest precedence)

```bash
nexus-v init --publisher my-org --variant webview --no-git
```

### 2. Environment variables

```bash
export NEXUSV_PUBLISHER="my-org"
export NEXUSV_DEFAULT_VARIANT="command"
export NEXUSV_TELEMETRY="on"
```

### 3. Configuration file (lowest precedence)

Place a `.nexusvrc.yaml` in your home directory or project root:

```yaml
# ~/.nexusvrc.yaml

defaults:
  publisher: "my-org"
  variant: "command"
  git: true
  license: "MIT"

telemetry:
  enabled: false

hooks:
  post_scaffold:
    - "npm install"
    - "code ."
```

---

## Hooks

NEXUS-V supports lifecycle hooks that run shell commands during scaffolding.

### Hook stages

| Stage | Trigger |
|---|---|
| `pre_scaffold` | After prompts resolve, before any files are written |
| `post_scaffold` | After all files are written to disk |

### Built-in hook shortcuts

```bash
nexus-v init --install    # runs "npm install" after scaffold
nexus-v init --open       # runs "code ." after scaffold
nexus-v init --git        # runs "git init" after scaffold
```

---

## Telemetry

NEXUS-V includes **optional, anonymous** usage telemetry. **Telemetry is off by default.**

### What is collected (when opted in)

- NEXUS-V version, OS and architecture
- Template variant selected
- Basename of the project directory (sanitized for privacy)

Telemetry respects the `DO_NOT_TRACK=1` environment variable.

---

## Project Structure

```
nexus-v/
├── cmd/
│   └── nexus-v/
│       └── main.go              # CLI entry point
├── internal/
│   ├── templates/
│   │   ├── templates.go         # Template engine and embedded FS
│   │   └── files/               # Embedded template files
│   ├── cli/
│   │   ├── cli.go               # Command handlers and aliases
│   │   └── spinner.go           # UX components
│   ├── config/
│   │   └── config.go            # YAML and Env configuration
│   ├── update/
│   │   └── update.go            # Self-update logic
│   └── version/
│       └── version.go           # Version tracking
└── README.md
```

---

## Next Steps & Roadmap

### Completed ✅
- [x] **CI/CD pipeline** — GitHub Actions for automated releases
- [x] **Cross-compilation matrix** — Prebuilt binaries for all major platforms
- [x] **`nexus-v update`** — Self-update from GitHub Releases
- [x] **Integration tests** — Go test suite for template validation
- [x] **Custom template support** — `--template-dir` flag

### Planned additions
- [ ] **Plugin system** — Community template packs
- [ ] **`nexus-v doctor`** — Environment diagnostic tool
- [ ] **Monorepo variant** — Multi-extension monorepo support

### Stretch goals
- [ ] **TUI mode** — Rich terminal UI (Bubble Tea)
- [ ] **Scaffold history** — `nexus-v history` / `nexus-v replay`
- [ ] **VS Code extension for NEXUS-V** — Meta extension

---

## License

MIT © [Geriatric Sailor](https://github.com/geriatric-sailor)