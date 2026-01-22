// =============================================================================
// Tests for Auto-Update Module
// =============================================================================
// These tests cover:
// - Semantic version comparison (CompareVersions, parseVersion)
// - UpdateInfo struct validation
// - GitHubRelease struct validation
// - Configuration validation
//
// NOTE: Tests that require network (CheckForUpdate) are skipped.
// Integration testing should be done manually.
// =============================================================================

package main

import (
	"context"
	"strings"
	"testing"
)

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name     string
		v1       string
		v2       string
		expected bool
	}{
		{"newer major", "v2.0.0", "v1.0.0", true},
		{"newer minor", "v1.2.0", "v1.1.0", true},
		{"newer patch", "v1.1.2", "v1.1.1", true},
		{"same version", "v1.0.0", "v1.0.0", false},
		{"older major", "v1.0.0", "v2.0.0", false},
		{"older minor", "v1.0.0", "v1.1.0", false},
		{"older patch", "v1.0.0", "v1.0.1", false},
		{"without v prefix", "1.2.0", "1.1.0", true},
		{"mixed prefix", "v1.2.0", "1.1.0", true},
		{"partial version", "1.2", "1.1.0", true},
		{"dev version", "v1.0.0", "dev", true},
		{"both dev", "dev", "dev", false},
		{"v0 versions", "v0.2.0", "v0.1.0", true},
		{"large numbers", "v10.20.30", "v10.20.29", true},
		{"prerelease tag", "v1.0.0-beta", "v0.9.9", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CompareVersions(tt.v1, tt.v2)
			if result != tt.expected {
				t.Errorf("CompareVersions(%q, %q) = %v, want %v", tt.v1, tt.v2, result, tt.expected)
			}
		})
	}
}

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		expected [3]int
	}{
		{"full version", "1.2.3", [3]int{1, 2, 3}},
		{"major minor only", "1.2", [3]int{1, 2, 0}},
		{"major only", "1", [3]int{1, 0, 0}},
		{"empty string", "", [3]int{0, 0, 0}},
		{"with trailing text", "1.2.3-beta", [3]int{1, 2, 3}},
		{"zeros", "0.0.0", [3]int{0, 0, 0}},
		{"large numbers", "100.200.300", [3]int{100, 200, 300}},
		{"with extra parts", "1.2.3.4.5", [3]int{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseVersion(tt.version)
			if result != tt.expected {
				t.Errorf("parseVersion(%q) = %v, want %v", tt.version, result, tt.expected)
			}
		})
	}
}

func TestUpdateInfoDefaults(t *testing.T) {
	info := UpdateInfo{
		Available:  false,
		CurrentVer: CurrentVersion,
	}

	if info.Available {
		t.Error("UpdateInfo.Available should default to false")
	}
	if info.CurrentVer != CurrentVersion {
		t.Errorf("UpdateInfo.CurrentVer = %q, want %q", info.CurrentVer, CurrentVersion)
	}
	if info.LatestVer != "" {
		t.Errorf("UpdateInfo.LatestVer should be empty, got %q", info.LatestVer)
	}
	if info.DownloadURL != "" {
		t.Errorf("UpdateInfo.DownloadURL should be empty, got %q", info.DownloadURL)
	}
	if info.ReleaseURL != "" {
		t.Errorf("UpdateInfo.ReleaseURL should be empty, got %q", info.ReleaseURL)
	}
}

func TestUpdateInfoWithValues(t *testing.T) {
	info := UpdateInfo{
		Available:   true,
		CurrentVer:  "v1.0.0",
		LatestVer:   "v2.0.0",
		DownloadURL: "https://example.com/download",
		ReleaseURL:  "https://example.com/release",
	}

	if !info.Available {
		t.Error("UpdateInfo.Available should be true")
	}
	if info.CurrentVer != "v1.0.0" {
		t.Errorf("UpdateInfo.CurrentVer = %q, want %q", info.CurrentVer, "v1.0.0")
	}
	if info.LatestVer != "v2.0.0" {
		t.Errorf("UpdateInfo.LatestVer = %q, want %q", info.LatestVer, "v2.0.0")
	}
}

func TestGitHubReleaseStruct(t *testing.T) {
	release := GitHubRelease{
		TagName: "v1.0.0",
		HTMLURL: "https://github.com/user/repo/releases/v1.0.0",
	}

	if release.TagName != "v1.0.0" {
		t.Errorf("GitHubRelease.TagName = %q, want %q", release.TagName, "v1.0.0")
	}
	if release.HTMLURL != "https://github.com/user/repo/releases/v1.0.0" {
		t.Errorf("GitHubRelease.HTMLURL = %q", release.HTMLURL)
	}
	if len(release.Assets) != 0 {
		t.Error("GitHubRelease.Assets should be empty")
	}
}

func TestGitHubConstants(t *testing.T) {
	// Skip if constants are still template placeholders
	if strings.HasPrefix(GitHubOwner, "{{") {
		t.Skip("GitHubOwner is still a template placeholder - update before testing")
	}
	if strings.HasPrefix(GitHubRepo, "{{") {
		t.Skip("GitHubRepo is still a template placeholder - update before testing")
	}

	if GitHubOwner == "" {
		t.Error("GitHubOwner should not be empty")
	}
	if GitHubRepo == "" {
		t.Error("GitHubRepo should not be empty")
	}
}

func TestCurrentVersionDefault(t *testing.T) {
	// CurrentVersion should have a default value
	if CurrentVersion == "" {
		t.Error("CurrentVersion should not be empty")
	}
}

func TestGetCurrentVersion(t *testing.T) {
	app := NewApp()
	app.startup(context.Background())

	version := app.GetCurrentVersion()
	if version != CurrentVersion {
		t.Errorf("GetCurrentVersion() = %q, want %q", version, CurrentVersion)
	}
}

func TestCompareVersionsEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		v1       string
		v2       string
		expected bool
	}{
		{"empty vs empty", "", "", false},
		{"v vs empty", "v", "", false},
		{"empty vs v1.0.0", "", "v1.0.0", false},
		{"dots only", "...", "1.0.0", false},
		{"letters", "abc", "1.0.0", false},
		{"negative", "-1.0.0", "0.0.0", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CompareVersions(tt.v1, tt.v2)
			if result != tt.expected {
				t.Errorf("CompareVersions(%q, %q) = %v, want %v", tt.v1, tt.v2, result, tt.expected)
			}
		})
	}
}
