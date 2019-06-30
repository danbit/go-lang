package main

import (
	"os"
	"time"
)

var game Game

func main() {
	game = Game{}
	game.Start("Go Snake 2D", ScreenWidth+BorderSize, ScreenHeight+BorderSize, false)

	ticker := time.NewTicker(time.Second / 10)

	for game.IsRunning() {

		game.HandleEvents()
		game.Update()
		game.Render()

		<-ticker.C
	}

	game.Destroy()
	os.Exit(0)
}
