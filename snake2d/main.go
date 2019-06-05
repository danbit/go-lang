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
	ScreenWidth   int32  = 640
	ScreenHeight  int32  = 480
	CellSize      int32  = 20
	ScoreText     string = "SCORE:"
	HighScoreText string = "HIGHSCORE:"
)

type GameState struct {
	To  string
	FSM *fsm.FSM
}

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
	gameState = NewGameState("snake")
	highScore = readHighscore()
	rand.Seed(time.Now().UnixNano())
	ticker := time.NewTicker(time.Second / 10)

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
					}

					snake.ChangeDirection(e)
				}
			}
		}

		<-ticker.C

		DrawBackground(renderer)
		renderer.Clear()

		switch cs := gameState.FSM.Current(); cs {
		case PLAYING.value():
			snake.Move()

			DrawFood(renderer, &food)
			DrawSnake(renderer, &snake)

			foodArea := sdl.Rect{X: food.position.X, Y: food.position.Y, H: CellSize, W: CellSize}
			if snake.CheckCollision(foodArea) {
				food.position = RandomPosition()
				snake.AddTail()
				score += 10
			}

			position := sdl.Point{X: 5, Y: 0}
			RenderText(renderer, &position, &scoreRect, ScoreText, font)
			position = sdl.Point{X: scoreRect.X + scoreRect.W + 5, Y: 0}
			RenderText(renderer, &position, &scoreRect, FormatInt32(score), font)

			position = sdl.Point{X: ScreenWidth - 200, Y: 0}
			RenderText(renderer, &position, &highScoreRect, HighScoreText, font)
			position = sdl.Point{X: highScoreRect.X + highScoreRect.W + 5, Y: 0}
			RenderText(renderer, &position, &highScoreRect, FormatInt32(highScore), font)

			if snake.IsTryingToEat() {
				gameState.changeState(GAME_OVER)
				sdl.Delay(500)
			}
		case ON_MENU.value():
			RenderText(renderer, nil, nil, "<Press ENTER to Start>", font)
		case PAUSED.value():
			RenderText(renderer, nil, nil, "Paused", font)
		case GAME_OVER.value():
			highScore = max(highScore, score)
			createHighscore(highScore)
			RenderText(renderer, nil, nil, "Game Over", font)
		}

		renderer.Present()
	}

	os.Exit(0)
}

func DrawBackground(renderer *sdl.Renderer) {
	renderer.SetDrawColor(0, 0, 0, 0)
}

func Draw(renderer *sdl.Renderer, rect *sdl.Rect) {
	renderer.SetDrawColor(255, 255, 255, 255)
	renderer.FillRect(rect)
}

func DrawSnake(renderer *sdl.Renderer, snake *Snake) {
	for _, element := range snake.positions {
		Draw(renderer, &sdl.Rect{
			X: element.X, Y: element.Y, W: snake.dimension.width, H: snake.dimension.height})
	}
}

func DrawFood(renderer *sdl.Renderer, food *Food) {
	Draw(renderer, &sdl.Rect{
		X: food.position.X, Y: food.position.Y, W: food.dimension.width, H: food.dimension.height})
}

func RenderText(renderer *sdl.Renderer, position *sdl.Point, rect *sdl.Rect, text string, font *ttf.Font) {

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
	return math.Round(float64(rand.Float64()*float64((max-min)+min)/float64(CellSize))) * float64(CellSize)
}

func RandomPosition() sdl.Point {
	newX := random(0, ScreenWidth-CellSize)
	newY := random(0, ScreenHeight-CellSize)

	newPosition := sdl.Point{
		X: int32(newX),
		Y: int32(newY),
	}

	return newPosition
}

func (s *Snake) Move() {
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

func (s *Snake) AddTail() {
	newPart := sdl.Point{X: 1, Y: 1}
	s.positions = append([]sdl.Point{newPart}, s.positions...)
}

func (s *Snake) ChangeDirection(e *sdl.KeyboardEvent) {
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

func (s *Snake) CheckCollision(area sdl.Rect) bool {
	if head := s.positions[len(s.positions)-1]; head.InRect(&area) {
		return true
	}

	return false
}

func (s *Snake) IsTryingToEat() bool {
	for i := 0; i < len(s.positions)-1; i++ {
		position := s.positions[i]
		area := sdl.Rect{X: position.X, Y: position.Y, H: s.dimension.height, W: s.dimension.width}

		if s.CheckCollision(area) {
			return true
		}
	}

	return false
}

func FormatInt32(n int32) string {
	return strconv.FormatInt(int64(n), 10)
}

func NewGameState(to string) *GameState {
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
		snake, food, score = CreateLevel()
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

func CreateLevel() (snake Snake, food Food, score int32) {
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
		position:  RandomPosition(),
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

	sEnc := base64.StdEncoding.EncodeToString([]byte(FormatInt32(hs)))

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
