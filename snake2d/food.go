package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Food struct {
	dimension Dimension
	position  sdl.Point
	color     Color
}

func (f *Food) Draw(renderer *sdl.Renderer) error {
	r := sdl.Rect{
		X: f.position.X, Y: f.position.Y, W: f.dimension.W, H: f.dimension.H}

	renderer.SetDrawColor(f.color.R, f.color.G, f.color.B, f.color.A)
	renderer.FillRect(&r)

	return nil
}

func createFoodPoint(l *Level, s *Snake) sdl.Point {
	newPos := l.randomPosition()

	for _, p := range s.positions {
		if p.InRect(&sdl.Rect{X: newPos.X, Y: newPos.Y, W: s.dimension.W, H: s.dimension.H}) {
			createFoodPoint(l, s)
		}
	}

	return newPos
}
