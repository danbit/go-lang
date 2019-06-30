package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Game struct{}

const (
	ScreenWidth   int32  = 688
	ScreenHeight  int32  = 496
	CellSize      int32  = 16
	BorderSize    int32  = 32
	ScoreText     string = "SCORE:"
	HighScoreText string = "HIGHSCORE:"
)

type Dimension struct {
	W int32
	H int32
}

type Food struct {
	dimension Dimension
	position  sdl.Point
}

type Color struct {
	R, G, B uint8
}

var isRunning = false
var level Level
var snake Snake
var food Food
var window *sdl.Window
var renderer *sdl.Renderer
var font *ttf.Font
var event sdl.Event
var err error

func (g Game) Start(title string, width int32, height int32, fullscreen bool) error {
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

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println("initializing renderer background:", err)
	}

	// Using the SDL_ttf library so need to initialize it before using it
	if err = ttf.Init(); err != nil {
		fmt.Printf("Failed to initialize TTF: %s\n", err)
	}

	if font, err = ttf.OpenFont("./fonts/Roboto-Regular.ttf", 18); err != nil {
		fmt.Printf("Failed to open font: %s\n", err)
	}

	isRunning = true

	return nil
}

func (g Game) Update() {
}

func (g Game) Render() {
	renderer.Clear()

	renderer.SetDrawColor(0, 0, 0, 255)

	renderer.Present()
}

func (g Game) HandleEvents() {
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

func (g Game) Destroy() {
	sdl.Quit()
	window.Destroy()
	renderer.Destroy()
	font.Close()
}

func (g Game) IsRunning() bool {
	return isRunning
}

func gameOver(gs *GameState) {
	gs.changeState(GAME_OVER)
	sdl.Delay(1000)
}

func gameWin(gs *GameState) {
	gs.changeState(GAME_WIN)
	sdl.Delay(1000)
}
