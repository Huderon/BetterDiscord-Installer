package main

import (
	"context"
	"embed"
	"fmt"
	"log"

	"installer/backend"
	"installer/backend/discord"
	"installer/backend/utils"

	"github.com/pkg/browser"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/build
var assets embed.FS

type App struct {
	ctx context.Context
}

// CreateApp creates a new App application struct
func CreateApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) SetContext(ctx context.Context) {
	a.ctx = ctx
}

// Override logger to send as events to the GUI
func (a *App) Write(p []byte) (n int, err error) {
	fmt.Println("CUSTOM WRITE")
	fmt.Println(string(p[:]))
	runtime.EventsEmit(a.ctx, "log", string(p[:]))
	return len(p), nil
}

func main() {
	// Create an instance of the app structure
	app := CreateApp()
	backend := backend.CreateBackend()

	bound := []interface{}{app, backend}
	others := backend.GetModules()
	bound = append(bound, others...)

	// Create application with options
	err := wails.Run(&options.App{
		Title:     "BetterDiscord Installer",
		Frameless: true,
		Width:     550,
		Height:    350,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			app.SetContext(ctx)
			backend.SetContext(ctx)
			app.CheckForUpdate()

			// Setup default logger to send data to GUI
			// TODO: create custom logger and use log.SetDefault or even slog.SetDefault
			log.SetOutput(app)
			log.SetFlags(0) // Don't add date/time
		},
		Bind: bound,
		EnumBind: []interface{}{
			discord.Channels,
			// backend.LogEvent,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

var version = "dev"

func (a *App) CheckForUpdate() {
	// If the version is "dev", it's a development build,
	// so we don't check for updates
	if version == "dev" {
		return
	}

	// Get latest installer version from GitHub API
	apiData, err := utils.DownloadJSON[utils.Release]("https://api.github.com/repos/BetterDiscord/Installer/releases/latest")
	if err != nil {
		return
	}

	// If the current version is greater than or equal
	// to the latest version, no update is needed
	if version >= apiData.TagName {
		return
	}

	result, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:          runtime.QuestionDialog,
		Title:         "Update Available",
		Message:       fmt.Sprintf("A new version (%s) of the installer is available. Would you like to download it now?", apiData.TagName),
		DefaultButton: "Yes",
	})

	if err != nil {
		return
	}

	if result == "Yes" {
		browser.OpenURL(apiData.HTMLURL)
	}
}
