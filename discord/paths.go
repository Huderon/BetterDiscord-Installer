package discord

import (
	"installer/types"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

var searchPaths []string
var versionRegex = regexp.MustCompile(`[0-9]+\.[0-9]+\.[0-9]+`)
var allDiscordInstalls map[types.DiscordChannel][]*DiscordInstall

func GetAllInstalls() map[types.DiscordChannel][]*DiscordInstall {
	var installs = map[types.DiscordChannel][]*DiscordInstall{}

	for _, path := range searchPaths {
		if result := Validate(path); result != nil {
			installs[result.channel] = append(installs[result.channel], result)
		}
	}

	sortInstalls()

	return installs
}

func GetVersion(proposed string) string {
	for _, folder := range strings.Split(proposed, string(filepath.Separator)) {
		if version := versionRegex.FindString(folder); version != "" {
			return version
		}
	}
	return ""
}

func GetChannel(proposed string) types.DiscordChannel {
	for _, folder := range strings.Split(proposed, string(filepath.Separator)) {
		for _, channel := range types.Channels {
			if strings.ToLower(folder) == strings.ReplaceAll(strings.ToLower(channel.Name()), " ", "") {
				return channel
			}
		}
	}
	return types.Stable
}

func GetSuggestedPath(channel types.DiscordChannel) string {
	if len(allDiscordInstalls[channel]) > 0 {
		return allDiscordInstalls[channel][0].corePath
	}
	return ""
}

func AddCustomPath(proposed string) *DiscordInstall {
	result := Validate(proposed)
	if result == nil {
		return nil
	}

	// Check if this already exists in our list and return reference
	index := slices.IndexFunc(allDiscordInstalls[result.channel], func(d *DiscordInstall) bool { return d.corePath == result.corePath })
	if index >= 0 {
		return allDiscordInstalls[result.channel][index]
	}

	allDiscordInstalls[result.channel] = append(allDiscordInstalls[result.channel], result)

	sortInstalls()

	return result
}

func ResolvePath(proposed string) *DiscordInstall {
	for channel := range allDiscordInstalls {
		index := slices.IndexFunc(allDiscordInstalls[channel], func(d *DiscordInstall) bool { return d.corePath == proposed })
		if index >= 0 {
			return allDiscordInstalls[channel][index]
		}
	}

	// If it wasn't found as an existing install, try to add it
	return AddCustomPath(proposed)
}

func sortInstalls() {
	for channel := range allDiscordInstalls {
		slices.SortFunc(allDiscordInstalls[channel], func(a, b *DiscordInstall) int {
			switch {
			case a.version > b.version:
				return -1
			case b.version > a.version:
				return 1
			}
			return 0
		})
	}
}
