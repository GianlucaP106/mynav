package tui

import (
	"math"
	"strings"
)

type FrameType = []rune

var (
	ThickFrame FrameType = []rune{'═', '║', '╔', '╗', '╚', '╝', '╠', '╣', '╦', '╩', '╬'}
	ThinFrame  FrameType = []rune{'─', '│', '╭', '╮', '╰', '╯', '├', '┤', '┬', '┴', '┼'}
)

type Alignment uint

const (
	LeftAlign Alignment = iota
	CenterAlign
)

func Pad(content string, size int) string {
	return withCharPadding(content, size, " ")
}

func TrimEnd(s string, n int) string {
	if n >= len(s) {
		return ""
	}
	return s[:len(s)-n]
}

func TrimStart(s string, n int) string {
	if n >= len(s) {
		return ""
	}
	return s[n:]
}

func withCharPadding(content string, size int, c string) string {
	repeat := size - len(content)
	if repeat <= 0 {
		repeat = int(math.Abs(float64(repeat)))
		return TrimEnd(content, repeat+3) + "..."
	}
	return content + strings.Repeat(c, repeat)
}
