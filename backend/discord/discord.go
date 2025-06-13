package discord

import (
	_ "embed"
	"installer/backend/utils"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type DiscordInstall struct {
	corePath  string         `json:"corePath"`
	channel   DiscordChannel `json:"channel"`
	version   string         `json:"version"`
	isFlatpak bool           `json:"isFlatpak"`
	isSnap    bool           `json:"isSnap"`
}

func (discord *DiscordInstall) IsChannel(channel DiscordChannel) bool {
	return discord.channel == channel
}

func (discord *DiscordInstall) GetPath() string {
	return discord.corePath
}

func (discord *DiscordInstall) InstallBD() error {
	// Gets the global BetterDiscord install
	betterdiscord := GetBetterDiscord()

	// Snaps get their own local BD install
	if discord.isSnap {
		betterdiscord = GetBetterDiscord(filepath.Clean(filepath.Join(discord.corePath, "..", "..", "..", "..")))
	}

	// Make BetterDiscord folders
	log.Printf("## Preparing BetterDiscord...")
	if err := betterdiscord.prepare(); err != nil {
		return err
	}
	log.Printf("✅ BetterDiscord prepared for install")
	log.Printf("")

	// Download and write betterdiscord.asar
	log.Printf("## Downloading BetterDiscord...")
	if err := betterdiscord.download(); err != nil {
		return err
	}
	log.Printf("✅ BetterDiscord downloaded")
	log.Printf("")

	// Write injection script to discord_desktop_core/index.js
	log.Printf("## Injecting into Discord...")
	if err := discord.inject(betterdiscord); err != nil {
		return err
	}
	log.Printf("✅ Injection successsful")
	log.Printf("")

	// Terminate and restart Discord if possible
	log.Printf("## Restarting %s...", discord.channel.Name())
	if err := discord.restart(); err != nil {
		return err
	}
	log.Printf("")

	return nil
}

func (discord *DiscordInstall) UninstallBD() error {
	log.Printf("## Removing injection...")
	if err := discord.uninject(); err != nil {
		return err
	}
	log.Printf("")

	log.Printf("## Restarting %s...", discord.channel.Name())
	if err := discord.restart(); err != nil {
		return err
	}
	log.Printf("")

	return nil
}

func (discord *DiscordInstall) RepairBD() error {
	if err := discord.UninstallBD(); err != nil {
		return err
	}

	// Gets the global BetterDiscord install
	betterdiscord := GetBetterDiscord()

	// Snaps get their own local BD install
	if discord.isSnap {
		betterdiscord = GetBetterDiscord(filepath.Clean(filepath.Join(discord.corePath, "..", "..", "..", "..")))
	}

	if err := betterdiscord.repair(discord.channel); err != nil {
		return nil
	}

	return nil
}

//go:embed injection.js
var injectionScript string

func (discord *DiscordInstall) inject(bd *BetterDiscord) error {
	if discord.isFlatpak {
		cmd := exec.Command("flatpak", "--user", "override", "com.discordapp."+discord.channel.Exe(), "--filesystem="+bd.root)
		if err := cmd.Run(); err != nil {
			log.Printf("❌ Could not give flatpak access to %s", bd.root)
			log.Printf("❌ %s", err.Error())
			return err
		}
	}

	if err := os.WriteFile(filepath.Join(discord.corePath, "index.js"), []byte(injectionScript), 0755); err != nil {
		log.Printf("❌ Unable to write index.js in %s", discord.corePath)
		log.Printf("❌ %s", err.Error())
		return err
	}

	log.Printf("✅ Injected into %s", discord.corePath)
	return nil
}

func (discord *DiscordInstall) uninject() error {
	indexFile := filepath.Join(discord.corePath, "index.js")

	contents, err := os.ReadFile(indexFile)

	// First try to check the file, but if there's an issue we try to blindly overwrite below
	if err == nil {
		if !strings.Contains(strings.ToLower(string(contents)), "betterdiscord") {
			log.Printf("✅ No injection found for %s", discord.channel.Name())
			return nil
		}
	}

	if err := os.WriteFile(indexFile, []byte(`module.exports = require("./core.asar");`), 0o644); err != nil {
		log.Printf("❌ Unable to write file %s", indexFile)
		log.Printf("❌ %s", err.Error())
		return err
	}
	log.Printf("✅ Removed from %s", discord.channel.Name())

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
	if discord.isFlatpak || discord.isSnap {
		cmd.Path, _ = os.UserHomeDir()
	}

	if err := cmd.Start(); err != nil {
		log.Printf("❌ Unable to restart %s, please do so manually!", discord.channel.Name())
		log.Printf("❌ %s", err.Error())
		return err
	}
	log.Printf("✅ Restarted %s", discord.channel.Name())
	return nil
}
