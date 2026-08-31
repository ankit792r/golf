package game

import (
	"embed"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func loadPNG(fs embed.FS, path string) rl.Texture2D {
	data, err := fs.ReadFile(path)
	if err != nil {
		panic("missing embedded asset " + path + ": " + err.Error())
	}

	img := rl.LoadImageFromMemory(".png", data, int32(len(data)))
	tex := rl.LoadTextureFromImage(img)
	rl.UnloadImage(img)
	return tex
}
