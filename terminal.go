package main

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"os"
)

type winsize struct {
	Rows 	uint16
	Cols 	uint16
	Width 	uint16
	Height 	uint16
}

type Terminal struct {
	Width	int
	Height 	int
	Ratio	float64
}

const defaultRatio float64 = 1.0 // The terminal's default cursor width/height ratio

// Screen buffer
var screen *bytes.Buffer = new(bytes.Buffer)
var output *bufio.Writer = bufio.NewWriter(os.Stdout)
var Window *Terminal = getTerminal()

func init() {
	// Clear console
	output.WriteString("\033[2J")
	// Remove blinking cursor
	output.WriteString("\033[?25l")
}

// Get terminal size
func getTerminal() (*Terminal) {
	var whRatio float64
	ws, err := getWinsize()
	if err != nil {
		panic(err)
	}
	whRatio = defaultRatio
	if ws.Width > 0 && ws.Height > 0 {
		whRatio = float64(ws.Height / ws.Rows) / float64(ws.Width / ws.Cols) * 0.5
	}
	return &Terminal{
		Width 	: int(ws.Cols),
		Height 	: int(ws.Rows),
		Ratio	: whRatio,
	}
}

// Flush buffer and ensure that it will not overflow screen
func (terminal *Terminal) Flush() {
	for idx, str := range strings.Split(screen.String(), "\n") {
		if idx > Window.Height {
			return
		}
		output.WriteString(str + "\n")
	}

	screen.Reset()
	output.Flush()
	terminal.MoveCursor(0,0)
}

// Move cursor to given position
func (terminal *Terminal) MoveCursor(x int, y int) {
	fmt.Fprintf(screen, "\033[%d;%dH", x, y)
}