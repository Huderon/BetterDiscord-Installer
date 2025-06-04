package backend

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"installer/backend/utils"

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

func (action *Actions) log(message string) {
	runtime.EventsEmit(action.ctx, "log", message)
}

func (action *Actions) makeDirectory(folder string) error {
	exists := utils.Exists(folder)

	if exists {
		action.log(fmt.Sprintf("✅ Directory exists: %s", folder))
		return nil
	}

	if err := os.MkdirAll(folder, 0755); err != nil {
		action.log(fmt.Sprintf("❌ Failed to create directory: %s", folder))
		action.log(fmt.Sprintf("❌ %s", err.Error()))
		return err
	}

	action.log(fmt.Sprintf("✅ Directory created: %s", folder))
	return nil
}

func (action *Actions) ensureBDFolders() error {
	if err := action.makeDirectory(utils.Data); err != nil {
		return err
	}
	if err := action.makeDirectory(utils.Plugins); err != nil {
		return err
	}
	if err := action.makeDirectory(utils.Themes); err != nil {
		return err
	}
	return nil
}

func (action *Actions) downloadAsar() error {

	resp, err := utils.DownloadFile("https://betterdiscord.app/Download/betterdiscord.asar", utils.Asar)
	if err == nil {
		version := resp.Header[http.CanonicalHeaderKey("x-bd-version")]
		action.log(fmt.Sprintf("✅ Downloaded BetterDiscord version %s from the official website", version))
		return nil
	} else {
		action.log("❌ Failed to download package from official website")
		action.log(fmt.Sprintf("❌ %s", err.Error()))
		action.log("#### Falling back to GitHub...")
	}

	// Get download URL from GitHub API
	apiData, err := utils.DownloadJSON[utils.Release]("https://api.github.com/repos/BetterDiscord/BetterDiscord/releases/latest")
	if err != nil {
		action.log("❌ Failed to get asset url from GitHub")
		action.log(fmt.Sprintf("❌ %s", err.Error()))
		return err
	}

	var index = 0
	for i, asset := range apiData.Assets {
		if asset.Name == "betterdiscord.asar" {
			index = i
			break
		}
	}

	var downloadUrl = apiData.Assets[index].URL
	var version = apiData.TagName

	if downloadUrl != "" {
		action.log(fmt.Sprintf("✅ Found BetterDiscord package: %s", downloadUrl))
	}

	// Download asar into the BD folder
	_, err = utils.DownloadFile(downloadUrl, utils.Asar)
	if err != nil {
		action.log("❌ Failed to download package from GitHub")
		action.log(fmt.Sprintf("❌ %s", err.Error()))
		return err
	}

	action.log(fmt.Sprintf("✅ Downloaded BetterDiscord version %s from GitHub", version))

	return nil
}

