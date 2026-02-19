package betterdiscord

import (
	"fmt"
	"installer/types"
	"installer/utils"
	"log"
)

func (i *BDInstall) download() error {
	if i.hasDownloaded {
		log.Printf("✅ Already downloaded to %s\n", i.asar)
		return nil
	}

	resp, err := utils.DownloadFile("https://betterdiscord.app/Download/betterdiscord.asar", i.asar)
	if err == nil {
		version := resp.Header.Get("x-bd-version")
		if version == "" {
			log.Println("✅ Downloaded BetterDiscord from the official website")
		} else {
			log.Printf("✅ Downloaded BetterDiscord version %s from the official website\n", utils.FormatVersion(version))
		}
		i.hasDownloaded = true
		return nil
	} else {
		log.Println("❌ Failed to download BetterDiscord from official website")
		log.Printf("❌ %s\n", err.Error())
		log.Println("")
		log.Println("🔁 Falling back to GitHub...")
	}

	// Get download URL from GitHub API
	apiData, err := utils.DownloadJSON[types.GitHubRelease]("https://api.github.com/repos/BetterDiscord/BetterDiscord/releases/latest")
	if err != nil {
		log.Println("❌ Failed to get asset url from GitHub")
		log.Printf("❌ %s\n", err.Error())
		return err
	}

	var index = -1
	for idx, asset := range apiData.Assets {
		if asset.Name == "betterdiscord.asar" {
			index = idx
			break
		}
	}

	if index == -1 {
		log.Println("❌ Failed to find the BetterDiscord asar on GitHub")
		return fmt.Errorf("failed to find betterdiscord.asar asset in GitHub release")
	}

	var downloadUrl = apiData.Assets[index].URL
	var version = apiData.TagName

	if downloadUrl != "" {
		log.Printf("✅ Found BetterDiscord: %s\n", downloadUrl)
	}

	// Download asar into the BD folder
	_, err = utils.DownloadFile(downloadUrl, i.asar)
	if err != nil {
		log.Println("❌ Failed to download BetterDiscord from GitHub")
		log.Printf("❌ %s\n", err.Error())
		return err
	}

	if version == "" {
		log.Println("✅ Downloaded BetterDiscord from GitHub")
	} else {
		log.Printf("✅ Downloaded BetterDiscord version %s from GitHub\n", utils.FormatVersion(version))
	}
	i.hasDownloaded = true

	return nil
}
