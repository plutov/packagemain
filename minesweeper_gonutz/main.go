package main

import (
	"fmt"
	"time"

	"golang.org/x/exp/rand"

	"github.com/gonutz/prototype/draw"
)

const (
	size      = 30
	winWidth  = 30 * size
	winHeight = 31 * size
)

type state struct {
	justStarted bool
	menu        bool
	gameOver    bool
	gameWon     bool
	rows        int32
	cols        int32
	mines       int32
	field       [][]point
	startedAt   time.Time
	finishedAt  time.Time
}

type point struct {
	hasMine     bool
	open        bool
	marked      bool
	minesAround int
}

func (s *state) markedMinesAround(x, y int) int {
	n := 0
	s.doForNeighbours(x, y, func(x, y int) {
		if s.field[x][y].marked {
			n++
		}
	})
	return n
}

func (s *state) openUnmarkedAround(x, y int) {
	s.doForNeighbours(x, y, func(x, y int) {
		if !s.field[x][y].marked {
			s.revealTile(x, y)
		}
	})
}

func (s *state) reset() {
	s.gameOver = false
	s.gameWon = false
	s.menu = true
}

func (s *state) getWidth() int {
	return size * int(s.cols)
}

func (s *state) getHeight() int {
	return size*int(s.rows) + size
}

func (s *state) getStatus() string {
	var elapsed time.Duration
	if s.gameOver || s.gameWon {
		elapsed = s.finishedAt.Sub(s.startedAt)
	} else {
		elapsed = time.Since(s.startedAt)
	}
	return fmt.Sprintf("Time: %.2fs", elapsed.Seconds())
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

	s.justStarted = true
	s.menu = false
	s.startedAt = time.Now()
}

