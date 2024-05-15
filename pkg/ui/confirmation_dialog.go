package ui

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
)

type ConfirmationDialogState struct {
	editor   Editor
	viewName string
	title    string
	active   bool
}

func newConfirmationDialogState() *ConfirmationDialogState {
	return &ConfirmationDialogState{
		viewName: "ConfirmationDialog",
	}
}

func (ui *UI) initConfirmationDialogView(sizeX int) *gocui.View {
	view := ui.setCenteredView(ui.confirmation.viewName, sizeX, 3, 0)
	view.Title = withSurroundingSpaces("Confirm")
	view.Wrap = true
	view.Editor = ui.confirmation.editor
	view.Editable = true
	return view
}

func (ui *UI) openConfirmationDialog(onConfirm func(bool), title string) {
	ui.confirmation.title = title
	ui.confirmation.editor = newConfirmationEditor(func() {
		ui.closeConfirmationDialog()
		onConfirm(true)
	}, func() {
		ui.closeConfirmationDialog()
		onConfirm(false)
	})
	ui.confirmation.active = true
}

func (ui *UI) closeConfirmationDialog() {
	ui.confirmation.active = false
	ui.gui.DeleteView(ui.confirmation.viewName)
}

func (ui *UI) getConfirmationDialogView() *gocui.View {
	return ui.getView(ui.confirmation.viewName)
}

func (ui *UI) renderConfirmationDialog() {
	if !ui.confirmation.active {
		return
	}
	titleSize := len(ui.confirmation.title)
	view := ui.initConfirmationDialogView(titleSize + 3)

	ui.focusView(ui.confirmation.viewName)

	sizeX, _ := view.Size()
	fmt.Fprintln(view, displayWhiteText(ui.confirmation.title, Left, sizeX))
}
