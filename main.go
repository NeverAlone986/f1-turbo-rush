package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"image/color"
	"math/rand"
	"sync"
	"time"
)

const (
	windowWidth  = 400
	windowHeight = 600
	carWidth     = 50
	carHeight    = 80
)

var (
	playerX = float32(windowWidth/2 - carWidth/2)
	playerY = float32(windowHeight - carHeight - 20)
	enemyCars   []*canvas.Rectangle
	enemySpeed = float32(5)
	gameMutex  sync.Mutex
)

func main() {
	app := app.New()
	w := app.NewWindow("F1 Turbo Rush")
	w.Resize(fyne.NewSize(windowWidth, windowHeight))

	// Фон трассы
	bg := canvas.NewRectangle(color.RGBA{30, 30, 30, 255})

	// Машина игрока
	playerCar := canvas.NewRectangle(color.RGBA{255, 0, 0, 255})
	playerCar.Resize(fyne.NewSize(carWidth, carHeight))
	playerCar.Move(fyne.NewPos(playerX, playerY))

	// Создание машин-соперников
	for i := 0; i < 3; i++ {
		enemy := canvas.NewRectangle(color.RGBA{0, 0, 255, 255})
		enemy.Resize(fyne.NewSize(carWidth, carHeight))
		enemy.Move(fyne.NewPos(float32(rand.Intn(windowWidth-carWidth)), float32(rand.Intn(300))))
		enemyCars = append(enemyCars, enemy)
	}

	// Основной игровой контейнер
	gameContainer := container.NewWithoutLayout(bg, playerCar)
	for _, enemy := range enemyCars {
		gameContainer.Add(enemy)
	}

	w.SetContent(gameContainer)

	// Горутина для обновления игры
	go func() {
		for {
			gameMutex.Lock()
			for _, enemy := range enemyCars {
				enemy.Move(fyne.NewPos(enemy.Position().X, enemy.Position().Y+enemySpeed))
				if enemy.Position().Y > windowHeight {
					enemy.Move(fyne.NewPos(float32(rand.Intn(windowWidth-carWidth)), -carHeight))
				}
			}
			gameMutex.Unlock()
			canvas.Refresh(gameContainer)
			time.Sleep(50 * time.Millisecond) // FPS-логика
		}
	}()

	// Управление машиной
	w.Canvas().SetOnTypedKey(func(e *fyne.KeyEvent) {
		gameMutex.Lock()
		if e.Name == fyne.KeyLeft && playerX > 0 {
			playerX -= 20
		} else if e.Name == fyne.KeyRight && playerX < windowWidth-carWidth {
			playerX += 20
		}
		playerCar.Move(fyne.NewPos(playerX, playerY))
		canvas.Refresh(playerCar)
		gameMutex.Unlock()
	})

	w.ShowAndRun()
}
