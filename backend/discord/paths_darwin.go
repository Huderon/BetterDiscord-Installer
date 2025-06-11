package discord

import (
	"installer/backend/utils"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var searchPaths []string
var versionRegex = regexp.MustCompile(`[0-9]+\.[0-9]+\.[0-9]+`)

func init() {
	config, _ := os.UserConfigDir()
	paths := []string{
		filepath.Join(config, "{channel}"),
	}

	for _, channel := range Channels {
		for _, path := range paths {
			folder := strings.ReplaceAll(strings.ToLower(channel.Name()), " ", "")
			searchPaths = append(
				searchPaths,
				strings.ReplaceAll(path, "{channel}", folder),
			)
		}
	}
}

func GetAllInstalls() []*DiscordInstall {
	var installs []*DiscordInstall

	for _, path := range searchPaths {
		if result := Validate(path); result != nil {
			installs = append(installs, result)
		}
	}

	return installs
}

func GetChannel(proposed string) DiscordChannel {
	for _, folder := range strings.Split(proposed, string(filepath.Separator)) {
		for _, channel := range Channels {
			if folder == strings.ReplaceAll(strings.ToLower(channel.Name()), " ", "") {
				return channel
			}
		}
	}
	return Stable
}

func GetVersion(proposed string) string {
	for _, folder := range strings.Split(proposed, string(filepath.Separator)) {
		if versionRegex.MatchString(folder) {
			return folder
		}
	}
	return ""
}

func Validate(proposed string) *DiscordInstall {
	var finalPath = ""
	var selected = filepath.Base(proposed)
	if strings.HasPrefix(selected, "discord") {
		// Get version dir like 1.0.9002
		var dFiles, err = os.ReadDir(proposed)
		if err != nil {
			return nil
		}

		var candidates = utils.Filter(dFiles, func(file fs.DirEntry) bool { return file.IsDir() && versionRegex.MatchString(file.Name()) })
		sort.Slice(candidates, func(i, j int) bool { return candidates[i].Name() < candidates[j].Name() })
		var versionDir = candidates[len(candidates)-1].Name()
		finalPath = filepath.Join(proposed, versionDir, "modules", "discord_desktop_core")
	}

	if len(strings.Split(selected, ".")) == 3 {
		finalPath = filepath.Join(proposed, "modules", "discord_desktop_core")
	}

	if selected == "modules" {
		finalPath = filepath.Join(proposed, "discord_desktop_core")
	}

	if selected == "discord_desktop_core" {
		finalPath = proposed
	}

	// If the path and the asar exist, all good
	if utils.Exists(finalPath) && utils.Exists(filepath.Join(finalPath, "core.asar")) {
		return &DiscordInstall{
			corePath:  finalPath,
			channel:   GetChannel(finalPath),
			version:   GetVersion(finalPath),
			isFlatpak: false,
			isSnap:    false,
		}
	}

	return nil
}
