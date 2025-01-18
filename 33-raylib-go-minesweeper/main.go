package main

import (
	"fmt"

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
	defaultWinHeight = 400
)

type state struct {
	menu     bool
	gameOver bool
	won      bool
	rows     int32
	cols     int32
	mines    int32
	// [x][y]
	field [][]point
}

type point struct {
	hasMine    bool
	open       bool
	marked     bool
	neighbours int
}

// buttons and paddings
func (s *state) getWidth() int32 {
	return size*s.cols + 2*padding
}

func (s *state) getHeight() int32 {
	// with status bar
	return size*s.rows + 2*padding + size
}

func (s *state) getStatus() string {
	fps := rl.GetFPS()

	return fmt.Sprintf("FPS: %d", fps)
}

func (s *state) reset() {
	s.gameOver = false
	s.won = false
	s.menu = true
	s.rows = 9
	s.cols = 9
	s.mines = 10
}

func (s *state) start() {
	s.field = make([][]point, s.rows)
	for x := range s.rows {
		s.field[x] = make([]point, s.cols)
		for y := range s.cols {
			s.field[x][y] = point{}
		}
	}

	// plant mines
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
}

func (s *state) doForNeighbours(x, y int, do func(x, y int)) {
	// with diagonals
	dx := []int{-1, -1, -1, 0, 0, 1, 1, 1}
	dy := []int{-1, 0, 1, -1, 1, -1, 0, 1}

	for i := 0; i < len(dx); i++ {
		nx := x + dx[i]
		ny := y + dy[i]

		if nx >= 0 && nx < int(s.cols) && ny >= 0 && ny < int(s.rows) {
			do(nx, ny)
		}
	}
}

func (s *state) drawMenu() {
	w, h := defaultWinWidth, defaultWinHeight
	colw := float32(w / 2)
	var lineHeight int32 = 50
	rl.SetWindowSize(w, h)

	rl.DrawText("ROWS:", padding, lineHeight, size, rl.White)
	s.rows = gui.Spinner(rl.NewRectangle(colw, float32(lineHeight), float32(colw-padding), size), "", &s.rows, minRowsCols, maxRowsCols, true)

	rl.DrawText("COLS:", padding, 2*lineHeight, size, rl.White)
	s.cols = gui.Spinner(rl.NewRectangle(colw, float32(2*lineHeight), float32(colw-padding), size), "", &s.cols, minRowsCols, maxRowsCols, true)

	rl.DrawText("MINES:", padding, 3*lineHeight, size, rl.White)
	s.mines = gui.Spinner(rl.NewRectangle(colw, float32(3*lineHeight), float32(colw-padding), size), "", &s.mines, minRowsCols, maxRowsCols, true)

	clicked := gui.Button(rl.NewRectangle(padding, float32(4*lineHeight), float32(w-2*padding), size), "START")
	if clicked {
		s.start()
	}
}

func (s *state) drawField() {
	w := float32(s.getWidth())
	h := float32(s.getHeight())

	rl.SetWindowSize(int(w), int(h))
	gui.StatusBar(rl.NewRectangle(0, h-size, w, size), s.getStatus())

	for x := range s.field {
		for y := range s.field[x] {
			if s.field[x][y].open {
				rl.DrawRectangle(padding+int32(x)*size, padding+int32(y)*size, size, size, rl.Gray)
				if s.field[x][y].hasMine {
					gui.DrawIcon(gui.ICON_STAR, 5+padding+int32(x)*size, 5+padding+int32(y)*size, 1, rl.Red)
					s.gameOver = true
				} else {
					rl.DrawText(fmt.Sprintf("%d", s.field[x][y].neighbours), 5+padding+int32(x)*size, 5+padding+int32(y)*size, 20, rl.DarkGreen)
				}
			} else {
				s.field[x][y].open = gui.Button(rl.NewRectangle(float32(padding+x*size), float32(padding+y*size), size, size), "")
			}
		}
	}
}

func (s *state) drawMessageWithRestart() {
	w, h := defaultWinWidth, defaultWinHeight
	var lineHeight int32 = 50
	rl.SetWindowSize(w, h)

	if s.gameOver {
		rl.DrawText("GAME OVER :(", padding, lineHeight, size, rl.White)
	}
	if s.won {
		rl.DrawText("WELL DONE !", padding, lineHeight, size, rl.White)
	}

	clicked := gui.Button(rl.NewRectangle(padding, float32(2*lineHeight), float32(w-2*padding), size), "ONE MORE")
	if clicked {
		s.reset()
	}
}

func main() {
	game := &state{}
	game.reset()

	rl.InitWindow(defaultWinWidth, defaultWinHeight, "minesweeper")
	rl.SetTargetFPS(1000)
	defer rl.CloseWindow()

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.DarkGray)

		if game.gameOver {
			game.drawMessageWithRestart()
		} else if game.menu {
			game.drawMenu()
		} else {
			game.drawField()
		}

		rl.EndDrawing()
	}
}
