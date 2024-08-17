package ui

import (
	"fmt"
	"mynav/pkg/constants"
	"mynav/pkg/tui"
)

type editorDialog struct {
	view *tui.View
}

type editorSize = uint

const (
	smallEditorSize editorSize = iota
	largeEditorSize
)

func openEditorDialog(onEnter func(string), onEsc func(), title string, size editorSize) *editorDialog {
	return openEditorDialogWithDefaultValue(onEnter, onEsc, title, size, "")
}

func openEditorDialogWithDefaultValue(onEnter func(string), onEsc func(), title string, size editorSize, defaultValue string) *editorDialog {
	ed := &editorDialog{}

	var height int
	switch size {
	case smallEditorSize:
		height = 3
	case largeEditorSize:
		height = 7
	}

	ed.view = tui.SetCenteredView(constants.EditorDialogName, 80, height, 0)
	ed.view.Editable = true
	ed.view.Title = tui.WithSurroundingSpaces(title)
	ed.view.Wrap = true
	ed.view.FrameColor = tui.OnFrameColor
	tui.StyleView(ed.view)
	tui.ToggleCursor(true)

	if defaultValue != "" {
		fmt.Fprint(ed.view, defaultValue)
		ed.view.MoveCursor(len(defaultValue), 0)
	}

	prevView := tui.GetFocusedView()
	ed.view.Editor = tui.NewSimpleEditor(func(s string) {
		ed.close()
		if prevView != nil {
			prevView.Focus()
		}
		onEnter(s)
	}, func() {
		ed.close()
		if prevView != nil {
			prevView.Focus()
		}
		onEsc()
	})

	ed.view.Focus()

	return ed
}

func (ed *editorDialog) close() {
	tui.ToggleCursor(false)
	ed.view.Delete()
}
