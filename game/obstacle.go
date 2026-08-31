package game

import (
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type obstacleKind int

const (
	kindSmall obstacleKind = iota
	kindMedium
	kindLarge
	kindXLarge
	kindCount
)

type obstacle struct {
	pos     rl.Vector2
	size    rl.Vector2
	texture rl.Texture2D
}

func (o *obstacle) bounds() rl.Rectangle {
	return rl.Rectangle{
		X:      o.pos.X - o.size.X/2,
		Y:      o.pos.Y - o.size.Y/2,
		Width:  o.size.X,
		Height: o.size.Y,
	}
}

func (o *obstacle) draw() {
	r := o.bounds()
	source := rl.Rectangle{Width: float32(o.texture.Width), Height: float32(o.texture.Height)}
	rl.DrawTexturePro(o.texture, source, r, rl.Vector2{}, 0, rl.White)
}

func (o *obstacle) bounce(b *ball) {
	r := o.bounds()
	closest := rl.Vector2{
		X: clamp(b.position.X, r.X, r.X+r.Width),
		Y: clamp(b.position.Y, r.Y, r.Y+r.Height),
	}
	delta := rl.Vector2Subtract(b.position, closest)
	dist := rl.Vector2Length(delta)

	if dist < 0.001 {
		left := b.position.X - r.X
		right := r.X + r.Width - b.position.X
		top := b.position.Y - r.Y
		bottom := r.Y + r.Height - b.position.Y

		n := rl.Vector2{X: -1, Y: 0}
		minEdge := left
		if right < minEdge {
			minEdge = right
			n = rl.Vector2{X: 1, Y: 0}
		}
		if top < minEdge {
			minEdge = top
			n = rl.Vector2{X: 0, Y: -1}
		}
		if bottom < minEdge {
			minEdge = bottom
			n = rl.Vector2{X: 0, Y: 1}
		}

		b.position = rl.Vector2Add(b.position, rl.Vector2Scale(n, minEdge+b.radius))
		reflectVelocity(b, n)
		return
	}

	if dist >= b.radius {
		return
	}

	n := rl.Vector2Normalize(delta)
	b.position = rl.Vector2Add(b.position, rl.Vector2Scale(n, b.radius-dist))
	reflectVelocity(b, n)
}

func reflectVelocity(b *ball, n rl.Vector2) {
	vn := rl.Vector2DotProduct(b.velocity, n)
	if vn < 0 {
		b.velocity = rl.Vector2Subtract(b.velocity, rl.Vector2Scale(n, 2*vn))
	}
}

func (o *obstacle) overlapsCircle(center rl.Vector2, radius float32) bool {
	r := o.bounds()
	closest := rl.Vector2{
		X: clamp(center.X, r.X, r.X+r.Width),
		Y: clamp(center.Y, r.Y, r.Y+r.Height),
	}
	return rl.Vector2Distance(center, closest) < radius
}

func (o *obstacle) overlapsObstacle(other *obstacle, pad float32) bool {
	a := o.bounds()
	b := other.bounds()
	return a.X < b.X+b.Width+pad &&
		a.X+a.Width+pad > b.X &&
		a.Y < b.Y+b.Height+pad &&
		a.Y+a.Height+pad > b.Y
}

func circleHitsObstacles(center rl.Vector2, radius float32, obstacles []*obstacle) bool {
	for _, o := range obstacles {
		if o.overlapsCircle(center, radius) {
			return true
		}
	}
	return false
}

func spawnObstacles(level int, textures [kindCount]rl.Texture2D, avoid rl.Vector2) []*obstacle {
	kinds := kindsForLevel(level)
	if len(kinds) == 0 {
		return nil
	}

	screenW := float32(rl.GetScreenWidth())
	screenH := float32(rl.GetScreenHeight())
	placed := make([]*obstacle, 0, len(kinds))

	for _, kind := range kinds {
		tex := textures[kind]
		size := rl.Vector2{X: float32(tex.Width), Y: float32(tex.Height)}
		marginX := size.X/2 + 24
		minY := 130 + size.Y/2
		maxY := screenH - 180 - size.Y/2
		if maxY <= minY {
			continue
		}

		var next *obstacle
		for range 80 {
			candidate := &obstacle{
				pos: rl.Vector2{
					X: marginX + rand.Float32()*(screenW-2*marginX),
					Y: minY + rand.Float32()*(maxY-minY),
				},
				size:    size,
				texture: tex,
			}
			if candidate.overlapsCircle(avoid, ballRadius*5) {
				continue
			}
			ok := true
			for _, existing := range placed {
				if candidate.overlapsObstacle(existing, 20) {
					ok = false
					break
				}
			}
			if !ok {
				continue
			}
			next = candidate
			break
		}
		if next != nil {
			placed = append(placed, next)
		}
	}

	return placed
}

func clamp(v, lo, hi float32) float32 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
