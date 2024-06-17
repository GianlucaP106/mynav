package ui

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
)

const ToastDialogStateName = "ToastDialogView"

type ToastDialog struct {
	editor  Editor
	message string
}

var _ Dialog = &ToastDialog{}

func newToastDialogState() *ToastDialog {
	return &ToastDialog{}
}

func (td *ToastDialog) Name() string {
	return ToastDialogStateName
}

func (td *ToastDialog) Open(message string, exit func()) {
	td.message = message
	messageLength := len(td.message)
	view := SetCenteredView(td.Name(), messageLength+5, 3, 0)
	view.FrameColor = gocui.ColorRed
	view.Title = withSurroundingSpaces("Error")
	view.Editable = true
	keys := []gocui.Key{
		gocui.KeyEnter,
		gocui.KeyBackspace,
		gocui.KeyBackspace2,
		gocui.KeyEsc,
	}

	prevView := GetFocusedView()
	td.editor = NewSingleActionEditor(keys, func() {
		td.Close()
		if prevView != nil {
			FocusView(prevView.Name())
		}
		exit()
	})
	view.Editor = td.editor
	FocusView(td.Name())
}

func (td *ToastDialog) Close() {
	td.message = ""
	DeleteView(td.Name())
}

func (td *ToastDialog) Render(ui *UI) error {
	view := GetInternalView(td.Name())
	if view == nil {
		return nil
	}

	sizeX, _ := view.Size()
	fmt.Fprintln(view, displayWhiteText(td.message, Left, sizeX))
	return nil
}
