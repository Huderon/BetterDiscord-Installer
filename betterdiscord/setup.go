package betterdiscord

import (
	"installer/types"
	"installer/utils"
	"log"
	"os"
	"path/filepath"
)

func makeDirectory(folder string) error {
	exists := utils.Exists(folder)

	if exists {
		log.Printf("✅ Directory exists: %s\n", folder)
		return nil
	}

	if err := os.MkdirAll(folder, 0755); err != nil {
		log.Printf("❌ Failed to create directory: %s\n", folder)
		log.Printf("   %s\n", err.Error())
		return err
	}

	log.Printf("✅ Directory created: %s\n", folder)
	return nil
}

func (i *BDInstall) prepare() error {
	if err := makeDirectory(i.data); err != nil {
		return err
	}
	if err := makeDirectory(i.plugins); err != nil {
		return err
	}
	if err := makeDirectory(i.themes); err != nil {
		return err
	}
	return nil
}

func (i *BDInstall) repair(channel types.DiscordChannel, options types.RepairOptions) error {
	channelFolder := filepath.Join(i.data, channel.String())

	if options.DisablePlugins {
		if err := removeFile(channelFolder, "plugins.json", "plugins", channel.Name()); err != nil {
			return err
		}
	}

	if options.DisableThemes {
		if err := removeFile(channelFolder, "themes.json", "themes", channel.Name()); err != nil {
			return err
		}
	}

	if options.ClearCustomCSS {
		if err := removeFile(channelFolder, "custom.css", "custom CSS", channel.Name()); err != nil {
			return err
		}
	}

	if options.ClearWebpackCache {
		if err := removeFile(channelFolder, "webpack.json", "webpack cache", channel.Name()); err != nil {
			return err
		}
	}

	if options.ClearAddonStoreCache {
		if err := removeFile(channelFolder, "addon-store.json", "addon store cache", channel.Name()); err != nil {
			return err
		}
	}

	if options.ResetSettings {
		if err := removeFile(channelFolder, "settings.json", "settings", channel.Name()); err != nil {
			return err
		}
	}

	return nil
}

func removeFile(basePath string, filename string, label string, channelName string) error {
	path := filepath.Join(basePath, filename)
	if !utils.Exists(path) {
		log.Printf("✅ No %s found for %s\n", label, channelName)
		return nil
	}

	if err := os.Remove(path); err != nil {
		log.Printf("❌ Unable to remove %s: %s\n", label, path)
		log.Printf("   %s\n", err.Error())
		return err
	}

	log.Printf("✅ Removed %s for %s\n", label, channelName)
	return nil
}

func removeDir(basePath string, dirname string, label string, channelName string) error {
	path := filepath.Join(basePath, dirname)
	if !utils.Exists(path) {
		log.Printf("✅ No %s found for %s\n", label, channelName)
		return nil
	}

	if err := os.RemoveAll(path); err != nil {
		log.Printf("❌ Unable to remove %s: %s\n", label, path)
		log.Printf("   %s\n", err.Error())
		return err
	}

	log.Printf("✅ Removed %s for %s\n", label, channelName)
	return nil
}
