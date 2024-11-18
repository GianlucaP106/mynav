package app

import (
	"fmt"
	"mynav/pkg/tui"
	"strconv"
	"time"

	"github.com/awesome-gocui/gocui"
)

// Toast is a view that shows a message at the corner of the screen for a period of time.
type Toast struct {
	view *tui.View
}

// ToastType dictates what style the toast will be.
type ToastType uint

const (
	toastInfo ToastType = iota
	toastError
	toastWarn
)

// toastCount ensures that each toast trigger has a different name.
// If we use the same name, the system will update the dimensions of an existing view.
var toastCount = 0

func toast(msg string, typ ToastType) *Toast {
	tcount := strconv.Itoa(toastCount)
	toastCount++
	td := &Toast{}
	msg = " " + msg + "  "
	toastSize := len(msg)
	_, maxY := a.ui.Size()
	td.view = a.ui.SetView(
		tui.NewViewPosition(
			ToastDialog+tcount,
			2,
			maxY-4,
			2+toastSize,
			maxY-2,
			0,
		))

	td.view.FrameRunes = tui.ThinFrame
	fmt.Fprintln(td.view, msg)

	switch typ {
	case toastError:
		td.view.Title = "Error"
		td.view.TitleColor = gocui.ColorRed
		td.view.FrameColor = gocui.ColorRed
	case toastWarn:
		td.view.Title = "Warning"
		td.view.TitleColor = gocui.ColorYellow
		td.view.FrameColor = gocui.ColorYellow
	case toastInfo:
		td.view.Title = "Info"
		td.view.TitleColor = gocui.ColorGreen
		td.view.FrameColor = gocui.ColorGreen
	}

	time.AfterFunc(5*time.Second, func() {
		a.ui.Update(func() {
			a.ui.DeleteView(td.view)
		})
	})

	return td
}
