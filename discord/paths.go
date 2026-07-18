package discord

import (
	"installer/types"
	"installer/utils"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

var searchPaths []string
var versionRegex = regexp.MustCompile(`[0-9]+\.[0-9]+\.[0-9]+`)
var allDiscordInstalls map[types.DiscordChannel][]*DiscordInstall

func GetAllInstalls() map[types.DiscordChannel][]*DiscordInstall {
	allDiscordInstalls = map[types.DiscordChannel][]*DiscordInstall{}

	for _, path := range searchPaths {
		if result := Validate(path); result != nil {
			allDiscordInstalls[result.Channel] = append(allDiscordInstalls[result.Channel], result)
		}
	}

	sortInstalls()

	return allDiscordInstalls
}

func GetVersion(proposed string) string {
	for folder := range strings.SplitSeq(proposed, string(filepath.Separator)) {
		if version := versionRegex.FindString(folder); version != "" {
			return version
		}
	}
	return ""
}

func GetChannel(proposed string) types.DiscordChannel {
	for folder := range strings.SplitSeq(proposed, string(filepath.Separator)) {
		// Normalize the segment so macOS bundle names ("Discord Canary.app")
		// match the space-stripped channel names ("discordcanary").
		normalized := strings.ReplaceAll(strings.ToLower(folder), " ", "")
		normalized = strings.TrimSuffix(normalized, ".app")
		for _, channel := range types.Channels {
			if normalized == strings.ReplaceAll(strings.ToLower(channel.Name()), " ", "") {
				return channel
			}
		}
	}
	return types.Stable
}

func GetSuggestedPath(channel types.DiscordChannel) string {
	if len(allDiscordInstalls[channel]) > 0 {
		return allDiscordInstalls[channel][0].ResourcesPath
	}
	return ""
}

func AddCustomPath(proposed string) *DiscordInstall {
	result := Validate(proposed)
	if result == nil {
		return nil
	}

	// Check if this already exists in our list and return reference
	index := slices.IndexFunc(allDiscordInstalls[result.Channel], func(d *DiscordInstall) bool { return d.ResourcesPath == result.ResourcesPath })
	if index >= 0 {
		return allDiscordInstalls[result.Channel][index]
	}

	allDiscordInstalls[result.Channel] = append(allDiscordInstalls[result.Channel], result)

	sortInstalls()

	return result
}

func ResolvePath(proposed string) *DiscordInstall {
	for channel := range allDiscordInstalls {
		index := slices.IndexFunc(allDiscordInstalls[channel], func(d *DiscordInstall) bool { return d.ResourcesPath == proposed })
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
			// Descending (highest version first) with a numeric compare so
			// e.g. 1.0.10000 sorts above 1.0.9999.
			return utils.CompareVersions(b.Version, a.Version)
		})
	}
}
