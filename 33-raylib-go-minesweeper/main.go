package main

import (
	"fmt"
	"time"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"golang.org/x/exp/rand"
)

const (
	size             = 30
	padding          = 20
	minRowsCols      = 3
	maxRowsCols      = 20
	defaultWinWidth  = 300
	defaultWinHeight = 450
)

type state struct {
	menu      bool
	gameOver  bool
	gameWon   bool
	startedAt time.Time
	rows      int32
	cols      int32
	mines     int32
	// [x][y]
	field [][]point
}

type point struct {
	hasMine    bool
	open       bool
	marked     bool
	neighbours int
}

func (s *state) getWidth() int32 {
	// buttons, paddings
	return size*s.cols + 2*padding
}

func (s *state) getHeight() int32 {
	// buttons, paddings, status bar
	return size*s.rows + 2*padding + size
}

func (s *state) getStatus() string {
	fps := rl.GetFPS()
	elapsed := time.Since(s.startedAt)

	return fmt.Sprintf("FPS: %d, TIME: %d", fps, int(elapsed.Seconds()))
}

func (s *state) reset() {
	s.gameOver = false
	s.gameWon = false
	s.menu = true
}

func (s *state) start() {
	// Build grid
	s.field = make([][]point, s.rows)
	for x := range s.rows {
		s.field[x] = make([]point, s.cols)
		for y := range s.cols {
			s.field[x][y] = point{}
		}
	}

	// Plant mines
	m := s.mines
	for m > 0 {
		x, y := rand.Intn(int(s.rows)), rand.Intn(int(s.cols))

		// make sure placements are unique
		if s.field[x][y].hasMine {
			continue
		}

		s.field[x][y].hasMine = true
		// mark neighbours
		s.doForNeighbours(x, y, func(x, y int) {
			s.field[x][y].neighbours++
		})
		m--
	}

	s.menu = false
	s.startedAt = time.Now()
}

func (s *state) doForNeighbours(x, y int, do func(x, y int)) {
	// with diagonals
	dx := []int{-1, 0, 1, -1, 1, -1, 0, 1}
	dy := []int{-1, -1, -1, 0, 0, 1, 1, 1}

	for i := 0; i < len(dx); i++ {
		nx := x + dx[i]
		ny := y + dy[i]

		if nx >= 0 && nx < int(s.rows) && ny >= 0 && ny < int(s.cols) {
			do(nx, ny)
		}
	}
}

func (s *state) checkIfGameWon() bool {
	open := 0
	total := int(s.rows * s.cols)

	for x := 0; x < int(s.rows); x++ {
		for y := 0; y < int(s.cols); y++ {
			if s.field[x][y].open {
				open++
			}
		}
	}

	return open == total-int(s.mines)
}

func (s *state) revealTile(x, y int) {
	if s.field[x][y].open {
		return
	}

	s.field[x][y].open = true

	if s.field[x][y].hasMine {
		s.gameOver = true
		return
	}

	s.gameWon = s.checkIfGameWon()

	// No neighbors, reveal all adjacent tiles recursively.
	if s.field[x][y].neighbours == 0 {
		s.doForNeighbours(x, y, func(nx, ny int) {
			s.revealTile(nx, ny)
		})
	}
}

func (s *state) drawMenu() {
	w, h := defaultWinWidth, defaultWinHeight
	colw := float32(w / 2)
	var rowh float32 = 50
	var fontSize int32 = 20
	buttonWidth := float32(w - 2*padding)

	rl.SetWindowSize(w, h)
	rl.SetWindowPosition((rl.GetMonitorWidth(0)-int(w))/2, (rl.GetMonitorHeight(0)-int(h))/2)

	if clicked := gui.Button(rl.NewRectangle(padding, rowh, buttonWidth, size), "BEGINNER"); clicked {
		s.rows = 9
		s.cols = 9
		s.mines = 10
	}
	if clicked := gui.Button(rl.NewRectangle(padding, 2*rowh, buttonWidth, size), "INTERMEDIATE"); clicked {
		s.rows = 16
		s.cols = 16
		s.mines = 40
	}
	if clicked := gui.Button(rl.NewRectangle(padding, 3*rowh, buttonWidth, size), "EXPERT"); clicked {
		s.rows = 30
		s.cols = 30
		s.mines = 99
	}

	rl.DrawText("ROWS:", padding, 4*int32(rowh)+5, fontSize, rl.White)
	s.rows = gui.Spinner(rl.NewRectangle(colw, 4*rowh, float32(colw-padding), size), "", &s.rows, minRowsCols, maxRowsCols, true)

	rl.DrawText("COLS:", padding, 5*int32(rowh)+5, fontSize, rl.White)
	s.cols = gui.Spinner(rl.NewRectangle(colw, 5*rowh, float32(colw-padding), size), "", &s.cols, minRowsCols, maxRowsCols, true)

	rl.DrawText("MINES:", padding, 6*int32(rowh)+5, fontSize, rl.White)
	s.mines = gui.Spinner(rl.NewRectangle(colw, 6*rowh, float32(colw-padding), size), "", &s.mines, 1, int(s.rows)*int(s.cols), true)

	if clicked := gui.Button(rl.NewRectangle(padding, 7*rowh, buttonWidth, size), "START"); clicked {
		s.start()
	}

}

