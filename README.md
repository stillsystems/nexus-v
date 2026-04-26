![Still Systems Banner](https://raw.githubusercontent.com/stillsystems/.github/main/profile/banner.png)

# ⚓ NEXUS-V

**Modern developer tooling engineered for real-world conditions.**  
Zero-install. Dependency-light. Predictable.

---

## 🧭 Philosophy

NEXUS-V is the flagship project of the **Still Systems** ecosystem and embodies our core principles:
- 🛡️ **Clarity over cleverness** — predictable behavior, no magic.
- 📦 **Portability over complexity** — single static binaries, zero runtime dependencies.
- ⚓ **Durability over trends** — built for long-term maintainability with minimal dependencies.

---

## 🚀 Features

| Feature | Description |
|---|---|
| **Zero-Install** | One executable, zero runtime dependencies — no Node.js required to run |
| **Interactive TUI** | Beautiful terminal UI for choosing template variants |
| **Remote Plugins** | Scaffold directly from any GitHub repository with pinning support |
| **Offline-First** | All core templates are bundled; no internet required for local use |
| **Doctor Command** | Diagnostic tool to verify your environment health |
| **Multi-Platform** | Native support for Windows, macOS, and Linux |
| **Secure Hooks** | Pre- and post-generation hooks with safety-first warnings |

---

## 📦 Installation

### **Homebrew (macOS / Linux)**
```bash
brew tap stillsystems/nexusv
brew install nexus-v
```

### **Scoop (Windows)**
```powershell
scoop bucket add stillsystems https://github.com/stillsystems/scoop-bucket
scoop install nexus-v
```

### **Winget (Windows)**
```powershell
winget install stillsystems.nexusv
```

---

## 🛠️ Usage

### **Interactive Mode**
```bash
nexus-v init  # or simply: nexus-v i
```

### **Remote Templates**
```bash
nexus-v init --template-dir https://github.com/user/my-template
```

### **Check Health**
```bash
nexus-v doctor
```

---

## 🔧 Command Options (`init`)

| Flag | Description |
|---|---|
| `--out <dir>` | Specify the output directory |
| `--variant <name>`| Select a specific template variant |
| `--license <type>`| `MIT`, `Apache-2.0`, `GPL-3.0`, `BSD-3-Clause`, `Unlicense` |
| `--dry-run` | Preview files without writing them |
| `--force` | Overwrite existing files |
| `--install` | Automatically run `npm install` after scaffold |
| `--git` | Automatically run `git init` after scaffold |

---

## 🤝 Contributing

We welcome contributions! Please see the [Still Systems Contribution Rules](https://github.com/stillsystems/.github/blob/main/brand/governance/contribution-rules.md) for our standards on clarity and dependency management.

---

⚓ **Still Systems** — Tools engineered for real-world conditions.  
MIT © [Still Systems](https://github.com/stillsystems)
