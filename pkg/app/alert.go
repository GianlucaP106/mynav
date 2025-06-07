package app

import (
	"fmt"

	"github.com/GianlucaP106/mynav/pkg/tui"
	"github.com/awesome-gocui/gocui"
	"github.com/gookit/color"
)

type Alert struct {
	view  *tui.View
	title string
}

func alert(onConfirm func(bool), title string) *Alert {
	// build alert dialog
	cd := &Alert{}
	cd.title = title
	cd.view = a.ui.SetCenteredView(ConfirmationDialog, len(title)+5, 4, 0, 0)
	cd.view.Wrap = true
	cd.view.FrameColor = gocui.ColorWhite

	cd.view.Title = " Confirm "
	a.styleView(cd.view)

	// set key bindings
	prevView := a.ui.FocusedView()
	a.ui.KeyBinding(cd.view).
		Set(gocui.KeyEsc, "Cancel", func() {
			a.ui.DeleteView(cd.view)
			if prevView != nil {
				a.ui.FocusView(prevView)
			}
			onConfirm(false)
		}).
		Set(gocui.KeyEnter, "Confirm", func() {
			a.ui.DeleteView(cd.view)
			if prevView != nil {
				a.ui.FocusView(prevView)
			}
			onConfirm(true)
		})

	a.ui.FocusView(cd.view)

	// write view content
	cd.view.Clear()
	fmt.Fprintln(cd.view, color.Note.Sprint(" "+cd.title))
	fmt.Fprintln(cd.view)

	line := fmt.Sprintf(" %s %s %s %s %s ",
		timestampColor.Sprint("Press"),
		sessionMarkerColor.Sprint("Enter"),
		timestampColor.Sprint("to confirm,"),
		color.Danger.Sprint("Esc"),
		timestampColor.Sprint("to cancel"),
	)
	fmt.Fprintln(cd.view, line)

	return cd
}
