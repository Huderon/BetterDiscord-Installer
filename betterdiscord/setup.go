package betterdiscord

import (
	"installer/types"
	"installer/utils"
	"log"
	"os"
	"path/filepath"
	"time"
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
	backupRoot := filepath.Join(channelFolder, "backup")
	backupStamp := time.Now().Format("20060102-150405")
	backupDir := filepath.Join(backupRoot, "repair-"+backupStamp)

	if options.DisablePlugins {
		if err := backupAndRemoveFile(channelFolder, "plugins.json", "plugins", channel.Name(), backupDir, true); err != nil {
			return err
		}
	}

	if options.DisableThemes {
		if err := backupAndRemoveFile(channelFolder, "themes.json", "themes", channel.Name(), backupDir, true); err != nil {
			return err
		}
	}

	if options.ClearCustomCSS {
		if err := backupAndRemoveFile(channelFolder, "custom.css", "custom CSS", channel.Name(), backupDir, false); err != nil {
			return err
		}
	}

	if options.ClearWebpackCache {
		if err := backupAndRemoveFile(channelFolder, "webpack.json", "webpack cache", channel.Name(), backupDir, true); err != nil {
			return err
		}
	}

	if options.ClearAddonStoreCache {
		if err := backupAndRemoveFile(channelFolder, "addon-store.json", "addon store cache", channel.Name(), backupDir, true); err != nil {
			return err
		}
	}

	if options.ResetSettings {
		if err := backupAndRemoveFile(channelFolder, "settings.json", "settings", channel.Name(), backupDir, false); err != nil {
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

func backupAndRemoveFile(basePath string, filename string, label string, channelName string, backupDir string, allowDeleteOnBackupFailure bool) error {
	path := filepath.Join(basePath, filename)
	if !utils.Exists(path) {
		log.Printf("✅ No %s found for %s\n", label, channelName)
		return nil
	}

	backupErr := backupFile(path, backupDir, label, channelName)
	if backupErr != nil && !allowDeleteOnBackupFailure {
		return backupErr
	}

	if backupErr != nil {
		log.Printf("⚠️  Proceeding without backup for %s\n", label)
	}

	if err := os.Remove(path); err != nil {
		log.Printf("❌ Unable to remove %s: %s\n", label, path)
		log.Printf("   %s\n", err.Error())
		return err
	}
	log.Printf("✅ Removed %s for %s\n", label, channelName)
	return nil
}

func backupFile(path string, backupDir string, label string, channelName string) error {
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		log.Printf("❌ Unable to create backup folder: %s\n", backupDir)
		log.Printf("   %s\n", err.Error())
		return err
	}

	contents, err := os.ReadFile(path)
	if err != nil {
		log.Printf("❌ Unable to read %s: %s\n", label, path)
		log.Printf("   %s\n", err.Error())
		return err
	}

	backupPath := filepath.Join(backupDir, filepath.Base(path))
	if err := os.WriteFile(backupPath, contents, 0644); err != nil {
		log.Printf("❌ Unable to backup %s to %s\n", label, backupPath)
		log.Printf("   %s\n", err.Error())
		return err
	}
	log.Printf("✅ Backed up %s for %s\n", label, channelName)
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
