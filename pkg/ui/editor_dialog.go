package ui

import "github.com/awesome-gocui/gocui"

type EditorDialogState struct {
	editor   Editor
	viewName string
	title    string
	active   bool
	height   int
}

type EditorSize = uint

const (
	Small EditorSize = iota
	Large
)

func newEditorDialogState() *EditorDialogState {
	editor := &EditorDialogState{
		viewName: "EditorDialog",
		active:   false,
		height:   3,
	}
	return editor
}

func (ui *UI) initEditorDialogView() *gocui.View {
	view := ui.setCenteredView(ui.editor.viewName, 80, ui.editor.height, 0)
	view.Editable = true
	view.Editor = ui.editor.editor
	view.Title = withSurroundingSpaces(ui.editor.title)
	view.Wrap = true
	ui.toggleCursor(true)
	return view
}

type OpenEditorRequest struct {
	onEnter func(string)
	onEsc   func()
	title   string
}

func (ui *UI) openEditorDialog(onEnter func(string), onEsc func(), title string, size EditorSize) {
	switch size {
	case Small:
		ui.editor.height = 3
	case Large:
		ui.editor.height = 7
	}
	ui.editor.title = title
	ui.editor.editor = newSimpleEditor(func(s string) {
		ui.closeEditorDialog()
		onEnter(s)
	}, func() {
		ui.closeEditorDialog()
		onEsc()
	})
	ui.editor.active = true
}

func (ui *UI) closeEditorDialog() {
	ui.editor.active = false
	ui.gui.Cursor = false
	ui.gui.DeleteView(ui.editor.viewName)
}

func (ui *UI) renderEditorDialog() {
	if !ui.editor.active {
		return
	}
	ui.initEditorDialogView()
	ui.focusView(ui.editor.viewName)
}
