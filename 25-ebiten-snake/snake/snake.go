package snake

type Coord struct {
	x, y int
}

const (
	RIGHT int = iota
	LEFT
	UP
	DOWN
)

type Snake struct {
	body      []Coord
	direction int
}

func NewSnake(body []Coord) *Snake {
	return &Snake{
		body: body,
	}
}

func (s *Snake) Head() Coord {
	return s.body[len(s.body)-1]
}

func (s *Snake) HeadHits(x, y int) bool {
	h := s.Head()

	return h.x == x && h.y == y
}

func (s *Snake) HeadHitsBody() bool {
	h := s.Head()
	bodyWithoutHead := s.body[:len(s.body)-1]

	for _, b := range bodyWithoutHead {
		if b.x == h.x && b.y == b.y {
			return true
		}
	}

	return false
}

func (s *Snake) Move() {
	h := s.Head()
	newHead := Coord{x: h.x, y: h.y}

	switch s.direction {
	case RIGHT:
		newHead.y++
	case LEFT:
		newHead.y--
	case UP:
		newHead.x--
	case DOWN:
		newHead.x++
	}

	s.body = append(s.body[1:], newHead)
}
