package ui

import (
	"fmt"
	"mynav/pkg/api"

	"github.com/awesome-gocui/gocui"
	"github.com/gookit/color"
)

const TmuxSessionViewName = "TmuxSessionView"

type TmuxSessionView struct {
	editor       Editor
	listRenderer *ListRenderer
	sessions     []*api.TmuxSession
}

var _ Dialog = &TmuxSessionView{}

func newTmuxSessionView() *TmuxSessionView {
	ts := &TmuxSessionView{}

	return ts
}

func (tv *TmuxSessionView) Open(exit func()) {
	view := SetViewLayout(tv.Name())

	view.Title = withSurroundingSpaces("TMUX Sessions")
	view.TitleColor = gocui.ColorBlue
	view.FrameColor = gocui.ColorBlue

	_, sizeY := view.Size()
	tv.listRenderer = newListRenderer(0, sizeY, 0)
	tv.refreshTmuxSessions()

	tv.editor = NewListRendererEditor(
		func() {
			tv.listRenderer.decrement()
		}, func() {
			tv.listRenderer.increment()
		}, func() {
		}, func() {
			tv.Close()
			exit()
		})

	view.Editor = tv.editor
	view.Editable = true
	FocusView(tv.Name())
}

func (tv *TmuxSessionView) Close() {
	DeleteView(tv.Name())
}

func (ts *TmuxSessionView) refreshTmuxSessions() {
	out := make([]*api.TmuxSession, 0)
	for _, session := range Api().TmuxSessionContainer {
		out = append(out, session)
	}

	ts.sessions = out

	if ts.listRenderer != nil {
		newListSize := len(ts.sessions)
		if ts.listRenderer.listSize != newListSize {
			ts.listRenderer.setListSize(newListSize)
		}
	}
}

func (tv *TmuxSessionView) Name() string {
	return TmuxSessionViewName
}

func (tv *TmuxSessionView) Init() *gocui.View {
	view := GetInternalView(tv.Name())

	view.Title = withSurroundingSpaces("TMUX Sessions")
	view.TitleColor = gocui.ColorBlue
	view.FrameColor = gocui.ColorBlue

	_, sizeY := view.Size()
	tv.listRenderer = newListRenderer(0, sizeY, 0)
	tv.refreshTmuxSessions()

	KeyBinding(tv.Name()).
		set('j', func() {
			tv.listRenderer.increment()
		}).
		set('k', func() {
			tv.listRenderer.decrement()
		}).
		set(gocui.KeyEsc, func() {
			tv.Close()
		})

	return view
}

func (tv *TmuxSessionView) Render(ui *UI) error {
	view := GetInternalView(tv.Name())
	if view == nil {
		return nil
	}

	view.Clear()
	sizeX, _ := view.Size()
	tv.listRenderer.forEach(func(idx int) {
		session := tv.sessions[idx]
		line := withSpacePadding(" "+session.Name+" ", sizeX)
		if idx == tv.listRenderer.selected {
			line = color.New(color.BgCyan).Sprint(line)
		}
		fmt.Fprintln(view, line)
	})

	return nil
}
