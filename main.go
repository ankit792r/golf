package main

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	SCREEN_H = 1024
	SCREEN_W = 720
)

func main() {
	rl.InitWindow(SCREEN_W, SCREEN_H, "Go Golf")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	ball := NewBall()

	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()
		ball.UpdateBall(dt)

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)
		ball.DrawBall()
		ball.DrawPowerBar()

		rl.EndDrawing()
	}
}
