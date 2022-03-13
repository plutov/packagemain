## packagemain #25: Snake game in Go using Ebiten

[Ebiten](https://ebiten.org/) is an open source game library in Go for building 2D games that can be ran across multiple platforms. Ebiten games work on desktop, web browsers (through WebAssembly), as well as on Mobile and even on Nintendo Switch.

In this video we'll give it a try and create a Snake game in Go, which we'll run in the browser using WebAssembly.

### Ebiten API

![overview](https://ebiten.org/images/overview2.2.png)

To create a new game we have to implement an `ebiten.Game` interface which consist of the following functions as it shown on the image above:
- `Update()` - update the logical state
- `Draw(screen *ebiten.Image)` - render the screen
- `Layout(outsideWidth, outsideHeight int)` - define the size of the screen

Let's create our `Snake` struct to implement these methods. For now it will do nothing, except setting the size of the screen and filling it with a background. We want to confirm that we can compile it and see the results in the browser.

**snake/game.go**

```go
package snake

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	ScreenWidth  = 600
	ScreenHeight = 600
	boardRows    = 20
	boardCols    = 20
)

var (
	backgroundColor = color.RGBA{50, 100, 50, 50}
)

type Game struct {}

func NewGame() *Game {
	return &Game{}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) Update() error {
	// todo
    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)
}
```

**main.go**

```go
package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/plutov/packagemain/25-ebiten-snake/snake"
)

func main() {
	game := snake.NewGame()

	ebiten.SetWindowSize(snake.ScreenWidth, snake.ScreenHeight)
	ebiten.SetWindowTitle("Snake")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
```

### Compiling into WebAssembly

As I mentioned before we'll run this game in a browser, as it's easy to compile it in WebAssembly. There are 2 options running this game in the browser.

1. We can just compile it manually and create our own HTML file
2. We can use [wasmserve](https://github.com/hajimehoshi/wasmserve) project that will compile a code and run it in the browser for us

For simplicity I'll use the second option.

```
go install github.com/hajimehoshi/wasmserve@latest
wasmserve .
```

Then if we open http://localhost:8080/ we should see a gray-ish square.

### Snake Logic

Now we can write a snake logic that we'll be using later in functions `Update` and `Draw`.

#### Food

`Food` struct simply describes the point on the board that snake should eat.

**snake/food.go**

```go
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
```

`Input` struct is using `inpututil` Ebiten's package to check the latest key pressed.

**snake/input.go**

```go
package snake

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Input struct{}

func NewInput() *Input {
	return &Input{}
}

func (i *Input) Dir() (ebiten.Key, bool) {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		return ebiten.KeyArrowUp, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		return ebiten.KeyArrowLeft, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		return ebiten.KeyArrowRight, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		return ebiten.KeyArrowDown, true
	}

	return 0, false
}
```

`Snake` has logic on how to move a snake and has few functions to check for collisions.

**snake/snake.go**

```go
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
```

Now we can combine these types in our new `Board` struct that will define a game board.

`Board` has an `Update` function that will be called from our `Game.Update()` method. Ebiten runs `Update()` aroudn 60/s that's why we need some intervals logic there to prevent snake moving too fast. For example we start with snake moving at move/200ms, and when it grows it will speed up to move/100s.

`Board` struct also holds information such as `points` and `gameOver` so in case snake hits the wall or itself, we can later just render a message.

**snake/board.go**

```go
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
```

What's left is to make few updates to our `game.go` file to initialize the `Board` and render the snake/food at each current state .

**snake/game.go**

```go
// ...

var (
	backgroundColor = color.RGBA{50, 100, 50, 50}
	snakeColor      = color.RGBA{200, 50, 150, 150}
	foodColor       = color.RGBA{200, 200, 50, 150}
)

type Game struct {
	input *Input
	board *Board
}

func NewGame() *Game {
	return &Game{
		input: NewInput(),
		board: NewBoard(boardRows, boardCols),
	}
}

func (g *Game) Update() error {
	return g.board.Update(g.input)
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)
	if g.board.gameOver {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("Game Over. Score: %d", g.board.points))
	} else {
		width := ScreenHeight / boardRows

		for _, p := range g.board.snake.body {
			ebitenutil.DrawRect(screen, float64(p.y*width), float64(p.x*width), float64(width), float64(width), snakeColor)
		}
		if g.board.food != nil {
			ebitenutil.DrawRect(screen, float64(g.board.food.y*width), float64(g.board.food.x*width), float64(width), float64(width), foodColor)
		}
		ebitenutil.DebugPrint(screen, fmt.Sprintf("Score: %d", g.board.points))
	}
}
```

### Playing the game

Repeat `wasmserve .` or just refresh the page and you'll be able to play the game.

![screenshot.png](https://raw.githubusercontent.com/plutov/packagemain/master/25-ebiten-snake/screenshot.png)