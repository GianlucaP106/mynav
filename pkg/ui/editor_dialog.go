package ui

import "fmt"

type EditorDialog struct {
	view *View
}

type EditorSize = uint

const EditorDialogName = "EditorDialog"

const (
	Small EditorSize = iota
	Large
)

func OpenEditorDialog(onEnter func(string), onEsc func(), title string, size EditorSize) *EditorDialog {
	return OpenEditorDialogWithDefaultValue(onEnter, onEsc, title, size, "")
}

func OpenEditorDialogWithDefaultValue(onEnter func(string), onEsc func(), title string, size EditorSize, defaultValue string) *EditorDialog {
	ed := &EditorDialog{}

	var height int
	switch size {
	case Small:
		height = 3
	case Large:
		height = 7
	}

	ed.view = SetCenteredView(EditorDialogName, 80, height, 0)
	ed.view.Editable = true
	ed.view.Title = withSurroundingSpaces(title)
	ed.view.Wrap = true
	ToggleCursor(true)

	if defaultValue != "" {
		fmt.Fprint(ed.view, defaultValue)
		ed.view.MoveCursor(len(defaultValue), 0)
	}

	prevView := GetFocusedView()
	ed.view.Editor = NewSimpleEditor(func(s string) {
		ed.Close()
		if prevView != nil {
			prevView.Focus()
		}
		onEnter(s)
	}, func() {
		ed.Close()
		if prevView != nil {
			prevView.Focus()
		}
		onEsc()
	})

	ed.view.Focus()

	return ed
}

func (ed *EditorDialog) Close() {
	ToggleCursor(false)
	ed.view.Delete()
}
