package main

import "fmt"

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

func main() {
	ng := NewGame(40, 20)
	fmt.Printf("Игра создана: поле %dx%d, змейка в (%d, %d), направление вправо, уровень %d", ng.width, ng.height, ng.snake[0].x, ng.snake[0].y, ng.level)
}
