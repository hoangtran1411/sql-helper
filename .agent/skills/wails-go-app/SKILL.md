---
name: Wails Go App
description: Create desktop applications using Go backend with Wails v3 framework and modern HTML/CSS/JS frontend. Lightweight alternative to Fyne.
---

# Wails Go Desktop Application Skill (Wails v3)

This skill provides instructions for building modern, lightweight desktop applications using **Wails v3** with a Go backend and HTML/CSS/JS frontend.

## When to Use This Skill

Use this skill when:
- Building a Go desktop application with GUI using Wails v3
- Need a lightweight alternative to Fyne (Fyne requires OpenGL compilation)
- Want to use modern web technologies (HTML/CSS/JS) for UI
- Building cross-platform desktop apps (Windows, macOS, Linux)

## Prerequisites

### 1. Install Wails v3 CLI
```bash
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
```

### 2. Verify Installation
```bash
wails3 doctor
```

This checks:
- Go version (Go 1.24+ recommended, project uses Go 1.27)
- Node.js (required for frontend dev tooling or `@wailsio/runtime`)
- WebView2 runtime (Windows) / WebKitGTK (Linux)

## Project Structure

A typical Wails v3 project structure:

```
project/
├── main.go              # Wails v3 entry point (application.New, window creation)
├── app.go               # Backend service logic (public Go methods exposed to frontend)
├── wails.json           # Wails configuration file
├── go.mod               # Go module
├── frontend/
│   ├── index.html       # Main HTML file (<script type="module" src="main.js">)
│   ├── style.css        # CSS styles
│   ├── main.js          # Frontend JavaScript using ES module bindings
│   └── bindings/        # Auto-generated Wails v3 JS/TS bindings
└── internal/            # Business logic (modular, UI-agnostic)
```

## Core Files

### 1. main.go - Application Entry Point

```go
package main

import (
	"embed"
	"io/fs"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/*
var assets embed.FS

func main() {
	frontendFS, err := fs.Sub(assets, "frontend")
	if err != nil {
		log.Fatalf("failed to create frontend sub filesystem: %v", err)
	}

	appService := NewApp()

	app := application.New(application.Options{
		Name:        "My Wails App",
		Description: "A desktop utility built with Go and Wails v3",
		Services: []application.Service{
			application.NewService(appService),
		},
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(frontendFS),
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "My Wails App",
		Width:  1024,
		Height: 768,
		URL:    "/",
		BackgroundColour: application.RGBA{
			Red:   248,
			Green: 249,
			Blue:  250,
			Alpha: 255,
		},
	})

	if err := app.Run(); err != nil {
		log.Fatal("Error:", err.Error())
	}
}
```

### 2. app.go - Backend Service Logic

```go
package main

import (
	"fmt"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type App struct{}

func NewApp() *App {
	return &App{}
}

// OpenFile opens a file selection dialog
func (a *App) OpenFile() (string, error) {
	dialog := application.Get().Dialog.OpenFileWithOptions(&application.OpenFileDialogOptions{
		Title: "Select File",
		Filters: []application.FileFilter{
			{
				DisplayName: "All Files (*.*)",
				Pattern:     "*.*",
			},
		},
	})

	path, err := dialog.PromptForSingleSelection()
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "cancel") {
			return "", nil // User cancelled
		}
		return "", err
	}
	return path, nil
}

// CopyText copies text to system clipboard
func (a *App) CopyText(text string) error {
	if !application.Get().Clipboard.SetText(text) {
		return fmt.Errorf("failed to copy text to clipboard")
	}
	return nil
}

// EmitProgress emits a real-time event to the frontend
func (a *App) EmitProgress(percent int, message string) {
	application.Get().Event.Emit("progress", map[string]interface{}{
		"percent": percent,
		"message": message,
	})
}
```

### 3. frontend/main.js - Frontend Logic with ES Bindings

```javascript
import * as App from "./bindings/github.com/yourname/project/app.js";
import { Events } from "@wailsio/runtime";

// Call Go backend service method
async function selectFile() {
    try {
        const path = await App.OpenFile();
        if (path) {
            console.log('Selected file:', path);
        }
    } catch (err) {
        console.error('Error selecting file:', err);
    }
}

// Listen for real-time events from Go
Events.On('progress', (event) => {
    const { percent, message } = event.data;
    console.log(`Progress: ${percent}% - ${message}`);
});
```

### 4. frontend/index.html

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>My Wails App</title>
    <link rel="stylesheet" href="style.css">
</head>
<body>
    <div id="app">
        <!-- UI Elements -->
    </div>
    <script type="module" src="main.js"></script>
</body>
</html>
```

## Commands

### Generate Bindings
```bash
wails3 generate bindings
```
Generates typed ES module bindings inside `frontend/bindings/`.

### Development Mode
```bash
wails3 dev
```
- Hot reload for frontend changes
- Auto-rebuild for Go backend changes

### Build Production
```bash
# Standard build
wails3 build

# Windows AMD64
wails3 build -platform windows/amd64

# Linux AMD64
wails3 build -platform linux/amd64

# macOS Universal
wails3 build -platform darwin/universal
```

## Go ↔ JavaScript Communication

### 1. Invoking Go from JavaScript
Public methods on registered services are generated into the bindings package:
```javascript
import * as App from "./bindings/github.com/hoangtran1411/sql-helper/app.js";

const result = await App.GenerateSQL(headers, dataRows, options);
```

### 2. Emitting Events from Go
```go
application.Get().Event.Emit("eventName", payload)
```

### 3. Listening to Events in Frontend
```javascript
import { Events } from "@wailsio/runtime";

Events.On("eventName", (event) => {
    console.log("Received data:", event.data);
});
```

## Common Wails v3 Runtime APIs

```go
import "github.com/wailsapp/wails/v3/pkg/application"

// Dialogs
application.Get().Dialog.OpenFileWithOptions(&application.OpenFileDialogOptions{...})
application.Get().Dialog.SaveFileWithOptions(&application.SaveFileDialogOptions{...})

// Clipboard
application.Get().Clipboard.SetText("text")
application.Get().Clipboard.Text()

// Browser
application.Get().Browser.OpenURL("https://example.com")

// Events
application.Get().Event.Emit("event-name", data)

// Application lifecycle
application.Get().Quit()
```

## Tips for High Performance

1. **O(1) Streaming**: Never read entire large files into frontend memory. Use streaming readers (`excelize.Rows`) and write buffered chunks directly to disk (`bufio.Writer`).
2. **Preview Pagination**: Only send preview subsets (e.g. 100 rows) across the bridge for DOM rendering.
3. **Vanilla Frontend**: For desktop utilities, vanilla JS with modern CSS glassmorphism avoids heavyweight bundlers and speeds up startup.
