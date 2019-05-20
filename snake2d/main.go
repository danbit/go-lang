package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	ScreenWidth  int32 = 640
	ScreenHeight int32 = 480
	CellSize     int32 = 15
	DELAY              = 300
)

type Vector2D struct {
	X, Y int32
}

type Snake struct {
	width  int32
	height int32
	parts  []Vector2D
}

// horizontalVelocity, verticalVelocity
var dX, dY int32

func main() {
	var winTitle = "Go Snake 2D"
	var window *sdl.Window
	var context sdl.GLContext
	var renderer *sdl.Renderer
	var event sdl.Event
	var snake Snake
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
		ScreenWidth, ScreenHeight,
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

	var rectFood = sdl.Rect{X: 0, Y: 0, W: 5, H: 5}

	snake = Snake{
		width:  15,
		height: 15,
		parts: []Vector2D{
			Vector2D{X: ScreenWidth/2 + CellSize*2 + 1, Y: ScreenHeight / 2},
			Vector2D{X: ScreenWidth/2 + CellSize + 1, Y: ScreenHeight / 2},
			Vector2D{X: ScreenWidth / 2, Y: ScreenHeight / 2}},
	}

	dX = 1
	dY = 0

	var nextTick uint32
	for running {

		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				if e.Type == sdl.KEYDOWN {
					changeDirection(e)
				}
			}
		}

		drawBackground(renderer)
		renderer.Clear()

		draw(renderer, &rectFood)

		moveSnake(&snake)
		drawSnake(renderer, &snake)

		renderer.Present()

		//improve preformance
		ticks := sdl.GetTicks()
		if ticks < nextTick {
			sdl.Delay(nextTick - ticks)
		}
		nextTick = ticks + (DELAY)
	}
}

func drawBackground(renderer *sdl.Renderer) {
	renderer.SetDrawColor(0, 0, 0, 0)
}

func draw(renderer *sdl.Renderer, rect *sdl.Rect) {
	renderer.SetDrawColor(255, 255, 255, 255)
	renderer.FillRect(rect)
}

func drawSnake(renderer *sdl.Renderer, snake *Snake) {
	for _, element := range snake.parts {
		draw(renderer, &sdl.Rect{
			X: element.X, Y: element.Y, W: snake.width, H: snake.height})
	}
}

func moveSnake(snake *Snake) {

	//TODO check collision

	for i := 0; i < len(snake.parts)-1; i++ {
		snake.parts[i] = snake.parts[i+1]
	}

	tail := snake.parts[len(snake.parts)-1]
	snake.parts[len(snake.parts)-1].X = tail.X + dX*CellSize
	snake.parts[len(snake.parts)-1].Y = tail.Y + dY*CellSize

	for i := range snake.parts {
		if headPos := snake.parts[i]; headPos.X > ScreenWidth {
			snake.parts[i].X = 0
		} else if headPos.X <= 0 {
			snake.parts[i].X = ScreenWidth
		} else if headPos.Y > ScreenHeight {
			snake.parts[i].Y = 0
		} else if headPos.Y <= 0 {
			snake.parts[i].Y = ScreenHeight
		}
	}
}

func changeDirection(e *sdl.KeyboardEvent) {

	if key := e.Keysym.Sym; key == sdl.K_RIGHT && dX != -1 {
		dX = 1
		dY = 0
	} else if key == sdl.K_LEFT && dX != 1 {
		dX = -1
		dY = 0
	} else if key == sdl.K_UP && dY != 1 {
		dX = 0
		dY = -1
	} else if key == sdl.K_DOWN && dY != -1 {
		dX = 0
		dY = 1
	}
}
