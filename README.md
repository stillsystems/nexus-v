# NEXUS-V

**A modern, zero-install VS Code extension scaffolder — built in Go, outputs TypeScript.**

NEXUS-V replaces the legacy Yeoman `yo code` generator with a single static binary that produces clean, modern VS Code extension projects. **The tool itself is dependency-free** — no Node.js or global packages required to run the scaffolder, and no hidden background state.

---

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Remote Plugins](#remote-templates-plugins)
- [Template Authoring](#template-authoring-plugins)
- [Template Variants](#template-variants)
- [Configuration](#configuration)
- [Hooks](#hooks)
- [Environment Health](#check-environment)
- [Telemetry](#telemetry)
- [Roadmap](#next-steps--roadmap)
- [License](#license)

---

## Overview

The VS Code extension ecosystem depends heavily on Yeoman — a scaffolding tool with deep Node.js dependency trees, global installs, and a maintenance surface area that drifts over time. NEXUS-V takes a different approach:

- **Written in Go** — compiles to a single static binary with zero runtime dependencies.
- **Outputs TypeScript** — every scaffolded project is a clean, modern TypeScript extension.
- **Interactive TUI** — uses **Bubble Tea** for a rich, visual template selection experience.
- **Portable by design** — no PATH pollution, no registry entries, and no hidden background state.
- **Offline-first** — all core templates are bundled; no internet required for local scaffolding.

---

## Features

| Feature | Description |
|---|---|
| **Single binary** | One executable, zero runtime dependencies — works on Windows, macOS, and Linux |
| **Interactive TUI** | Searchable, stylized menu for choosing template variants |
| **Remote Plugins** | Scaffold directly from any GitHub repository with optional pinning (`--template-ref`) |
| **Offline mode** | Works without an internet connection using built-in templates |
| **Doctor Command** | Diagnostic tool to verify your local environment (`node`, `npm`, `vsce`) |
| **GoReleaser Pipeline** | Automated builds and distribution via Homebrew, Scoop, and Winget |
| **Self-Update** | Built-in `update` command with package-manager awareness |
| **Hook system** | Pre- and post-scaffold hooks for custom automation (`--install`, `--git`) |
| **Opt-in telemetry** | Anonymous, minimal usage telemetry — off by default |

---

## Installation

### Homebrew (macOS / Linux)
```bash
brew tap billy-kidd-dev/nexusv
brew install nexus-v
```

### Scoop (Windows)
```powershell
scoop bucket add nexusv https://github.com/billy-kidd-dev/scoop-bucket
scoop install nexus-v
```

### Winget (Windows)
```powershell
winget install BillyKidd.NexusV
```

### One-liner (Unix)
```bash
curl -fsSL https://raw.githubusercontent.com/billy-kidd-dev/nexus-v/main/install.sh | bash
```

---

## Usage

### Interactive Mode (TUI)

```bash
nexus-v init  # or: nexus-v i
```

NEXUS-V will launch a beautiful interactive menu for selecting your extension type.

### Remote Templates (Plugins)

NEXUS-V supports scaffolding directly from remote Git repositories. This allows you to use community-created templates as plugins:

```bash
nexus-v init --template-dir https://github.com/user/my-custom-template
```

### Check Environment

Ensure you have all the necessary tools installed for VS Code development:

```bash
nexus-v doctor  # or: nexus-v dr
```

### Update NEXUS-V

```bash
nexus-v update  # or: nexus-v u
```

> [!NOTE]
> If you installed NEXUS-V via a package manager (Homebrew, Scoop, or Winget), it is recommended to update using that manager instead (e.g., `brew upgrade nexus-v`) to ensure your system state remains consistent.

### List Variants

List available local templates:

```bash
nexus-v list  # or: nexus-v ls
```

List variants available in a remote template:

```bash
nexus-v list --template-dir https://github.com/user/my-custom-template
```

---

## Command Options (`init`)

| Flag | Description |
|---|---|
| `--out <dir>` | Specify the output directory (defaults to extension ID) |
| `--variant <name>` | Select a specific template variant |
| `--template-dir <url>` | Use a remote Git repository (HTTPS, SSH, or `file://`) as a template |
| `--template-ref <ref>` | Pin a remote template to a specific branch, tag, or commit SHA |
| `--license <type>` | `MIT`, `Apache-2.0`, `GPL-3.0`, `BSD-3-Clause`, `Unlicense` |
| `--dry-run` | Preview the file structure without writing any files |
| `--force` | Overwrite existing files in the target directory |
| `--install` | Automatically run `npm install` after scaffolding |
| `--git` | Automatically run `git init` and create an initial commit |

---

## Template Authoring (Plugins)

NEXUS-V supports external templates via Git URLs. For a repository to be compatible with `nexus-v list` and `nexus-v init --variant`, it must follow this structure:

```text
my-template-repo/
├── files/
│   ├── default/          # Base template (fallback)
│   │   ├── package.json.tmpl
│   │   └── ...
│   ├── webview/          # A variant
│   │   └── ...
│   └── custom-variant/   # Another variant
└── ...
```

- The `files/` directory is mandatory for variant discovery.
- Use `.tmpl` suffix for files that require variable interpolation (Go template syntax).
- The `default` variant is used as a fallback if the requested variant is missing files.

### Available Template Variables

Authors can use the following variables in their `.tmpl` files:

| Variable | Description |
|---|---|
| `{{ .Name }}` | Human-readable project name |
| `{{ .Identifier }}` | Extension identifier (e.g., `my-extension`) |
| `{{ .Publisher }}` | Publisher ID |
| `{{ .Description }}` | Short project description |
| `{{ .License }}` | Selected license identifier |
| `{{ .Template }}` | Selected variant name |
| `{{ .CommandName }}` | Auto-generated command ID (e.g., `my-extension.helloWorld`) |
| `{{ currentYear }}` | Helper function to return the current year |

---

## Security & Trust

### Remote Templates
When using `--template-dir` with a remote URL, NEXUS-V clones the repository to a temporary directory. **Only use remote templates from sources you trust.** Remote templates can contain hooks that execute shell commands on your machine.

### Zero Runtime Dependencies
While the NEXUS-V binary itself requires no dependencies to run, the **scaffolded projects** are TypeScript-based and typically require Node.js and `npm` for development and compilation. Use `nexus-v doctor` to ensure your environment is ready for VS Code extension development.

| Variant | Description |
|---|---|
| `command` | Basic extension with a registered command and activation event |
| `webview` | Extension with a webview panel boilerplate |
| `language` | Language support with syntax highlighting and config |
| `theme` | Color theme extension with a base theme JSON |

---

## Configuration

NEXUS-V is "stateless" by default, but you can provide explicit configuration via a `.nexusvrc.yaml` file in your home directory or project root. This is the only form of persistence the tool acknowledges.

Place a `.nexusvrc.yaml` in your home directory or project root:

```yaml
defaults:
  publisher: "my-org"
  variant: "command"
  git: true
  license: "MIT"

hooks:
  post_scaffold:
    - "npm install"
    - "code ."
```

---

## Next Steps & Roadmap

### Completed ✅
- [x] **TUI Mode** — Rich terminal UI for variant selection
- [x] **Plugin System** — Remote Git template support
- [x] **`nexus-v doctor`** — Environment diagnostic tool
- [x] **Multi-Channel Distribution** — Homebrew, Scoop, and Winget
- [x] **Self-Update** — Built-in `update` command with package-manager awareness
- [x] **CI/CD Pipeline** — GoReleaser + GitHub Actions
- [x] **Plugin Discovery** — `nexus-v list` for remote repositories
- [x] **Pinned Templates** — Support for branch/tag/SHA refs

### Planned additions
- [ ] **Monorepo variant** — Multi-extension monorepo support
- [ ] **Scaffold history** — `nexus-v history` / `nexus-v replay`
- [ ] **VS Code Meta Extension** — Native UI wrapper for NEXUS-V

---

## License

MIT © [Billy Kidd](https://github.com/billy-kidd-dev)