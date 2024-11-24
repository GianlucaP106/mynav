package app

import (
	"fmt"
	"mynav/pkg/tui"
)

type Preview struct {
	view *tui.View
}

func newPreview() *Preview {
	p := &Preview{}
	return p
}

func (p *Preview) init() {
	p.view = a.ui.SetView(getViewPosition(PreviewView))
	p.view.Title = " Preview "
	a.styleView(p.view)
}

func (p *Preview) show(content string) {
	p.view.Clear()
	p.view = a.ui.SetView(getViewPosition(p.view.Name()))
	a.styleView(p.view)
	fmt.Fprintln(p.view, content)
}
