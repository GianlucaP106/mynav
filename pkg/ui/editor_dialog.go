package ui

type EditorDialog struct {
	editor Editor
	title  string
	height int
}

var _ Dialog = &EditorDialog{}

type EditorSize = uint

const EditorDialogStateName = "EditorDialog"

const (
	Small EditorSize = iota
	Large
)

func newEditorDialogState() *EditorDialog {
	editor := &EditorDialog{
		height: 3,
	}
	return editor
}

func (eds *EditorDialog) Name() string {
	return EditorDialogStateName
}

func (ed *EditorDialog) Open(onEnter func(string), onEsc func(), title string, size EditorSize) {
	switch size {
	case Small:
		ed.height = 3
	case Large:
		ed.height = 7
	}
	ed.title = title

	prevView := GetFocusedView()
	ed.editor = NewSimpleEditor(func(s string) {
		ed.Close()
		if prevView != nil {
			FocusView(prevView.Name())
		}
		onEnter(s)
	}, func() {
		ed.Close()
		if prevView != nil {
			FocusView(prevView.Name())
		}
		onEsc()
	})

	view := SetCenteredView(ed.Name(), 80, ed.height, 0)
	FocusView(ed.Name())
	view.Editable = true
	view.Editor = ed.editor
	view.Title = withSurroundingSpaces(ed.title)
	view.Wrap = true
	ToggleCursor(true)
}

func (ed *EditorDialog) Close() {
	ToggleCursor(false)
	DeleteView(ed.Name())
}

func (eds *EditorDialog) Render(ui *UI) error {
	return nil
}
