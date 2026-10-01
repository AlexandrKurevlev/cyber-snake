package main

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

type Point struct {
	x, y int
}

func (p Point) ToRune() rune {
	if p.x == 1 && p.y == 0 {
		return '▶'
	} else if p.x == 0 && p.y == 1 {
		return '▼'
	} else if p.x == -1 && p.y == 0 {
		return '◀'
	} else if p.x == 0 && p.y == -1 {
		return '▲'
	} else {
		return '●'
	}
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
	termbox.SetCell(g.snake[0].x, g.snake[0].y, g.dir.ToRune(), termbox.ColorWhite, termbox.ColorDefault)

	termbox.Flush()
}

func (g *Game) handleInput(ev termbox.Event) {
	if ev.Type != termbox.EventKey {
		return
	}

	if ev.Ch == 0 {
		switch ev.Key {
		case termbox.KeyArrowUp:
			g.handleInputUp()
		case termbox.KeyArrowRight:
			g.handleInputRight()
		case termbox.KeyArrowDown:
			g.handleInputDown()
		case termbox.KeyArrowLeft:
			g.handleInputLeft()
		case termbox.KeyEsc:
			close(g.quit)
		}
	} else {
		switch ev.Ch {
		case 'w':
			g.handleInputUp()
		case 'd':
			g.handleInputRight()
		case 's':
			g.handleInputDown()
		case 'a':
			g.handleInputLeft()
		case 'q':
			close(g.quit)
		}
	}
}

func (g *Game) handleInputUp() {
	if g.dir.y == 1 {
		return
	}

	g.dir.x = 0
	g.dir.y = -1
}

func (g *Game) handleInputRight() {
	if g.dir.x == -1 {
		return
	}

	g.dir.x = 1
	g.dir.y = 0
}

func (g *Game) handleInputDown() {
	if g.dir.y == -1 {
		return
	}

	g.dir.x = 0
	g.dir.y = 1
}

func (g *Game) handleInputLeft() {
	if g.dir.x == 1 {
		return
	}

	g.dir.x = -1
	g.dir.y = 0
}

func (g *Game) isOnSnake(p Point) bool {
	for _, sp := range g.snake {
		if p.x == sp.x && p.y == sp.y {
			return true
		}
	}

	return false
}

func (g *Game) isOnMalware(p Point) bool {
	for _, mp := range g.malware {
		if p.x == mp.x && p.y == mp.y {
			return true
		}
	}

	return false
}

func (g *Game) isOutOfBounds(p Point) bool {
	return !(p.x > 1 && p.x < g.width-2 && p.y > 1 && p.y < g.height-2)
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

	eventCh := make(chan termbox.Event)

	go func() {
		for {
			eventCh <- termbox.PollEvent()
		}
	}()

	for {
		select {
		case ev := <-eventCh:
			ng.handleInput(ev)
			ng.draw()
		case <-ng.quit:
			return
		}
	}
}
