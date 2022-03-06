package snake

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 600
	ScreenHeight = 600
	boardRows    = 20
	boardCols    = 20
)

// Game represents a game state.
type Game struct {
	input *Input
	board *Board
}

// NewGame generates a new Game object.
func NewGame() *Game {
	return &Game{
		input: NewInput(),
		board: NewBoard(boardRows, boardCols),
	}
}

// Layout implements ebiten.Game's Layout.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

// Update updates the current game state.
func (g *Game) Update() error {
	return g.board.Update(g.input)
}

// Draw draws the current game to the given screen.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)
}
