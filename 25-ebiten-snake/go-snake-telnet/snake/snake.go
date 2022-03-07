// Copyright (c) 2017 Alex Pliutau

package snake

import (
	"errors"
)

const (
	// RIGHT const
	RIGHT direction = 1 + iota
	// LEFT const
	LEFT
	// UP const
	UP
	// DOWN const
	DOWN
)

type direction int

type snake struct {
	body      []coord
	direction direction
	length    int
}

func newSnake(d direction, b []coord) *snake {
	return &snake{
		length:    len(b),
		body:      b,
		direction: d,
	}
}

func (s *snake) changeDirection(d direction) {
	opposites := map[direction]direction{
		RIGHT: LEFT,
		LEFT:  RIGHT,
		UP:    DOWN,
		DOWN:  UP,
	}

	if o := opposites[d]; o != 0 && o != s.direction {
		s.direction = d
	}
}

func (s *snake) head() coord {
	return s.body[len(s.body)-1]
}

func (s *snake) die() error {
	return errors.New("Game over")
}

func (s *snake) move() error {
	h := s.head()
	c := coord{x: h.x, y: h.y}

	switch s.direction {
	case RIGHT:
		c.y++
	case LEFT:
		c.y--
	case UP:
		c.x--
	case DOWN:
		c.x++
	}

	if s.hits(c) {
		return s.die()
	}

	if s.length > len(s.body) {
		s.body = append(s.body, c)
	} else {
		s.body = append(s.body[1:], c)
	}

	return nil
}

func (s *snake) hits(c coord) bool {
	for _, b := range s.body {
		if b.x == c.x && b.y == c.y {
			return true
		}
	}

	return false
}
