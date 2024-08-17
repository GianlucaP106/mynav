package ui

import (
	"mynav/pkg/tui"

	"github.com/awesome-gocui/gocui"
)

const (
	offFrameColor = gocui.AttrDim | gocui.ColorWhite
	onFrameColor  = gocui.ColorWhite
)

func styleView(v *tui.View) {
	v.FrameRunes = tui.ThickFrame
	v.TitleColor = gocui.AttrBold | gocui.ColorYellow
}
