package ui

import (
	"fmt"
	"mynav/pkg/tui"

	"github.com/awesome-gocui/gocui"
)

type toastDialog struct {
	view *tui.View
}

type toastDialogType uint

const (
	toastDialogErrType toastDialogType = iota
	toastDialogSuccessType
	toastDialogNeutralType
)

func openToastDialog(message string, dialogType toastDialogType, title string, exit func()) *toastDialog {
	td := &toastDialog{}

	td.view = tui.SetCenteredView(ToastDialog, max(len(message), len(title))+5, 3, 0)
	td.view.Title = tui.WithSurroundingSpaces(title)
	td.view.Editable = true
	styleView(td.view)
	switch dialogType {
	case toastDialogErrType:
		td.view.FrameColor = gocui.ColorRed
	case toastDialogSuccessType:
		td.view.FrameColor = gocui.ColorGreen
	default:
		td.view.FrameColor = onFrameColor
	}

	keys := []gocui.Key{
		gocui.KeyEnter,
		gocui.KeyBackspace,
		gocui.KeyBackspace2,
		gocui.KeyEsc,
	}

	prevView := tui.GetFocusedView()
	td.view.Editor = tui.NewSingleActionEditor(keys, func() {
		td.close()
		if prevView != nil {
			prevView.Focus()
		}
		exit()
	})

	sizeX, _ := td.view.Size()
	fmt.Fprintln(td.view, tui.DisplayWhite(message, tui.LeftAlign, sizeX))

	td.view.Focus()

	return td
}

func openToastDialogError(message string) *toastDialog {
	return openToastDialog(message, toastDialogErrType, "Error", func() {})
}

func (td *toastDialog) close() {
	td.view.Delete()
}
