package snake

import "github.com/hajimehoshi/ebiten/v2"

type Coord struct {
	x, y int
}

type Snake struct {
	body      []Coord
	direction ebiten.Key
	justAte   bool
}

func NewSnake(body []Coord, direction ebiten.Key) *Snake {
	return &Snake{
		body:      body,
		direction: direction,
	}
}

func (s *Snake) Head() Coord {
	return s.body[len(s.body)-1]
}

func (s *Snake) ChangeDirection(newDir ebiten.Key) {
	opposites := map[ebiten.Key]ebiten.Key{
		ebiten.KeyArrowUp:    ebiten.KeyArrowDown,
		ebiten.KeyArrowRight: ebiten.KeyArrowLeft,
		ebiten.KeyArrowDown:  ebiten.KeyArrowUp,
		ebiten.KeyArrowLeft:  ebiten.KeyArrowRight,
	}

	// don't allow changing direction to opposite
	if o, ok := opposites[newDir]; ok && o != s.direction {
		s.direction = newDir
	}
}

func (s *Snake) HeadHits(x, y int) bool {
	h := s.Head()

	return h.x == x && h.y == y
}

func (s *Snake) HeadHitsBody() bool {
	h := s.Head()
	bodyWithoutHead := s.body[:len(s.body)-1]

	for _, b := range bodyWithoutHead {
		if b.x == h.x && b.y == h.y {
			return true
		}
	}

	return false
}

func (s *Snake) Move() {
	h := s.Head()
	newHead := Coord{x: h.x, y: h.y}

	switch s.direction {
	case ebiten.KeyArrowUp:
		newHead.x--
	case ebiten.KeyArrowRight:
		newHead.y++
	case ebiten.KeyArrowDown:
		newHead.x++
	case ebiten.KeyArrowLeft:
		newHead.y--
	}

	if s.justAte {
		s.body = append(s.body, newHead)
		s.justAte = false
	} else {
		s.body = append(s.body[1:], newHead)
	}
}
