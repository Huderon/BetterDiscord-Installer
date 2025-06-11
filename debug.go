package main

import (
	"fmt"
	"installer/backend/discord"
)

func main() {
	// TODO: de-dup the code across the 3 files
	installs := discord.GetAllInstalls()

	for _, install := range installs {
		fmt.Println(install)
	}
}
