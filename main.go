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
	carSpeed     = 4
	enemySpeed   = 3
	updateDelay  = 16 * time.Millisecond
	recordFile   = "record.txt"
	minDistance  = 60 // Минимальное расстояние между объектами
)

var score int
var record int
var playerCar *canvas.Image
var enemyCars []*canvas.Image
var obstacles []*canvas.Image
var scoreLabel *widget.Label
var recordLabel *widget.Label
var gameContent *fyne.Container
var playerX float32 = 175
var playerY float32 = 500
var keysPressed = make(map[fyne.KeyName]bool)

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
	keysPressed = make(map[fyne.KeyName]bool)

	background := canvas.NewImageFromFile("assets/track.png")
	background.Resize(fyne.NewSize(windowWidth, windowHeight))

	playerCar = canvas.NewImageFromFile("assets/player_car.png")
	playerCar.Resize(fyne.NewSize(50, 100))
	playerCar.Move(fyne.NewPos(playerX, playerY))

	enemyCars = spawnUniqueObjects(2, "assets/enemy_car", 50, 100)
	obstacles = spawnUniqueObjects(3, "assets/obstacle", 50, 50)

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

func spawnUniqueObjects(count int, assetPrefix string, width, height float32) []*canvas.Image {
	var objects []*canvas.Image
	for i := 0; i < count; i++ {
		var obj *canvas.Image
		for {
			x := float32(rand.Intn(350))
			y := float32(rand.Intn(300))

			// Проверяем, не пересекается ли объект с уже созданными
			if !checkOverlap(x, y, width, height, objects) {
				obj = canvas.NewImageFromFile(assetPrefix + strconv.Itoa(i+1) + ".png")
				obj.Resize(fyne.NewSize(width, height))
				obj.Move(fyne.NewPos(x, y))
				objects = append(objects, obj)
				break
			}
		}
	}
	return objects
}

func checkOverlap(x, y, width, height float32, objects []*canvas.Image) bool {
	for _, obj := range objects {
		pos := obj.Position()
		if x < pos.X+width+minDistance && x+width > pos.X-minDistance &&
			y < pos.Y+height+minDistance && y+height > pos.Y-minDistance {
			return true
		}
	}
	return false
}

func addKeyboardControl(window fyne.Window) {
	if deskCanvas, ok := window.Canvas().(desktop.Canvas); ok {
		deskCanvas.SetOnKeyDown(func(e *fyne.KeyEvent) {
			keysPressed[e.Name] = true
		})
		deskCanvas.SetOnKeyUp(func(e *fyne.KeyEvent) {
			keysPressed[e.Name] = false
		})
	}
}

func updatePlayerPosition() {
	if keysPressed[fyne.KeyLeft] || keysPressed[fyne.KeyA] {
		if playerX > 10 {
			playerX -= carSpeed
		}
	}
	if keysPressed[fyne.KeyRight] || keysPressed[fyne.KeyD] {
		if playerX < windowWidth-60 {
			playerX += carSpeed
		}
	}
	if keysPressed[fyne.KeyUp] || keysPressed[fyne.KeyW] {
		if playerY > 10 {
			playerY -= carSpeed
		}
	}
	if keysPressed[fyne.KeyDown] || keysPressed[fyne.KeyS] {
		if playerY < windowHeight-110 {
			playerY += carSpeed
		}
	}
	playerCar.Move(fyne.NewPos(playerX, playerY))
}

func gameLoop(window fyne.Window) {
	ticker := time.NewTicker(updateDelay)
	defer ticker.Stop()

	for range ticker.C {
		updatePlayerPosition()

		score++
		scoreLabel.SetText("Score: " + strconv.Itoa(score))

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
