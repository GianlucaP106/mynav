package ui

import (
	"fmt"
)

type ConfirmationDialog struct {
	view  *View
	title string
}

const ConfirmationDialogName = "ConfirmationDialog"

func OpenConfirmationDialog(onConfirm func(bool), title string) *ConfirmationDialog {
	cd := &ConfirmationDialog{}
	cd.title = title
	prevView := GetFocusedView()
	cd.view = SetCenteredView(ConfirmationDialogName, len(title)+5, 3, 0)
	cd.view.Title = withSurroundingSpaces("Confirm")
	cd.view.Wrap = true
	cd.view.Editable = true
	cd.view.Editor = NewConfirmationEditor(func() {
		cd.Close()
		if prevView != nil {
			prevView.Focus()
		}
		onConfirm(true)
	}, func() {
		cd.Close()
		if prevView != nil {
			prevView.Focus()
		}
		onConfirm(false)
	})

	sizeX, _ := cd.view.Size()
	prevView.Focus()
	cd.view.Clear()
	fmt.Fprintln(cd.view, displayWhiteText(cd.title, Left, sizeX))
	return cd
}

func (cd *ConfirmationDialog) Close() {
	cd.view.Delete()
}
