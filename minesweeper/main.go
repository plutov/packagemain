package main

import (
	"fmt"
	"time"

	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
	"golang.org/x/exp/rand"
)

const (
	size      = 30
	padding   = 20
	minSize   = 3
	maxSize   = 20
	winWidth  = 300
	winHeight = 450
)

type state struct {
	menu       bool
	gameOver   bool
	gameWon    bool
	startedAt  time.Time
	finishedAt time.Time
	rows       int32
	cols       int32
	mines      int32
	field      [][]point
}

type point struct {
	hasMine    bool
	open       bool
	marked     bool
	neighbours int
}

func (s *state) reset() {
	s.gameOver = false
	s.gameWon = false
	s.menu = true
}

func (s *state) getWidth() int32 {
	return size*s.cols + 2*padding
}

func (s *state) getHeight() int32 {
	return size*s.rows + 2*padding + size
}

func (s *state) getStatus() string {
	fps := rl.GetFPS()
	var elapsed time.Duration
	if s.gameOver || s.gameWon {
		elapsed = s.finishedAt.Sub(s.startedAt)
	} else {
		elapsed = time.Since(s.startedAt)
	}

	return fmt.Sprintf("FPS: %d, TIME: %.2f", fps, elapsed.Seconds())
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

	for i := range len(dx) {
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

	for x := range s.rows {
		for y := range s.cols {
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
		s.finishedAt = time.Now()
		return
	}

	s.gameWon = s.checkIfGameWon()

	// No neighbors, reveal all adjacent tiles recursively
	if s.field[x][y].neighbours == 0 {
		s.doForNeighbours(x, y, func(nx, ny int) {
			s.revealTile(nx, ny)
		})
	}
}

func (s *state) drawMenu() {
	w := winWidth
	colw := float32(w / 2)
	rowSpacing := float32(50)
	baseY := rowSpacing
	var fontSize int32 = 20
	buttonWidth := float32(w - 2*padding)

	if clicked := gui.Button(rl.NewRectangle(padding, baseY, buttonWidth, size), "BEGINNER"); clicked {
		s.rows = 9
		s.cols = 9
		s.mines = 10
	}
	baseY += rowSpacing

	if clicked := gui.Button(rl.NewRectangle(padding, baseY, buttonWidth, size), "INTERMEDIATE"); clicked {
		s.rows = 16
		s.cols = 16
		s.mines = 40
	}
	baseY += rowSpacing

	if clicked := gui.Button(rl.NewRectangle(padding, baseY, buttonWidth, size), "EXPERT"); clicked {
		s.rows = 30
		s.cols = 30
		s.mines = 99
	}
	baseY += rowSpacing

	rl.DrawText("ROWS:", padding, int32(baseY)+5, fontSize, rl.White)
	s.rows = gui.Spinner(rl.NewRectangle(colw, baseY, float32(colw-padding), size), "", &s.rows, minSize, maxSize, true)
	baseY += rowSpacing

	rl.DrawText("COLS:", padding, int32(baseY)+5, fontSize, rl.White)
	s.cols = gui.Spinner(rl.NewRectangle(colw, baseY, float32(colw-padding), size), "", &s.cols, minSize, maxSize, true)
	baseY += rowSpacing

	rl.DrawText("MINES:", padding, int32(baseY)+5, fontSize, rl.White)
	s.mines = gui.Spinner(rl.NewRectangle(colw, baseY, float32(colw-padding), size), "", &s.mines, 1, int(s.rows)*int(s.cols), true)
	baseY += rowSpacing

	if clicked := gui.Button(rl.NewRectangle(padding, baseY, buttonWidth, size), "START"); clicked {
		s.start()
	}
}

func (s *state) drawField() {
	w := float32(s.getWidth())
	h := float32(s.getHeight())

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
	w, h := winWidth, winHeight
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

	rl.InitWindow(winWidth, winHeight, "minesweeper")
	rl.SetTargetFPS(60)
	defer rl.CloseWindow()

	for !rl.WindowShouldClose() {
		if game.menu {
			centerWindow(winWidth, winHeight)
		} else if game.gameWon {
			centerWindow(winWidth, winHeight)
		} else {
			centerWindow(game.getWidth(), game.getHeight())
		}

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
