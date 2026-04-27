# <picture><source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/stillsystems/.github/main/brand/logo.png"><img alt="Still Systems" src="https://raw.githubusercontent.com/stillsystems/.github/main/brand/logo.png" width="32" height="32"></picture> 🏗️ NEXUS-V

**Modern developer tooling engineered for real-world conditions.**  
Zero runtime dependencies. Dependency-light. Predictable.

[![CI](https://github.com/stillsystems/nexus-v/actions/workflows/ci.yml/badge.svg)](https://github.com/stillsystems/nexus-v/actions)
[![License](https://img.shields.io/github/license/stillsystems/nexus-v?style=flat-square&color=111827)](LICENSE)
[![Release](https://img.shields.io/github/v/release/stillsystems/nexus-v?style=flat-square&color=111827)](https://github.com/stillsystems/nexus-v/releases)

![NEXUS-V Launch Demo](https://vhs.charm.sh/vhs-DEqFCb9FeLANXwRzxNFqw.gif)

## Overview

NEXUS-V is the flagship project of the **Still Systems** ecosystem. It is a lightweight, high-utility scaffolding engine designed to provide software that "just works"—allowing you to focus on your build rather than troubleshooting your tools.

It embodies our core principles:
- 🛡️ **Clarity over cleverness** — predictable behavior, no magic.
- 📦 **Portability over complexity** — single static binaries, zero runtime dependencies.
- ⚓ **Durability over trends** — built for long-term maintainability.

## Installation / Quickstart

### Homebrew (macOS / Linux)
```bash
brew tap stillsystems/nexusv
brew install nexus-v
```

### Scoop (Windows)
```powershell
scoop bucket add stillsystems https://github.com/stillsystems/scoop-bucket
scoop install nexus-v
```

### WinGet (Windows)
```powershell
winget install stillsystems.nexusv
```

## Usage

### Interactive Mode
```bash
nexus-v init  # or simply: nexus-v i
```

### Remote Templates
```bash
nexus-v init --template-dir https://github.com/user/my-template
```

### Check Health
```bash
nexus-v doctor
```

## Configuration

| Flag | Description |
|---|---|
| `--out <dir>` | Specify the output directory |
| `--variant <name>`| Select a specific template variant |
| `--license <type>`| `MIT`, `Apache-2.0`, `GPL-3.0`, `BSD-3-Clause`, `Unlicense` |
| `--dry-run` | Preview files without writing them |
| `--force` | Overwrite existing files |
| `--install` | Automatically run `npm install` after scaffold |
| `--git` | Automatically run `git init` after scaffold |

## Contributing

Please refer to our [Global Contributing Guidelines](https://github.com/stillsystems/.github/blob/main/CONTRIBUTING.md).

## License

This project is licensed under the MIT License.

## Roadmap

We are committed to the long-term stability of the Still Systems ecosystem. See our [Roadmap](ROADMAP.md) for planned features and architectural goals.

Documentation • [Issues](https://github.com/stillsystems/nexus-v/issues) • [Support](https://github.com/stillsystems/.github/blob/main/SUPPORT.md)

---
🧱 **Still Systems** — Tools engineered for real-world conditions.
