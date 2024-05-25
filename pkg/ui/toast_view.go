package ui

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
)

const ToastDialogStateName = "ToastDialogView"

type ToastDialogState struct {
	editor  Editor
	message string
}

var _ Dialog = &ToastDialogState{}

func newToastDialogState() *ToastDialogState {
	return &ToastDialogState{}
}

func (td *ToastDialogState) Name() string {
	return ToastDialogStateName
}

func (td *ToastDialogState) Open(message string, exit func()) {
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

	td.editor = NewSingleActionEditor(keys, func() {
		td.Close()
		exit()
	})
	view.Editor = td.editor
	FocusView(td.Name())
}

func (td *ToastDialogState) Close() {
	td.message = ""
	DeleteView(td.Name())
}

func (td *ToastDialogState) Render(ui *UI) error {
	view := GetInternalView(td.Name())
	if view == nil {
		return nil
	}

	sizeX, _ := view.Size()
	fmt.Fprintln(view, displayWhiteText(td.message, Left, sizeX))
	return nil
}
