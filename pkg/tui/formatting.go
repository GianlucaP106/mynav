package tui

import (
	"math"
	"strings"

	"github.com/awesome-gocui/gocui"
	"github.com/gookit/color"
)

type FrameType = []rune

var ThickFrame FrameType = []rune{'═', '║', '╔', '╗', '╚', '╝', '╠', '╣', '╦', '╩', '╬'}

type Alignment uint

const (
	LeftAlign Alignment = iota
	CenterAlign
)

const (
	OffFrameColor = gocui.AttrDim | gocui.ColorWhite
	OnFrameColor  = gocui.ColorWhite
)

func StyleView(v *View) {
	v.FrameRunes = ThickFrame
	v.TitleColor = gocui.AttrBold | gocui.ColorYellow
}

func DisplayColored(content string, alignment Alignment, maxWidth int, color color.Style) string {
	var line string
	switch alignment {
	case LeftAlign:
		line = color.Sprint(WithSpaces(content, maxWidth))
	case CenterAlign:
		line = color.Sprint(BlankLine((maxWidth*2)/5)) + color.Sprint(WithSpaces(content, (maxWidth*3)/5))
	}
	return line
}

func Display(content string, alignment Alignment, maxWidth int) string {
	var line string
	switch alignment {
	case LeftAlign:
		line = WithSpaces(content, maxWidth)
	case CenterAlign:
		line = BlankLine((maxWidth/2)-(len(content)/2)) + content
	}
	return line
}

func BlankLine(size int) string {
	return WithSpaces("", size)
}

func DisplayWhite(content string, alignment Alignment, size int) string {
	white := color.New(color.White)
	return DisplayColored(content, alignment, size, white)
}

func WithSpaces(content string, size int) string {
	return withCharPadding(content, size, " ")
}

func WithSurroundingSpaces(s string) string {
	return " " + s + " "
}

func trimEnd(s string, n int) string {
	if n >= len(s) {
		return ""
	}
	return s[:len(s)-n]
}

func withCharPadding(content string, size int, c string) string {
	repeat := size - len(content)
	if repeat <= 0 {
		repeat = int(math.Abs(float64(repeat)))
		return trimEnd(content, repeat+4) + "... "
	}
	return content + strings.Repeat(c, repeat)
}
