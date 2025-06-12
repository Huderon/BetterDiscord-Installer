package backend

import (
	"context"

	"installer/backend/discord"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type Actions struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewActionsManager() *Actions {
	return &Actions{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (action *Actions) SetContext(ctx context.Context) {
	action.ctx = ctx
}

func (action *Actions) Install(corePaths ...string) {
	for i := range corePaths {
		install := discord.ResolvePath(corePaths[i])
		if install == nil {
			continue
		}

		if err := install.InstallBD(); err != nil {
			runtime.EventsEmit(action.ctx, "failure")
			return
		}
	}

	runtime.EventsEmit(action.ctx, "success")
}

func (action *Actions) Uninstall(corePaths ...string) {
	for i := range corePaths {
		install := discord.ResolvePath(corePaths[i])
		if install == nil {
			continue
		}

		if err := install.UninstallBD(); err != nil {
			runtime.EventsEmit(action.ctx, "failure")
			return
		}
	}

	runtime.EventsEmit(action.ctx, "success")
}

func (action *Actions) Repair(corePaths ...string) {
	for i := range corePaths {
		install := discord.ResolvePath(corePaths[i])
		if install == nil {
			continue
		}

		if err := install.RepairBD(); err != nil {
			runtime.EventsEmit(action.ctx, "failure")
			return
		}
	}

	runtime.EventsEmit(action.ctx, "success")

	result, err := runtime.MessageDialog(action.ctx, runtime.MessageDialogOptions{
		Type:          runtime.QuestionDialog,
		Title:         "Reinstall BetterDiscord?",
		Message:       "After repairing, you need to reinstall BetterDiscord. Would you like to do that now?",
		DefaultButton: "No",
	})

	if err != nil {
		return
	}

	if result == "Yes" {
		runtime.EventsEmit(action.ctx, "reset")
		action.Install(corePaths...)
	}
}

func (action *Actions) Other(install discord.DiscordInstall, channel discord.DiscordChannel) {

}
