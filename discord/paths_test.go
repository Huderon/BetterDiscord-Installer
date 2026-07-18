package discord

import (
	"installer/types"
	"path/filepath"
	"testing"
)

func TestGetVersion(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "Finds version in Windows style path",
			path:     filepath.Join("C:", "Users", "Me", "AppData", "Local", "Discord", "app-1.0.9002", "modules", "discord_desktop_core-1", "discord_desktop_core"),
			expected: "1.0.9002",
		},
		{
			name:     "Finds version in Unix style path",
			path:     filepath.Join("/home/me/.config/discord", "0.0.35", "modules", "discord_desktop_core"),
			expected: "0.0.35",
		},
		{
			name:     "No version",
			path:     filepath.Join("/home/me/.config/discord", "modules", "discord_desktop_core"),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetVersion(tt.path)
			if result != tt.expected {
				t.Errorf("GetVersion(%q) = %q, expected %q", tt.path, result, tt.expected)
			}
		})
	}
}

func TestGetChannel(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected types.DiscordChannel
	}{
		{
			name:     "Stable path",
			path:     filepath.Join("/home/me/.config/discord", "0.0.35", "modules", "discord_desktop_core"),
			expected: types.Stable,
		},
		{
			name:     "Canary path",
			path:     filepath.Join("/home/me/.config/discordcanary", "0.0.90", "modules", "discord_desktop_core"),
			expected: types.Canary,
		},
		{
			name:     "PTB path",
			path:     filepath.Join("/home/me/.config/discordptb", "0.0.56", "modules", "discord_desktop_core"),
			expected: types.PTB,
		},
		{
			name:     "Unknown defaults to stable",
			path:     filepath.Join("/home/me/.config/discordunknown", "0.0.56", "modules", "discord_desktop_core"),
			expected: types.Stable,
		},
		{
			name:     "Windows path names",
			path:     filepath.Join("C:", "Users", "Me", "AppData", "Local", "DiscordCanary", "app-1.0.9002", "modules", "discord_desktop_core-1", "discord_desktop_core"),
			expected: types.Canary,
		},
		{
			name:     "macOS bundle name",
			path:     filepath.Join("/Applications", "Discord Canary.app", "Contents", "Resources"),
			expected: types.Canary,
		},
		{
			name:     "macOS stable bundle name",
			path:     filepath.Join("/Applications", "Discord.app", "Contents", "Resources"),
			expected: types.Stable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetChannel(tt.path)
			if result != tt.expected {
				t.Errorf("GetChannel(%q) = %v, expected %v", tt.path, result, tt.expected)
			}
		})
	}
}
