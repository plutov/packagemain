package snake

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Input represents the current key states.
type Input struct{}

// NewInput generates a new Input object.
func NewInput() *Input {
	return &Input{}
}

// Dir returns a currently pressed direction.
// Dir returns false if no direction key is pressed.
func (i *Input) Dir() (ebiten.Key, bool) {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		return ebiten.KeyArrowUp, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		return ebiten.KeyArrowLeft, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		return ebiten.KeyArrowRight, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		return ebiten.KeyArrowDown, true
	}

	return 0, false
}
