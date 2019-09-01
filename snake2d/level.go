package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Level struct {
	dimension Dimension
	snake     Snake
	food      Food
	score     int32
}

// NewLevel return a level struct. Thread safety singleton.
func NewLevel() *Level {
	l := &Level{}

	d := Dimension{
		W: ScreenWidth - BorderSize,
		H: ScreenHeight - BorderSize,
	}
	l.dimension = d

	screenCenterW := (d.W/CellSize)/2*CellSize + CellSize
	screenCenterH := (d.H/CellSize)/2*CellSize + CellSize

	s := Snake{
		direction: sdl.Point{X: 1, Y: 0},
		dimension: Dimension{W: CellSize, H: CellSize},
		positions: []sdl.Point{
			sdl.Point{X: screenCenterW, Y: screenCenterH},
			sdl.Point{X: screenCenterW + CellSize, Y: screenCenterH},
			sdl.Point{X: screenCenterW + CellSize*2, Y: screenCenterH},
		},
		color: Color{R: 255, G: 255, B: 255, A: 255},
	}
	l.snake = s

	f := Food{
		dimension: Dimension{W: CellSize, H: CellSize},
		position:  createFoodPoint(l, &snake),
		color:     Color{R: 255, G: 255, B: 255, A: 255},
	}
	l.food = f

	return l
}

func (l *Level) randomPosition() sdl.Point {
	newX := random(BorderSize, l.dimension.W)
	newY := random(BorderSize, l.dimension.H)

	newPosition := sdl.Point{
		X: int32(newX),
		Y: int32(newY),
	}

	return newPosition
}
