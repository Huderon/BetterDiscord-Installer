package discord

import (
	"runtime"
	"strings"
)

type DiscordChannel int64

const (
	Stable DiscordChannel = iota
	Canary
	PTB
)

var Channels = []DiscordChannel{Stable, Canary, PTB}

func (channel DiscordChannel) String() string {
	switch channel {
	case Stable:
		return "stable"
	case Canary:
		return "canary"
	case PTB:
		return "ptb"
	}
	return ""
}

func (channel DiscordChannel) TSName() string {
	return strings.ToUpper(channel.String())
}

func (channel DiscordChannel) Name() string {
	switch channel {
	case Stable:
		return "Discord"
	case Canary:
		return "Discord Canary"
	case PTB:
		return "Discord PTB"
	}
	return ""
}

func (channel DiscordChannel) Exe() string {
	name := channel.Name()

	if runtime.GOOS != "darwin" {
		name = strings.ReplaceAll(name, " ", "")
	}

	if runtime.GOOS == "windows" {
		name = name + ".exe"
	}

	return name
}

func ParseChannel(input string) DiscordChannel {
	switch strings.ToLower(input) {
	case "stable":
		return Stable
	case "canary":
		return Canary
	case "ptb":
		return PTB
	}
	return Stable
}
