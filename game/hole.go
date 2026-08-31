package game

import (
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	sinkSpeed      = 220
	holeMinPadding = 40
)

type hole struct {
	position rl.Vector2
	radius   float32
	texture  rl.Texture2D
}

func newHole(texture rl.Texture2D) *hole {
	return &hole{
		radius:  float32(texture.Width) / 2,
		texture: texture,
	}
}

func (h *hole) contains(b *ball) bool {
	dist := rl.Vector2Distance(h.position, b.position)
	speed := rl.Vector2Length(b.velocity)
	return dist < h.radius*0.7 && speed < sinkSpeed
}

func (h *hole) randomizePosition(avoid rl.Vector2, obstacles []*obstacle) {
	screenW := float32(rl.GetScreenWidth())
	screenH := float32(rl.GetScreenHeight())
	margin := h.radius + holeMinPadding
	minDist := h.radius + ballRadius*4

	for range 80 {
		pos := rl.Vector2{
			X: margin + rand.Float32()*(screenW-2*margin),
			Y: margin + rand.Float32()*(screenH*0.5),
		}
		if rl.Vector2Distance(pos, avoid) < minDist {
			continue
		}
		if circleHitsObstacles(pos, h.radius+12, obstacles) {
			continue
		}
		h.position = pos
		return
	}

	h.position = rl.Vector2{X: screenW / 2, Y: margin + 80}
}

func (h *hole) draw() {
	source := rl.Rectangle{Width: float32(h.texture.Width), Height: float32(h.texture.Height)}
	dest := rl.Rectangle{
		X:      h.position.X,
		Y:      h.position.Y,
		Width:  float32(h.texture.Width),
		Height: float32(h.texture.Height),
	}
	origin := rl.Vector2{X: float32(h.texture.Width) / 2, Y: float32(h.texture.Height) / 2}

	rl.DrawTexturePro(h.texture, source, dest, origin, 0, rl.White)
}
