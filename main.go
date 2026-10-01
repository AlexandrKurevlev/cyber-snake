package main

import (
	"fmt"
	"math/rand"
	"time"

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
	ng := &Game{
		snake:   []Point{{x: width / 2, y: height / 2}},
		malware: make([]Point, 0),
		dir:     Point{x: 1, y: 0},
		level:   1,
		width:   width,
		height:  height,
		quit:    make(chan struct{}),
	}

	ng.placeFood()
	ng.placeMalware()

	return ng
}

func (g *Game) draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	//border
	termbox.SetCell(0, 0, '┌', termbox.ColorWhite, termbox.ColorDefault)
	for col := 1; col < g.width-1; col++ {
		termbox.SetCell(col, 0, '─', termbox.ColorWhite, termbox.ColorDefault)
	}

	termbox.SetCell(g.width-1, 0, '┐', termbox.ColorWhite, termbox.ColorDefault)
	for row := 1; row < g.height-1; row++ {
		termbox.SetCell(g.width-1, row, '│', termbox.ColorWhite, termbox.ColorDefault)
	}

	termbox.SetCell(g.width-1, g.height-1, '┘', termbox.ColorWhite, termbox.ColorDefault)
	for col := 1; col < g.width-1; col++ {
		termbox.SetCell(col, g.height-1, '─', termbox.ColorWhite, termbox.ColorDefault)
	}

	termbox.SetCell(0, g.height-1, '└', termbox.ColorWhite, termbox.ColorDefault)
	for row := 1; row < g.height-1; row++ {
		termbox.SetCell(0, row, '│', termbox.ColorWhite, termbox.ColorDefault)
	}

	//snake head
	termbox.SetCell(g.snake[0].x, g.snake[0].y, g.dir.ToRune(), termbox.ColorWhite, termbox.ColorDefault)

	//snake body
	for i := 1; i < len(g.snake); i++ {
		termbox.SetCell(g.snake[i].x, g.snake[i].y, '○', termbox.ColorWhite, termbox.ColorDefault)
	}

	//score and level
	for i, ch := range fmt.Sprintf("Score: %d Level: %d", g.score, g.level) {
		termbox.SetCell(i+5, g.height, ch, termbox.ColorWhite, termbox.ColorDefault)
	}

	//food
	termbox.SetCell(g.food.x, g.food.y, '●', termbox.ColorGreen, termbox.ColorDefault)

	//malware
	for _, mp := range g.malware {
		termbox.SetCell(mp.x, mp.y, '✗', termbox.ColorRed, termbox.ColorDefault)
	}

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
	return !(p.x >= 1 && p.x <= g.width-2 && p.y >= 1 && p.y <= g.height-2)
}

func (g *Game) isOnFood(p Point) bool {
	return p.x == g.food.x && p.y == g.food.y
}

func (g *Game) placeFood() {
	for {
		fp := Point{x: rand.Intn(g.width-3) + 1, y: rand.Intn(g.height-3) + 1}
		if !g.isOnSnake(fp) && !g.isOnMalware(fp) {
			g.food = fp
			break
		}
	}
}

func (g *Game) placeMalware() {
	for {
		mp := Point{x: rand.Intn(g.width-3) + 1, y: rand.Intn(g.height-3) + 1}
		if !g.isOnSnake(mp) && !g.isOnMalware(mp) && !g.isOnFood(mp) {
			g.malware = append(g.malware, mp)
			break
		}
	}
}

func (g *Game) move() {
	newHead := Point{x: g.snake[0].x + g.dir.x, y: g.snake[0].y + g.dir.y}
	if g.isOutOfBounds(newHead) || g.isOnSnake(newHead) || g.isOnMalware(newHead) {
		g.gameOver = true
		return
	}

	if g.isOnFood(newHead) {
		g.score++
		g.snake = append([]Point{newHead}, g.snake...)
		g.placeFood()
	} else {
		g.snake = append([]Point{newHead}, g.snake[:len(g.snake)-1]...)
	}
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

	ticker := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case ev := <-eventCh:
			ng.handleInput(ev)
		case <-ticker.C:
			if !ng.gameOver {
				ng.move()
				ng.draw()
			}
		case <-ng.quit:
			return
		}
	}
}