func (action *Actions) inject(corePath string) error {
	var indString = `require("` + utils.Asar + `");`
	indString = strings.ReplaceAll(indString, `\`, "/")
	indString = indString + "\nmodule.exports = require(\"./core.asar\");"

	if err := os.WriteFile(path.Join(corePath, "index.js"), []byte(indString), 0755); err != nil {
		action.log(fmt.Sprintf("❌ Unable to write index.js in %s", corePath))
		action.log(fmt.Sprintf("❌ %s", err.Error()))
		return err
	}

	action.log(fmt.Sprintf("✅ Injected into %s", corePath))
	return nil
}

func (action *Actions) restartDiscord(channel string) error {
	name := utils.GetExecutableName(channel)
	exeName := utils.GetProcessExe(name)

	if running, _ := utils.IsRunning(name); !running {
		action.log(fmt.Sprintf("✅ %s not running", utils.GetChannelName(channel)))
		return nil
	}

	if err := utils.KillProcess(name); err != nil {
		action.log(fmt.Sprintf("❌ Unable to restart %s, please do so manually!", utils.GetChannelName(channel)))
		action.log(fmt.Sprintf("❌ %s", err.Error()))
		return err
	}

	cmd := exec.Command(exeName)
	cmd.Start()
	action.log(fmt.Sprintf("✅ Restarted %s", utils.GetChannelName(channel)))
	return nil
}

func (action *Actions) Install(channels []string, corePaths []string) {
	// Make BD directories
	action.log("## Creating required folders...")
	if err := action.ensureBDFolders(); err != nil {
		runtime.EventsEmit(action.ctx, "failure")
	}
	action.log("✅ Directories created")
	action.log("")

	action.log("## Downloading BetterDiscord package...")
	if err := action.downloadAsar(); err != nil {
		runtime.EventsEmit(action.ctx, "failure")
	}
	action.log("✅ Package downloaded")
	action.log("")

	action.log("## Injecting into Discord...")
	for _, core := range corePaths {
		if err := action.inject(core); err != nil {
			runtime.EventsEmit(action.ctx, "failure")
		}
	}
	action.log("✅ Injections successful")
	action.log("")

	action.log("## Restarting Discord...")
	for _, channel := range channels {
		action.restartDiscord(channel)
	}
	action.log("")

	runtime.EventsEmit(action.ctx, "success")
}

func (action *Actions) deleteInjection(channel string, corePath string) error {
	indexFile := path.Join(corePath, "index.js")

	contents, err := os.ReadFile(indexFile)

	// First try to check the file, but if there's an issue we try to blindly overwrite below
	if err == nil {
		if !strings.Contains(strings.ToLower(string(contents)), "betterdiscord") {
			action.log(fmt.Sprintf("✅ No injection found for %s", utils.GetChannelName(channel)))
			return nil
		}
	}

	if err := os.WriteFile(indexFile, []byte(`module.exports = require("./core.asar");`), 0o644); err != nil {
		action.log(fmt.Sprintf("❌ Unable to write file %s", indexFile))
		action.log(fmt.Sprintf("❌ %s", err.Error()))
		return err
	}
	action.log(fmt.Sprintf("✅ Removed from %s", utils.GetChannelName(channel)))

	return nil
}

func (action *Actions) Uninstall(channels []string, corePaths []string) {
	action.log("## Removing Discord injection...")
	for i, core := range corePaths {
		if err := action.deleteInjection(channels[i], core); err != nil {
			runtime.EventsEmit(action.ctx, "failure")
		}
	}
	// action.log("✅ Injections removed")
	action.log("")

	action.log("## Restarting Discord...")
	for _, channel := range channels {
		action.restartDiscord(channel)
	}
	action.log("")

	runtime.EventsEmit(action.ctx, "success")
}

func (action *Actions) disablePlugins(channel string) error {
	channelFolder := path.Join(utils.Data, channel)
	pluginsJson := path.Join(channelFolder, "plugins.json")

	if !utils.Exists(pluginsJson) {
		action.log(fmt.Sprintf("✅ No plugins enabled for %s", utils.GetChannelName(channel)))
		return nil
	}

	if err := os.Remove(pluginsJson); err != nil {
		action.log(fmt.Sprintf("❌ Unable to remove file %s", pluginsJson))
		action.log(fmt.Sprintf("❌ %s", err.Error()))
		return err
	}

	action.log(fmt.Sprintf("✅ Plugins disabled for %s", utils.GetChannelName(channel)))
	return nil
}

func (action *Actions) Repair(channels []string, corePaths []string) {
	action.log("## Removing Discord injection...")
	for i, core := range corePaths {
		if err := action.deleteInjection(channels[i], core); err != nil {
			runtime.EventsEmit(action.ctx, "failure")
		}
	}
	// action.log("✅ Injections removed")
	action.log("")

	action.log("## Restarting Discord...")
	for _, channel := range channels {
		action.restartDiscord(channel)
	}
	action.log("")

	action.log("## Disabling all plugins...")
	for _, channel := range channels {
		if err := action.disablePlugins(channel); err != nil {
			runtime.EventsEmit(action.ctx, "failure")
		}
	}
	action.log("")

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
		action.Install(channels, corePaths)
	}
}
