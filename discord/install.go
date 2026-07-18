package discord

import (
	"log"

	"installer/betterdiscord"
	"installer/types"
	"installer/utils"
)

type DiscordInstall struct {
	// ResourcesPath is the Discord install's `resources` directory (the one
	// holding app.asar). The json tag stays `corePath` so the existing Wails
	// bindings and frontend keep working; a full frontend rename is deferred.
	ResourcesPath string               `json:"corePath"`
	Channel       types.DiscordChannel `json:"channel"`
	Version       string               `json:"version"`
	IsFlatpak     bool                 `json:"isFlatpak"`
	IsSnap        bool                 `json:"isSnap"`
}

// InstallBD installs BetterDiscord into this Discord installation
func (discord *DiscordInstall) InstallBD(options types.InstallOptions) error {
	bd, err := discord.GetBetterDiscordInstall()
	if err != nil {
		return err
	}

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

	if options.RestartDiscord {
		// Terminate and restart Discord if possible
		log.Printf("🔄 Restarting %s...\n", discord.Channel.Name())
		if err := discord.restart(); err != nil {
			return err
		}
		log.Println("")
	}

	return nil
}

// UninstallBD removes BetterDiscord from this Discord installation
func (discord *DiscordInstall) UninstallBD(options types.UninstallOptions) error {
	log.Println("🧹 Removing injection...")
	if err := discord.uninject(); err != nil {
		return err
	}
	log.Println("")

	if options.FullUninstall {
		install, err := discord.GetBetterDiscordInstall()
		if err != nil {
			return err
		}
		if err := install.RemoveAll(); err != nil {
			return err
		}
		log.Println("")
	}

	if options.RestartDiscord {
		log.Printf("🔄 Restarting %s...\n", discord.Channel.Name())
		if err := discord.restart(); err != nil {
			return err
		}
		log.Println("")
	}

	return nil
}

// RepairBD repairs BetterDiscord for this Discord installation
func (discord *DiscordInstall) RepairBD(options types.RepairOptions) error {
	if err := discord.UninstallBD(types.UninstallOptions{FullUninstall: false}); err != nil {
		return err
	}

	bd, err := discord.GetBetterDiscordInstall()
	if err != nil {
		return err
	}

	if err := bd.Repair(discord.Channel, options); err != nil {
		return err
	}

	return nil
}

func (discord *DiscordInstall) GetBetterDiscordInstall() (*betterdiscord.BDInstall, error) {
	// Gets the global BetterDiscord install
	bd := betterdiscord.GetInstallation()

	// Snaps and flatpaks get their own local BD install
	if discord.IsSnap || discord.IsFlatpak {
		segment := "config"
		if discord.IsSnap {
			segment = ".config"
		}

		configPath, err := utils.FindSegment(discord.ResourcesPath, segment)
		if err != nil {
			return nil, err
		}
		bd = betterdiscord.GetInstallation(configPath)
	}

	return bd, nil
}
