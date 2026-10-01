package main

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

type Point struct {
	x, y int
}

type Game struct {
	snake         []Point
	food          Point
	malware       []Point
	dir           Point
	score         int
	level         int
	gameOver      bool
	width, height int
	quit          chan struct{}
}

func NewGame(width, height int) *Game {
	return &Game{
		snake:   []Point{{x: width / 2, y: height / 2}},
		malware: make([]Point, 0),
		dir:     Point{x: 1, y: 0},
		level:   1,
		width:   width,
		height:  height,
		quit:    make(chan struct{}),
	}
}

func (g *Game) draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	//score and level
	for i, ch := range fmt.Sprintf("Score: %d Level: %d", g.score, g.level) {
		termbox.SetCell(i+5, 0, ch, termbox.ColorWhite, termbox.ColorDefault)
	}

	//border
	termbox.SetCell(0, 1, '┌', termbox.ColorWhite, termbox.ColorDefault)
	for col := 1; col < g.width-1; col++ {
		termbox.SetCell(col, 1, '─', termbox.ColorWhite, termbox.ColorDefault)
	}

	termbox.SetCell(g.width-1, 1, '┐', termbox.ColorWhite, termbox.ColorDefault)
	for row := 1; row < g.height-1; row++ {
		termbox.SetCell(g.width-1, row+1, '│', termbox.ColorWhite, termbox.ColorDefault)
	}

	termbox.SetCell(g.width-1, g.height, '┘', termbox.ColorWhite, termbox.ColorDefault)
	for col := 1; col < g.width-1; col++ {
		termbox.SetCell(col, g.height, '─', termbox.ColorWhite, termbox.ColorDefault)
	}

	termbox.SetCell(0, g.height, '└', termbox.ColorWhite, termbox.ColorDefault)
	for row := 1; row < g.height-1; row++ {
		termbox.SetCell(0, row+1, '│', termbox.ColorWhite, termbox.ColorDefault)
	}

	//snake
	termbox.SetCell(g.snake[0].x, g.snake[0].y, 'o', termbox.ColorWhite, termbox.ColorDefault)

	termbox.Flush()
}

func main() {
	err := termbox.Init()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer termbox.Close()

	ng := NewGame(40, 20)
	ng.draw()
}
