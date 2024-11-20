package main

import (
	"github.com/mattn/go-tty"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"time"
)

type config struct {
	LeftBorder, TopBorder     int
	RightBorder, BottomBorder int
	SnakeSpeed                time.Duration
}

type game struct {
	score  int
	snake  *snake
	food   position
	config config
}

func newGame() *game {
	rand.Seed(time.Now().UnixNano())
	cfg := config{
		LeftBorder:   5,
		RightBorder:  40,
		TopBorder:    5,
		BottomBorder: 30,

		SnakeSpeed: time.Millisecond * 130,
	}

	snake := newSnake(cfg)

	game := &game{
		score:  0,
		snake:  snake,
		food:   position{10, 10},
		config: cfg,
	}

	go game.listenForKeyPress()

	return game
}

func (g *game) drawBorder() {

	// top & bottom
	for i := g.config.LeftBorder; i <= g.config.RightBorder; i++ {
		moveCursor(position{i, g.config.TopBorder})
		draw("#")
	}

	for i := g.config.LeftBorder; i <= g.config.RightBorder; i++ {
		moveCursor(position{i, g.config.BottomBorder})
		draw("#")
	}

	// left & right
	for j := g.config.TopBorder; j <= g.config.BottomBorder; j++ {
		moveCursor(position{g.config.LeftBorder, j})
		draw("#")
	}

	for j := g.config.TopBorder; j <= g.config.BottomBorder; j++ {
		moveCursor(position{g.config.RightBorder, j})
		draw("#")
	}
}

func (g *game) draw() {
	clear()
	//maxX, maxY := getSize()

	g.drawBorder()

	status := "score " + strconv.Itoa(g.score)
	statusXPos := g.config.RightBorder - len(status)/2

	moveCursor(position{statusXPos, 0})
	draw(status)

	moveCursor(g.food)
	draw("*")

	for i, pos := range g.snake.body {
		moveCursor(pos)

		if i == 0 {
			draw("0")
		} else {
			draw("o")
		}
	}

	render()
	// Snake moves on each render
	time.Sleep(g.config.SnakeSpeed)
}

func (g *game) over() {
	clear()
	showCursor()

	moveCursor(position{1, 1})
	draw("game over. score: " + strconv.Itoa(g.score))

	render()

	os.Exit(0)
}

func (g *game) beforeGame() {
	hideCursor()

	// handle CTRL C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			g.over()
		}
	}()
}

func (g *game) placeNewFood() {
	for {
		newFoodPosition := randomPosition(g)
		if positionAreSame(newFoodPosition, g.food) {
			continue
		}

		for _, pos := range g.snake.body {
			if positionAreSame(newFoodPosition, pos) {
				continue
			}
		}

		g.food = newFoodPosition

		break
	}
}

func (g *game) listenForKeyPress() {
	terminalAPI, err := tty.Open()
	if err != nil {
		panic(err)
	}
	defer terminalAPI.Close()

	for {
		char, err := terminalAPI.ReadRune()
		if err != nil {
			panic(err)
		}

		switch char {
		case 'A':
			g.snake.direction = north
		case 'B':
			g.snake.direction = south
		case 'C':
			g.snake.direction = east
		case 'D':
			g.snake.direction = west
		}
	}
}
