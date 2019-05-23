package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	ScreenWidth  int32 = 640
	ScreenHeight int32 = 480
	CellSize     int32 = 15
)

type Snake struct {
	dX, dY    int32 //direction X and Y
	width     int32
	height    int32
	positions []sdl.Point
}

type Food struct {
	width    int32
	height   int32
	position sdl.Point
}

var random = rand.New(rand.NewSource(time.Now().UnixNano()))
var font *ttf.Font
var scoreSurface *sdl.Surface
var scoreTexture *sdl.Texture

func main() {
	var winTitle = "Go Snake 2D"
	var window *sdl.Window
	var renderer *sdl.Renderer
	var event sdl.Event
	var snake Snake
	var food Food
	var scoreText = "SCORE: %d"
	var score int32
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

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println("initializing renderer background:", err)
		return
	}
	//sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")
	defer renderer.Destroy()

	// Using the SDL_ttf library so need to initialize it before using it
	if err = ttf.Init(); err != nil {
		fmt.Printf("Failed to initialize TTF: %s\n", err)
	}

	if font, err = ttf.OpenFont("./fonts/Roboto-Regular.ttf", 14); err != nil {
		fmt.Printf("Failed to open font: %s\n", err)
	}

	if scoreSurface, err = font.RenderUTF8Solid(fmt.Sprintf(scoreText, score), sdl.Color{R: 255, G: 255, B: 255, A: 255}); err != nil {
		fmt.Printf("Failed to render text: %s\n", err)
	}

	if scoreTexture, err = renderer.CreateTextureFromSurface(scoreSurface); err != nil {
		fmt.Printf("Failed to create texture: %s\n", err)
	}

	defer scoreTexture.Destroy()
	defer font.Close()
	scoreSurface.Free()

	running = true

	snake = Snake{
		dX:     1,
		dY:     0,
		width:  CellSize,
		height: CellSize,
		positions: []sdl.Point{
			sdl.Point{X: ScreenWidth / 2, Y: ScreenHeight / 2},
			sdl.Point{X: ScreenWidth/2 + CellSize + 1, Y: ScreenHeight / 2},
			sdl.Point{X: ScreenWidth/2 + CellSize*2 + 1, Y: ScreenHeight / 2},
		},
	}

	food = Food{
		width:    CellSize,
		height:   CellSize,
		position: randomPosition(),
	}

	fmt.Printf("food.position= %d, %d\n", food.position.X, food.position.Y)
	fmt.Printf("snake.position= %d, %d\n", snake.positions[0].X, snake.positions[0].Y)

	ticker := time.NewTicker(time.Second / 10)

	for running {

		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				if e.Type == sdl.KEYDOWN {
					changeDirection(e, &snake)

					if key := e.Keysym.Sym; key == sdl.K_SPACE {
						food.position = sdl.Point{X: -1, Y: -100}
					}
				}
			}
		}

		moveSnake(&snake)
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
		drawScore(renderer, scoreTexture, fmt.Sprintf(scoreText, score))

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
	if food.position.X == -1 {
		food.position = randomPosition()
		fmt.Println(food.position)
	}

	draw(renderer, &sdl.Rect{
		X: food.position.X, Y: food.position.Y, W: food.width, H: food.height})
}

func drawScore(renderer *sdl.Renderer, scoreTexture *sdl.Texture, scoreText string) {
	//TODO refactory
	var err error
	if scoreSurface, err = font.RenderUTF8Solid(scoreText, sdl.Color{R: 255, G: 255, B: 255, A: 255}); err != nil {
		fmt.Printf("Failed to render text: %s\n", err)
	}

	if scoreTexture, err = renderer.CreateTextureFromSurface(scoreSurface); err != nil {
		fmt.Printf("Failed to create texture: %s\n", err)
	}
	renderer.Copy(scoreTexture, nil, &sdl.Rect{W: (ScreenWidth / 7), H: 35, X: 5, Y: 0})
}

func randomPosition() sdl.Point {
	newX := math.Floor(random.Float64()*(float64(ScreenWidth/CellSize))) * float64(CellSize)
	newY := math.Floor(random.Float64()*(float64(ScreenHeight/CellSize))) * float64(CellSize)

	newPosition := sdl.Point{
		X: int32(newX),
		Y: int32(newY),
	}

	return newPosition
}

func moveSnake(snake *Snake) {
	for i := 0; i < len(snake.positions)-1; i++ {
		snake.positions[i] = snake.positions[i+1]
	}

	head := snake.positions[len(snake.positions)-1]
	snake.positions[len(snake.positions)-1].X = head.X + snake.dX*CellSize
	snake.positions[len(snake.positions)-1].Y = head.Y + snake.dY*CellSize

	if snakeLength := len(snake.positions) - 1; head.X > ScreenWidth {
		snake.positions[snakeLength].X = 0
	} else if head.X < 0 {
		snake.positions[snakeLength].X = ScreenWidth
	} else if head.Y > ScreenHeight {
		snake.positions[snakeLength].Y = 0
	} else if head.Y < 0 {
		snake.positions[snakeLength].Y = ScreenHeight
	}
}

func (s *Snake) addTail() {
	newPart := sdl.Point{X: 1, Y: 1}
	s.positions = append([]sdl.Point{newPart}, s.positions...)
}

func changeDirection(e *sdl.KeyboardEvent, snake *Snake) {
	if key := e.Keysym.Sym; key == sdl.K_RIGHT && !(snake.dX == -1) {
		snake.dX = 1
		snake.dY = 0
	} else if key == sdl.K_LEFT && !(snake.dX == 1) {
		snake.dX = -1
		snake.dY = 0
	} else if key == sdl.K_UP && !(snake.dY == 1) {
		snake.dX = 0
		snake.dY = -1
	} else if key == sdl.K_DOWN && !(snake.dY == -1) {
		snake.dX = 0
		snake.dY = 1
	}
}

func (s *Snake) checkCollision(area sdl.Rect) bool {
	if head := s.positions[len(s.positions)-1]; head.InRect(&area) {
		return true
	}

	return false
}
