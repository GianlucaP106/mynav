package ui

import (
	"fmt"
	"mynav/pkg/constants"
	"mynav/pkg/tui"

	"github.com/awesome-gocui/gocui"
)

type toastDialog struct {
	view *tui.View
}

func openToastDialog(message string, error bool, title string, exit func()) *toastDialog {
	td := &toastDialog{}

	td.view = tui.SetCenteredView(constants.ToastDialogName, max(len(message), len(title))+5, 3, 0)
	td.view.Title = tui.WithSurroundingSpaces(title)
	td.view.Editable = true
	tui.StyleView(td.view)
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
	return openToastDialog(message, true, "Error", func() {})
}

func (td *toastDialog) close() {
	td.view.Delete()
}
