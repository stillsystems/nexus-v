# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.8] - 2026-05-01
### Fixed
- Restored cross-repo publishing for Homebrew, Scoop, and Winget.
- Finalized CI/CD workflow with correct permissions.

## [0.2.5] - 2026-04-27
### Fixed
- Release pipeline authentication and token scopes.

## [0.2.4] - 2026-04-27
### Added
- SBOM (Software Bill of Materials) generation using Syft.

## [0.2.3] - 2026-04-26
### Added
- Interactive **License** prompt during `init` (defaults to MIT).
- Brand-new **Launch Demo GIF** showcasing the full scaffolding flow.
- Modernized `package.json` templates (removed redundant `activationEvents`).
### Changed
- Rebranded from "Zero-install" to **"Zero runtime dependencies"** across all documentation and manifests for better technical accuracy.
- Refined terminal pacing in the launch recording.

## [0.2.2] - 2026-04-22
### Changed
- Transferred repository and all distribution channels to the **Still Systems** organization.
- Updated internal update logic, installation scripts, and package manager manifests to point to the new organization.
- Rebranded Winget publisher and package identifier to `stillsystems`.

## [0.2.1] - 2026-04-22
### Fixed
- Cleaned up template "bloat": variants like `theme` and `language` no longer include unwanted `src/` boilerplate from the default template.
- Fixed broken installation URL in `install.ps1`.
- Filtered internal implementation templates (`default`, `.vscode`) from the `list` command output.
### Changed
- Standardized README examples to use `./nexus-v` for easier local execution and copy-paste.
- Improved TUI variant selection to be dynamic based on available templates.
- Simplified `dist` binary naming for better ease of use.
### Added
- Local directory support for the `list` command (`--template-dir <path>`).

## [0.2.0] - 2026-04-22
### Added
- Remote template variant discovery (`nexus-v list --template-dir <URL>`)
- Template pinning support with `--template-ref` (branch, tag, or SHA)
- License validation and automatic `LICENSE` file generation
- Template authoring guide and variable reference in README
- Interactive TUI search and filtering improvements
### Changed
- Upgraded remote cloning strategy to handle SHA pinning reliably
- Improved update command with package-manager awareness (Homebrew, Scoop, Winget)
- Clarified "zero runtime dependency" and "stateless" claims in documentation
### Fixed
- Cleaned up temporary clone directories on all error paths

## [0.1.0] - 2026-04-21
### Added
- Initial release of NEXUS-V
- Template variants system (command, webview, language, theme)
- Dynamic template rendering with Go `text/template`
- Filesystem walker for embedded and local templates
- Interactive CLI prompts with `bufio`
- YAML-based configuration (`.nexusvrc.yaml`) with Environment Variable support
- Self-update functionality via GitHub Releases
- Multi-platform support (Windows, Linux, macOS)
- Opt-in telemetry system
- GitHub Actions for CI/CD and automated releases