func getTextColor(neighbors int) rl.Color {
	switch neighbors {
	case 1:
		return rl.Blue
	case 2:
		return rl.Green
	case 3:
		return rl.Red
	default:
		return rl.Black
	}
}

func (s *state) drawField() {
	w := float32(s.getWidth())
	h := float32(s.getHeight())

	rl.SetWindowSize(int(w), int(h))
	rl.SetWindowPosition((rl.GetMonitorWidth(0)-int(w))/2, (rl.GetMonitorHeight(0)-int(h))/2)

	gui.StatusBar(rl.NewRectangle(0, h-size, w, size), s.getStatus())
	if restart := gui.Button(rl.NewRectangle(w-65, h-size+5, 60, size-10), "RESTART"); restart {
		s.reset()
		return
	}

	for x := range s.field {
		for y := range s.field[x] {
			rect := rl.NewRectangle(float32(padding+x*size), float32(padding+y*size), size, size)

			if s.gameOver {
				// reveal current state
				if s.field[x][y].hasMine {
					rl.DrawText("*", 5+padding+int32(x)*size, 5+padding+int32(y)*size, 20, rl.Red)
				} else {
					text := ""
					if s.field[x][y].neighbours > 0 {
						text = fmt.Sprintf("%d", s.field[x][y].neighbours)
					}

					rl.DrawText(text, 5+padding+int32(x)*size, 5+padding+int32(y)*size, 20, getTextColor(s.field[x][y].neighbours))
				}
				continue
			}

			// Mark on right mouse button
			if rl.IsMouseButtonPressed(rl.MouseButtonRight) {
				if rl.CheckCollisionPointRec(rl.GetMousePosition(), rect) {
					if !s.field[x][y].open {
						s.field[x][y].marked = !s.field[x][y].marked
					}
				}
			}

			if s.field[x][y].marked {
				rl.DrawText("M", 5+padding+int32(x)*size, 5+padding+int32(y)*size, 20, rl.Violet)
			} else if s.field[x][y].open {
				text := ""
				if s.field[x][y].neighbours > 0 {
					text = fmt.Sprintf("%d", s.field[x][y].neighbours)
				}

				rl.DrawText(text, 5+padding+int32(x)*size, 5+padding+int32(y)*size, 20, getTextColor(s.field[x][y].neighbours))
			} else {
				if open := gui.Button(rect, ""); open {
					s.revealTile(x, y)
				}
			}
		}
	}
}

func (s *state) congrats() {
	w, h := defaultWinWidth, defaultWinHeight
	var lineHeight int32 = 50
	rl.SetWindowSize(w, h)

	if s.gameWon {
		rl.DrawText("WELL DONE !", padding, lineHeight, size, rl.White)
	}

	clicked := gui.Button(rl.NewRectangle(padding, float32(2*lineHeight), float32(w-2*padding), size), "PLAY AGAIN")
	if clicked {
		s.reset()
	}
}

func main() {
	game := &state{
		rows:  9,
		cols:  9,
		mines: 10,
	}
	game.reset()

	rl.InitWindow(defaultWinWidth, defaultWinHeight, "minesweeper")
	rl.SetTargetFPS(1000)
	defer rl.CloseWindow()

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.DarkGray)

		if game.gameWon {
			game.congrats()
		} else if game.menu {
			game.drawMenu()
		} else {
			game.drawField()
		}

		rl.EndDrawing()
	}
}
