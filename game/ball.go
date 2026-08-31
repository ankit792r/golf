package game

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	ballRadius   = 20
	ballMaxPower = 800
	friction     = 100

	sinkDuration = 1.0
	sinkShakeAmp = 8
)

type ballState int

const (
	aiming ballState = iota
	directionLocked
	charging
	moving
	stopped
	sinking
)

type arrow struct {
	active  bool
	angle   float32
	rect    rl.Rectangle
	texture rl.Texture2D
}

func newArrow(pos rl.Vector2, texture rl.Texture2D) *arrow {
	return &arrow{
		active:  true,
		rect:    rl.Rectangle{X: pos.X, Y: pos.Y, Width: float32(texture.Width), Height: float32(texture.Height)},
		texture: texture,
	}
}

type ball struct {
	position rl.Vector2
	startPos rl.Vector2
	radius   float32

	power    float32
	maxPower float32
	chargeY  float32

	velocity rl.Vector2
	friction float32

	arrow   *arrow
	state   ballState
	texture rl.Texture2D

	scale      float32
	sinkTime   float32
	sinkFrom   rl.Vector2
	sinkTarget rl.Vector2
	shakeOff   rl.Vector2
	shakeRot   float32
}

func newBall(texture, arrowTexture rl.Texture2D) *ball {
	screenH := rl.GetScreenHeight()
	screenW := rl.GetScreenWidth()

	pos := rl.Vector2{X: float32(screenW / 2), Y: float32(screenH - ballRadius - 40)}

	return &ball{
		position: pos,
		startPos: pos,
		radius:   float32(texture.Height / 2),
		maxPower: ballMaxPower,
		friction: friction,
		arrow:    newArrow(pos, arrowTexture),
		state:    aiming,
		texture:  texture,
		scale:    1,
	}
}

func (b *ball) update(dt float32, obstacles []*obstacle) {
	mousePos := rl.GetMousePosition()

	switch b.state {
	case aiming:
		b.arrow.angle = float32(math.Atan2(
			float64(mousePos.Y-b.position.Y),
			float64(mousePos.X-b.position.X),
		))

		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			b.state = directionLocked
		}

	case directionLocked:
		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			b.state = charging
			b.chargeY = mousePos.Y
			b.power = 0
		}

	case charging:
		deltaY := b.chargeY - mousePos.Y
		b.power = float32(math.Max(0, math.Min(float64(deltaY), float64(b.maxPower))))

		if rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
			b.launch()
		}

	case moving:
		b.move(dt, obstacles)
	case stopped:
		b.resetAim()
	case sinking:
		b.updateSink(dt)
	}
}

func (b *ball) startSink(target rl.Vector2) {
	b.state = sinking
	b.arrow.active = false
	b.velocity = rl.Vector2Zero()
	b.power = 0
	b.sinkTime = 0
	b.sinkFrom = b.position
	b.sinkTarget = target
	b.scale = 1
	b.shakeOff = rl.Vector2{}
	b.shakeRot = 0
}

func (b *ball) updateSink(dt float32) {
	b.sinkTime += dt
	t := b.sinkTime / sinkDuration
	if t > 1 {
		t = 1
	}

	eased := 1 - (1-t)*(1-t)*(1-t)
	b.position = rl.Vector2Lerp(b.sinkFrom, b.sinkTarget, eased)
	b.scale = (1 - t) * (1 - t)

	amp := sinkShakeAmp * (1 - t)
	w := float64(b.sinkTime * 2 * math.Pi)
	b.shakeOff = rl.Vector2{
		X: float32(math.Sin(w*26)) * amp,
		Y: float32(math.Cos(w*21)) * amp,
	}
	b.shakeRot = float32(math.Sin(w*18)) * 22 * (1 - t)
}

func (b *ball) finishedSink() bool {
	return b.state == sinking && b.sinkTime >= sinkDuration
}

