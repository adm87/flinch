package flinch

import "github.com/hajimehoshi/ebiten/v2"

type Screen interface {
	Buffer() *ebiten.Image
	SetSize(width, height int)
}

type screen struct {
	buffer *ebiten.Image
}

func NewScreen() Screen {
	return &screen{}
}

func (s *screen) Buffer() *ebiten.Image {
	return s.buffer
}

func (s *screen) SetSize(width, height int) {
	if s.buffer == nil {
		s.buffer = ebiten.NewImage(width, height)
		return
	}

	if s.buffer.Bounds().Dx() != width || s.buffer.Bounds().Dy() != height {
		s.buffer.Deallocate()
		s.buffer = ebiten.NewImage(width, height)
	}
}
