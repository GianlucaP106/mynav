package ui

import (
	"fmt"
	"mynav/pkg/core"
	"mynav/pkg/tui"
)

type tmuxPreviewView struct {
	view    *tui.View
	content *core.Value[string]
}

var _ viewable = new(tmuxPreviewView)

func newTmuxPreviewView() *tmuxPreviewView {
	return &tmuxPreviewView{}
}

func getTmuxPreviewView() *tmuxPreviewView {
	return getViewable[*tmuxPreviewView]()
}

func (t *tmuxPreviewView) Focus() {
	focusView(t.getView().Name())
}

func (t *tmuxPreviewView) getView() *tui.View {
	return t.view
}

func (t *tmuxPreviewView) init() {
	t.view = getViewPosition(TmuxPreviewView).Set()

	t.view.Title = tui.WithSurroundingSpaces("Tmux Preview")
	styleView(t.view)
	t.view.Wrap = false

	t.content = core.NewValue("")

	t.refresh()
}

func (t *tmuxPreviewView) refresh() {
	pane := getTmuxPaneView().getSelectedPane()
	if pane == nil {
		t.content.Set("")
		return
	}

	content, err := pane.Capture()
	if err != nil {
		return
	}

	t.content.Set(content)
}

func (t *tmuxPreviewView) render() error {
	t.view.Clear()
	t.view.Resize(getViewPosition(t.view.Name()))
	fmt.Fprintln(t.view, t.content.Get())
	return nil
}