func (b *ball) launch() {
	direction := rl.Vector2{
		X: float32(math.Cos(float64(b.arrow.angle))),
		Y: float32(math.Sin(float64(b.arrow.angle))),
	}

	b.velocity = rl.Vector2Scale(direction, b.power)
	b.arrow.active = false
	b.state = moving
}

func (b *ball) move(dt float32, obstacles []*obstacle) {
	b.position = rl.Vector2Add(b.position, rl.Vector2Scale(b.velocity, dt))
	b.checkWallCollision()

	for _, o := range obstacles {
		o.bounce(b)
	}

	speed := rl.Vector2Length(b.velocity)
	if speed > 0 {
		speed -= b.friction * dt
		if speed < 0 {
			speed = 0
		}
		b.velocity = rl.Vector2Scale(rl.Vector2Normalize(b.velocity), speed)
	}

	if speed < 5 {
		b.velocity = rl.Vector2Zero()
		b.state = stopped
	}
}

func (b *ball) resetAim() {
	b.state = aiming
	b.arrow.active = true
	b.arrow.angle = 0
	b.arrow.rect.X = b.position.X
	b.arrow.rect.Y = b.position.Y
	b.power = 0
	b.velocity = rl.Vector2Zero()
}

func (b *ball) resetToStart() {
	b.position = b.startPos
	b.scale = 1
	b.shakeOff = rl.Vector2{}
	b.shakeRot = 0
	b.sinkTime = 0
	b.resetAim()
}

func (b *ball) checkWallCollision() {
	screenW := float32(rl.GetScreenWidth())
	screenH := float32(rl.GetScreenHeight())

	if b.position.X-b.radius <= 0 {
		b.position.X = b.radius
		b.velocity.X *= -1
	}
	if b.position.X+b.radius >= screenW {
		b.position.X = screenW - b.radius
		b.velocity.X *= -1
	}
	if b.position.Y-b.radius <= 0 {
		b.position.Y = b.radius
		b.velocity.Y *= -1
	}
	if b.position.Y+b.radius >= screenH {
		b.position.Y = screenH - b.radius
		b.velocity.Y *= -1
	}
}

func (b *ball) draw() {
	if b.scale <= 0.01 {
		return
	}

	pos := rl.Vector2Add(b.position, b.shakeOff)
	w := float32(b.texture.Width) * b.scale
	h := float32(b.texture.Height) * b.scale
	source := rl.Rectangle{Width: float32(b.texture.Width), Height: float32(b.texture.Height)}
	dest := rl.Rectangle{X: pos.X, Y: pos.Y, Width: w, Height: h}
	origin := rl.Vector2{X: w / 2, Y: h / 2}

	rl.DrawTexturePro(b.texture, source, dest, origin, b.shakeRot, rl.White)

	if b.arrow.active {
		b.drawArrow()
	}
}

func (b *ball) drawArrow() {
	texture := b.arrow.texture
	source := rl.Rectangle{Width: float32(texture.Width), Height: float32(texture.Height)}
	dest := rl.Rectangle{
		X:      b.position.X,
		Y:      b.position.Y,
		Width:  float32(texture.Width),
		Height: float32(texture.Height),
	}
	origin := rl.Vector2{X: -float32(b.texture.Height), Y: float32(texture.Height) / 2}

	rl.DrawTexturePro(texture, source, dest, origin, b.arrow.angle*rl.Rad2deg, rl.White)
}

func (b *ball) drawPowerBar() {
	if b.state != charging {
		return
	}

	screenW := float32(rl.GetScreenWidth())
	barWidth := float32(200)
	barHeight := float32(20)
	x := (screenW - barWidth) / 2
	y := float32(30)

	rl.DrawRectangle(int32(x), int32(y), int32(barWidth), int32(barHeight), rl.LightGray)
	rl.DrawRectangle(int32(x), int32(y), int32((b.power/b.maxPower)*barWidth), int32(barHeight), rl.Red)
}
