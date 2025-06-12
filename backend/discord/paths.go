package discord

import (
	"path/filepath"
	"regexp"
	"strings"
)

var searchPaths []string
var versionRegex = regexp.MustCompile(`[0-9]+\.[0-9]+\.[0-9]+`)

func GetAllInstalls() map[DiscordChannel][]*DiscordInstall {
	var installs = map[DiscordChannel][]*DiscordInstall{}

	for _, path := range searchPaths {
		if result := Validate(path); result != nil {
			installs[result.channel] = append(installs[result.channel], result)
		}
	}

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

func GetChannel(proposed string) DiscordChannel {
	for _, folder := range strings.Split(proposed, string(filepath.Separator)) {
		for _, channel := range Channels {
			if strings.ToLower(folder) == strings.ReplaceAll(strings.ToLower(channel.Name()), " ", "") {
				return channel
			}
		}
	}
	return Stable
}
