package discord

import (
	"installer/types"
	"os"
	"path/filepath"
	"testing"
)

func TestIsInjected(t *testing.T) {
	tmpDir := t.TempDir()
	corePath := filepath.Join(tmpDir, "discord_desktop_core")
	if err := os.MkdirAll(corePath, 0755); err != nil {
		t.Fatalf("Failed to create core path: %v", err)
	}
	indexFile := filepath.Join(corePath, "index.js")

	install := &DiscordInstall{
		CorePath: corePath,
		Channel:  types.Stable,
	}

	if install.IsInjected() {
		t.Fatalf("Expected IsInjected to be false with missing index.js")
	}

	if err := os.WriteFile(indexFile, []byte(`module.exports = require("./core.asar");`), 0644); err != nil {
		t.Fatalf("Failed to write index.js: %v", err)
	}
	if install.IsInjected() {
		t.Fatalf("Expected IsInjected to be false for default index.js")
	}

	if err := os.WriteFile(indexFile, []byte(`// BetterDiscord injected`), 0644); err != nil {
		t.Fatalf("Failed to write injection index.js: %v", err)
	}
	if !install.IsInjected() {
		t.Fatalf("Expected IsInjected to be true when BetterDiscord is present")
	}
}
