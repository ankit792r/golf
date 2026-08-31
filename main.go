package main

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	SCREEN_H = 1024
	SCREEN_W = 720
)

func main() {
	rl.InitWindow(SCREEN_W, SCREEN_H, "Go Golf")
	defer rl.CloseWindow()

	background := rl.LoadTexture("assets/ground.png")
	defer rl.UnloadTexture(background)

	rl.SetTargetFPS(60)

	ball := NewBall()
	hole := NewHole()

	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()
		ball.UpdateBall(dt)
		hole.UpdateHole(dt)

		rl.BeginDrawing()
		rl.ClearBackground(rl.DarkGray)
		rl.DrawTexture(background, 0, 0, rl.White)

		ball.DrawBall()
		ball.DrawPowerBar()
		hole.DrawHole()

		rl.EndDrawing()
	}
}
