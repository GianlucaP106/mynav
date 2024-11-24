package app

import (
	"fmt"
	"mynav/pkg/core"
	"mynav/pkg/tui"
)

type Preview struct {
	view *tui.View

	// session to show
	previews []string

	// idx of preview to show
	idx int
}

func newPreview() *Preview {
	p := &Preview{}
	return p
}

func (p *Preview) init() {
	p.view = a.ui.SetView(getViewPosition(PreviewView))
	p.view.Title = " Preview "
	a.styleView(p.view)
	p.setPreviews(nil)
}

func (p *Preview) show(session *core.Session) {
	// if session is nil we set a list of 1
	if session == nil {
		p.setPreviews(nil)
		p.render()
		return
	}

	// get all windows for this session
	windows, _ := session.ListWindows()

	// collect all previews (one per pane)
	previews := make([]string, 0)
	for _, w := range windows {
		panes, _ := w.ListPanes()
		for _, p := range panes {
			preview, _ := p.Capture()
			previews = append(previews, preview)
		}
	}

	p.setPreviews(previews)
	p.render()
}

func (p *Preview) render() {
	p.view.Clear()
	a.ui.Resize(p.view, getViewPosition(p.view.Name()))

	if len(p.previews) == 0 {
		p.view.Subtitle = ""
		return
	}

	s := p.previews[p.idx]
	p.view.Subtitle = fmt.Sprintf(" %d / %d ", min(p.idx+1, len(p.previews)), len(p.previews))
	fmt.Fprintln(p.view, s)
}

func (p *Preview) setPreviews(previews []string) {
	if len(previews) == 0 {
		p.idx = 0
		p.previews = previews
		return
	}

	if p.idx >= len(previews) {
		p.idx = len(previews) - 1
	}
	p.previews = previews
}

func (p *Preview) increment() {
	if p.idx == len(p.previews)-1 {
		p.idx = 0
	} else {
		p.idx++
	}
	p.render()
}

func (p *Preview) decrement() {
	if p.idx == 0 {
		p.idx = len(p.previews) - 1
	} else {
		p.idx--
	}
	p.render()
}
