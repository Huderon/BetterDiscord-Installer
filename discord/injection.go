package discord

import (
	_ "embed"
	"log"
	"os"
	"path/filepath"
	"strings"

	"installer/betterdiscord"
)

//go:embed assets/injection.js
var injectionScript string

func (discord *DiscordInstall) inject(bd *betterdiscord.BDInstall) error {

	// 0o644: index.js is a require target, it doesn't need the executable bit
	// (matches the mode uninject writes).
	if err := os.WriteFile(filepath.Join(discord.CorePath, "index.js"), []byte(injectionScript), 0o644); err != nil {
		log.Printf("❌ Unable to write index.js in %s\n", discord.CorePath)
		log.Printf("   %s\n", err.Error())
		return err
	}

	log.Printf("✅ Injected into %s\n", discord.CorePath)
	return nil
}

func (discord *DiscordInstall) uninject() error {
	indexFile := filepath.Join(discord.CorePath, "index.js")

	contents, err := os.ReadFile(indexFile)

	// First try to check the file, but if there's an issue we try to blindly overwrite below
	if err == nil {
		if !strings.Contains(strings.ToLower(string(contents)), "betterdiscord") {
			log.Printf("✅ No injection found for %s\n", discord.Channel.Name())
			return nil
		}
	}

	if err := os.WriteFile(indexFile, []byte(`module.exports = require("./core.asar");`), 0o644); err != nil {
		log.Printf("❌ Unable to write file %s\n", indexFile)
		log.Printf("   %s\n", err.Error())
		return err
	}
	log.Printf("✅ Removed from %s\n", discord.Channel.Name())

	return nil
}

// TODO: consider putting this in the betterdiscord package
func (discord *DiscordInstall) IsInjected() bool {
	indexFile := filepath.Join(discord.CorePath, "index.js")
	contents, err := os.ReadFile(indexFile)
	if err != nil {
		return false
	}
	lower := strings.ToLower(string(contents))
	return strings.Contains(lower, "betterdiscord")
}
