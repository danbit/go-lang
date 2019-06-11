package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"time"

	"github.com/looplab/fsm"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type State int

const (
	ON_MENU State = iota
	PLAYING
	PAUSED
	GAME_OVER
)

const (
	ScreenWidth   int32  = 688
	ScreenHeight  int32  = 496
	CellSize      int32  = 16
	BorderSize    int32  = 32
	ScoreText     string = "SCORE:"
	HighScoreText string = "HIGHSCORE:"
)

type GameState struct {
	To  string
	FSM *fsm.FSM
}

type Dimension struct {
	W int32
	H int32
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

type Color struct {
	R, G, B uint8
}

type Level struct {
	dimension Dimension
}

var level Level
var snake Snake
var food Food
var highScore, score int32

func main() {
	var winTitle = "Go Snake 2D"
	var window *sdl.Window
	var renderer *sdl.Renderer
	var font *ttf.Font
	var event sdl.Event
	var scoreRect sdl.Rect
	var highScoreRect sdl.Rect
	var gameState *GameState
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
		ScreenWidth+BorderSize, ScreenHeight+BorderSize,
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

	gameState = newGameState("snake")
	highScore = readHighscore()
	rand.Seed(time.Now().UnixNano())

	showGrid := false
	ticker := time.NewTicker(time.Second / 10)
	running := true

	for running {

		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				if e.Type == sdl.KEYDOWN {
					if key := e.Keysym.Sym; key == sdl.K_RETURN {
						gameState.changeState(PLAYING)
						gameState.changeState(ON_MENU)
					} else if key == sdl.K_SPACE {
						gameState.changeState(PAUSED)
					} else if key == sdl.K_g {
						showGrid = !showGrid
					}

					if gameState.FSM.Current() == PLAYING.value() {
						snake.changeDirection(e)
					}
				}
			}
		}

		<-ticker.C

		drawBackground(renderer)
		renderer.Clear()

		switch cs := gameState.FSM.Current(); cs {
		case PLAYING.value():
			snake.move()

			drawBorder(renderer)

			if showGrid {
				drawGrid(renderer)
			}

			drawSnake(renderer, &snake)

			foodArea := drawFood(renderer, &food)
			if snake.checkCollision(foodArea) {
				food.position = createFoodPoint(&snake)
				snake.addTail()
				score += 10
			}

			if snake.checkBorderCollision() {
				gameOver(gameState, highScore, score)
			}

			position := sdl.Point{X: 5, Y: 0}
			renderText(renderer, &position, &scoreRect, ScoreText, font)
			position = sdl.Point{X: scoreRect.X + scoreRect.W + 5, Y: 0}
			renderText(renderer, &position, &scoreRect, formatInt32(score), font)

			position = sdl.Point{X: level.dimension.W - 200, Y: 0}
			renderText(renderer, &position, &highScoreRect, HighScoreText, font)
			position = sdl.Point{X: highScoreRect.X + highScoreRect.W + 5, Y: 0}
			renderText(renderer, &position, &highScoreRect, formatInt32(highScore), font)

			if snake.isTryingToEat() {
				gameOver(gameState, highScore, score)
			}
		case ON_MENU.value():
			renderText(renderer, nil, nil, "<Press ENTER to Start>", font)
		case PAUSED.value():
			renderText(renderer, nil, nil, "Paused", font)
		case GAME_OVER.value():
			renderText(renderer, nil, nil, "Game Over", font)
		}

		renderer.Present()
	}

	os.Exit(0)
}

func draw(renderer *sdl.Renderer, rect *sdl.Rect, color Color) {
	renderer.SetDrawColor(color.R, color.G, color.B, 255)
	if rect != nil {
		renderer.FillRect(rect)
	}
}

func drawBackground(renderer *sdl.Renderer) {
	draw(renderer, nil, Color{R: 0, G: 0, B: 0})
}

func drawBorder(renderer *sdl.Renderer) {
	r := sdl.Rect{
		X: BorderSize, Y: BorderSize, W: level.dimension.W, H: level.dimension.H,
	}
	draw(renderer, nil, Color{R: 255, G: 255, B: 255})
	renderer.DrawRect(&r)
}

func drawGrid(renderer *sdl.Renderer) {
	draw(renderer, nil, Color{R: 255, G: 255, B: 255})

	var i int32
	for i = BorderSize; i < ScreenWidth; i += CellSize {
		renderer.DrawLine(i, BorderSize, i, ScreenHeight)
	}

	for i = BorderSize; i < ScreenHeight; i += CellSize {
		renderer.DrawLine(BorderSize, i, ScreenWidth, i)
	}
}

