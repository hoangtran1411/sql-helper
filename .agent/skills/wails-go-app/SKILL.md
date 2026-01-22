---
name: Wails Go App
description: Create desktop applications using Go backend with Wails framework and modern HTML/CSS/JS frontend. Lightweight alternative to Fyne.
---

# Wails Go Desktop Application Skill

This skill provides instructions for building modern, lightweight desktop applications using **Wails v2** with Go backend and HTML/CSS/JS frontend.

## When to Use This Skill

Use this skill when:
- Building a Go desktop application with GUI
- Need a lightweight alternative to Fyne (Fyne requires OpenGL compilation)
- Want to use web technologies (HTML/CSS/JS) for UI
- Building cross-platform apps (Windows, macOS, Linux)

## Prerequisites

### 1. Install Wails CLI
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### 2. Verify Installation
```bash
wails doctor
```

This checks:
- Go version
- Node.js (required for frontend bundling if using frameworks)
- WebView2 runtime (Windows) - usually pre-installed on Windows 10/11

## Project Structure

A typical Wails project structure:

```
project/
├── main.go              # Wails entry point with app configuration
├── app.go               # Backend logic (Go methods exposed to frontend)
├── wails.json           # Wails configuration file
├── go.mod               # Go module
├── frontend/
│   └── dist/
│       ├── index.html   # Main HTML file
│       ├── style.css    # CSS styles
│       └── app.js       # Frontend JavaScript logic
└── internal/            # Business logic (optional)
```

## Core Files

### 1. main.go - Application Entry Point

```go
package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:     "My App Name",
		Width:     900,
		Height:    720,
		MinWidth:  700,
		MinHeight: 600,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
```

### 2. app.go - Backend Logic

```go
package main

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Exposed methods - callable from JavaScript
func (a *App) SelectFile() (string, error) {
	return runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select File",
		Filters: []runtime.FileFilter{
			{DisplayName: "All Files", Pattern: "*.*"},
		},
	})
}

func (a *App) SelectFolder() (string, error) {
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Folder",
	})
}

// Emit events to frontend
func (a *App) EmitProgress(progress float64) {
	runtime.EventsEmit(a.ctx, "progress", progress)
}
```

### 3. wails.json - Configuration

```json
{
  "$schema": "https://wails.io/schemas/config.v2.json",
  "name": "myapp",
  "outputfilename": "MyApp",
  "frontend:install": "",
  "frontend:build": "",
  "author": {
    "name": "Your Name",
    "email": "your@email.com"
  }
}
```

### 4. frontend/dist/index.html

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>My App</title>
    <link rel="stylesheet" href="style.css">
</head>
<body>
    <div id="app">
        <!-- Your UI here -->
    </div>
    <script src="wails.js"></script>
    <script src="app.js"></script>
</body>
</html>
```

### 5. frontend/dist/app.js - Frontend Logic

```javascript
// Call Go methods
async function selectFile() {
    try {
        const path = await window.go.main.App.SelectFile();
        console.log('Selected:', path);
    } catch (err) {
        console.error(err);
    }
}

// Listen for events from Go
document.addEventListener('DOMContentLoaded', function() {
    if (typeof runtime !== 'undefined') {
        runtime.EventsOn('progress', function(value) {
            console.log('Progress:', value);
        });
    }
});
```

## Commands

### Development Mode
```bash
wails dev
```
- Hot reload for frontend changes
- Auto-rebuild for Go changes
- Access via browser at displayed URL

### Build Production
```bash
wails build
```

### Build with Compression (requires UPX)
```bash
wails build -upx
```

### Build Windows Installer (requires NSIS)
```bash
wails build -nsis
```

## Go ↔ JavaScript Communication

### Calling Go from JavaScript
```javascript
// All public methods on bound structs are available
const result = await window.go.main.App.MethodName(arg1, arg2);
```

### Emitting Events from Go to JavaScript
```go
runtime.EventsEmit(a.ctx, "eventName", data)
```

### Listening for Events in JavaScript
```javascript
runtime.EventsOn("eventName", (data) => {
    console.log(data);
});
```

## Common Runtime Methods

```go
import "github.com/wailsapp/wails/v2/pkg/runtime"

// Dialogs
runtime.OpenFileDialog(ctx, options)
runtime.OpenDirectoryDialog(ctx, options)
runtime.SaveFileDialog(ctx, options)
runtime.MessageDialog(ctx, options)

// Events
runtime.EventsEmit(ctx, eventName, data)
runtime.EventsOn(ctx, eventName, callback)
runtime.EventsOff(ctx, eventName)

// Window Control
runtime.WindowMinimise(ctx)
runtime.WindowMaximise(ctx)
runtime.WindowSetTitle(ctx, title)
runtime.WindowSetSize(ctx, width, height)
runtime.Quit(ctx)
```

## CSS Dark Theme Template

See `examples/dark-theme.css` for a complete premium dark theme with:
- CSS variables for easy customization
- Card components
- Form inputs
- Buttons with hover effects
- Progress bars
- Status messages
- Smooth animations

## Tips

1. **Keep frontend simple**: For simple apps, use vanilla HTML/CSS/JS without bundlers
2. **Use embed.FS**: Always embed frontend assets for single-binary distribution
3. **Struct binding**: Only public methods (uppercase) are exposed to JavaScript
4. **Error handling**: Go errors are returned as JavaScript promise rejections
5. **Progress updates**: Use events for real-time progress reporting

## Comparison: Wails vs Fyne

| Aspect | Wails | Fyne |
|--------|-------|------|
| Binary Size | ~8-10MB | ~25-30MB |
| Compile Time | Fast | Slow (OpenGL) |
| UI Technology | HTML/CSS/JS | Go widgets |
| Styling | Full CSS control | Limited |
| Learning Curve | Web devs friendly | Go-only |
| Hot Reload | Yes (dev mode) | No |

## Resources

- [Wails Documentation](https://wails.io/docs/introduction)
- [Wails GitHub](https://github.com/wailsapp/wails)
- [Wails Templates](https://wails.io/docs/community/templates)
