// Copyright (c) 2017 Alex Pliutau

package snake

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strconv"
	"time"
)

var (
	topScoreFile = "/tmp/snake.score"
	topScoreChan chan int
	topScoreVal  int
)

func init() {
	// Get top score from the file
	line, readErr := ioutil.ReadFile(topScoreFile)
	if readErr == nil {
		var castErr error
		topScoreVal, castErr = strconv.Atoi(string(line))
		if castErr != nil {
			log.Printf("can't cast score: %v", castErr)
		}
	}

	topScoreChan = make(chan int)
	go func() {
		for {
			s := <-topScoreChan
			if s > topScoreVal {
				topScoreVal = s
				ioutil.WriteFile(topScoreFile, []byte(fmt.Sprintf("%d", topScoreVal)), 0777)
			}
		}
	}()
}

// Game type
type Game struct {
	KeyboardEventsChan chan KeyboardEvent
	PointsChan         chan int
	arena              *arena
	score              int
	IsOver             bool
}

// NewGame returns Game obj
func NewGame() *Game {
	return &Game{
		arena: initialArena(),
		score: initialScore(),
	}
}

// Start game func
func (g *Game) Start() {
	g.KeyboardEventsChan = make(chan KeyboardEvent)
	g.PointsChan = make(chan int)
	g.arena.pointsChan = g.PointsChan

	for {
		select {
		case p := <-g.PointsChan:
			g.addPoints(p)
			topScoreChan <- g.score
		case e := <-g.KeyboardEventsChan:
			d := keyToDirection(e.Key)
			if d > 0 {
				g.arena.snake.changeDirection(d)
			}
		default:
			if g.IsOver {
				log.Printf("Game over, score: %d\n", g.score)
				return
			}

			if err := g.arena.moveSnake(); err != nil {
				g.IsOver = true
			}

			time.Sleep(g.moveInterval())
		}
	}
}

func initialSnake() *snake {
	return newSnake(RIGHT, []coord{
		coord{x: 1, y: 1},
		coord{x: 1, y: 2},
		coord{x: 1, y: 3},
		coord{x: 1, y: 4},
	})
}

func initialScore() int {
	return 0
}

func initialArena() *arena {
	return newArena(initialSnake(), 20, 20)
}

func (g *Game) moveInterval() time.Duration {
	ms := 400 - math.Max(float64(g.score), 100)
	return time.Duration(ms) * time.Millisecond
}

func (g *Game) addPoints(p int) {
	g.score += p
}