func drawSnake(renderer *sdl.Renderer, snake *Snake) {
	var sArea sdl.Rect
	c := Color{R: 255, G: 255, B: 255}

	for i, s := range snake.positions {
		sArea = sdl.Rect{
			X: s.X, Y: s.Y, W: snake.dimension.W, H: snake.dimension.H}

		if i == len(snake.positions)-1 {
			c = Color{R: 255, G: 0, B: 0}
		}

		draw(renderer, &sArea, c)
	}
}

func drawFood(renderer *sdl.Renderer, food *Food) (r sdl.Rect) {
	r = sdl.Rect{
		X: food.position.X, Y: food.position.Y, W: food.dimension.W, H: food.dimension.H}

	draw(renderer, &r,
		Color{R: 255, G: 255, B: 255})

	return
}

func renderText(renderer *sdl.Renderer, position *sdl.Point, rect *sdl.Rect, text string, font *ttf.Font) {

	solidSurface, err := font.RenderUTF8Solid(text, sdl.Color{R: 255, G: 255, B: 255, A: 255})
	if err != nil {
		fmt.Printf("Failed to render text: %s\n", err)
	}

	texture, err := renderer.CreateTextureFromSurface(solidSurface)
	if err != nil {
		fmt.Printf("Failed to create texture: %s\n", err)
	}

	solidSurface.Free()

	if rect == nil {
		rect = &sdl.Rect{
			X: (ScreenWidth - solidSurface.W) / 2,
			Y: (ScreenHeight - solidSurface.H) / 2,
			W: solidSurface.W,
			H: solidSurface.H,
		}
	} else {
		rect.X = position.X
		rect.Y = position.Y
		rect.W = solidSurface.W
		rect.H = solidSurface.H
	}

	renderer.Copy(texture, nil, rect)
	texture.Destroy()
}

func random(min, max int32) float64 {
	fmin := float64(min)
	fmax := float64(max)
	return math.Round((rand.Float64()*(fmax-fmin)+fmin)/float64(CellSize)) * float64(CellSize)
}

func randomPosition() sdl.Point {
	newX := random(BorderSize, level.dimension.W)
	newY := random(BorderSize, level.dimension.H)

	newPosition := sdl.Point{
		X: int32(newX),
		Y: int32(newY),
	}

	return newPosition
}

func createFoodPoint(s *Snake) sdl.Point {
	newPos := randomPosition()

	for _, p := range s.positions {
		if p.InRect(&sdl.Rect{X: newPos.X, Y: newPos.Y, W: s.dimension.W, H: s.dimension.H}) {
			fmt.Println("ignored food pos", newPos)
			createFoodPoint(s)
		}
	}

	fmt.Println("new food pos", newPos)

	return newPos
}

func (s *Snake) move() {
	for i := 0; i < len(s.positions)-1; i++ {
		s.positions[i] = s.positions[i+1]
	}

	tail := s.positions[len(s.positions)-1]
	s.positions[len(s.positions)-1].X = tail.X + s.direction.X*CellSize
	s.positions[len(s.positions)-1].Y = tail.Y + s.direction.Y*CellSize
}

func (s *Snake) addTail() {
	newPart := sdl.Point{X: 1, Y: 1}
	s.positions = append([]sdl.Point{newPart}, s.positions...)
}

func (s *Snake) changeDirection(e *sdl.KeyboardEvent) {
	dX := s.direction.X
	dY := s.direction.Y
	if key := e.Keysym.Sym; key == sdl.K_RIGHT && !(dX == -1) {
		s.direction.X = 1
		s.direction.Y = 0
	} else if key == sdl.K_LEFT && !(dX == 1) {
		s.direction.X = -1
		s.direction.Y = 0
	} else if key == sdl.K_UP && !(dY == 1) {
		s.direction.X = 0
		s.direction.Y = -1
	} else if key == sdl.K_DOWN && !(dY == -1) {
		s.direction.X = 0
		s.direction.Y = 1
	}
}

func (s *Snake) checkCollision(area sdl.Rect) bool {
	if head := s.positions[len(s.positions)-1]; head.InRect(&area) {
		return true
	}

	return false
}

func (s *Snake) checkBorderCollision() bool {
	head := s.positions[len(s.positions)-1]

	var collided bool
	if head.X-s.dimension.W > level.dimension.W {
		collided = true
	} else if head.X < BorderSize {
		collided = true
	} else if head.Y-s.dimension.H > level.dimension.H {
		collided = true
	} else if head.Y < BorderSize {
		collided = true
	}

	return collided
}

