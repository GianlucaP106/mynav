package ui

import (
	"github.com/awesome-gocui/gocui"
)

type FsState struct {
	viewName   string
	focusedTab string
}

func newFsState() *FsState {
	return &FsState{
		viewName: "FsView",
	}
}

func (ui *UI) initFsView() *gocui.View {
	view := ui.getView(ui.fs.viewName)
	if view != nil {
		return view
	}
	sizeX, sizeY := ui.gui.Size()
	view = ui.setCenteredView(ui.fs.viewName, sizeX, sizeY-7, 3)
	view.Frame = false
	ui.fs.focusedTab = ui.topics.viewName
	return view
}

func (ui *UI) setFocusedFsView(focusedTab string) {
	if ui.fs.focusedTab != focusedTab {
		ui.fs.focusedTab = focusedTab
	}
	ui.focusView(ui.fs.focusedTab)

	wv := ui.getView(ui.workspaces.viewName)
	tv := ui.getView(ui.topics.viewName)
	off := gocui.ColorBlue
	on := gocui.ColorGreen
	tab := ui.fs.focusedTab
	switch tab {
	case ui.workspaces.viewName:
		wv.FrameColor = on
		tv.FrameColor = off
	case ui.topics.viewName:
		wv.FrameColor = off
		tv.FrameColor = on

	}
}

func (ui *UI) renderFsView() error {
	if !ui.controller.IsConfigInitialized() {
		ui.openConfirmationDialog(func(b bool) {
			if b {
				ui.controller.InitConfiguration()
			}
		}, "No workspace configuration found. Would you like to initialize this directory?")
		return nil
	}
	ui.initFsView()
	ui.renderTopicsView()
	ui.renderWorkspacesView()
	ui.setFocusedFsView(ui.fs.focusedTab)
	return nil
}
