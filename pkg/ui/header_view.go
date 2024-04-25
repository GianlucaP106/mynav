package ui

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
)

type HeaderState struct {
	viewName string
}

func newHeaderState() *HeaderState {
	titleView := &HeaderState{
		viewName: "HeaderView",
	}

	return titleView
}

func (ui *UI) initHeaderView() *gocui.View {
	view := ui.setView(ui.header.viewName)
	return view
}

func (ui *UI) renderHeaderView() {
	currentView := ui.getView(ui.header.viewName)
	if currentView == nil {
		currentView = ui.initHeaderView()
	}
	sizeX, _ := currentView.Size()
	fmt.Fprintln(currentView, blankLine(sizeX))
	fmt.Fprintln(currentView, displayLineNormal("Welcome to mynav, a workspace manager", Center, sizeX))
	fmt.Fprintln(currentView, blankLine(sizeX))
}
