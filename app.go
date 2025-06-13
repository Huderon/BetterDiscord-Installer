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
