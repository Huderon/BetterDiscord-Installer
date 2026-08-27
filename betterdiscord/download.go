package betterdiscord

import (
	"fmt"
	"log"

	"installer/types"
	"installer/utils"
)

// Endpoints for fetching the BetterDiscord asar. Declared as package vars so
// tests can point them at a local httptest server.
var (
	websiteAsarURL         = "https://betterdiscord.app/Download/betterdiscord.asar"
	githubLatestReleaseURL = "https://api.github.com/repos/BetterDiscord/BetterDiscord/releases/latest"
	// githubCanaryReleaseURL is the rolling pre-release tagged "canary" (rebuilt
	// on every merge to the development branch). It is GitHub-only — the website
	// has no mirror — so the dev-build path fetches it by tag rather than via the
	// "latest" endpoint, which excludes pre-releases by design.
	githubCanaryReleaseURL = "https://api.github.com/repos/BetterDiscord/BetterDiscord/releases/tags/canary"
)

func (i *BDInstall) download(useDevBuild bool) error {
	if i.hasDownloaded {
		log.Printf("✅ Already downloaded to %s\n", i.asar)
		return nil
	}

	// The development build lives only on GitHub, so skip the website leg and go
	// straight to the canary release. A failure here must NOT fall back to the
	// stable asar: a developer who asked for the dev build silently receiving
	// stable is a confusing, near-undetectable footgun.
	if useDevBuild {
		return i.downloadFromGitHubRelease(githubCanaryReleaseURL, "GitHub (development build)")
	}

	resp, err := utils.DownloadFile(websiteAsarURL, i.asar)
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

	return i.downloadFromGitHubRelease(githubLatestReleaseURL, "GitHub")
}

// downloadFromGitHubRelease fetches the release metadata at apiURL, locates the
// betterdiscord.asar asset, and downloads it into the BD folder. sourceLabel is
// used only for logging (e.g. "GitHub" or "GitHub (development build)").
func (i *BDInstall) downloadFromGitHubRelease(apiURL, sourceLabel string) error {
	apiData, err := utils.DownloadJSON[types.GitHubRelease](apiURL)
	if err != nil {
		log.Printf("❌ Failed to get asset url from %s\n", sourceLabel)
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
		log.Printf("❌ Failed to find the BetterDiscord asar on %s\n", sourceLabel)
		return fmt.Errorf("failed to find betterdiscord.asar asset in %s release", sourceLabel)
	}

	var downloadUrl = apiData.Assets[index].URL
	var version = apiData.TagName

	if downloadUrl != "" {
		log.Printf("✅ Found BetterDiscord: %s\n", downloadUrl)
	}

	// Download asar into the BD folder
	_, err = utils.DownloadFile(downloadUrl, i.asar)
	if err != nil {
		log.Printf("❌ Failed to download BetterDiscord from %s\n", sourceLabel)
		log.Printf("❌ %s\n", err.Error())
		return err
	}

	if version == "" {
		log.Printf("✅ Downloaded BetterDiscord from %s\n", sourceLabel)
	} else {
		log.Printf("✅ Downloaded BetterDiscord version %s from %s\n", utils.FormatVersion(version), sourceLabel)
	}
	i.hasDownloaded = true

	return nil
}
