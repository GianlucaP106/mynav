package ui

import (
	"fmt"
	"mynav/pkg/constants"
	"mynav/pkg/events"
	"mynav/pkg/persistence"

	"github.com/awesome-gocui/gocui"
)

type TmuxPreviewView struct {
	view    *View
	content *persistence.Value[string]
}

var _ Viewable = new(TmuxPreviewView)

func NewTmuxPreviewView() *TmuxPreviewView {
	return &TmuxPreviewView{}
}

func GetTmuxPreviewView() *TmuxPreviewView {
	return GetViewable[*TmuxPreviewView]()
}

func (t *TmuxPreviewView) Focus() {
	FocusView(t.View().Name())
}

func (t *TmuxPreviewView) View() *View {
	return t.view
}

func (t *TmuxPreviewView) Init() {
	t.view = GetViewPosition(constants.TmuxPreviewViewName).Set()

	t.view.Title = withSurroundingSpaces("Tmux Preview")
	t.view.TitleColor = gocui.ColorBlue
	t.view.FrameColor = gocui.ColorGreen
	t.view.Wrap = false

	t.content = persistence.NewValue("")

	events.AddEventListener(constants.TmuxPreviewChangeEventName, func(s string) {
		t.refresh()
		RenderView(t)
	})

	t.refresh()
}

func (t *TmuxPreviewView) refresh() {
	pane := GetTmuxPaneView().getSelectedPane()
	if pane == nil {
		return
	}

	content, err := pane.Capture()
	if err != nil {
		return
	}

	t.content.Set(content)
}

func (t *TmuxPreviewView) Render() error {
	t.view.Clear()
	fmt.Fprintln(t.view, t.content.Get())
	return nil
}
