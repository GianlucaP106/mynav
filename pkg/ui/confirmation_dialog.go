package ui

import (
	"fmt"
)

type ConfirmationDialog struct {
	editor Editor
	title  string
}

var _ Dialog = &ConfirmationDialog{}

const ConfirmationDialogStateName = "ConfirmationDialog"

func newConfirmationDialogState() *ConfirmationDialog {
	return &ConfirmationDialog{}
}

func (cd *ConfirmationDialog) Open(onConfirm func(bool), title string) {
	cd.title = title
	cd.editor = NewConfirmationEditor(func() {
		cd.Close()
		onConfirm(true)
	}, func() {
		cd.Close()
		onConfirm(false)
	})

	sizeX := len(cd.title)
	view := SetCenteredView(cd.Name(), sizeX+5, 3, 0)
	FocusView(cd.Name())
	view.Title = withSurroundingSpaces("Confirm")
	view.Wrap = true
	view.Editor = cd.editor
	view.Editable = true
}

func (cd *ConfirmationDialog) Name() string {
	return ConfirmationDialogStateName
}

func (cd *ConfirmationDialog) Close() {
	DeleteView(cd.Name())
}

func (cd *ConfirmationDialog) Render(ui *UI) error {
	view := GetInternalView(cd.Name())
	if view == nil {
		return nil
	}

	sizeX, _ := view.Size()
	fmt.Fprintln(view, displayWhiteText(cd.title, Left, sizeX))
	return nil
}
