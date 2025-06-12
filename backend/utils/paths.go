package utils

import (
	"os"
	"path/filepath"
)

var Roaming string
var BetterDiscord string
var Data string
var Asar string
var Plugins string
var Themes string

func init() {
	var configDir, err = os.UserConfigDir()
	if err != nil {
		return
	}
	Roaming = configDir
	BetterDiscord = filepath.Join(configDir, "BetterDiscord")
	Data = filepath.Join(BetterDiscord, "data")
	Asar = filepath.Join(Data, "betterdiscord.asar")
	Plugins = filepath.Join(BetterDiscord, "plugins")
	Themes = filepath.Join(BetterDiscord, "themes")
}

func Exists(path string) bool {
	var _, err = os.Stat(path)
	return err == nil
}

func Filter[T any](source []T, filterFunc func(T) bool) (ret []T) {
	var returnArray = []T{}
	for _, s := range source {
		if filterFunc(s) {
			returnArray = append(ret, s)
		}
	}
	return returnArray
}
