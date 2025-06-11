package discord

import (
	"installer/backend/utils"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type DiscordChannel int64

const (
	Stable DiscordChannel = iota
	Canary
	PTB
)

var Channels = []DiscordChannel{Stable, Canary, PTB}

func (channel *DiscordChannel) String() string {
	switch *channel {
	case Stable:
		return "stable"
	case Canary:
		return "canary"
	case PTB:
		return "ptb"
	}
	return ""
}

func (channel *DiscordChannel) Name() string {
	switch *channel {
	case Stable:
		return "Discord"
	case Canary:
		return "Discord Canary"
	case PTB:
		return "Discord PTB"
	}
	return ""
}

func (channel *DiscordChannel) Exe() string {
	name := channel.Name()

	if runtime.GOOS != "darwin" {
		name = strings.ReplaceAll(name, " ", "")
	}

	if runtime.GOOS == "windows" {
		name = name + ".exe"
	}

	return name
}

type DiscordInstall struct {
	corePath  string
	channel   DiscordChannel
	version   string
	isFlatpak bool
	isSnap    bool
}

func (discord *DiscordInstall) patch() {
	// Gets the global BetterDiscord install
	betterdiscord := GetBetterDiscord()

	// Snaps get their own local BD install
	if discord.isSnap {
		betterdiscord = GetBetterDiscord(filepath.Clean(filepath.Join(discord.corePath, "..", "..", "..", "..", "Betterdiscord")))
	}

	// Make BetterDiscord folders
	if err := betterdiscord.prepare(); err != nil {
		return // TODO: do better?
	}

	// Download and write betterdiscord.asar
	if err := betterdiscord.download(); err != nil {
		return // TODO: do better?
	}

	// Write injection script to discord_desktop_core/index.js
	if err := discord.inject(betterdiscord); err != nil {
		return // TODO: do better?
	}

	// Terminate and restart Discord if possible
	if err := discord.restart(); err != nil {
		return // TODO: do better?
	}
}

func (discord *DiscordInstall) inject(bd *BetterDiscord) error {
	if discord.isFlatpak {
		cmd := exec.Command("flatpak", "--user", "override", "com.discordapp."+discord.channel.Exe(), "--filesystem="+bd.root)
		if err := cmd.Run(); err != nil {
			log.Printf("❌ Could not give flatpak access to %s", bd.root)
			log.Printf("❌ %s", err.Error())
			return err
		}
	}

	var injectionScript = `// BetterDiscord's Injection Script
const path = require("path");
const electron = require("electron");

// Windows and macOS both use the fixed global BetterDiscord folder but
// Electron gives the postfixed version of userData, so go up a directory
let userConfig = path.join(electron.app.getPath("userData"), "..");

// If we're on Linux there are a couple cases to deal with
if (process.platform !== "win32" && process.platform !== "darwin") {
    // Use || instead of ?? because a falsey value of "" is invalid per XDG spec
    userConfig = process.env.XDG_CONFIG_HOME || path.join(process.env.HOME, ".config");

    // HOST_XDG_CONFIG_HOME is set by flatpak, so use without validation if set
    if (process.env.HOST_XDG_CONFIG_HOME) userConfig = process.env.HOST_XDG_CONFIG_HOME;
}

require(path.join(userConfig, "BetterDiscord", "data", "betterdiscord.asar"));

// Discord's Default Export
module.exports = require("./core.asar");
`

	if err := os.WriteFile(filepath.Join(discord.corePath, "index.js"), []byte(injectionScript), 0755); err != nil {
		log.Printf("❌ Unable to write index.js in %s", discord.corePath)
		log.Printf("❌ %s", err.Error())
		return err
	}

	log.Printf("✅ Injected into %s", discord.corePath)
	return nil
}

func (discord *DiscordInstall) restart() error {
	name := discord.channel.Exe()
	exeName := utils.GetProcessExe(name)

	if running, _ := utils.IsRunning(name); !running {
		log.Printf("✅ %s not running", discord.channel.Name())
		return nil
	}

	if err := utils.KillProcess(name); err != nil {
		log.Printf("❌ Unable to restart %s, please do so manually!", discord.channel.Name())
		log.Printf("❌ %s", err.Error())
		return err
	}

	// Use binary found in killing process
	cmd := exec.Command(exeName)
	if discord.isFlatpak {
		cmd = exec.Command("flatpak", "run", "com.discordapp."+discord.channel.Exe())
	} else if discord.isSnap {
		cmd = exec.Command("snap", "run", discord.channel.Exe())
	}

	// Set working directory to user home
	cmd.Path, _ = os.UserHomeDir()
	if err := cmd.Run(); err != nil {
		log.Printf("❌ Unable to restart %s, please do so manually!", discord.channel.Name())
		log.Printf("❌ %s", err.Error())
		return err
	}
	log.Printf("✅ Restarted %s", discord.channel.Name())
	return nil
}
