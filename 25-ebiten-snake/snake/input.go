package snake

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Dir represents a direction.
type Dir int

const (
	DirUp Dir = iota
	DirRight
	DirDown
	DirLeft
)

// Input represents the current key states.
type Input struct{}

// NewInput generates a new Input object.
func NewInput() *Input {
	return &Input{}
}

// Dir returns a currently pressed direction.
// Dir returns false if no direction key is pressed.
func (i *Input) Dir() (Dir, bool) {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		return DirUp, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		return DirLeft, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		return DirRight, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		return DirDown, true
	}

	return 0, false
}
