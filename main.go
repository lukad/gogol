package main

import (
	"github.com/nsf/termbox-go"
	flag "github.com/ogier/pflag"
)

var (
	toroidal = flag.BoolP("toroidal", "t", false, "Use a toroidal array")
)

var (
	lastState []bool
	curState  []bool
	width     int
	height    int
)

func draw(board []bool) {
	termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if get(board, x, y) {
				termbox.SetCell(x, y, '#', termbox.ColorWhite, termbox.ColorBlack)
			}
		}
	}
	termbox.Flush()
}

func get(board []bool, x, y int) bool {
	if *toroidal {
		if x < 0 {
			x = width - 1
		} else if x >= width {
			x = 0
		}
		if y < 0 {
			y = height - 1
		} else if y >= height {
			y = 0
		}
	} else if x < 0 || x >= width || y < 0 || y >= height {
		return false
	}
	return board[y*width+x]
}

func set(board []bool, x, y int, v bool) {
	if *toroidal {
		if x < 0 {
			x = width - 1
		} else if x >= width {
			x = 0
		}
		if y < 0 {
			y = height - 1
		} else if y >= height {
			y = 0
		}
	}
	board[y*width+x] = v
}

func neighborCount(board []bool, x, y int) (count int) {
	if get(board, x-1, y) { // left
		count++
	}
	if get(board, x-1, y-1) { // top left
		count++
	}
	if get(board, x, y-1) { // top
		count++
	}
	if get(board, x+1, y-1) { // top right
		count++
	}
	if get(board, x+1, y) { // right
		count++
	}
	if get(board, x+1, y+1) { // bottom right
		count++
	}
	if get(board, x, y+1) { // bottom
		count++
	}
	if get(board, x-1, y+1) { // bottom left
		count++
	}
	return count
}

// Goes through one iteration and copies the current state into last
func step(last, current []bool) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			set(current, x, y, get(last, x, y) && neighborCount(last, x, y) == 2 || neighborCount(last, x, y) == 3)
		}
	}
	copy(last, current)
}

func loop() {
	for {
		switch event := termbox.PollEvent(); event.Type {
		case termbox.EventKey:
			switch event.Key {
			case termbox.KeyEsc, termbox.KeyCtrlC, termbox.KeyCtrlD:
				return
			case termbox.KeySpace:
				step(lastState, curState)
				draw(curState)
			}
		}
	}
	termbox.SetCursor(5, 4)
}

func init() {
	flag.Parse()
}

func main() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	// Setup slices
	width, height = termbox.Size()
	lastState = make([]bool, width*height)
	curState = make([]bool, width*height)

	// Glider
	set(lastState, 5, 5, true)
	set(lastState, 6, 5, true)
	set(lastState, 7, 5, true)
	set(lastState, 7, 4, true)
	set(lastState, 6, 3, true)

	// Blinker
	set(lastState, 20, 4, true)
	set(lastState, 20, 5, true)
	set(lastState, 20, 6, true)

	copy(curState, lastState)
	draw(curState)

	loop()
}
