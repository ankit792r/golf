package game

import (
	"embed"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	ScreenWidth  = 720
	ScreenHeight = 1024
)

type scene int

const (
	sceneMenu scene = iota
	scenePlay
)

type Game struct {
	scene scene
	quit  bool

	background  rl.Texture2D
	ballTex     rl.Texture2D
	holeTex     rl.Texture2D
	arrowTex    rl.Texture2D
	obstacleTex [kindCount]rl.Texture2D

	ball      *ball
	hole      *hole
	obstacles []*obstacle

	level   int
	score   int
	strokes int

	playBtn button
	quitBtn button
}

func New(assets embed.FS) *Game {
	g := &Game{
		scene:      sceneMenu,
		background: loadPNG(assets, "assets/ground.png"),
		ballTex:    loadPNG(assets, "assets/ball.png"),
		holeTex:    loadPNG(assets, "assets/hole.png"),
		arrowTex:   loadPNG(assets, "assets/arrow.png"),
	}
	g.obstacleTex = [kindCount]rl.Texture2D{
		loadPNG(assets, "assets/obstical-small.png"),
		loadPNG(assets, "assets/obstical-medium.png"),
		loadPNG(assets, "assets/obstical-large.png"),
		loadPNG(assets, "assets/obstical-xlarge.png"),
	}

	g.ball = newBall(g.ballTex, g.arrowTex)
	g.hole = newHole(g.holeTex)

	cx := float32(ScreenWidth) / 2
	cy := float32(ScreenHeight) / 2
	g.playBtn = newButton(cx, cy, 240, 64, "Play")
	g.quitBtn = newButton(cx, cy+90, 240, 64, "Quit")

	return g
}

func Run(assets embed.FS) {
	rl.InitWindow(ScreenWidth, ScreenHeight, "Go Golf")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	g := New(assets)
	defer g.Close()

	for !rl.WindowShouldClose() && !g.quit {
		g.Update(rl.GetFrameTime())
		g.Draw()
	}
}

func (g *Game) Close() {
	rl.UnloadTexture(g.background)
	rl.UnloadTexture(g.ballTex)
	rl.UnloadTexture(g.holeTex)
	rl.UnloadTexture(g.arrowTex)
	for i := range g.obstacleTex {
		rl.UnloadTexture(g.obstacleTex[i])
	}
}

func (g *Game) Update(dt float32) {
	switch g.scene {
	case sceneMenu:
		g.updateMenu()
	case scenePlay:
		g.updatePlay(dt)
	}
}

func (g *Game) updateMenu() {
	if g.playBtn.released() {
		g.startGame()
	}
	if g.quitBtn.released() {
		g.quit = true
	}
}

func (g *Game) updatePlay(dt float32) {
	if rl.IsKeyPressed(rl.KeyEscape) {
		g.scene = sceneMenu
		return
	}

	if g.ball.state != sinking && g.hole.contains(g.ball) {
		g.ball.startSink(g.hole.position)
	}

	wasMoving := g.ball.state == moving
	g.ball.update(dt, g.obstacles)
	if !wasMoving && g.ball.state == moving {
		g.strokes++
	}

	if g.ball.finishedSink() {
		g.completeLevel()
	}
}

func (g *Game) startGame() {
	g.score = 0
	g.level = 1
	g.scene = scenePlay
	g.startLevel()
}

func (g *Game) startLevel() {
	g.strokes = 0
	g.ball.resetToStart()
	g.obstacles = spawnObstacles(g.level, g.obstacleTex, g.ball.startPos)
	g.hole.randomizePosition(g.ball.startPos, g.obstacles)
}

func (g *Game) completeLevel() {
	g.score += scoreForHole(g.level, g.strokes)
	g.level++
	g.startLevel()
}

func (g *Game) Draw() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.DarkGray)
	rl.DrawTexture(g.background, 0, 0, rl.White)

	switch g.scene {
	case sceneMenu:
		g.drawMenu()
	case scenePlay:
		g.drawPlay()
	}

	rl.EndDrawing()
}

func (g *Game) drawMenu() {
	rl.DrawRectangle(0, 0, ScreenWidth, ScreenHeight, rl.Fade(rl.Black, 0.45))

	title := "GO GOLF"
	titleSize := int32(64)
	tw := rl.MeasureText(title, titleSize)
	rl.DrawText(title, (ScreenWidth-tw)/2+2, 222, titleSize, rl.Black)
	rl.DrawText(title, (ScreenWidth-tw)/2, 220, titleSize, rl.RayWhite)

	g.playBtn.draw()
	g.quitBtn.draw()
}

func (g *Game) drawPlay() {
	g.hole.draw()
	for _, o := range g.obstacles {
		o.draw()
	}
	g.ball.draw()
	g.ball.drawPowerBar()
	g.drawHUD()
}

func (g *Game) drawHUD() {
	drawOutlinedText(fmt.Sprintf("Level %d", g.level), 16, 12, 28, rl.White)
	drawOutlinedText(fmt.Sprintf("Score %d", g.score), 16, 44, 28, rl.White)
}

func drawOutlinedText(text string, x, y, size int32, color rl.Color) {
	rl.DrawText(text, x+2, y+2, size, rl.Black)
	rl.DrawText(text, x, y, size, color)
}
