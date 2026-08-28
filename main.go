package main

import (
	"embed"
	"io/fs"
	"log"
	"net"
	"os"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed frontend/*
var assets embed.FS

func main() {
	// If FRONTEND_DEVSERVER_URL is set but not reachable, unset it to avoid dev mode blocking
	if devServer := os.Getenv("FRONTEND_DEVSERVER_URL"); devServer != "" {
		if conn, err := net.DialTimeout("tcp", "localhost:9245", 150*time.Millisecond); err != nil {
			_ = os.Unsetenv("FRONTEND_DEVSERVER_URL")
		} else {
			_ = conn.Close()
		}
	}

	frontendFS, err := fs.Sub(assets, "frontend")
	if err != nil {
		log.Fatalf("failed to create frontend sub filesystem: %v", err)
	}

	appService := NewApp()

	app := application.New(application.Options{
		Name:        "SQL Helper",
		Description: "A desktop utility to convert Excel data to SQL INSERT statements",
		Services: []application.Service{
			application.NewService(appService),
		},
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(frontendFS),
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "SQL Helper",
		Width:  1200,
		Height: 800,
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
