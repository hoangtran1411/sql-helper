package main

import (
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
}
