package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	ScreenWidth       int32  = 640
	ScreenHeight      int32  = 480
	ScreenScoreOffset int32  = 30
	CellSize          int32  = 15
	ScoreText         string = "SCORE:"
)

type Dimension struct {
	width  int32
	height int32
}

type Snake struct {
	direction sdl.Point //direction X and Y
	dimension Dimension
	positions []sdl.Point
}

type Food struct {
	dimension Dimension
	position  sdl.Point
}

func main() {
	var winTitle = "Go Snake 2D"
	var window *sdl.Window
	var renderer *sdl.Renderer
	var font *ttf.Font
	var event sdl.Event
	var rect sdl.Rect
	var snake Snake
	var food Food
	var score int32
	var running bool
	var err error

	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Println("initializing SDL:", err)
		return
	}
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")
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

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println("initializing renderer background:", err)
		return
	}
	defer renderer.Destroy()

	// Using the SDL_ttf library so need to initialize it before using it
	if err = ttf.Init(); err != nil {
		fmt.Printf("Failed to initialize TTF: %s\n", err)
	}

	if font, err = ttf.OpenFont("./fonts/Roboto-Regular.ttf", 18); err != nil {
		fmt.Printf("Failed to open font: %s\n", err)
	}
	defer font.Close()

	running = true

	snake = Snake{
		direction: sdl.Point{X: 1, Y: 0},
		dimension: Dimension{width: CellSize, height: CellSize},
		positions: []sdl.Point{
			sdl.Point{X: ScreenWidth / 2, Y: ScreenHeight / 2},
			sdl.Point{X: ScreenWidth/2 + CellSize + 1, Y: ScreenHeight / 2},
			sdl.Point{X: ScreenWidth/2 + CellSize*2 + 1, Y: ScreenHeight / 2},
		},
	}

	food = Food{
		dimension: Dimension{width: CellSize, height: CellSize},
		position:  randomPosition(),
	}

	rand.Seed(time.Now().UnixNano())
	ticker := time.NewTicker(time.Second / 10)

	for running {

		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				if e.Type == sdl.KEYDOWN {
					changeDirection(e, &snake)
				}
			}
		}

		<-ticker.C

		snake.move()
		snakeArea := sdl.Rect{X: snake.positions[0].X, Y: snake.positions[0].Y, H: CellSize, W: CellSize}
		if snake.checkCollision(snakeArea) {
			fmt.Println("Game Over")
			running = false
			return
		}

		drawBackground(renderer)
		renderer.Clear()

		drawFood(renderer, &food)

		foodArea := sdl.Rect{X: food.position.X, Y: food.position.Y, H: CellSize, W: CellSize}
		if snake.checkCollision(foodArea) {
			food.position = sdl.Point{X: -1, Y: -100}
			snake.addTail()
			score += 10
		}

		drawSnake(renderer, &snake)
		renderText(renderer, sdl.Point{X: 5, Y: 0}, &rect, ScoreText, font)
		renderText(renderer, sdl.Point{X: rect.X + rect.W + 5, Y: 0}, &rect, formatInt32(score), font)

		renderer.Present()

	}

	os.Exit(0)
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
			X: element.X, Y: element.Y, W: snake.dimension.width, H: snake.dimension.height})
	}
}

func drawFood(renderer *sdl.Renderer, food *Food) {
	if food.position.X == -1 {
		food.position = randomPosition()
	}

	draw(renderer, &sdl.Rect{
		X: food.position.X, Y: food.position.Y, W: food.dimension.width, H: food.dimension.height})
}

func renderText(renderer *sdl.Renderer, position sdl.Point, rect *sdl.Rect, text string, font *ttf.Font) {

	solidSurface, err := font.RenderUTF8Solid(text, sdl.Color{R: 255, G: 255, B: 255, A: 255})
	if err != nil {
		fmt.Printf("Failed to render text: %s\n", err)
	}

	texture, err := renderer.CreateTextureFromSurface(solidSurface)
	if err != nil {
		fmt.Printf("Failed to create texture: %s\n", err)
	}

	solidSurface.Free()

	rect.X = position.X
	rect.Y = position.Y
	rect.W = solidSurface.W
	rect.H = solidSurface.H

	renderer.Copy(texture, nil, rect)
	texture.Destroy()
}

func randomPosition() sdl.Point {
	newX := math.Floor(rand.Float64()*(float64(ScreenWidth/CellSize))) * float64(CellSize)
	newY := math.Floor(rand.Float64()*(float64(ScreenHeight/CellSize))) * float64(CellSize)

	newPosition := sdl.Point{
		X: int32(newX),
		Y: int32(newY),
	}

	return newPosition
}

// todo
func (s *Snake) move() {
	for i := 0; i < len(s.positions)-1; i++ {
		s.positions[i] = s.positions[i+1]
	}

	head := s.positions[len(s.positions)-1]
	s.positions[len(s.positions)-1].X = head.X + s.direction.X*CellSize
	s.positions[len(s.positions)-1].Y = head.Y + s.direction.Y*CellSize

	if snakeLength := len(s.positions) - 1; head.X > ScreenWidth {
		s.positions[snakeLength].X = 0
	} else if head.X < 0 {
		s.positions[snakeLength].X = ScreenWidth
	} else if head.Y > ScreenHeight {
		s.positions[snakeLength].Y = 0
	} else if head.Y < 0 {
		s.positions[snakeLength].Y = ScreenHeight
	}
}

func (s *Snake) addTail() {
	newPart := sdl.Point{X: 1, Y: 1}
	s.positions = append([]sdl.Point{newPart}, s.positions...)
}

func changeDirection(e *sdl.KeyboardEvent, snake *Snake) {
	dX := snake.direction.X
	dY := snake.direction.Y
	if key := e.Keysym.Sym; key == sdl.K_RIGHT && !(dX == -1) {
		snake.direction.X = 1
		snake.direction.Y = 0
	} else if key == sdl.K_LEFT && !(dX == 1) {
		snake.direction.X = -1
		snake.direction.Y = 0
	} else if key == sdl.K_UP && !(dY == 1) {
		snake.direction.X = 0
		snake.direction.Y = -1
	} else if key == sdl.K_DOWN && !(dY == -1) {
		snake.direction.X = 0
		snake.direction.Y = 1
	}
}

func (s *Snake) checkCollision(area sdl.Rect) bool {
	if head := s.positions[len(s.positions)-1]; head.InRect(&area) {
		return true
	}

	return false
}

func formatInt32(n int32) string {
	return strconv.FormatInt(int64(n), 10)
}
