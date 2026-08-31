package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	BallRadius  = 20
	ArrowHeight = 10
	ArrowWidth  = 60

	BallMaxPower = 800
	Friction     = 100
)

type BallState int

const (
	Aiming BallState = iota
	DirectionLocked
	Charging
	Moving
	Stopped
)

type Arrow struct {
	active bool
	angle  float32
	rect   rl.Rectangle

	texture rl.Texture2D
}

func NewArrow(pos *rl.Vector2) *Arrow {
	texture := rl.LoadTexture("assets/arrow.png")

	rect := rl.Rectangle{
		X:      pos.X,
		Y:      pos.Y,
		Width:  float32(texture.Width),
		Height: float32(texture.Height),
	}
	return &Arrow{
		active:  true,
		rect:    rect,
		texture: texture,
	}
}

type Ball struct {
	position rl.Vector2
	radius   float32

	// Charging
	power    float32
	maxPower float32
	chargeY  float32

	// Movement
	velocity rl.Vector2
	friction float32

	arrow *Arrow

	state   BallState
	texture rl.Texture2D
}

func NewBall() *Ball {
	screenH := rl.GetScreenHeight()
	screenW := rl.GetScreenWidth()

	pos := rl.Vector2{X: float32(screenW / 2), Y: float32(screenH - BallRadius - 40)}

	arrow := NewArrow(&pos)
	texture := rl.LoadTexture("assets/ball.png")

	return &Ball{
		position: pos,
		radius:   float32(texture.Height / 2),

		maxPower: BallMaxPower,
		friction: Friction,
		arrow:    arrow,

		state:   Aiming,
		texture: texture,
	}
}

func (b *Ball) UpdateBall(dt float32) {
	mousePos := rl.GetMousePosition()

	switch b.state {
	case Aiming:
		b.arrow.angle = float32(math.Atan2(
			float64(mousePos.Y-b.position.Y),
			float64(mousePos.X-b.position.X),
		))

		// First click locks direction.
		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			b.state = DirectionLocked
		}

	case DirectionLocked:
		// Do NOT update angle here.
		// It stays at the angle we locked.

		// Second click starts charging.
		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			b.state = Charging

			b.chargeY = mousePos.Y
			b.power = 0
		}

	case Charging:
		// Mouse UP = more power
		//
		// Mouse DOWN = less power
		deltaY := b.chargeY - mousePos.Y

		b.power = deltaY

		// Clamp
		b.power = float32(math.Max(
			0,
			math.Min(float64(b.power), float64(b.maxPower)),
		))

		// Release mouse = launch
		if rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
			b.LaunchBall()
		}

	case Moving:
		b.moveBall(dt)
	case Stopped:
		b.resetBall()
	}

}

func (b *Ball) LaunchBall() {
	direction := rl.Vector2{
		X: float32(math.Cos(float64(b.arrow.angle))),
		Y: float32(math.Sin(float64(b.arrow.angle))),
	}

	b.velocity = rl.Vector2Scale(direction, b.power)

	b.arrow.active = false
	b.state = Moving
}

func (b *Ball) moveBall(dt float32) {
	b.position = rl.Vector2Add(
		b.position,
		rl.Vector2Scale(b.velocity, dt),
	)

	b.CheckWallCollision()

	speed := rl.Vector2Length(b.velocity)

	if speed > 0 {
		speed -= b.friction * dt

		if speed < 0 {
			speed = 0
		}

		b.velocity = rl.Vector2Scale(
			rl.Vector2Normalize(b.velocity),
			speed,
		)
	}

	if speed < 5 {
		b.velocity = rl.Vector2Zero()
		b.state = Stopped
	}
}

func (b *Ball) resetBall() {
	b.state = Aiming

	b.arrow.active = true
	b.arrow.angle = 0
	b.arrow.rect.X = b.position.X
	b.arrow.rect.Y = b.position.Y

	b.power = 0
	b.velocity = rl.Vector2Zero()
}

func (b *Ball) CheckWallCollision() {
	screenW := float32(rl.GetScreenWidth())
	screenH := float32(rl.GetScreenHeight())

	// Left wall
	if b.position.X-b.radius <= 0 {
		b.position.X = b.radius
		b.velocity.X *= -1
	}

	// Right wall
	if b.position.X+b.radius >= screenW {
		b.position.X = screenW - b.radius
		b.velocity.X *= -1
	}

	// Top wall
	if b.position.Y-b.radius <= 0 {
		b.position.Y = b.radius
		b.velocity.Y *= -1
	}

	// Bottom wall
	if b.position.Y+b.radius >= screenH {
		b.position.Y = screenH - b.radius
		b.velocity.Y *= -1
	}
}

func (b *Ball) DrawBall() {
	source := rl.Rectangle{
		X:      0,
		Y:      0,
		Width:  float32(b.texture.Width),
		Height: float32(b.texture.Height),
	}

	dest := rl.Rectangle{
		X:      b.position.X,
		Y:      b.position.Y,
		Width:  float32(b.texture.Width),
		Height: float32(b.texture.Height),
	}

	origin := rl.Vector2{
		X: float32(b.texture.Width) / 2,
		Y: float32(b.texture.Height) / 2,
	}

	rl.DrawTexturePro(
		b.texture,
		source,
		dest,
		origin,
		0,
		rl.White,
	)

	// Draw arrow
	if b.arrow.active {
		b.DrawArrow()
	}
}

func (b *Ball) DrawArrow() {
	texture := b.arrow.texture

	source := rl.Rectangle{
		X:      0,
		Y:      0,
		Width:  float32(texture.Width),
		Height: float32(texture.Height),
	}

	dest := rl.Rectangle{
		X:      b.position.X,
		Y:      b.position.Y,
		Width:  float32(texture.Width),
		Height: float32(texture.Height),
	}

	origin := rl.Vector2{
		X: -float32(b.texture.Height),
		Y: float32(texture.Height) / 2,
	}

	rl.DrawTexturePro(
		texture,
		source,
		dest,
		origin,
		b.arrow.angle*rl.Rad2deg,
		rl.White,
	)
}

func (b *Ball) DrawPowerBar() {
	if b.state != Charging {
		return
	}

	screenW := float32(rl.GetScreenWidth())

	barWidth := float32(200)
	barHeight := float32(20)

	x := (screenW - barWidth) / 2
	y := float32(30)

	rl.DrawRectangle(
		int32(x),
		int32(y),
		int32(barWidth),
		int32(barHeight),
		rl.LightGray,
	)

	powerWidth := (b.power / b.maxPower) * barWidth

	rl.DrawRectangle(
		int32(x),
		int32(y),
		int32(powerWidth),
		int32(barHeight),
		rl.Red,
	)
}
