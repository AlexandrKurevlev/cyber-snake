package main

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

func NewGame() *Game {
	return &Game{
		snake:   make([]Point, 0),
		malware: make([]Point, 0),
		quit:    make(chan struct{}),
	}
}

func main() {}
