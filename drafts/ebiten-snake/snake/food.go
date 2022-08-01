package snake

type Food struct {
	x, y int
}

func NewFood(x, y int) *Food {
	return &Food{
		x: x,
		y: y,
	}
}
