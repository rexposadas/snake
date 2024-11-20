package main

import "math/rand"

func main() {
	game := newGame()
	game.beforeGame()

	for {
		maxX := game.config.RightBorder
		maxY := game.config.BottomBorder

		newHeadPos := game.snake.body[0]

		switch game.snake.direction {
		case north:
			newHeadPos[1]--
		case east:
			newHeadPos[0]++
		case south:
			newHeadPos[1]++
		case west:
			newHeadPos[0]--
		}

		hitWall := newHeadPos[0] <= game.config.LeftBorder || newHeadPos[1] <= game.config.TopBorder ||
			newHeadPos[0] >= maxX || newHeadPos[1] >= maxY
		if hitWall {
			game.over()
		}

		// Running into yourself, game over
		for _, pos := range game.snake.body {
			if positionAreSame(newHeadPos, pos) {
				game.over()
			}
		}

		// add head to body
		game.snake.body = append([]position{newHeadPos}, game.snake.body...)

		ateFood := positionAreSame(game.food, newHeadPos)
		if ateFood {
			game.score++
			game.placeNewFood()
		} else {
			game.snake.body = game.snake.body[:len(game.snake.body)-1]
		}

		game.draw()
	}
}

func positionAreSame(a, b position) bool {
	return a[0] == b[0] && a[1] == b[1]
}

func randomPosition(g *game) position {
	// create a border within the game border where the food can spawn

	leftBuffer := g.config.LeftBorder + 3
	rightBuffer := g.config.RightBorder - leftBuffer - 3

	topBuffer := g.config.TopBorder + 3
	bottomBuffer := g.config.BottomBorder - topBuffer - 3

	x := leftBuffer + rand.Intn(rightBuffer)
	y := topBuffer + rand.Intn(bottomBuffer)

	return position{x, y}
}