func (s *Snake) isTryingToEat() bool {
	for i := 0; i < len(s.positions)-1; i++ {
		position := s.positions[i]
		area := sdl.Rect{X: position.X, Y: position.Y, H: s.dimension.H, W: s.dimension.W}

		if s.checkCollision(area) {
			return true
		}
	}

	return false
}

func formatInt32(n int32) string {
	return strconv.FormatInt(int64(n), 10)
}

func newGameState(to string) *GameState {
	gameState := &GameState{
		To: to,
	}

	gameState.FSM = fsm.NewFSM(
		ON_MENU.value(),
		fsm.Events{
			{Name: PLAYING.value(), Src: []string{ON_MENU.value()}, Dst: PLAYING.value()},
			{Name: PAUSED.value(), Src: []string{PLAYING.value()}, Dst: PAUSED.value()},
			{Name: PAUSED.value(), Src: []string{PAUSED.value()}, Dst: PLAYING.value()},
			{Name: GAME_OVER.value(), Src: []string{PLAYING.value()}, Dst: GAME_OVER.value()},
			{Name: GAME_OVER.value(), Src: []string{GAME_OVER.value()}, Dst: ON_MENU.value()},
			{Name: ON_MENU.value(), Src: []string{GAME_OVER.value()}, Dst: ON_MENU.value()},
		},
		fsm.Callbacks{
			"enter_state": func(e *fsm.Event) { gameState.enterState(e) },
		},
	)

	return gameState
}

func (gs *GameState) enterState(e *fsm.Event) {
	fmt.Printf("The gameState to %s is %s\n", gs.To, e.Dst)

	if e.Src == ON_MENU.value() && e.Dst == PLAYING.value() {
		snake, food, score = createLevel()
	}
}

func (gs *GameState) changeState(state State) {
	err := gs.FSM.Event(state.value())
	if err != nil {
		fmt.Printf("Failed to change state: %s\n", err)
	}
}

func (s State) value() string {
	return [...]string{"on_menu", "playing", "paused", "game_over"}[s]
}

func createLevel() (snake Snake, food Food, score int32) {
	level = Level{
		dimension: Dimension{
			W: ScreenWidth - BorderSize,
			H: ScreenHeight - BorderSize,
		},
	}

	screenCenterW := (level.dimension.W/CellSize)/2*CellSize + CellSize
	screenCenterH := (level.dimension.H/CellSize)/2*CellSize + CellSize

	snake = Snake{
		direction: sdl.Point{X: 1, Y: 0},
		dimension: Dimension{W: CellSize, H: CellSize},
		positions: []sdl.Point{
			sdl.Point{X: screenCenterW, Y: screenCenterH},
			sdl.Point{X: screenCenterW + CellSize, Y: screenCenterH},
		},
	}

	food = Food{
		dimension: Dimension{W: CellSize, H: CellSize},
		position:  createFoodPoint(&snake),
	}

	score = 0

	return
}

func max(x, y int32) int32 {
	if x >= y {
		return x
	}
	return y
}

func createHighscore(hs int32) {
	f, err := os.Create(getScorePath())
	check(err)

	sEnc := base64.StdEncoding.EncodeToString([]byte(formatInt32(hs)))

	_, err = f.WriteString(sEnc)
	f.Sync()
}

func readHighscore() int32 {
	var f *os.File
	var err error
	p := getScorePath()

	if _, ferr := os.Stat(p); os.IsNotExist(ferr) {
		fmt.Println("creating score file")
		f, err = os.Create(p)
		check(err)
	} else {
		fmt.Println("opening score file")
		f, err = os.Open(p)
		check(err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil || fi.Size() <= 0 {
		fmt.Printf("Failed to read score file status %s\n", err)
		return 0
	}

	b := make([]byte, fi.Size())
	_, err = f.Read(b)
	check(err)

	sDec, err := base64.StdEncoding.DecodeString(string(b))
	if err != nil {
		fmt.Println("decode error:", err)
		return 0
	}

	r, _ := strconv.Atoi(string(sDec))
	return int32(r)
}

func getScorePath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	p := filepath.FromSlash(fmt.Sprintf("%s/.gosnake", usr.HomeDir))
	return p
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func gameOver(gs *GameState, highScore int32, score int32) {
	gs.changeState(GAME_OVER)
	highScore = max(highScore, score)
	createHighscore(highScore)
	sdl.Delay(1000)
}