func (s *state) plantMines(notX, notY int) {
	m := s.mines
	for m > 0 {
		x, y := rand.Intn(int(s.rows)), rand.Intn(int(s.cols))
		if abs(x-notX) <= 1 && abs(y-notY) <= 1 {
			// Leave a 3x3 square around the first click empty so the user does
			// not hit a mine and also hits a nice open field at the start.
			continue
		}

		// make sure placements are unique
		if s.field[x][y].hasMine {
			continue
		}

		s.field[x][y].hasMine = true
		// mark neighbours
		s.doForNeighbours(x, y, func(x, y int) {
			s.field[x][y].minesAround++
		})
		m--
	}

	s.justStarted = false
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
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

func (s *state) isGameWon() bool {
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

	s.gameWon = s.isGameWon()
	if s.gameWon {
		s.finishedAt = time.Now()
	}

	// No neighbors, reveal all adjacent tiles recursively
	if s.field[x][y].minesAround == 0 {
		s.doForNeighbours(x, y, func(nx, ny int) {
			s.revealTile(nx, ny)
		})
	}
}

func textBox(window draw.Window, x, y, w, h int, text string, textColor, borderColor, backColor draw.Color) {
	window.FillRect(x, y, w, h, backColor)
	window.DrawRect(x, y, w, h, borderColor)
	textW, textH := window.GetTextSize(text)
	window.DrawText(text, x+(w-textW)/2, y+(h-textH)/2, textColor)
}

func button(window draw.Window, x, y, w, h int, text string) bool {
	mx, my := window.MousePosition()
	backColor := draw.LightGray
	if x <= mx && mx < x+w && y <= my && my < y+h {
		backColor = draw.LightBlue
	}
	textBox(window, x, y, w, h, text, draw.Black, draw.Black, backColor)
	for _, c := range window.Clicks() {
		if c.Button == draw.LeftButton &&
			x <= c.X && c.X < x+w &&
			y <= c.Y && c.Y < y+h {
			return true
		}
	}
	return false
}

func (s *state) drawMenu(window draw.Window) {
	var (
		rowSpacing = 50
		baseY      = 50
	)

	if clicked := button(window, 0, baseY, winWidth, size, "Beginner"); clicked {
		s.rows = 9
		s.cols = 9
		s.mines = 10
	}
	baseY += rowSpacing

	if clicked := button(window, 0, baseY, winWidth, size, "Intermediate"); clicked {
		s.rows = 16
		s.cols = 16
		s.mines = 40
	}
	baseY += rowSpacing

	if clicked := button(window, 0, baseY, winWidth, size, "Expert"); clicked {
		s.rows = 30
		s.cols = 30
		s.mines = 99
	}
	baseY += rowSpacing * 2

	if clicked := button(window, 0, baseY, winWidth, size, "Start"); clicked {
		s.start()
	}
}

func (s *state) drawField(window draw.Window) {
	w := s.getWidth()
	h := s.getHeight()

	offsetX := (winWidth - w) / 2
	offsetY := (winHeight - h) / 2

	textBox(window, offsetX, offsetY+h-size, w-85, size, s.getStatus(), draw.Black, draw.Black, draw.LightGray)
	if restart := button(window, offsetX+w-85, offsetY+h-size, 85, size, "Restart"); restart {
		s.reset()
		return
	}

	for x := range s.field {
		for y := range s.field[x] {
			if s.gameOver {
				var (
					text  string
					color draw.Color
				)

				if s.field[x][y].hasMine {
					text = "*"
					color = draw.LightRed
				} else if s.field[x][y].minesAround > 0 {
					color = getTextColor(s.field[x][y].minesAround)
					text = fmt.Sprintf("%d", s.field[x][y].minesAround)
				}

				textBox(window, offsetX+x*size, offsetY+y*size, size, size, text, color, draw.DarkGray, draw.DarkGray)
				continue
			}

			fx := offsetX + x*size
			fy := offsetY + y*size

			// Mark on right mouse button
			for _, c := range window.Clicks() {
				if fx <= c.X && c.X < fx+size &&
					fy <= c.Y && c.Y < fy+size {
					if c.Button == draw.RightButton && window.IsMouseDown(draw.LeftButton) ||
						c.Button == draw.LeftButton && window.IsMouseDown(draw.RightButton) {
						// Clicking both left and right mouse buttons on a
						// number with all its mines will open all closed
						// fields. Say you click both buttons on a 2 and there
						// are exaclty 2 mines marked around that 2. Then all
						// other closed fields will be opened.
						if s.field[x][y].open && s.field[x][y].minesAround == s.markedMinesAround(x, y) {
							s.openUnmarkedAround(x, y)
						}
					} else if c.Button == draw.RightButton && !s.field[x][y].open {
						s.field[x][y].marked = !s.field[x][y].marked
					}
				}
			}

			if s.field[x][y].marked {
				textBox(window, fx+1, fy+1, size-2, size-2, "M", draw.DarkPurple, draw.LightPurple, draw.LightPurple)
			} else if s.field[x][y].open {
				text := ""
				if s.field[x][y].minesAround > 0 {
					text = fmt.Sprintf("%d", s.field[x][y].minesAround)
				}
				textBox(window, fx, fy, size, size, text, getTextColor(s.field[x][y].minesAround), draw.DarkGray, draw.DarkGray)
			} else {
				if open := button(window, fx, fy, size, size, ""); open {
					if s.justStarted {
						s.plantMines(x, y)
					}
					s.revealTile(x, y)
				}
			}
		}
	}
}

func (s *state) drawCongrats(window draw.Window) {
	w := winWidth
	var lineHeight = 50

	if s.gameWon {
		window.DrawText("Well Done! "+s.getStatus(), 0, lineHeight, draw.White)
	}

	clicked := button(window, 0, 2*lineHeight, w, size, "Play Again")
	if clicked {
		s.reset()
	}
}

func main() {
	rand.Seed(uint64(time.Now().UnixNano()))
	game := &state{
		rows:  9,
		cols:  9,
		mines: 10,
	}
	game.reset()

	draw.RunWindow("Minesweeper", winWidth, winHeight, func(window draw.Window) {
		window.FillRect(0, 0, winWidth, winHeight, draw.DarkGray)

		if game.gameWon {
			game.drawCongrats(window)
		} else if game.menu {
			game.drawMenu(window)
		} else {
			game.drawField(window)
		}
	})
}
