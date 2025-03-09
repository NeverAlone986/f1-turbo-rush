package main

import (
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

const (
	windowWidth  = 400
	windowHeight = 600
	carSpeed     = 6
	enemySpeed   = 5
	updateDelay  = 15 * time.Millisecond // Ускорение темпа игры
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("F1 Turbo Rush")
	myWindow.Resize(fyne.NewSize(windowWidth, windowHeight))

	background := canvas.NewImageFromFile("assets/track.png") // Фон трассы F1
	background.Resize(fyne.NewSize(windowWidth, windowHeight))

	playerCar := canvas.NewImageFromFile("assets/player_car.png")
	playerCar.Move(fyne.NewPos(170, 500))

	enemyCars := []*canvas.Image{}
	for i := 0; i < 3; i++ {
		enemyCar := canvas.NewImageFromFile("assets/enemy_car.png")
		enemyCar.Move(fyne.NewPos(float32(rand.Intn(350)), float32(rand.Intn(300))))
		enemyCars = append(enemyCars, enemyCar)
	}

	obstacles := []*canvas.Image{}
	for i := 0; i < 2; i++ {
		obstacle := canvas.NewImageFromFile("assets/obstacle.png") // Разные препятствия
		obstacle.Move(fyne.NewPos(float32(rand.Intn(350)), float32(rand.Intn(300))))
		obstacles = append(obstacles, obstacle)
	}

	content := container.NewWithoutLayout(background, playerCar)
	for _, car := range enemyCars {
		content.Add(car)
	}
	for _, obstacle := range obstacles {
		content.Add(obstacle)
	}

	myWindow.SetContent(content)

	go func() {
		for {
			time.Sleep(updateDelay)
			for _, car := range enemyCars {
				pos := car.Position()
				car.Move(fyne.NewPos(pos.X, pos.Y+enemySpeed))
				if pos.Y > windowHeight {
					car.Move(fyne.NewPos(float32(rand.Intn(350)), 0))
				}
			}
			content.Refresh()
		}
	}()

	myWindow.ShowAndRun()
}
