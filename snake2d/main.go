package main

import (
	"os"
	"time"
)

var game Game

func main() {
	game = Game{}
	defer game.Destroy()
	if err := game.Start("Go Snake 2D", ScreenWidth+BorderSize, ScreenHeight+BorderSize, false); err != nil {
		os.Exit(-1)
	}

	ticker := time.NewTicker(time.Second / 10)

	for game.IsRunning() {
		game.HandleEvents()
		game.Update()
		game.Render()

		<-ticker.C
	}

	os.Exit(0)
}
