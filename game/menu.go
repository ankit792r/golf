package game

import rl "github.com/gen2brain/raylib-go/raylib"

type button struct {
	bounds rl.Rectangle
	label  string
}

func newButton(cx, cy, w, h float32, label string) button {
	return button{
		bounds: rl.Rectangle{X: cx - w/2, Y: cy - h/2, Width: w, Height: h},
		label:  label,
	}
}

func (b *button) hovered() bool {
	return rl.CheckCollisionPointRec(rl.GetMousePosition(), b.bounds)
}

func (b *button) released() bool {
	return b.hovered() && rl.IsMouseButtonReleased(rl.MouseButtonLeft)
}

func (b *button) draw() {
	bg := rl.DarkGreen
	fg := rl.White
	if b.hovered() {
		bg = rl.Lime
		fg = rl.DarkGreen
	}

	rl.DrawRectangleRec(b.bounds, bg)
	rl.DrawRectangleLinesEx(b.bounds, 3, rl.RayWhite)

	size := int32(32)
	tw := rl.MeasureText(b.label, size)
	x := int32(b.bounds.X + (b.bounds.Width-float32(tw))/2)
	y := int32(b.bounds.Y + (b.bounds.Height-float32(size))/2)
	rl.DrawText(b.label, x, y, size, fg)
}
