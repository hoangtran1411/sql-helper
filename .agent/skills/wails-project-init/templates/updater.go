// =============================================================================
// Auto-Update Module for Wails Desktop Applications
// =============================================================================
// This module provides:
// - Version checking against GitHub Releases API
// - Semantic version comparison (major.minor.patch)
// - Windows self-update via batch script (download, replace, restart)
// - Progress events for frontend UI
//
// SETUP:
// 1. Replace {{GITHUB_OWNER}} with your GitHub username
// 2. Replace {{GITHUB_REPO}} with your repository name
// 3. Replace {{PROJECT_NAME}} with your project name (for temp files)
//
// BUILD:
// Use ldflags to inject version: -ldflags "-X main.CurrentVersion=v1.0.0"
//
// FRONTEND EVENTS:
// - Listen for "updateProgress" event to show download/install status
// =============================================================================

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// CurrentVersion is the application version
// Set at build time via ldflags: -X main.CurrentVersion=v1.x.x
// Defaults to "dev" for development builds
var CurrentVersion = "dev"

// GitHub repository configuration
// TODO: Update these constants with your GitHub repository info
const (
	GitHubOwner = "{{GITHUB_OWNER}}"
	GitHubRepo  = "{{GITHUB_REPO}}"
)

// UpdateInfo holds information about available updates
// Exported to frontend via JSON
type UpdateInfo struct {
	Available   bool   `json:"available"`
	CurrentVer  string `json:"currentVersion"`
	LatestVer   string `json:"latestVersion"`
	DownloadURL string `json:"downloadUrl"`
	ReleaseURL  string `json:"releaseUrl"`
}

// GitHubRelease represents a GitHub release API response
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// GetCurrentVersion returns the current app version
// Exposed to frontend for display
func (a *App) GetCurrentVersion() string {
	return CurrentVersion
}

// CheckForUpdate checks GitHub for newer versions
// Returns UpdateInfo with availability status and download URL
func (a *App) CheckForUpdate() UpdateInfo {
	info := UpdateInfo{
		Available:  false,
		CurrentVer: CurrentVersion,
	}

	// Call GitHub Releases API
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", GitHubOwner, GitHubRepo)
	resp, err := http.Get(url)
	if err != nil {
		return info
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return info
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return info
	}

	info.LatestVer = release.TagName
	info.ReleaseURL = release.HTMLURL

	// Find Windows zip asset
	// Release workflow creates: {{PROJECT_NAME}}-windows-amd64.zip
	for _, asset := range release.Assets {
		assetName := strings.ToLower(asset.Name)
		if strings.Contains(assetName, "windows") && strings.HasSuffix(assetName, ".zip") {
			info.DownloadURL = asset.BrowserDownloadURL
			break
		}
	}

	// Compare versions using semantic versioning
	if info.LatestVer != "" && CompareVersions(info.LatestVer, CurrentVersion) {
		info.Available = true
	}

	return info
}

// CompareVersions returns true if v1 is newer than v2
// Uses semantic versioning: major.minor.patch
func CompareVersions(v1, v2 string) bool {
	// Remove 'v' prefix
	v1 = strings.TrimPrefix(v1, "v")
	v2 = strings.TrimPrefix(v2, "v")

	// Parse version parts
	parts1 := parseVersion(v1)
	parts2 := parseVersion(v2)

	// Compare major, minor, patch in order
	for i := 0; i < 3; i++ {
		if parts1[i] > parts2[i] {
			return true
		}
		if parts1[i] < parts2[i] {
			return false
		}
	}
	return false // Equal versions
}

// parseVersion splits version string into [major, minor, patch] integers
func parseVersion(v string) [3]int {
	var result [3]int
	parts := strings.Split(v, ".")

	for i := 0; i < len(parts) && i < 3; i++ {
		// Parse integer, ignore errors (defaults to 0)
		//nolint:errcheck // intentionally ignore parse errors, default to 0
		fmt.Sscanf(parts[i], "%d", &result[i])
	}
	return result
}

// PerformUpdate downloads and installs the new version (Windows only)
// Uses a batch script approach:
// 1. Download EXE to temp directory
// 2. Create batch script that waits for app to exit
// 3. Batch deletes old exe, moves new exe in place
// 4. Batch restarts app and deletes itself
func (a *App) PerformUpdate(downloadURL string) (bool, error) {
	if downloadURL == "" {
		return false, fmt.Errorf("no download URL provided")
	}

	// Get current executable path
	exePath, err := os.Executable()
	if err != nil {
		return false, fmt.Errorf("failed to get executable path: %w", err)
	}
	exePath, err = filepath.Abs(exePath)
	if err != nil {
		return false, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Create temp path for download
	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, "{{PROJECT_NAME}}_update.exe")

	// Emit progress event to frontend
	runtime.EventsEmit(a.ctx, "updateProgress", "Downloading update...")

	// Download new version
	resp, err := http.Get(downloadURL)
	if err != nil {
		return false, fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false, fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	// Create output file
	out, err := os.Create(tempFile)
	if err != nil {
		return false, fmt.Errorf("failed to create temp file: %w", err)
	}

	// Download
	_, err = io.Copy(out, resp.Body)
	out.Close()
	if err != nil {
		return false, fmt.Errorf("failed to save update: %w", err)
	}

	runtime.EventsEmit(a.ctx, "updateProgress", "Installing update...")

	// Create update batch script
	// This runs after the app exits and:
	// 1. Waits 2 seconds for app to fully exit
	// 2. Deletes old executable
	// 3. Moves new executable to original location
	// 4. Starts new version
	// 5. Deletes itself
	batchPath := filepath.Join(tempDir, "update_{{PROJECT_NAME}}.bat")
	batchContent := fmt.Sprintf(`@echo off
timeout /t 2 /nobreak >nul
del "%s"
move /y "%s" "%s"
start "" "%s"
del "%%~f0"
`, exePath, tempFile, exePath, exePath)

	if err := os.WriteFile(batchPath, []byte(batchContent), 0644); err != nil {
		return false, fmt.Errorf("failed to create update script: %w", err)
	}

	// Run the batch script (hidden/minimized)
	cmd := exec.Command("cmd", "/c", "start", "/min", "", batchPath)
	if err := cmd.Start(); err != nil {
		return false, fmt.Errorf("failed to start update script: %w", err)
	}

	// Quit the application to allow batch to replace exe
	runtime.Quit(a.ctx)

	return true, nil
}

// OpenReleaseURL opens the release page in the default browser
// Use this as a fallback if auto-update fails
func (a *App) OpenReleaseURL(url string) {
	runtime.BrowserOpenURL(a.ctx, url)
}
