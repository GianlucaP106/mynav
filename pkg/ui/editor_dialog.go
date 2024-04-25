package ui

import "github.com/awesome-gocui/gocui"

type EditorDialogState struct {
	editor   Editor
	viewName string
	title    string
	active   bool
}

func newEditorDialogState() *EditorDialogState {
	editor := &EditorDialogState{
		viewName: "EditorDialog",
		active:   false,
	}
	return editor
}

func (ui *UI) initEditorDialogView() *gocui.View {
	view := ui.setCenteredView(ui.editor.viewName, 80, 3, 0)
	view.Editable = true
	view.Editor = ui.editor.editor
	view.Title = withSurroundingSpaces(ui.editor.title)
	ui.toggleCursor(true)
	return view
}

type OpenEditorRequest struct {
	onEnter func(string)
	onEsc   func()
	title   string
}

func (ui *UI) openEditorDialog(onEnter func(string), onEsc func(), title string) {
	ui.editor.title = title
	ui.editor.editor = gocui.EditorFunc(simpleEditorFactory(func(s string) {
		ui.closeEditorDialog()
		onEnter(s)
	}, func() {
		ui.closeEditorDialog()
		onEsc()
	}))
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
