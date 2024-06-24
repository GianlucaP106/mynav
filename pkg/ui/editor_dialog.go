package ui

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

	prevView := GetFocusedView()
	ed.view.Editor = NewSimpleEditor(func(s string) {
		ed.Close()
		if prevView != nil {
			FocusViewInternal(prevView.Name())
		}
		onEnter(s)
	}, func() {
		ed.Close()
		if prevView != nil {
			FocusViewInternal(prevView.Name())
		}
		onEsc()
	})

	FocusViewInternal(ed.view.Name())

	return ed
}

func (ed *EditorDialog) Close() {
	ToggleCursor(false)
	DeleteView(ed.view.Name())
}
