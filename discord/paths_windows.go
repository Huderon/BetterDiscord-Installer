package discord

import (
	"installer/types"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	paths := []string{
		filepath.Join(os.Getenv("LOCALAPPDATA"), "{channel}"),
		filepath.Join(os.Getenv("PROGRAMDATA"), os.Getenv("USERNAME"), "{channel}"),
	}

	for _, channel := range types.Channels {
		for _, path := range paths {
			searchPaths = append(
				searchPaths,
				strings.ReplaceAll(path, "{channel}", strings.ReplaceAll(channel.Name(), " ", "")),
			)
		}
	}

	allDiscordInstalls = GetAllInstalls()
}

func Validate(proposed string) *DiscordInstall {
	return validateWindowsStyleInstall(proposed)
}
