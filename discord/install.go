package discord

import (
	"installer/betterdiscord"
	"installer/types"
	"log"
	"path/filepath"
)

type DiscordInstall struct {
	CorePath  string                `json:"corePath"`
	Channel   types.DiscordChannel  `json:"channel"`
	Version   string                `json:"version"`
	IsFlatpak bool                  `json:"isFlatpak"`
	IsSnap    bool                  `json:"isSnap"`
}

// InstallBD installs BetterDiscord into this Discord installation
func (discord *DiscordInstall) InstallBD() error {
	bd := discord.GetBetterDiscordInstall()

	// Make BetterDiscord folders
	log.Println("🛠 Preparing BetterDiscord...")
	if err := bd.Prepare(); err != nil {
		return err
	}
	log.Println("✅ BetterDiscord prepared for install")
	log.Println("")

	// Download and write betterdiscord.asar
	log.Println("📥 Downloading BetterDiscord...")
	if err := bd.Download(); err != nil {
		return err
	}
	log.Println("✅ BetterDiscord downloaded")
	log.Println("")

	// Write injection script to discord_desktop_core/index.js
	log.Println("🔌 Injecting into Discord...")
	if err := discord.inject(bd); err != nil {
		return err
	}
	log.Println("✅ Injection successful")
	log.Println("")

	// Terminate and restart Discord if possible
	log.Printf("🔄 Restarting %s...\n", discord.Channel.Name())
	if err := discord.restart(); err != nil {
		return err
	}
	log.Println("")

	return nil
}

// UninstallBD removes BetterDiscord from this Discord installation
func (discord *DiscordInstall) UninstallBD() error {
	log.Println("🧹 Removing injection...")
	if err := discord.uninject(); err != nil {
		return err
	}
	log.Println("")

	log.Printf("🔄 Restarting %s...\n", discord.Channel.Name())
	if err := discord.restart(); err != nil {
		return err
	}
	log.Println("")

	return nil
}

// RepairBD repairs BetterDiscord for this Discord installation
func (discord *DiscordInstall) RepairBD() error {
	if err := discord.UninstallBD(); err != nil {
		return err
	}

	// Gets the global BetterDiscord install
	bd := betterdiscord.GetInstallation()

	// Snaps get their own local BD install
	if discord.IsSnap {
		bd = betterdiscord.GetInstallation(filepath.Clean(filepath.Join(discord.CorePath, "..", "..", "..", "..")))
	}

	if err := bd.Repair(discord.Channel); err != nil {
		return err
	}

	return nil
}

func (discord *DiscordInstall) GetBetterDiscordInstall() *betterdiscord.BDInstall {
	// Gets the global BetterDiscord install
	bd := betterdiscord.GetInstallation()

	// Snaps get their own local BD install
	if discord.IsSnap {
		bd = betterdiscord.GetInstallation(filepath.Clean(filepath.Join(discord.CorePath, "..", "..", "..", "..")))
	}

	return bd
}
