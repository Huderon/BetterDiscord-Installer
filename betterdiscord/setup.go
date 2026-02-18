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

func (i *BDInstall) repair(channel types.DiscordChannel) error {
	channelFolder := filepath.Join(i.data, channel.String())
	pluginsJson := filepath.Join(channelFolder, "plugins.json")

	if !utils.Exists(pluginsJson) {
		log.Printf("✅ No plugins enabled for %s\n", channel.Name())
		return nil
	}

	if err := os.Remove(pluginsJson); err != nil {
		log.Printf("❌ Unable to remove file %s\n", pluginsJson)
		log.Printf("   %s\n", err.Error())
		return err
	}

	log.Printf("✅ Plugins disabled for %s\n", channel.Name())
	return nil
}
