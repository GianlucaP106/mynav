package ui

import (
	"log"

	"github.com/awesome-gocui/gocui"
)

func (ui *UI) setView(viewName string) *gocui.View {
	maxX, maxY := ui.gui.Size()
	views := map[string]func() *gocui.View{}
	views[ui.workspaces.viewName] = func() *gocui.View {
		view, _ := ui.gui.SetView(ui.workspaces.viewName, (maxX)/3+1, 8, maxX-2, maxY-4, 0)
		return view
	}

	views[ui.topics.viewName] = func() *gocui.View {
		view, _ := ui.gui.SetView(ui.topics.viewName, 2, 8, maxX/3-1, maxY-4, 0)
		return view
	}

	views[ui.header.viewName] = func() *gocui.View {
		view, _ := ui.gui.SetView(ui.header.viewName, 2, 1, maxX-2, 6, 0)
		return view
	}

	f := views[viewName]
	if f == nil {
		log.Panicln("invalid view")
	}

	return f()
}

func (ui *UI) setCenteredView(name string, sizeX int, sizeY int, verticalOffset int) *gocui.View {
	maxX, maxY := ui.gui.Size()
	view, _ := ui.gui.SetView(name, maxX/2-sizeX/2, maxY/2-sizeY/2+verticalOffset, maxX/2+sizeX/2, maxY/2+sizeY/2+verticalOffset, 0)
	return view
}
