package main

// import (
// 	"fmt"
// 	"installer/backend/discord"
// )

// func main() {
// 	// TODO: de-dup the code across the 3 files
// 	installs := discord.GetAllInstalls()

// 	for c, install := range installs {
// 		fmt.Println(c.Name())
// 		for _, idk := range install {
// 			fmt.Println(idk)
// 		}
// 		fmt.Println("")
// 	}

// 	// installs := discord.GetAllInstalls()
// 	// for _, channel := range discord.Channels {
// 	// 	subset := slices.Collect(func(yield func(*discord.DiscordInstall) bool) {
// 	// 		for _, install := range installs {
// 	// 			if install.IsChannel(channel) {
// 	// 				if !yield(install) {
// 	// 					return // triggered in "break"
// 	// 				}
// 	// 			}
// 	// 		}
// 	// 	})
// 	// 	fmt.Println(channel.Name())
// 	// 	fmt.Println(subset[0])
// 	// }
// }
