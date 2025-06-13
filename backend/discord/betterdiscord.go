package discord

import (
	"installer/backend/utils"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type BetterDiscord struct {
	root          string
	data          string
	asar          string
	plugins       string
	themes        string
	hasDownloaded bool
}

var lock = &sync.Mutex{}
var globalBetterDiscord *BetterDiscord

func getGlobalInstance() *BetterDiscord {
	if globalBetterDiscord != nil {
		return globalBetterDiscord
	}

	lock.Lock()
	defer lock.Unlock()
	if globalBetterDiscord != nil {
		return globalBetterDiscord
	}

	configDir, _ := os.UserConfigDir()
	globalBetterDiscord = GetBetterDiscord(configDir)

	return globalBetterDiscord
}

func CreateBetterDiscord(root string) *BetterDiscord {
	return &BetterDiscord{
		root:          root,
		data:          filepath.Join(root, "data"),
		asar:          filepath.Join(root, "data", "betterdiscord.asar"),
		plugins:       filepath.Join(root, "plugins"),
		themes:        filepath.Join(root, "themes"),
		hasDownloaded: false,
	}
}

func GetBetterDiscord(base ...string) *BetterDiscord {
	if len(base) == 0 {
		return getGlobalInstance()
	}
	return CreateBetterDiscord(filepath.Join(base[0], "BetterDiscord"))
}

func makeDirectory(folder string) error {
	exists := utils.Exists(folder)

	if exists {
		log.Printf("✅ Directory exists: %s", folder)
		return nil
	}

	if err := os.MkdirAll(folder, 0755); err != nil {
		log.Printf("❌ Failed to create directory: %s", folder)
		log.Printf("❌ %s", err.Error())
		return err
	}

	log.Printf("✅ Directory created: %s", folder)
	return nil
}

func (bd *BetterDiscord) prepare() error {
	if err := makeDirectory(bd.data); err != nil {
		return err
	}
	if err := makeDirectory(bd.plugins); err != nil {
		return err
	}
	if err := makeDirectory(bd.themes); err != nil {
		return err
	}
	return nil
}

func (bd *BetterDiscord) download() error {

	if bd.hasDownloaded {
		log.Printf("✅ Already downloaded to %s", bd.asar)
		return nil
	}

	resp, err := utils.DownloadFile("https://betterdiscord.app/Download/betterdiscord.asar", bd.asar)
	if err == nil {
		version := resp.Header.Get("x-bd-version")
		log.Printf("✅ Downloaded BetterDiscord version %s from the official website", version)
		return nil
	} else {
		log.Printf("❌ Failed to download BetterDiscord from official website")
		log.Printf("❌ %s", err.Error())
		log.Printf("")
		log.Printf("#### Falling back to GitHub...")
	}

	// Get download URL from GitHub API
	apiData, err := utils.DownloadJSON[utils.Release]("https://api.github.com/repos/BetterDiscord/BetterDiscord/releases/latest")
	if err != nil {
		log.Printf("❌ Failed to get asset url from GitHub")
		log.Printf("❌ %s", err.Error())
		return err
	}

	var index = 0
	for i, asset := range apiData.Assets {
		if asset.Name == "betterdiscord.asar" {
			index = i
			break
		}
	}

	var downloadUrl = apiData.Assets[index].URL
	var version = apiData.TagName

	if downloadUrl != "" {
		log.Printf("✅ Found BetterDiscord: %s", downloadUrl)
	}

	// Download asar into the BD folder
	_, err = utils.DownloadFile(downloadUrl, bd.asar)
	if err != nil {
		log.Printf("❌ Failed to download BetterDiscord from GitHub")
		log.Printf("❌ %s", err.Error())
		return err
	}

	log.Printf("✅ Downloaded BetterDiscord version %s from GitHub", version)
	bd.hasDownloaded = true

	return nil
}

func (bd *BetterDiscord) repair(channel DiscordChannel) error {
	channelFolder := filepath.Join(bd.data, channel.String())
	pluginsJson := filepath.Join(channelFolder, "plugins.json")

	if !utils.Exists(pluginsJson) {
		log.Printf("✅ No plugins enabled for %s", channel.Name())
		return nil
	}

	if err := os.Remove(pluginsJson); err != nil {
		log.Printf("❌ Unable to remove file %s", pluginsJson)
		log.Printf("❌ %s", err.Error())
		return err
	}

	log.Printf("✅ Plugins disabled for %s", channel.Name())
	return nil
}
