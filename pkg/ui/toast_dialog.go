package ui

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
)

type ToastDialog struct {
	view *View
}

const ToastDialogName = "ToastDialogView"

func OpenToastDialog(message string, error bool, title string, exit func()) *ToastDialog {
	td := &ToastDialog{}

	td.view = SetCenteredView(ToastDialogName, max(len(message), len(title))+5, 3, 0)
	td.view.Title = withSurroundingSpaces(title)
	td.view.Editable = true
	if error {
		td.view.FrameColor = gocui.ColorRed
	} else {
		td.view.FrameColor = gocui.ColorGreen
	}

	keys := []gocui.Key{
		gocui.KeyEnter,
		gocui.KeyBackspace,
		gocui.KeyBackspace2,
		gocui.KeyEsc,
	}

	prevView := GetFocusedView()
	td.view.Editor = NewSingleActionEditor(keys, func() {
		td.Close()
		if prevView != nil {
			SetCurrentView(prevView.Name())
		}
		exit()
	})

	sizeX, _ := td.view.Size()
	fmt.Fprintln(td.view, displayWhiteText(message, Left, sizeX))

	SetCurrentView(td.view.Name())

	return td
}

func OpenToastDialogError(message string) *ToastDialog {
	return OpenToastDialog(message, true, "Error", func() {})
}

func (td *ToastDialog) Close() {
	DeleteView(td.view.Name())
}
