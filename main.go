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
	carSpeed     = 5 // Скорость машины игрока
	enemySpeed   = 2 // Скорость врагов
	updateDelay  = 10 * time.Millisecond
	recordFile   = "record.txt"
	objectSize   = 50
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

var keysPressed = map[fyne.KeyName]bool{}

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

// Проверка на наложение объектов
func isOverlapping(x, y float32, objects []*canvas.Image) bool {
	for _, obj := range objects {
		pos := obj.Position()
		if x < pos.X+objectSize && x+objectSize > pos.X && y < pos.Y+objectSize && y+objectSize > pos.Y {
			return true
		}
	}
	return false
}

// Генерация позиций без наложений
func getValidPosition(objects []*canvas.Image) (float32, float32) {
	var x, y float32
	for {
		x = float32(rand.Intn(windowWidth - objectSize))
		y = float32(rand.Intn(windowHeight / 2))
		if !isOverlapping(x, y, objects) {
			break
		}
	}
	return x, y
}

func setupGame(window fyne.Window) {
	score = 0

	background := canvas.NewImageFromFile("assets/track.png")
	background.Resize(fyne.NewSize(windowWidth, windowHeight))
	background.Move(fyne.NewPos(0, 0))

	playerCar = canvas.NewImageFromFile("assets/player_car.png")
	playerCar.Resize(fyne.NewSize(objectSize, 100))
	playerCar.Move(fyne.NewPos(playerX, playerY))

	enemyCars = []*canvas.Image{}
	for i := 0; i < 2; i++ {
		x, y := getValidPosition(enemyCars)
		enemyCar := canvas.NewImageFromFile("assets/enemy_car" + strconv.Itoa(i+1) + ".png")
		enemyCar.Resize(fyne.NewSize(objectSize, 100))
		enemyCar.Move(fyne.NewPos(x, y))
		enemyCars = append(enemyCars, enemyCar)
	}

	obstacles = []*canvas.Image{}
	for i := 0; i < 3; i++ {
		x, y := getValidPosition(append(enemyCars, obstacles...))
		obstacle := canvas.NewImageFromFile("assets/obstacle" + strconv.Itoa(i+1) + ".png")
		obstacle.Resize(fyne.NewSize(objectSize, objectSize))
		obstacle.Move(fyne.NewPos(x, y))
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

// Добавляем управление с удержанием клавиш
func addKeyboardControl(window fyne.Window) {
	if deskCanvas, ok := window.Canvas().(desktop.Canvas); ok {
		deskCanvas.SetOnKeyDown(func(e *fyne.KeyEvent) {
			keysPressed[e.Name] = true
		})
		deskCanvas.SetOnKeyUp(func(e *fyne.KeyEvent) {
			delete(keysPressed, e.Name)
		})
	}

	go func() {
		for {
			time.Sleep(updateDelay)
			if keysPressed[fyne.KeyLeft] && playerX > 0 {
				playerX -= carSpeed
			}
			if keysPressed[fyne.KeyRight] && playerX < windowWidth-objectSize {
				playerX += carSpeed
			}
			if keysPressed[fyne.KeyUp] && playerY > 0 {
				playerY -= carSpeed
			}
			if keysPressed[fyne.KeyDown] && playerY < windowHeight-objectSize {
				playerY += carSpeed
			}

			playerCar.Move(fyne.NewPos(playerX, playerY))
			gameContent.Refresh()
		}
	}()
}

func gameLoop(window fyne.Window) {
	ticker := time.NewTicker(updateDelay)
	defer ticker.Stop()

	for range ticker.C {
		score++
		scoreLabel.SetText("Score: " + strconv.Itoa(score))

		for _, car := range enemyCars {
			pos := car.Position()
			car.Move(fyne.NewPos(pos.X, pos.Y+enemySpeed))
			if pos.Y > windowHeight {
				x, y := getValidPosition(enemyCars)
				car.Move(fyne.NewPos(x, y))
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
				x, y := getValidPosition(append(enemyCars, obstacles...))
				obstacle.Move(fyne.NewPos(x, y))
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
	return posA.X < posB.X+objectSize && posA.X+objectSize > posB.X &&
		posA.Y < posB.Y+objectSize && posA.Y+100 > posB.Y
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
