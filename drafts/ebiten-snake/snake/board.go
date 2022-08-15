package snake

import (
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Board struct {
	rows     int
	cols     int
	food     *Food
	snake    *Snake
	points   int
	gameOver bool
	timer    time.Time
}

func NewBoard(rows int, cols int) *Board {
	rand.Seed(time.Now().UnixNano())

	board := &Board{
		rows:  rows,
		cols:  cols,
		timer: time.Now(),
	}
	// start in top-left corner
	board.snake = NewSnake([]Coord{{0, 0}, {0, 1}, {0, 2}, {0, 3}}, ebiten.KeyArrowRight)
	board.placeFood()

	return board
}

func (b *Board) Update(input *Input) error {
	if b.gameOver {
		return nil
	}

	// snake goes faster when there are more points
	interval := time.Millisecond * 200
	if b.points > 10 {
		interval = time.Millisecond * 150
	} else if b.points > 20 {
		interval = time.Millisecond * 100
	}

	if newDir, ok := input.Dir(); ok {
		b.snake.ChangeDirection(newDir)
	}

	if time.Since(b.timer) >= interval {
		if err := b.moveSnake(); err != nil {
			return err
		}

		b.timer = time.Now()
	}

	return nil
}

func (b *Board) placeFood() {
	var x, y int

	for {
		x = rand.Intn(b.cols)
		y = rand.Intn(b.rows)

		// make sure we don't put a food on a snake
		if !b.snake.HeadHits(x, y) {
			break
		}
	}

	b.food = NewFood(x, y)
}

func (b *Board) moveSnake() error {
	// remove tail first, add 1 in front
	b.snake.Move()

	if b.snakeLeftBoard() || b.snake.HeadHitsBody() {
		b.gameOver = true
		return nil
	}

	if b.snake.HeadHits(b.food.x, b.food.y) {
		// the snake grows on the next move
		b.snake.justAte = true

		b.placeFood()
		b.points++
	}

	return nil
}

func (b *Board) snakeLeftBoard() bool {
	head := b.snake.Head()
	return head.x > b.cols-1 || head.y > b.rows-1 || head.x < 0 || head.y < 0
}
