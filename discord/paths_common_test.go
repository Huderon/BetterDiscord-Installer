package discord

import (
	"os"
	"path/filepath"
	"testing"
)

func writeCoreAsar(t *testing.T, corePath string) {
	t.Helper()
	if err := os.MkdirAll(corePath, 0755); err != nil {
		t.Fatalf("Failed to create core path: %v", err)
	}
	if err := os.WriteFile(filepath.Join(corePath, "core.asar"), []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to write core.asar: %v", err)
	}
}

func TestValidateWindowsStyleInstall_FromDiscordRoot(t *testing.T) {
	tmpDir := t.TempDir()
	root := filepath.Join(tmpDir, "Discord")
	versionDir := filepath.Join(root, "app-1.0.9002")
	coreWrap := filepath.Join(versionDir, "modules", "discord_desktop_core-1", "discord_desktop_core")

	writeCoreAsar(t, coreWrap)

	result := validateWindowsStyleInstall(root)
	if result == nil {
		t.Fatalf("Expected install for %s", root)
	}
	if result.CorePath != coreWrap {
		t.Errorf("CorePath = %s, expected %s", result.CorePath, coreWrap)
	}
}

func TestValidateWindowsStyleInstall_FromAppFolder(t *testing.T) {
	tmpDir := t.TempDir()
	root := filepath.Join(tmpDir, "Discord")
	versionDir := filepath.Join(root, "app-1.0.9002")
	coreWrap := filepath.Join(versionDir, "modules", "discord_desktop_core-1", "discord_desktop_core")

	writeCoreAsar(t, coreWrap)

	result := validateWindowsStyleInstall(versionDir)
	if result == nil {
		t.Fatalf("Expected install for %s", versionDir)
	}
	if result.CorePath != coreWrap {
		t.Errorf("CorePath = %s, expected %s", result.CorePath, coreWrap)
	}
}

func TestValidateWindowsStyleInstall_FromCoreFolder(t *testing.T) {
	tmpDir := t.TempDir()
	corePath := filepath.Join(tmpDir, "discord_desktop_core")
	writeCoreAsar(t, corePath)

	result := validateWindowsStyleInstall(corePath)
	if result == nil {
		t.Fatalf("Expected install for %s", corePath)
	}
	if result.CorePath != corePath {
		t.Errorf("CorePath = %s, expected %s", result.CorePath, corePath)
	}
}

func TestValidateWindowsStyleInstall_MissingAsar(t *testing.T) {
	tmpDir := t.TempDir()
	root := filepath.Join(tmpDir, "Discord")
	versionDir := filepath.Join(root, "app-1.0.9002")
	coreWrap := filepath.Join(versionDir, "modules", "discord_desktop_core-1", "discord_desktop_core")

	if err := os.MkdirAll(coreWrap, 0755); err != nil {
		t.Fatalf("Failed to create core path: %v", err)
	}

	result := validateWindowsStyleInstall(root)
	if result != nil {
		t.Fatalf("Expected no install when core.asar is missing")
	}
}

func TestValidateUnixStyleInstall_FromDiscordRoot(t *testing.T) {
	tmpDir := t.TempDir()
	root := filepath.Join(tmpDir, "discord")
	corePath := filepath.Join(root, "0.0.35", "modules", "discord_desktop_core")

	writeCoreAsar(t, corePath)

	result := validateUnixStyleInstall(root, true, true)
	if result == nil {
		t.Fatalf("Expected install for %s", root)
	}
	if result.CorePath != corePath {
		t.Errorf("CorePath = %s, expected %s", result.CorePath, corePath)
	}
}

func TestValidateUnixStyleInstall_FromVersionFolder(t *testing.T) {
	tmpDir := t.TempDir()
	root := filepath.Join(tmpDir, "discord")
	versionDir := filepath.Join(root, "0.0.35")
	corePath := filepath.Join(versionDir, "modules", "discord_desktop_core")

	writeCoreAsar(t, corePath)

	result := validateUnixStyleInstall(versionDir, true, true)
	if result == nil {
		t.Fatalf("Expected install for %s", versionDir)
	}
	if result.CorePath != corePath {
		t.Errorf("CorePath = %s, expected %s", result.CorePath, corePath)
	}
}

func TestValidateUnixStyleInstall_FromModulesFolder(t *testing.T) {
	tmpDir := t.TempDir()
	corePath := filepath.Join(tmpDir, "modules", "discord_desktop_core")

	writeCoreAsar(t, corePath)

	result := validateUnixStyleInstall(filepath.Join(tmpDir, "modules"), true, true)
	if result == nil {
		t.Fatalf("Expected install for %s", filepath.Join(tmpDir, "modules"))
	}
	if result.CorePath != corePath {
		t.Errorf("CorePath = %s, expected %s", result.CorePath, corePath)
	}
}

func TestValidateUnixStyleInstall_FlatpakDetection(t *testing.T) {
	tmpDir := t.TempDir()
	root := filepath.Join(tmpDir, "com.discordapp.Discord", "config", "discord")
	corePath := filepath.Join(root, "0.0.35", "modules", "discord_desktop_core")

	writeCoreAsar(t, corePath)

	result := validateUnixStyleInstall(root, true, false)
	if result == nil {
		t.Fatalf("Expected install for %s", root)
	}
	if !result.IsFlatpak {
		t.Fatalf("Expected flatpak detection")
	}
	if result.IsSnap {
		t.Fatalf("Did not expect snap detection")
	}
}

func TestValidateUnixStyleInstall_SnapDetection(t *testing.T) {
	tmpDir := t.TempDir()
	root := filepath.Join(tmpDir, "snap", "discord", "current", ".config", "discord")
	corePath := filepath.Join(root, "0.0.35", "modules", "discord_desktop_core")

	writeCoreAsar(t, corePath)

	result := validateUnixStyleInstall(root, false, true)
	if result == nil {
		t.Fatalf("Expected install for %s", root)
	}
	if !result.IsSnap {
		t.Fatalf("Expected snap detection")
	}
	if result.IsFlatpak {
		t.Fatalf("Did not expect flatpak detection")
	}
}
