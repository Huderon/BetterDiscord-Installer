package discord

import (
	"os"
	"path/filepath"

	"installer/types"
)

func init() {
	home, err := os.UserHomeDir()

	// On macOS the app.asar lives inside the application bundle
	// (Discord.app/Contents/Resources), not under Application Support. Search the
	// standard install locations for each channel's bundle.
	bases := []string{
		filepath.Join("/", "Applications"),
	}
	// Only add ~/Applications when the home dir resolved; otherwise the join would
	// produce a relative "Applications" and search the current working directory.
	if err == nil && home != "" {
		bases = append(bases, filepath.Join(home, "Applications"))
	}

	for _, channel := range types.Channels {
		bundle := channel.Name() + ".app"
		for _, base := range bases {
			searchPaths = append(searchPaths, filepath.Join(base, bundle))
		}
	}

	allDiscordInstalls = GetAllInstalls()
}

func Validate(proposed string) *DiscordInstall {
	return validateUnixStyleInstall(proposed, false, false)
}

// DefaultBrowseDir is where the "browse for Discord" dialog should open: app
// bundles live in /Applications.
func DefaultBrowseDir() string {
	return filepath.Join("/", "Applications")
}
