package betterdiscord

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"installer/types"
	"installer/utils"
	"installer/wsl"
)

type BDInstall struct {
	root          string
	data          string
	asar          string
	plugins       string
	themes        string
	hasDownloaded bool
}

// Root returns the root directory path of the BetterDiscord installation
func (i *BDInstall) Root() string {
	return i.root
}

// Data returns the data directory path
func (i *BDInstall) Data() string {
	return i.data
}

// Asar returns the path to the BetterDiscord asar file
func (i *BDInstall) Asar() string {
	return i.asar
}

// Plugins returns the plugins directory path
func (i *BDInstall) Plugins() string {
	return i.plugins
}

// Themes returns the themes directory path
func (i *BDInstall) Themes() string {
	return i.themes
}

// HasDownloaded returns whether BetterDiscord has been downloaded
func (i *BDInstall) HasDownloaded() bool {
	return i.hasDownloaded
}

// Download downloads the BetterDiscord asar file. When useDevBuild is true it
// pulls the rolling "canary" pre-release from GitHub instead of the stable asar.
func (i *BDInstall) Download(useDevBuild bool) error {
	return i.download(useDevBuild)
}

// Prepare creates all necessary directories for BetterDiscord
func (i *BDInstall) Prepare() error {
	return i.prepare()
}

// Repair repairs the BetterDiscord installation with specific options
func (i *BDInstall) Repair(channel types.DiscordChannel, options types.RepairOptions) error {
	return i.repair(channel, options)
}

func (i *BDInstall) IsAsarInstalled() bool {
	return utils.Exists(i.asar)
}

// RemoveAll deletes the BetterDiscord installation directory and all contents.
func (i *BDInstall) RemoveAll() error {
	if !utils.Exists(i.root) {
		log.Printf("✅ BetterDiscord folder not found: %s\n", i.root)
		return nil
	}

	if err := os.RemoveAll(i.root); err != nil {
		log.Printf("❌ Failed to remove BetterDiscord folder: %s\n", i.root)
		log.Printf("   %s\n", err.Error())
		return err
	}

	log.Printf("✅ Removed BetterDiscord folder: %s\n", i.root)
	return nil
}

var lock = &sync.Mutex{}
var globalInstance *BDInstall

func GetInstallation(base ...string) *BDInstall {
	if len(base) == 0 {
		if globalInstance != nil {
			return globalInstance
		}

		lock.Lock()
		defer lock.Unlock()
		if globalInstance != nil {
			return globalInstance
		}

		// Default to user config directory
		configDir, _ := os.UserConfigDir()

		// Handle WSL with Windows home directory
		if wsl.IsWSL() {
			winHome, err := wsl.WindowsHome()
			if err == nil && winHome != "" {
				configDir = filepath.Join(winHome, "AppData", "Roaming")
			}
		}

		globalInstance = GetInstallation(configDir)

		return globalInstance
	}

	return New(filepath.Join(base[0], "BetterDiscord"))
}

func New(root string) *BDInstall {
	return &BDInstall{
		root:          root,
		data:          filepath.Join(root, "data"),
		asar:          filepath.Join(root, "data", "betterdiscord.asar"),
		plugins:       filepath.Join(root, "plugins"),
		themes:        filepath.Join(root, "themes"),
		hasDownloaded: false,
	}
}
