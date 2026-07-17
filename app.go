package main

import (
	"context"
	"fmt"
	"strings"

	"installer/types"
	"installer/utils"

	"github.com/pkg/browser"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) SetContext(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetVersion() string {
	return version
}

func (a *App) CheckForUpdate() {
	// If the version doesn't start with "v" then it's a
	// development build, so we don't check for updates
	if !strings.HasPrefix(version, "v") {
		return
	}

	// Get latest installer version from GitHub API
	apiData, err := utils.DownloadJSON[types.GitHubRelease]("https://api.github.com/repos/BetterDiscord/Installer/releases/latest")
	if err != nil {
		// Offline or GitHub API error — skip the update check silently.
		return
	}

	// If the current version is greater than or equal
	// to the latest version, no update is needed
	if utils.CompareVersions(version, apiData.TagName) >= 0 {
		return
	}

	result, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    runtime.QuestionDialog,
		Title:   "Update Available",
		Message: fmt.Sprintf("A new version (%s) of the installer is available. Would you like to download it now?", apiData.TagName),
		// Explicit buttons are required on macOS: its dialog maps the result
		// back through Buttons, returning "" when none are set. The other
		// platforms ignore Buttons and return Yes/No on their own.
		Buttons:       []string{"Yes", "No"},
		DefaultButton: "Yes",
		CancelButton:  "No",
	})

	if err != nil {
		// Dialog failed to display — don't block startup over it.
		return
	}

	if result == "Yes" {
		browser.OpenURL(apiData.HTMLURL)
	}
}
