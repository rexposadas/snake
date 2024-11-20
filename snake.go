package main

type position [2]int
type direction int

type snake struct {
	body      []position
	direction direction
}

func newSnake(cfg config) *snake {
	maxX := cfg.RightBorder
	maxY := cfg.BottomBorder

	pos := position{maxX / 2, maxY / 2}

	return &snake{
		body:      []position{pos},
		direction: north,
	}
}
