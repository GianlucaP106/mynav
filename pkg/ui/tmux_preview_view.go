package ui

import (
	"fmt"
	"mynav/pkg/events"
	"mynav/pkg/persistence"
	"mynav/pkg/tui"
)

type tmuxPreviewView struct {
	view    *tui.View
	content *persistence.Value[string]
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

	t.content = persistence.NewValue("")

	events.AddEventListener(events.TmuxPreviewChangeEvent, func(s string) {
		t.refresh()
		renderView(t)
	})

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
	t.view = getViewPosition(t.view.Name()).Set()
	fmt.Fprintln(t.view, t.content.Get())
	return nil
}
