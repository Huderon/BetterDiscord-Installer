package main

// func main() {
// TODO: de-dup the code across the 3 files
// installs := discord.GetAllInstalls()

// for c, install := range installs {
// 	fmt.Println(c.Name())
// 	for _, idk := range install {
// 		fmt.Println(idk)
// 	}
// 	fmt.Println("")
// }
// versions := []string{"1.0.324", "1.0.9524", "0.0.97"}
// slices.SortFunc(versions, func(a, b string) int {
// 	switch {
// 	case a > b:
// 		return -1
// 	case b > a:
// 		return 1
// 	}
// 	return 0
// })
// fmt.Println(versions)

// installs := discord.GetAllInstalls()
// for _, channel := range discord.Channels {
// 	subset := slices.Collect(func(yield func(*discord.DiscordInstall) bool) {
// 		for _, install := range installs {
// 			if install.IsChannel(channel) {
// 				if !yield(install) {
// 					return // triggered in "break"
// 				}
// 			}
// 		}
// 	})
// 	fmt.Println(channel.Name())
// 	fmt.Println(subset[0])
// }
// }
