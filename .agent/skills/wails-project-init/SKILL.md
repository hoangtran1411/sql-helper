---
name: wails-project-init
description: Initialize a production-ready Go + Wails v3 desktop application with CI/CD, GitHub Release automation, and auto-update functionality.
---

# Wails Project Init (Wails v3)

## Overview

This skill sets up a complete Go + **Wails v3** desktop project with:
- **CI Pipeline**: Lint, Test with 70% coverage gate on `internal/`, Build on all platforms (Linux, Windows, macOS)
- **Release Workflow**: Automated GitHub Releases with Windows direct single-executable deployment
- **Auto-Update**: In-app update checking and self-installation from GitHub Releases via `application.Get()`
- **Makefile**: Development commands with `wails3 dev`, `wails3 build`, `wails3 generate bindings`
- **golangci-lint**: Configured cross-platform linting

## Prerequisites

- Go 1.27+ installed
- Node.js 20+ / Node 24 installed
- Wails v3 CLI: `go install github.com/wailsapp/wails/v3/cmd/wails3@latest`
- golangci-lint: `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`

---

## Step-by-Step Instructions

### 1. Create Wails v3 Project (if new)

```bash
wails3 init -n <project-name>
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
| `go-style-guide.md` | `.agent/rules/go-style-guide.md` |
| `go-idioms-reference.md` | `.agent/rules/go-idioms-reference.md` |

### 4. Update Project-Specific Values

After copying, update these placeholders:

| Placeholder | Replace With | Used In |
|-------------|--------------|----------|
| `{{PROJECT_NAME}}` | Your project name (e.g., `sql-helper`) | Both style guide files, Makefile, CI/Release workflows, updater.go |
| `{{PROJECT_DESCRIPTION}}` | Brief project description | go-style-guide.md |
| `{{DOMAIN}}` | Main domain package name (e.g., `excel`, `sql`) | go-style-guide.md |
| `{{GITHUB_OWNER}}` | Your GitHub username | updater.go |
| `{{GITHUB_REPO}}` | Your repository name | updater.go |

### 5. Integrate Updater with App

In your `app.go` or main service, ensure the App struct methods are available:
- `GetCurrentVersion() string`
- `CheckForUpdate() UpdateInfo`
- `PerformUpdate(downloadURL string) (bool, error)`
- `OpenReleaseURL(url string) error`

### 6. Frontend Integration

Add update UI to your frontend using ES module bindings:

```javascript
import * as App from "./bindings/github.com/owner/repo/app.js";
import { Events } from "@wailsio/runtime";

// Listen for progress
Events.On("updateProgress", (event) => {
    console.log("Update progress:", event.data);
});

// Check for updates on app load
async function checkForUpdates() {
    const updateInfo = await App.CheckForUpdate();
    if (updateInfo.available) {
        showUpdateButton(updateInfo);
    }
}

// Perform update
async function performUpdate(downloadUrl) {
    await App.PerformUpdate(downloadUrl);
}
```

### 7. Build with Version Injection

Use ldflags to inject version at build time:

```bash
# Development
wails3 dev

# Release (version from git tag)
go build -ldflags "-s -w -X main.CurrentVersion=v1.0.0" -o build/bin/myproject.exe .
```

---

## CI/CD Pipeline Details

### CI Workflow (ci.yml)

Triggers on: `push` and `pull_request` to `main`/`master`

| Job | Description |
|-----|-------------|
| **lint** | Runs `golangci-lint` on Windows |
| **test** | Runs tests with coverage on all 3 platforms, enforcing 70% threshold on `internal/` |
| **build** | Installs `wails3`, generates bindings, and builds for Linux, Windows, macOS |

### Release Workflow (release.yml)

Triggers on: `push` tags matching `v*`

| Step | Output |
|------|--------|
| Unit Tests | Verifies test suite passes |
| Generate Bindings | `wails3 generate bindings` |
| Build Windows EXE | `build/bin/<project>-windows-amd64.exe` with version injected |
| Create Release | GitHub Release with executable asset via `softprops/action-gh-release@v3` |

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
┌──────────────────────────────────────────────┐
│  PerformUpdate(downloadURL)                  │
│  1. Download EXE to temp                     │
│  2. Emit progress: application.Get().Event   │
│  3. Create update batch script               │
│  4. Start script & application.Get().Quit()  │
│  5. Script replaces exe and restarts app     │
└──────────────────────────────────────────────┘
```

---

## Makefile Commands

| Command | Description |
|---------|-------------|
| `make dev` | Run in development mode (`wails3 dev`) |
| `make bindings` | Generate frontend bindings (`wails3 generate bindings`) |
| `make build` | Build production binary (`wails3 build`) |
| `make test` | Run all tests |
| `make lint` | Run golangci-lint |
| `make coverage` | Run tests with coverage report |
| `make fmt` | Format code with gofmt/goimports |
| `make fmt-check` | Check code formatting |
| `make check` | Run lint + test |
| `make clean` | Remove build artifacts |

---

## ⚠️ Known CI Pitfalls & Fixes

### 1. Windows PowerShell Argument Parsing

**Problem**: PowerShell incorrectly parses `-coverprofile=coverage.out`, treating `.out` as a separate package name.

**Solution**: Use **space** instead of `=`:

```yaml
# ❌ WRONG - Fails on Windows PowerShell
run: go test -coverprofile=coverage.out ./...

# ✅ CORRECT - Works on all platforms
run: go test -coverprofile coverage.out ./...
```

### 2. Linux Wails v3 Dependencies

**Problem**: Linux builds for Wails v3 require GTK4 and WebKitGTK 6.0 development headers.

**Solution**: Install `libgtk-4-dev` and `libwebkitgtk-6.0-dev`:

```bash
sudo apt-get install -y libgtk-4-dev libwebkitgtk-6.0-dev
```

### 3. Wails v3 Service Architecture

**Problem**: Wails v2 used `context.Context` saved on startup for dialogs/events (`runtime.EventsEmit(ctx, ...)`). In Wails v3 this fails or is missing.

**Solution**: Use `application.Get()` singleton:
```go
application.Get().Event.Emit("eventName", data)
application.Get().Dialog.OpenFileWithOptions(...)
application.Get().Clipboard.SetText(text)
application.Get().Quit()
```
