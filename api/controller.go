package api

import (
	"context"
	"fmt"
	"log"

	"installer/discord"
	"installer/types"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Controller struct {
	ctx context.Context
}

func NewController() *Controller {
	return &Controller{}
}

// SetContext is called when the app starts. The context is saved
// so we can call the runtime methods
func (action *Controller) SetContext(ctx context.Context) {
	action.ctx = ctx
}

func (a *Controller) Write(p []byte) (n int, err error) {
	fmt.Println(string(p[:]))
	runtime.EventsEmit(a.ctx, "log", string(p[:]))
	return len(p), nil
}

func (d *Controller) GetDiscordPath(channel string) string {
	return discord.GetSuggestedPath(types.ParseChannel(channel))
}

// #region Actions
func (action *Controller) Install(resourcePaths []string, options types.InstallOptions) {
	for i := range resourcePaths {
		install := discord.ResolvePath(resourcePaths[i])
		if install == nil {
			continue
		}

		if err := install.InstallBD(options); err != nil {
			// Surface the reason: some failures (e.g. GetBetterDiscordInstall's
			// path resolution) return an error the inner calls didn't log.
			log.Printf("❌ %s\n", err.Error())
			runtime.EventsEmit(action.ctx, "failure")
			return
		}
	}

	runtime.EventsEmit(action.ctx, "success")
}

func (action *Controller) Uninstall(resourcePaths []string, options types.UninstallOptions) {
	for i := range resourcePaths {
		install := discord.ResolvePath(resourcePaths[i])
		if install == nil {
			continue
		}

		if err := install.UninstallBD(options); err != nil {
			log.Printf("❌ %s\n", err.Error())
			runtime.EventsEmit(action.ctx, "failure")
			return
		}
	}

	runtime.EventsEmit(action.ctx, "success")
}

func (action *Controller) Repair(resourcePaths []string, options types.RepairOptions) {
	for i := range resourcePaths {
		install := discord.ResolvePath(resourcePaths[i])
		if install == nil {
			continue
		}

		if err := install.RepairBD(options); err != nil {
			log.Printf("❌ %s\n", err.Error())
			runtime.EventsEmit(action.ctx, "failure")
			return
		}
	}

	result, err := runtime.MessageDialog(action.ctx, runtime.MessageDialogOptions{
		Type:    runtime.QuestionDialog,
		Title:   "Repair Complete",
		Message: "Repair is complete. Would you like to reinstall BetterDiscord now?",
		// Explicit buttons are required on macOS: its dialog maps the result
		// back through Buttons, returning "" when none are set. The other
		// platforms ignore Buttons and return Yes/No on their own.
		Buttons:       []string{"Yes", "No"},
		DefaultButton: "Yes",
		CancelButton:  "No",
	})

	if err != nil {
		runtime.EventsEmit(action.ctx, "success")
		return
	}

	if result == "Yes" {
		runtime.EventsEmit(action.ctx, "reset")
		runtime.EventsEmit(action.ctx, "navigate", map[string]string{
			"action": "install",
		})
		return
	}

	runtime.EventsEmit(action.ctx, "success")
}

// #endregion

// #region Dialogs
func (d *Controller) BrowseForDiscord(schannel string) string {
	channel := types.ParseChannel(schannel)

	selection, err := runtime.OpenDirectoryDialog(d.ctx, runtime.OpenDialogOptions{
		Title:                      "Browsing to " + channel.Name(),
		DefaultDirectory:           discord.DefaultBrowseDir(),
		ShowHiddenFiles:            true,
		TreatPackagesAsDirectories: true,
	})

	if err != nil || selection == "" {
		return ""
	}

	if result := discord.ResolvePath(selection); result != nil {
		return result.ResourcesPath
	}

	return ""
}

func (d *Controller) ConfirmAction(title string, message string) string {
	result, err := runtime.MessageDialog(d.ctx, runtime.MessageDialogOptions{
		Type:    runtime.QuestionDialog,
		Title:   title,
		Message: message,
		// See the note on explicit Buttons in RepairAction above.
		Buttons:       []string{"Yes", "No"},
		DefaultButton: "No",
		CancelButton:  "No",
	})

	if err != nil {
		return ""
	}

	return result
}

func (d *Controller) ShowNotice(dialog string, title string, message string) string {

	dialogType := runtime.InfoDialog
	if dialog == "error" {
		dialogType = runtime.ErrorDialog
	}

	result, err := runtime.MessageDialog(d.ctx, runtime.MessageDialogOptions{
		Type:    dialogType,
		Title:   title,
		Message: message,
	})

	if err != nil {
		return ""
	}

	return result
}

// #endregion
