# 🏗️ Nexus-V Roadmap

This document outlines the planned trajectory for **Nexus-V**. Our priority is always **Stability** and **Structural Clarity**, but we are excited to expand the engine's capabilities.

## 🟢 Phase 1: Foundation (Current)
*   [x] Single-binary Go execution engine.
*   [x] Zero runtime dependency architecture.
*   [x] Interactive TUI and CLI prompt system.
*   [x] Support for remote Git-based templates.
*   [x] Multi-platform distribution (WinGet, Scoop, Homebrew).

## 🟡 Phase 2: Refinement (Q2 2026)
*   **Template Plugin System**: Allow templates to run lightweight pre/post-scaffold scripts.
*   **Variable Persistence**: Save common configuration (Publisher, Author) in a global config file to skip redundant prompts.
*   **Improved Validation**: Deep-linting of `package.json` and `tsconfig.json` during template rendering.
*   **NPM Wrapper**: A lightweight shim for users who prefer `npx nexus-v`.

## 🟠 Phase 3: Expansion (Q3 2026)
*   **Language-Specific Optimizations**: Native support for Go and Rust extension templates.
*   **Visual Scaffolders**: A minimal web-based UI for users who prefer a graphical interface.
*   **Template Gallery**: An official index of community-contributed Still Systems templates.

---

🧱 **Still Systems** — Modern developer tooling engineered for real-world conditions.
