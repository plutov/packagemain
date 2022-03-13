## packagemain #25: Snake game in Go using Ebiten

[Ebiten](https://ebiten.org/) is an open source game library in Go for building 2D games that can be ran across multiple platforms. Ebiten games work on desktop, web browsers (through WebAssembly), as well as on Mobile and even on Nintendo Switch.

In this video we'll give it a try and create a Snake game in Go, which we'll run in the browser using WebAssembly.

### Ebiten API

~[overview](https://ebiten.org/images/overview2.2.png)

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

As I mentioned before we'll run this game in the browser, as it's easy to compile it in WebAssembly. There are 2 options running this game in the browser.

1. We can just compile it manually and create our own HTML file
2. We can use [wasmserve](https://github.com/hajimehoshi/wasmserve) project that will compile a code and run it in the browser for us

For simplicity I'll use the second option.

```
go install github.com/hajimehoshi/wasmserve@latest
wasmserve .
```

Then if we open http://localhost:8080/ we should see a gray-ish square.