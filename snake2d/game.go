package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var isRunning = false
var window *sdl.Window
var renderer *sdl.Renderer
var font *ttf.Font
var event sdl.Event
var err error

func Start(title string, width int32, height int32, fullscreen bool) error {
	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Println("initializing SDL:", err)
		return err
	}
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	var flags uint32
	if fullscreen {
		flags = sdl.WINDOW_FULLSCREEN_DESKTOP
	}

	window, err = sdl.CreateWindow(
		title,
		sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		width, height,
		sdl.WINDOW_OPENGL|flags)
	if err != nil {
		fmt.Println("initializing window:", err)
		return err
	}

	// Using the SDL_ttf library so need to initialize it before using it
	if err = ttf.Init(); err != nil {
		fmt.Printf("Failed to initialize TTF: %s\n", err)
	}

	if font, err = ttf.OpenFont("./fonts/Roboto-Regular.ttf", 18); err != nil {
		fmt.Printf("Failed to open font: %s\n", err)
	}

	return nil
}

func Update() {

}

func Render() {
	renderer.Clear()

	renderer.Present()
}

func HandleEvents() {
	for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch e := event.(type) {
		case *sdl.QuitEvent:
			isRunning = false
		case *sdl.KeyboardEvent:
			if e.Type == sdl.KEYDOWN {
				if key := e.Keysym.Sym; key == sdl.K_RETURN {
					fmt.Println("gameState.changeState(PLAYING)")
					fmt.Println("gameState.changeState(ON_MENU)")
				} else if key == sdl.K_SPACE {
					fmt.Println("gameState.changeState(PAUSED)")
				} else if key == sdl.K_g {
					//showGrid = !showGrid
				}

				// if gameState.FSM.Current() == PLAYING.value() {
				// 	snake.changeDirection(e)
				// }
			}
		}
	}

}

func Destroy() {
	sdl.Quit()
	window.Destroy()
	renderer.Destroy()
	font.Close()
}
