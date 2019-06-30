package main

import "github.com/veandco/go-sdl2/sdl"

type Level struct {
	dimension Dimension
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
			sdl.Point{X: screenCenterW + CellSize*2, Y: screenCenterH},
		},
	}

	food = Food{
		dimension: Dimension{W: CellSize, H: CellSize},
		position:  createFoodPoint(&snake),
	}

	score = 0

	return
}

func createFoodPoint(s *Snake) sdl.Point {
	newPos := randomPosition()

	for _, p := range s.positions {
		if p.InRect(&sdl.Rect{X: newPos.X, Y: newPos.Y, W: s.dimension.W, H: s.dimension.H}) {
			createFoodPoint(s)
		}
	}

	return newPos
}
