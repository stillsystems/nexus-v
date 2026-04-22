package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"

	"nexus-v/internal/version"
)

const repo = "billy-kidd-dev/nexus-v"

type release struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name string `json:"name"`
		URL  string `json:"browser_download_url"`
	} `json:"assets"`
}

func CheckAndApply() error {
	exePath, _ := os.Executable()
	if isManagedByPkgManager(exePath) {
		fmt.Println("⚠️  WARNING: nexus-v appears to be managed by a package manager (Homebrew, Scoop, or Winget).")
		fmt.Println("   Updating via 'nexus-v update' may conflict with your package manager's state.")
		fmt.Println("   It is recommended to use your package manager to update (e.g., 'brew upgrade nexus-v').")
		fmt.Print("   Do you want to proceed anyway? (y/N): ")
		var answer string
		fmt.Scanln(&answer)
		if strings.ToLower(answer) != "y" {
			fmt.Println("Update cancelled.")
			return nil
		}
	}

	fmt.Printf("Checking for updates for %s...\n", repo)

	resp, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo))
	if err != nil {
		return fmt.Errorf("failed to fetch latest release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("github api returned status: %s", resp.Status)
	}

	var rel release
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return fmt.Errorf("failed to decode release metadata: %w", err)
	}

	current := version.Version
	if !strings.HasPrefix(current, "v") {
		current = "v" + current
	}
	latest := rel.TagName
	if !strings.HasPrefix(latest, "v") {
		latest = "v" + latest
	}

	if current == latest {
		fmt.Println("✔ You are already on the latest version (" + current + ").")
		return nil
	}

	fmt.Printf("✨ New version found: %s (Current: %s)\n", latest, current)

	var downloadURL string
	ext := ""
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}

	// Match pattern: nexus-v-<os>-<arch><ext>
	targetName := fmt.Sprintf("nexus-v-%s-%s%s", runtime.GOOS, runtime.GOARCH, ext)
	for _, asset := range rel.Assets {
		if asset.Name == targetName {
			downloadURL = asset.URL
			break
		}
	}

	if downloadURL == "" {
		return fmt.Errorf("could not find binary for %s/%s in release %s", runtime.GOOS, runtime.GOARCH, latest)
	}

	return applyUpdate(downloadURL)
}

// isManagedByPkgManager uses a best-effort heuristic to detect if the binary
// was installed via a package manager. Note that this can fail if the user
// uses custom prefixes or symlinks that don't follow standard patterns.
func isManagedByPkgManager(path string) bool {
	path = strings.ToLower(path)
	return strings.Contains(path, "homebrew") ||
		strings.Contains(path, "scoop") ||
		strings.Contains(path, "winget") ||
		strings.Contains(path, "cellar")
}

func applyUpdate(url string) error {
	fmt.Println("Downloading update...")

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download update: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download update: %s", resp.Status)
	}

	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Create temp file in the same directory as the executable
	tmpPath := exePath + ".new"
	tmpFile, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}

	success := false
	defer func() {
		tmpFile.Close()
		if !success {
			os.Remove(tmpPath)
		}
	}()

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		return fmt.Errorf("failed to save update: %w", err)
	}
	tmpFile.Close()

	// Windows-safe update strategy:
	// 1. Rename running nexus-v.exe to nexus-v.exe.old
	// 2. Rename nexus-v.exe.new to nexus-v.exe

	oldPath := exePath + ".old"
	_ = os.Remove(oldPath) // Remove previous backup if exists

	if err := os.Rename(exePath, oldPath); err != nil {
		return fmt.Errorf("failed to backup current version: %w. Try running with elevated permissions.", err)
	}

	if err := os.Rename(tmpPath, exePath); err != nil {
		// Restore backup if replacement fails
		_ = os.Rename(oldPath, exePath)
		return fmt.Errorf("failed to install new version: %w", err)
	}

	success = true
	fmt.Println("🚀 Update successful! Please restart nexus-v.")

	// Proactively suggest deleting the .old file or do it on next run?
	// We'll leave it for now.
	return nil
}
