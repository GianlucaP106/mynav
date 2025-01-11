package app

import (
	"fmt"
	"mynav/pkg/tui"
)

type Editor struct {
	view *tui.View
}

type editorSize = uint

const (
	smallEditorSize editorSize = iota
	largeEditorSize
)

func editor(onEnter func(string), onEsc func(), title string, size editorSize, defaultValue string) *Editor {
	ed := &Editor{}

	var height int
	switch size {
	case smallEditorSize:
		height = 3
	case largeEditorSize:
		height = 7
	}

	ed.view = a.ui.SetCenteredView(EditorDialog, 80, height, 0)
	ed.view.Editable = true
	ed.view.Title = fmt.Sprintf(" %s ", title)
	ed.view.Wrap = true
	a.styleView(ed.view)
	a.ui.Cursor = true

	if defaultValue != "" {
		fmt.Fprint(ed.view, defaultValue)
		ed.view.MoveCursor(len(defaultValue), 0)
	}

	prevView := a.ui.FocusedView()
	ed.view.Editor = tui.NewSimpleEditor(func(s string) {
		ed.close()
		if prevView != nil {
			a.ui.FocusView(prevView)
		}
		onEnter(s)
	}, func() {
		ed.close()
		if prevView != nil {
			a.ui.FocusView(prevView)
		}
		onEsc()
	}, nil)

	a.ui.FocusView(ed.view)

	return ed
}

func (ed *Editor) close() {
	a.ui.Cursor = false
	a.ui.DeleteView(ed.view)
}
