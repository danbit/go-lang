package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	ScreenWidth  int32 = 640
	ScreenHeight int32 = 480
	CellSize     int32 = 15
)

type Snake struct {
	width     int32
	height    int32
	positions []sdl.Point
}

type Food struct {
	width    int32
	height   int32
	position sdl.Point
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
	var food Food
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

	snake = Snake{
		width:  CellSize,
		height: CellSize,
		positions: []sdl.Point{
			sdl.Point{X: ScreenWidth/2 + CellSize*4 + 1, Y: ScreenHeight / 2},
			sdl.Point{X: ScreenWidth/2 + CellSize*3 + 1, Y: ScreenHeight / 2},
			sdl.Point{X: ScreenWidth/2 + CellSize*2 + 1, Y: ScreenHeight / 2},
			sdl.Point{X: ScreenWidth/2 + CellSize + 1, Y: ScreenHeight / 2},
			sdl.Point{X: ScreenWidth / 2, Y: ScreenHeight / 2}},
	}

	food = Food{
		width:    CellSize,
		height:   CellSize,
		position: randomPosition(),
	}

	fmt.Printf("food.position= %d, %d\n", food.position.X, food.position.Y)

	dX = 1
	dY = 0

	ticker := time.NewTicker(time.Second / 20)

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

		moveSnake(&snake)

		drawFood(renderer, &food)
		drawSnake(renderer, &snake)

		renderer.Present()

		<-ticker.C
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
	for _, element := range snake.positions {
		draw(renderer, &sdl.Rect{
			X: element.X, Y: element.Y, W: snake.width, H: snake.height})
	}
}

func drawFood(renderer *sdl.Renderer, food *Food) {
	draw(renderer, &sdl.Rect{
		X: food.position.X, Y: food.position.Y, W: food.width, H: food.height})
}

func randomPosition() sdl.Point {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	newX := math.Floor(r1.Float64()*(float64(ScreenWidth/CellSize))) * float64(CellSize)
	newY := math.Floor(r1.Float64()*(float64(ScreenWidth/CellSize))) * float64(CellSize)

	newPosition := sdl.Point{
		X: int32(newX),
		Y: int32(newY),
	}

	return newPosition
}

func moveSnake(snake *Snake) {

	//TODO check collision

	for i := 0; i < len(snake.positions)-1; i++ {
		snake.positions[i] = snake.positions[i+1]
	}

	head := snake.positions[len(snake.positions)-1]
	snake.positions[len(snake.positions)-1].X = head.X + dX*CellSize
	snake.positions[len(snake.positions)-1].Y = head.Y + dY*CellSize

	for i := range snake.positions {
		if headPos := snake.positions[i]; headPos.X > ScreenWidth {
			snake.positions[i].X = 0
		} else if headPos.X < 0 {
			snake.positions[i].X = ScreenWidth
		} else if headPos.Y > ScreenHeight {
			snake.positions[i].Y = 0
		} else if headPos.Y < 0 {
			snake.positions[i].Y = ScreenHeight
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
