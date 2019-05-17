package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	screenWidth  = 640
	screenHeight = 480
	fps          = 30
)

func main() {
	var winTitle = "Snake 2D"
	var window *sdl.Window
	var context sdl.GLContext
	var renderer *sdl.Renderer
	var event sdl.Event
	var running bool
	var err error

	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Println("initializing SDL:", err)
		return
	}
	defer sdl.Quit()

	window, err = sdl.CreateWindow(
		winTitle,
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		screenWidth, screenHeight,
		sdl.WINDOW_OPENGL)
	if err != nil {
		fmt.Println("initializing window:", err)
		return
	}
	defer window.Destroy()

	context, err = window.GLCreateContext()
	if err != nil {
		panic(err)
	}
	defer sdl.GLDeleteContext(context)

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println("initializing renderer background:", err)
		return
	}
	defer renderer.Destroy()

	running = true

	var rectSnake = sdl.Rect{X: screenWidth / 2, Y: screenHeight / 2, W: 15, H: 15}
	///var rectFood = sdl.Rect{X: 0, Y: 0, W: 5, H: 5}

	var nextTick uint32
	for running {

		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				fmt.Println(e)
			}
		}

		drawBackground(renderer)
		renderer.Clear()

		draw(renderer, &rectSnake)
		move(renderer, &rectSnake)

		renderer.Present()

		//improve preformance
		ticks := sdl.GetTicks()

		if ticks < nextTick {
			sdl.Delay(nextTick - ticks)
		}
		nextTick = ticks + (1000 / fps)

	}
}

func update() {}

func drawBackground(renderer *sdl.Renderer) {
	renderer.SetDrawColor(0, 0, 0, 0)
}

func draw(renderer *sdl.Renderer, rect *sdl.Rect) {
	renderer.SetDrawColor(255, 255, 255, 255)
	renderer.FillRect(rect)
}

func move(renderer *sdl.Renderer, rect *sdl.Rect) {
	renderer.SetDrawColor(255, 255, 255, 255)

	rect.X++

	if rect.X > screenWidth {
		rect.X = 0
	}

	renderer.FillRect(rect)
}
