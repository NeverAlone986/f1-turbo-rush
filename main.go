package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Game struct {
	Car           *canvas.Rectangle
	Obstacle      *canvas.Rectangle
	ScoreLabel    *widget.Label
	GameOverLabel *widget.Label
	RestartButton *widget.Button
	Score         int
	Running       bool
}

func main() {
	a := app.New()
	w := a.NewWindow("F1 Turbo Rush")
	w.Resize(fyne.NewSize(400, 600))

	game := &Game{
		Car:        canvas.NewRectangle(color.RGBA{0, 255, 0, 255}),
		Obstacle:   canvas.NewRectangle(color.RGBA{255, 0, 0, 255}),
		ScoreLabel: widget.NewLabel("Score: 0"),
		Running:    true,
	}

	game.Car.Resize(fyne.NewSize(50, 100))
	game.Car.Move(fyne.NewPos(175, 450))

	game.Obstacle.Resize(fyne.NewSize(50, 100))
	game.Obstacle.Move(fyne.NewPos(float32(rand.Intn(350)), -100))

	game.GameOverLabel = widget.NewLabel("Game Over! Your Score: 0")
	game.GameOverLabel.Hide()

	game.RestartButton = widget.NewButton("Restart", func() {
		game.Restart()
	})
	game.RestartButton.Hide()

	content := container.NewWithoutLayout(
		game.Car,
		game.Obstacle,
	)
	w.SetContent(container.NewVBox(
		game.ScoreLabel,
		content,
		game.GameOverLabel,
		game.RestartButton,
	))

	go game.StartGameLoop()
	w.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		game.HandleInput(k)
	})

	w.ShowAndRun()
}

func (g *Game) StartGameLoop() {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		if !g.Running {
			return
		}
		g.Update()
	}
}

func (g *Game) Update() {
	pos := g.Obstacle.Position()
	pos.Y += 10
	if pos.Y > 600 {
		pos.Y = -100
		pos.X = float32(rand.Intn(350))
		g.Score++
		g.ScoreLabel.SetText("Score: " + strconv.Itoa(g.Score))
	}
	g.Obstacle.Move(pos)
	canvas.Refresh(g.Obstacle)

	if g.CheckCollision() {
		g.GameOver()
	}
}

func (g *Game) HandleInput(k *fyne.KeyEvent) {
	pos := g.Car.Position()
	if k.Name == fyne.KeyLeft && pos.X > 0 {
		pos.X -= 20
	} else if k.Name == fyne.KeyRight && pos.X < 350 {
		pos.X += 20
	}
	g.Car.Move(pos)
	canvas.Refresh(g.Car)
}

func (g *Game) CheckCollision() bool {
	carPos := g.Car.Position()
	obstaclePos := g.Obstacle.Position()

	if carPos.X < obstaclePos.X+50 && carPos.X+50 > obstaclePos.X &&
		carPos.Y < obstaclePos.Y+100 && carPos.Y+100 > obstaclePos.Y {
		return true
	}
	return false
}

func (g *Game) GameOver() {
	g.Running = false
	g.GameOverLabel.SetText(fmt.Sprintf("Game Over! Your Score: %d", g.Score))
	g.GameOverLabel.Show()
	g.RestartButton.Show()
}

func (g *Game) Restart() {
	g.Score = 0
	g.ScoreLabel.SetText("Score: 0")
	g.Car.Move(fyne.NewPos(175, 450))
	g.Obstacle.Move(fyne.NewPos(float32(rand.Intn(350)), -100))
	g.GameOverLabel.Hide()
	g.RestartButton.Hide()
	g.Running = true
	go g.StartGameLoop()
}
