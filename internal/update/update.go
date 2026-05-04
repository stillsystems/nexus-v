package update

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/stillsystems/nexus-v/internal/version"
)

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

func validateURL(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	allowedHosts := []string{
		"api.github.com",
		"github.com",
		"objects.githubusercontent.com",
	}

	for _, host := range allowedHosts {
		if u.Host == host {
			return nil
		}
	}

	return fmt.Errorf("untrusted download host: %s", u.Host)
}

func verifyChecksum(data []byte, url string, targetName string) error {
	if err := validateURL(url); err != nil {
		return err
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download checksums: %w", err)
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read checksums: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	var expectedHash string
	for _, line := range lines {
		if strings.Contains(line, targetName) {
			fields := strings.Fields(line)
			if len(fields) > 0 {
				expectedHash = fields[0]
				break
			}
		}
	}

	if expectedHash == "" {
		return fmt.Errorf("no checksum found for %s", targetName)
	}

	sum := sha256.Sum256(data)
	actualHash := hex.EncodeToString(sum[:])

	if actualHash != expectedHash {
		return fmt.Errorf("hash mismatch: expected %s, got %s", expectedHash, actualHash)
	}

	return nil
}

const repo = "stillsystems/nexus-v"

type release struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name string `json:"name"`
		URL  string `json:"browser_download_url"`
	} `json:"assets"`
}

func CheckAndApply() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to determine executable path: %w", err)
	}
	if isManagedByPkgManager(exePath) {
		fmt.Println("⚠️  WARNING: nexus-v appears to be managed by a package manager (Homebrew, Scoop, or Winget).")
		fmt.Println("   Updating via 'nexus-v update' may conflict with your package manager's state.")
		fmt.Println("   It is recommended to use your package manager to update (e.g., 'brew upgrade nexus-v').")
		fmt.Print("   Do you want to proceed anyway? (y/N): ")
		var answer string
		if _, err := fmt.Scanln(&answer); err != nil && err != io.EOF {
			return fmt.Errorf("failed to read user input: %w", err)
		}
		if strings.ToLower(answer) != "y" {
			fmt.Println("Update cancelled.")
			return nil
		}
	}

	fmt.Printf("Checking for updates for %s...\n", repo)

	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	if err := validateURL(apiURL); err != nil {
		return err
	}

	resp, err := httpClient.Get(apiURL)
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

	targetName := fmt.Sprintf("nexus-v_%s_%s%s", runtime.GOOS, runtime.GOARCH, ext)
	checksumName := "checksums.txt"
	var checksumURL string

	for _, asset := range rel.Assets {
		if asset.Name == targetName {
			downloadURL = asset.URL
		}
		if asset.Name == checksumName {
			checksumURL = asset.URL
		}
	}

	if downloadURL == "" {
		return fmt.Errorf("could not find binary for %s/%s in release %s", runtime.GOOS, runtime.GOARCH, latest)
	}
	if checksumURL == "" {
		return fmt.Errorf("could not find checksums.txt in release %s", latest)
	}

	return applyUpdate(downloadURL, checksumURL, targetName)
}

func isManagedByPkgManager(path string) bool {
	path = strings.ToLower(path)
	return strings.Contains(path, "homebrew") ||
		strings.Contains(path, "scoop") ||
		strings.Contains(path, "winget") ||
		strings.Contains(path, "cellar")
}

func applyUpdate(url, checksumURL, targetName string) error {
	fmt.Println("Downloading update...")

	if err := validateURL(url); err != nil {
		return err
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download update: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download update: %s", resp.Status)
	}

	binaryData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read update data: %w", err)
	}

	fmt.Println("Verifying checksum...")
	if err := verifyChecksum(binaryData, checksumURL, targetName); err != nil {
		return fmt.Errorf("checksum verification failed: %w", err)
	}

	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	tmpPath := exePath + ".new"
	tmpFile, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}

	success := false
	defer func() {
		_ = tmpFile.Close()
		if !success {
			_ = os.Remove(tmpPath)
		}
	}()

	if _, err := tmpFile.Write(binaryData); err != nil {
		return fmt.Errorf("failed to save update: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to finalize update file: %w", err)
	}

	oldPath := exePath + ".old"
	_ = os.Remove(oldPath)

	if err := os.Rename(exePath, oldPath); err != nil {
		return fmt.Errorf("failed to backup current version: %w. Try running with elevated permissions.", err)
	}

	if err := os.Rename(tmpPath, exePath); err != nil {
		_ = os.Rename(oldPath, exePath)
		return fmt.Errorf("failed to install new version: %w", err)
	}

	success = true
	fmt.Println("🚀 Update successful! Please restart nexus-v.")

	return nil
}

