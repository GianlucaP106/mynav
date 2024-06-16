package ui

import (
	"os"

	"github.com/awesome-gocui/gocui"
)

type MainView struct {
	wv          *WorkspacesView
	tv          *TopicsView
	pv          *PortView
	tmv         *TmuxSessionView
	configAsked bool
}

var _ View = &MainView{}

func newMainView(wv *WorkspacesView, tv *TopicsView, pv *PortView, tmv *TmuxSessionView) *MainView {
	return &MainView{
		wv:          wv,
		tv:          tv,
		pv:          pv,
		tmv:         tmv,
		configAsked: false,
	}
}

func (mv *MainView) Name() string {
	return "MainView"
}

func (mv *MainView) RequiresManager() bool {
	return true
}

func (ui *UI) FocusTopicsView() {
	ui.focusMainView(TopicViewName)
}

func (ui *UI) FocusWorkspacesView() {
	ui.focusMainView(WorkspacesViewName)
}

func (ui *UI) FocusPortView() {
	ui.focusMainView(PortViewName)
}

func (ui *UI) FocusTmuxView() {
	ui.focusMainView(TmuxSessionViewName)
}

func (ui *UI) focusMainView(window string) {
	FocusView(window)

	wv := GetInternalView(WorkspacesViewName)
	tv := GetInternalView(TopicViewName)
	pv := GetInternalView(PortViewName)
	tmv := GetInternalView(TmuxSessionViewName)

	off := gocui.ColorBlue
	on := gocui.ColorGreen

	switch window {
	case WorkspacesViewName:
		wv.FrameColor = on
		tv.FrameColor = off
		pv.FrameColor = off
		tmv.FrameColor = off
	case TopicViewName:
		tv.FrameColor = on
		wv.FrameColor = off
		pv.FrameColor = off
		tmv.FrameColor = off
	case PortViewName:
		pv.FrameColor = on
		tv.FrameColor = off
		wv.FrameColor = off
		tmv.FrameColor = off
	case TmuxSessionViewName:
		tmv.FrameColor = on
		pv.FrameColor = off
		tv.FrameColor = off
		wv.FrameColor = off
	}
}

func (mv *MainView) Init(ui *UI) {
	mv.tv.Init(ui)
	mv.pv.Init(ui)
	mv.tmv.Init(ui)
	mv.wv.Init(ui)
}

func (mv *MainView) Render(ui *UI) error {
	if !Api().IsConfigInitialized && !mv.configAsked {
		mv.configAsked = true

		homeDir, _ := os.UserHomeDir()
		cwd, _ := os.Getwd()
		if homeDir == cwd {
			mv.tmv.standalone = true
			return nil
		}

		GetDialog[*ConfirmationDialog](ui).Open(func(b bool) {
			if !b {
				mv.tmv.standalone = true
				return
			}

			Api().InitConfiguration()
			mv.Init(ui)
			ui.FocusTopicsView()
			ui.handleUpdate()
		}, "No configuration found. Would you like to initialize this directory?")
		return nil
	}

	if Api().IsConfigInitialized {
		mv.tv.Render(ui)
		mv.pv.Render(ui)
		mv.tmv.Render(ui)
		mv.wv.Render(ui)
	} else if mv.tmv.standalone {
		mv.tmv.Render(ui)
	}

	if ui.action.Command != nil {
		return gocui.ErrQuit
	}

	return nil
}

func (ui *UI) RefreshMainView() {
	tv := GetView[*TopicsView](ui)
	wv := GetView[*WorkspacesView](ui)
	pv := GetView[*PortView](ui)
	tv.refreshTopics()
	wv.refreshWorkspaces()
	pv.refreshPorts()
}

func (ui *UI) RefreshWorkspaces() {
	wv := GetView[*WorkspacesView](ui)
	wv.refreshWorkspaces()
}
