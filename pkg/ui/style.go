package ui

import (
	"math"
	"strings"

	"github.com/gookit/color"
)

type FrameType = []rune

var ThickFrame FrameType = []rune{'═', '║', '╔', '╗', '╚', '╝'}

type Alignment uint

const (
	Left Alignment = iota
	Center
	Right
)

func displayLine(content string, alignment Alignment, maxWidth int, color color.Style) string {
	numSpaces := maxWidth - len(content)
	if numSpaces <= 0 {
		numSpaces = int(math.Abs(float64(numSpaces)))
		return color.Sprint(" " + removeLastNChars(content, numSpaces+4) + "...")
	}
	numSpaces = numSpaces/2 + 1
	spaces := strings.Repeat(color.Sprint(" "), numSpaces)
	content = color.Sprint(" " + content)
	var line string
	switch alignment {
	case Left:
		line = content + spaces + spaces
	case Center:
		line = spaces + content + spaces
	case Right:
		line = spaces + spaces + content
	}
	return line
}

func removeLastNChars(s string, n int) string {
	if n >= len(s) {
		return ""
	}
	return s[:len(s)-n]
}

func blankLine(size int) string {
	white := color.New(color.White)
	return displayLine("", Center, size, white)
}

func highlightedBlankLine(size int) string {
	white := color.New(color.White, color.BgCyan)
	return displayLine("", Center, size, white)
}

func displayLineNormal(content string, alignment Alignment, size int) string {
	white := color.New(color.White)
	return displayLine(content, alignment, size, white)
}

func withSpacePadding(content string, size int) string {
	repeat := size - len(content)
	if repeat <= 0 {
		repeat = int(math.Abs(float64(repeat)))
		return color.Sprint(removeLastNChars(content, repeat+3) + "...")
	}
	return content + strings.Repeat(" ", size-len(content))
}

func withSurroundingSpaces(s string) string {
	return " " + s + " "
}
