package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Scene struct {
	GameObjects []GameObject
}

func (s *Scene) Draw(renderer *sdl.Renderer) error {
	for _, gob := range s.GameObjects {
		if err := gob.Draw(renderer); err != nil {
			return err
		}
	}

	return nil
}
