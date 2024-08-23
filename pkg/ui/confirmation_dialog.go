package ui

import (
	"fmt"
	"mynav/pkg/tui"
)

type confirmationDialog struct {
	view  *tui.View
	title string
}

func openConfirmationDialog(onConfirm func(bool), title string) *confirmationDialog {
	cd := &confirmationDialog{}
	cd.title = title
	prevView := tui.GetFocusedView()
	cd.view = tui.SetCenteredView(ConfirmationDialog, len(title)+5, 5, 0)
	cd.view.Title = tui.WithSurroundingSpaces("Confirm")
	cd.view.Wrap = true
	cd.view.Editable = true
	cd.view.FrameColor = onFrameColor
	styleView(cd.view)

	cd.view.Editor = tui.NewConfirmationEditor(func() {
		cd.close()
		if prevView != nil {
			prevView.Focus()
		}
		onConfirm(true)
	}, func() {
		cd.close()
		if prevView != nil {
			prevView.Focus()
		}
		onConfirm(false)
	})

	sizeX, _ := cd.view.Size()
	cd.view.Focus()
	cd.view.Clear()
	fmt.Fprintln(cd.view, tui.DisplayWhite(cd.title, tui.LeftAlign, sizeX))
	fmt.Fprintln(cd.view, tui.BlankLine(sizeX))
	fmt.Fprintln(cd.view, tui.DisplayWhite("Cancel [esc] | Confirm [enter]", tui.LeftAlign, sizeX))
	return cd
}

func (cd *confirmationDialog) close() {
	cd.view.Delete()
}
