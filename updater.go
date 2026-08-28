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
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// CurrentVersion is the application version
// This is set at build time via ldflags: -X main.CurrentVersion=v1.x.x
// If not set at build time, defaults to "dev"
var CurrentVersion = "dev"

// GitHub repository info
const (
	GitHubOwner = "hoangtran1411"
	GitHubRepo  = "sql-helper"
)

// UpdateInfo holds information about available updates
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
func (a *App) GetCurrentVersion() string {
	return CurrentVersion
}

// CheckForUpdate checks GitHub for newer versions
func (a *App) CheckForUpdate() UpdateInfo {
	info := UpdateInfo{
		Available:  false,
		CurrentVer: CurrentVersion,
	}

	// Call GitHub API
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", GitHubOwner, GitHubRepo)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(url)
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

	// Find Windows exe asset (direct EXE release)
	for _, asset := range release.Assets {
		assetName := strings.ToLower(asset.Name)
		if strings.Contains(assetName, "windows") && strings.HasSuffix(assetName, ".exe") {
			info.DownloadURL = asset.BrowserDownloadURL
			break
		}
	}

	// Compare versions
	if info.LatestVer != "" && CompareVersions(info.LatestVer, CurrentVersion) {
		info.Available = true
	}

	return info
}

// CompareVersions returns true if v1 is newer than v2
func CompareVersions(v1, v2 string) bool {
	v1 = strings.TrimPrefix(v1, "v")
	v2 = strings.TrimPrefix(v2, "v")

	parts1 := parseVersion(v1)
	parts2 := parseVersion(v2)

	for i := range 3 {
		if parts1[i] > parts2[i] {
			return true
		}
		if parts1[i] < parts2[i] {
			return false
		}
	}
	return false
}

// parseVersion splits version string into [major, minor, patch] integers
func parseVersion(v string) [3]int {
	var result [3]int
	parts := strings.Split(v, ".")
	for i := 0; i < len(parts) && i < 3; i++ {
		//nolint:errcheck // default to 0 on parse failure
		fmt.Sscanf(parts[i], "%d", &result[i])
	}
	return result
}

// PerformUpdate downloads and installs the new version (Windows only)
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

	// Create temp file for download
	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, "sql_helper_update.exe")

	// Emit progress event
	application.Get().Event.Emit("updateProgress", "Downloading update...")

	// Download new version
	client := &http.Client{
		Timeout: 30 * time.Minute,
	}
	resp, err := client.Get(downloadURL)
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

	_, err = io.Copy(out, resp.Body)
	out.Close()
	if err != nil {
		return false, fmt.Errorf("failed to save update: %w", err)
	}

	application.Get().Event.Emit("updateProgress", "Installing update...")

	// Create update batch script
	batchPath := filepath.Join(tempDir, "update_sql_helper.bat")
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

	// Run the batch script (hidden)
	cmd := exec.Command("cmd", "/c", "start", "/min", "", batchPath)
	if err := cmd.Start(); err != nil {
		return false, fmt.Errorf("failed to start update script: %w", err)
	}

	// Quit the application
	application.Get().Quit()

	return true, nil
}

// OpenReleaseURL opens the release page in the default browser
func (a *App) OpenReleaseURL(url string) {
	_ = application.Get().Browser.OpenURL(url)
}
