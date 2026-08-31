package main

import (
	"embed"

	"golf/game"
)

//go:embed assets/*.png
var assets embed.FS

func main() {
	game.Run(assets)
}
