package ui

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
)

type ToastDialogState struct {
	viewName string
	message  string
	active   bool
}

func newToastDialogState() *ToastDialogState {
	return &ToastDialogState{
		viewName: "ToastDialogView",
		active:   false,
	}
}

func (ui *UI) initToastDialogView(sizeX int, sizeY int) *gocui.View {
	view := ui.setCenteredView(ui.toast.viewName, sizeX, sizeY, 0)
	view.FrameColor = gocui.ColorRed
	view.Title = withSurroundingSpaces("Error")
	for _, key := range []gocui.Key{
		gocui.KeyEnter,
		gocui.KeyBackspace,
		gocui.KeyBackspace2,
		gocui.KeyEsc,
	} {
		ui.keyBinding(ui.toast.viewName).set(key, func() {
			ui.closeToastDialog()
		})
	}

	return view
}

func (ui *UI) openToastDialog(message string) {
	ui.toast.active = true
	ui.toast.message = message
}

func (ui *UI) closeToastDialog() {
	ui.toast.active = false
	ui.gui.DeleteView(ui.toast.viewName)
}

func (ui *UI) getToastDialogView() *gocui.View {
	return ui.getView(ui.toast.viewName)
}

func (ui *UI) renderToastDialog() {
	if !ui.toast.active {
		return
	}

	messageLength := len(ui.toast.message)
	view := ui.initToastDialogView(messageLength+5, 3)

	ui.focusView(ui.toast.viewName)

	sizeX, _ := view.Size()
	fmt.Fprintln(view, displayLineNormal(ui.toast.message, Left, sizeX))
}
