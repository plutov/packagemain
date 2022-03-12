package snake

type Coord struct {
	x, y int
}

type Snake struct {
	body      []Coord
	direction Dir
	justAte   bool
}

func NewSnake(body []Coord, direction Dir) *Snake {
	return &Snake{
		body:      body,
		direction: direction,
	}
}

func (s *Snake) Head() Coord {
	return s.body[len(s.body)-1]
}

func (s *Snake) ChangeDirection(newDir Dir) {
	opposites := map[Dir]Dir{
		DirUp:    DirDown,
		DirRight: DirLeft,
		DirDown:  DirUp,
		DirLeft:  DirRight,
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
	case DirUp:
		newHead.x--
	case DirRight:
		newHead.y++
	case DirDown:
		newHead.x++
	case DirLeft:
		newHead.y--
	}

	if s.justAte {
		s.body = append(s.body, newHead)
		s.justAte = false
	} else {
		s.body = append(s.body[1:], newHead)
	}
}
