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
	gprv        *GithubPrView
	configAsked bool
}

var _ View = &MainView{}

func newMainView(wv *WorkspacesView, tv *TopicsView, pv *PortView, tmv *TmuxSessionView, gprv *GithubPrView) *MainView {
	return &MainView{
		wv:          wv,
		tv:          tv,
		pv:          pv,
		tmv:         tmv,
		gprv:        gprv,
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

func (ui *UI) FocusPrView() {
	ui.focusMainView(GithubPrViewName)
}

func (ui *UI) focusMainView(viewName string) {
	FocusView(viewName)

	wv := GetInternalView(WorkspacesViewName)
	tv := GetInternalView(TopicViewName)
	pv := GetInternalView(PortViewName)
	tmv := GetInternalView(TmuxSessionViewName)
	gprv := GetInternalView(GithubPrViewName)
	views := []*gocui.View{wv, tv, pv, tmv, gprv}

	off := gocui.ColorBlue
	on := gocui.ColorGreen

	for _, v := range views {
		if v.Name() == viewName {
			v.FrameColor = on
		} else {
			v.FrameColor = off
		}
	}
}

func (mv *MainView) Init(ui *UI) {
	mv.tv.Init(ui)
	mv.pv.Init(ui)
	mv.gprv.Init(ui)
	mv.tmv.Init(ui)
	mv.wv.Init(ui)
}

func (mv *MainView) Render(ui *UI) error {
	if !Api().Core.IsConfigInitialized && !mv.configAsked {
		mv.configAsked = true

		homeDir, _ := os.UserHomeDir()
		cwd, _ := os.Getwd()
		if homeDir == cwd {
			mv.tmv.standalone = true
		}

		GetDialog[*ConfirmationDialog](ui).Open(func(b bool) {
			if !b {
				mv.tmv.standalone = true
				return
			}

			Api().InitConfiguration()
			mv.Init(ui)
			ui.handleUpdate()
		}, "No configuration found. Would you like to initialize this directory?")
		return nil
	}

	if Api().Core.IsConfigInitialized {
		mv.tv.Render(ui)
		mv.pv.Render(ui)
		mv.gprv.Render(ui)
		mv.tmv.Render(ui)
		mv.wv.Render(ui)
	} else if mv.tmv.standalone {
		mv.tmv.Render(ui)
	}

	ui.handleUpdate()

	if ui.action.Command != nil {
		return gocui.ErrQuit
	}

	return nil
}

func (ui *UI) RefreshMainView() {
	tv := GetView[*TopicsView](ui)
	wv := GetView[*WorkspacesView](ui)
	pv := GetView[*PortView](ui)
	tmv := GetView[*TmuxSessionView](ui)
	if !tmv.standalone {
		tv.refreshTopics()
		pv.refreshPorts()
		wv.refreshWorkspaces()
	}
	tmv.refreshTmuxSessions()
}
