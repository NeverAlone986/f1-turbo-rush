package main

import (
	"math/rand"
	"os"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

const (
	windowWidth  = 400
	windowHeight = 600
	carSpeed     = 20 // Скорость машины игрока
	enemySpeed   = 6  // Скорость врагов
	updateDelay  = 16 * time.Millisecond // 60 FPS
	recordFile   = "record.txt"
)

var score int
var record int
var playerCar *canvas.Image
var enemyCars []*canvas.Image
var obstacles []*canvas.Image
var scoreLabel *widget.Label
var recordLabel *widget.Label
var gameContent *fyne.Container
var playerX float32 = 175 // Начальная позиция игрока
var playerY float32 = 500 // Начальная вертикальная позиция игрока

func loadRecord() {
	if data, err := os.ReadFile(recordFile); err == nil {
		if r, err := strconv.Atoi(string(data)); err == nil {
			record = r
		}
	}
}

func saveRecord() {
	if score > record {
		record = score
		os.WriteFile(recordFile, []byte(strconv.Itoa(record)), 0644)
	}
}

func setupGame(window fyne.Window) {
	score = 0
	playerX = 175
	playerY = 500

	background := canvas.NewImageFromFile("assets/track.png")
	background.Resize(fyne.NewSize(windowWidth, windowHeight))

	playerCar = canvas.NewImageFromFile("assets/player_car.png")
	playerCar.Resize(fyne.NewSize(50, 100))
	playerCar.Move(fyne.NewPos(playerX, playerY))

	enemyCars = []*canvas.Image{}
	for i := 0; i < 2; i++ {
		enemyCar := canvas.NewImageFromFile("assets/enemy_car" + strconv.Itoa(i+1) + ".png")
		enemyCar.Resize(fyne.NewSize(50, 100))
		enemyCar.Move(fyne.NewPos(float32(rand.Intn(350)), float32(rand.Intn(300))))
		enemyCars = append(enemyCars, enemyCar)
	}

	obstacles = []*canvas.Image{}
	for i := 0; i < 3; i++ {
		obstacle := canvas.NewImageFromFile("assets/obstacle" + strconv.Itoa(i+1) + ".png")
		obstacle.Resize(fyne.NewSize(50, 50))
		obstacle.Move(fyne.NewPos(float32(rand.Intn(350)), float32(rand.Intn(300))))
		obstacles = append(obstacles, obstacle)
	}

	scoreLabel = widget.NewLabelWithStyle("Score: 0", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	scoreLabel.Move(fyne.NewPos(10, 10))
	recordLabel = widget.NewLabelWithStyle("Record: "+strconv.Itoa(record), fyne.TextAlignTrailing, fyne.TextStyle{Bold: true})
	recordLabel.Move(fyne.NewPos(105, 30))

	gameContent = container.NewWithoutLayout(background, playerCar, scoreLabel, recordLabel)
	for _, car := range enemyCars {
		gameContent.Add(car)
	}
	for _, obstacle := range obstacles {
		gameContent.Add(obstacle)
	}

	window.SetContent(container.NewStack(background, gameContent))
	addKeyboardControl(window)
	go gameLoop(window)
}

func addKeyboardControl(window fyne.Window) {
	if deskCanvas, ok := window.Canvas().(desktop.Canvas); ok {
		deskCanvas.SetOnKeyDown(func(e *fyne.KeyEvent) {
			switch e.Name {
			case fyne.KeyLeft, fyne.KeyA:
				if playerX > 10 {
					playerX -= carSpeed
				}
			case fyne.KeyRight, fyne.KeyD:
				if playerX < windowWidth-60 {
					playerX += carSpeed
				}
			case fyne.KeyUp, fyne.KeyW:
				if playerY > 300 {
					playerY -= carSpeed
				}
			case fyne.KeyDown, fyne.KeyS:
				if playerY < windowHeight-120 {
					playerY += carSpeed
				}
			}
			playerCar.Move(fyne.NewPos(playerX, playerY))
			gameContent.Refresh()
		})
	}
}

func gameLoop(window fyne.Window) {
	ticker := time.NewTicker(updateDelay)
	defer ticker.Stop()

	for range ticker.C {
		score++
		scoreLabel.SetText("Score: " + strconv.Itoa(score))

		// Двигаем врагов
		for _, car := range enemyCars {
			pos := car.Position()
			car.Move(fyne.NewPos(pos.X, pos.Y+enemySpeed))
			if pos.Y > windowHeight {
				car.Move(fyne.NewPos(float32(rand.Intn(350)), -50))
			}
			if checkCollision(playerCar, car) {
				gameOver(window)
				return
			}
		}

		// Двигаем препятствия
		for _, obstacle := range obstacles {
			pos := obstacle.Position()
			obstacle.Move(fyne.NewPos(pos.X, pos.Y+enemySpeed))
			if pos.Y > windowHeight {
				obstacle.Move(fyne.NewPos(float32(rand.Intn(350)), -50))
			}
			if checkCollision(playerCar, obstacle) {
				gameOver(window)
				return
			}
		}

		gameContent.Refresh()
	}
}

func checkCollision(a, b *canvas.Image) bool {
	posA, posB := a.Position(), b.Position()
	return posA.X < posB.X+50 && posA.X+50 > posB.X && posA.Y < posB.Y+50 && posA.Y+100 > posB.Y
}

func gameOver(window fyne.Window) {
	saveRecord()
	msg := "Game Over!\nYour Score: " + strconv.Itoa(score) + "\nRecord: " + strconv.Itoa(record)
	dialog := widget.NewLabel(msg)
	restartBtn := widget.NewButton("Restart", func() {
		setupGame(window)
	})

	box := container.NewVBox(dialog, restartBtn)
	window.SetContent(box)
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("F1 Turbo Rush")
	myWindow.Resize(fyne.NewSize(windowWidth, windowHeight))

	loadRecord()
	setupGame(myWindow)

	myWindow.ShowAndRun()
}
