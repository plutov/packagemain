// Copyright (c) 2017 Alex Pliutau

package snake

import (
	"math/rand"
)

type coord struct {
	x, y int
}

type arena struct {
	food       *food
	snake      *snake
	hasFood    func(*arena, coord) bool
	height     int
	width      int
	pointsChan chan (int)
}

func newArena(s *snake, h, w int) *arena {
	a := &arena{
		snake:   s,
		height:  h,
		width:   w,
		hasFood: hasFood,
	}

	a.placeFood()

	return a
}

func (a *arena) placeFood() {
	var x, y int

	for {
		x = rand.Intn(a.width)
		y = rand.Intn(a.height)

		if !a.snake.hits(coord{x: x, y: y}) {
			break
		}
	}

	a.food = newFood(x, y)
}

func (a *arena) moveSnake() error {
	if err := a.snake.move(); err != nil {
		return err
	}

	if a.snakeLeftArena() {
		return a.snake.die()
	}

	if a.hasFood(a, a.snake.head()) {
		go func() {
			a.pointsChan <- a.food.points
		}()
		a.snake.length++
		a.placeFood()
	}

	return nil
}

func (a *arena) snakeLeftArena() bool {
	h := a.snake.head()
	return h.x > a.width-1 || h.y > a.height-1 || h.x < 0 || h.y < 0
}

func hasFood(a *arena, c coord) bool {
	return c.x == a.food.x && c.y == a.food.y
}
