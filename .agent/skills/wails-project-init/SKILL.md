# Wails Project Init

---
name: wails-project-init
description: Initialize a production-ready Go + Wails desktop application with CI/CD, GitHub Release automation, and auto-update functionality.
---

## Overview

This skill sets up a complete Go + Wails project with:
- **CI Pipeline**: Lint, Test, Build on all platforms (Linux, Windows, macOS)
- **Release Workflow**: Automated GitHub Releases with multi-platform binaries
- **Auto-Update**: In-app update checking and installation from GitHub Releases
- **Makefile**: Standard development commands
- **golangci-lint**: Configured linting

## Prerequisites

- Go 1.23+ installed
- Node.js 20+ installed
- Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- golangci-lint: `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`

---

## Step-by-Step Instructions

### 1. Create Wails Project (if new)

```bash
wails init -n <project-name> -t vanilla
cd <project-name>
```

### 2. Set GitHub Repository Info

Update constants in `updater.go` (or create it):

```go
const (
    GitHubOwner = "<your-github-username>"
    GitHubRepo  = "<your-repo-name>"
)
```

### 3. Copy Template Files

Copy the following files from `templates/` directory into your project:

| Template File | Target Location |
|---------------|-----------------|
| `ci.yml` | `.github/workflows/ci.yml` |
| `release.yml` | `.github/workflows/release.yml` |
| `golangci.yml` | `.golangci.yml` |
| `Makefile` | `Makefile` |
| `updater.go` | `updater.go` |
| `updater_test.go` | `updater_test.go` |

### 4. Update Project-Specific Values

After copying, update these placeholders:

| Placeholder | Replace With |
|-------------|--------------|
| `{{PROJECT_NAME}}` | Your project name (e.g., `sql-helper`) |
| `{{GITHUB_OWNER}}` | Your GitHub username |
| `{{GITHUB_REPO}}` | Your repository name |
| `{{GO_VERSION}}` | Go version (default: `1.23`) |

### 5. Integrate Updater with App

In your `app.go`, ensure the App struct has these methods available:
- `GetCurrentVersion() string`
- `CheckForUpdate() UpdateInfo`
- `PerformUpdate(downloadURL string) (bool, error)`
- `OpenReleaseURL(url string)`

### 6. Frontend Integration

Add update UI to your frontend. Example JavaScript:

```javascript
// Check for updates on app load
async function checkForUpdates() {
    const updateInfo = await window.go.main.App.CheckForUpdate();
    if (updateInfo.available) {
        // Show update notification
        showUpdateButton(updateInfo);
    }
}

// Perform update
async function performUpdate(downloadUrl) {
    await window.go.main.App.PerformUpdate(downloadUrl);
}
```

### 7. Build with Version Injection

Use ldflags to inject version at build time:

```bash
# Development
wails build

# Release (version from git tag)
wails build -ldflags "-s -w -X main.CurrentVersion=v1.0.0"
```

---

## CI/CD Pipeline Details

### CI Workflow (ci.yml)

Triggers on: `push` and `pull_request` to `main`/`master`

| Job | Description |
|-----|-------------|
| **lint** | Runs `golangci-lint` on Ubuntu |
| **test** | Runs tests with coverage on all 3 platforms |
| **build** | Builds Wails app for Linux, Windows, macOS |

Coverage threshold: **70%** (configurable)

### Release Workflow (release.yml)

Triggers on: `push` tags matching `v*`

| Step | Output |
|------|--------|
| Build Linux | `<project>-linux-amd64.tar.gz` |
| Build Windows | `<project>-windows-amd64.zip` |
| Build macOS | `<project>-darwin-universal.tar.gz` |
| Create Release | GitHub Release with all assets |

---

## Auto-Update Flow

```
┌─────────────────┐
│  App Startup    │
└────────┬────────┘
         ▼
┌─────────────────────────────────┐
│  CheckForUpdate()               │
│  GET /repos/owner/repo/releases │
└────────┬────────────────────────┘
         ▼
┌─────────────────────────────────┐
│  Compare CurrentVersion vs      │
│  release.tag_name               │
└────────┬────────────────────────┘
         ▼
   ┌─────┴─────┐
   │ Newer?    │
   └─────┬─────┘
    Yes  │  No
    ▼    ▼
┌────────┐ ┌──────────┐
│Show UI │ │ Hide UI  │
└───┬────┘ └──────────┘
    ▼
┌─────────────────────────────────┐
│  PerformUpdate(downloadURL)     │
│  1. Download ZIP to temp        │
│  2. Create update batch script  │
│  3. Start script & quit app     │
│  4. Script replaces exe         │
│  5. Script restarts app         │
└─────────────────────────────────┘
```

---

## Makefile Commands

| Command | Description |
|---------|-------------|
| `make dev` | Run in development mode |
| `make build` | Build production binary |
| `make test` | Run all tests |
| `make lint` | Run golangci-lint |
| `make coverage` | Run tests with coverage report |
| `make fmt` | Format code with gofmt/goimports |
| `make check` | Run lint + test |
| `make clean` | Remove build artifacts |

---

## Common Issues & Fixes

### 1. golangci-lint gofmt errors

```bash
gofmt -w .
```

### 2. gosimple: nil slice check

```go
// Bad
if slice != nil && len(slice) != 0 { }

// Good
if len(slice) != 0 { }
```

### 3. Windows race detector not working

Windows CI uses `-race` flag only on Unix. Already handled in `ci.yml`.

### 4. Version shows "dev"

Ensure `-ldflags "-X main.CurrentVersion=vX.X.X"` is passed at build time.

---

## Checklist for New Project

- [ ] Create Wails project
- [ ] Copy template files
- [ ] Update `GitHubOwner` and `GitHubRepo` constants
- [ ] Integrate updater methods in `app.go`
- [ ] Add update UI to frontend
- [ ] Create initial commit
- [ ] Push to GitHub
- [ ] Create first release tag: `git tag v0.1.0 && git push origin v0.1.0`
- [ ] Verify CI passes
- [ ] Verify Release workflow creates assets
