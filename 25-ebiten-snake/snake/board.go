package snake

import "fmt"

// Board represents the game board.
type Board struct {
	rows int
	cols int
}

// NewBoard generates a new Board with giving a size.
func NewBoard(rows int, cols int) *Board {
	return &Board{
		rows: rows,
		cols: cols,
	}
}

// Update updates the board state.
func (b *Board) Update(input *Input) error {
	if dir, ok := input.Dir(); ok {
		fmt.Println(dir)
	}
	return nil
}
