package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Hole struct {
	position rl.Vector2
	texture  rl.Texture2D
}

func NewHole() *Hole {
	pos := rl.Vector2{X: 100, Y: 300}

	texture := rl.LoadTexture("assets/hole.png")

	return &Hole{
		position: pos,
		texture:  texture,
	}
}

func (h *Hole) UpdateHole(dt float32) {

}

func (h *Hole) DrawHole() {
	rl.DrawTexture(h.texture, int32(h.position.X), int32(h.position.Y), rl.White)
}
