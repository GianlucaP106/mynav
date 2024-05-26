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
)

func displayLine(content string, alignment Alignment, maxWidth int, color color.Style) string {
	var line string
	switch alignment {
	case Left:
		line = color.Sprint(withSpacePadding(content, maxWidth))
	case Center:
		line = color.Sprint(blankLine((maxWidth*2)/5)) + color.Sprint(withSpacePadding(content, (maxWidth*3)/5))
	}
	return line
}

func display(content string, alignment Alignment, maxWidth int) string {
	var line string
	switch alignment {
	case Left:
		line = withSpacePadding(content, maxWidth)
	case Center:
		line = blankLine((maxWidth/2)-(len(content)/2)) + content
	}
	return line
}

func trimEnd(s string, n int) string {
	if n >= len(s) {
		return ""
	}
	return s[:len(s)-n]
}

func blankLine(size int) string {
	return withSpacePadding("", size)
}

func highlightedBlankLine(size int) string {
	white := color.New(color.White, color.BgCyan)
	return displayLine("", Center, size, white)
}

func displayWhiteText(content string, alignment Alignment, size int) string {
	white := color.New(color.White)
	return displayLine(content, alignment, size, white)
}

func withSpacePadding(content string, size int) string {
	return withCharPadding(content, size, " ")
}

func withCharPadding(content string, size int, c string) string {
	repeat := size - len(content)
	if repeat <= 0 {
		repeat = int(math.Abs(float64(repeat)))
		return trimEnd(content, repeat+4) + "... "
	}
	return content + strings.Repeat(c, repeat)
}

func withSurroundingSpaces(s string) string {
	return " " + s + " "
}
