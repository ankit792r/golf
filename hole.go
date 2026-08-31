package main

import (
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	SinkSpeed      = 220
	HoleMinPadding = 40
)

type Hole struct {
	position rl.Vector2
	radius   float32
	texture  rl.Texture2D
}

func NewHole() *Hole {
	texture := rl.LoadTexture("assets/hole.png")
	hole := &Hole{
		radius:  float32(texture.Width) / 2,
		texture: texture,
	}
	hole.RandomizePosition(rl.Vector2{
		X: float32(rl.GetScreenWidth()) / 2,
		Y: float32(rl.GetScreenHeight()),
	})
	return hole
}

func (h *Hole) UpdateHole(dt float32) {
}

func (h *Hole) Contains(ball *Ball) bool {
	dist := rl.Vector2Distance(h.position, ball.position)
	speed := rl.Vector2Length(ball.velocity)
	return dist < h.radius*0.7 && speed < SinkSpeed
}

func (h *Hole) RandomizePosition(avoid rl.Vector2) {
	screenW := float32(rl.GetScreenWidth())
	screenH := float32(rl.GetScreenHeight())
	margin := h.radius + HoleMinPadding
	minDist := h.radius + BallRadius*4

	for range 50 {
		pos := rl.Vector2{
			X: margin + rand.Float32()*(screenW-2*margin),
			Y: margin + rand.Float32()*(screenH*0.55),
		}
		if rl.Vector2Distance(pos, avoid) >= minDist {
			h.position = pos
			return
		}
	}

	h.position = rl.Vector2{X: screenW / 2, Y: margin + 80}
}

func (h *Hole) DrawHole() {
	source := rl.Rectangle{
		Width:  float32(h.texture.Width),
		Height: float32(h.texture.Height),
	}

	dest := rl.Rectangle{
		X:      h.position.X,
		Y:      h.position.Y,
		Width:  float32(h.texture.Width),
		Height: float32(h.texture.Height),
	}

	origin := rl.Vector2{
		X: float32(h.texture.Width) / 2,
		Y: float32(h.texture.Height) / 2,
	}

	rl.DrawTexturePro(h.texture, source, dest, origin, 0, rl.White)
}
