package main

import "github.com/veandco/go-sdl2/sdl"

type Snake struct {
	direction sdl.Point //direction X and Y
	dimension Dimension
	positions []sdl.Point
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
