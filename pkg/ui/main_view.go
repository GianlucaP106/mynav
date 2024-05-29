package ui

import "github.com/awesome-gocui/gocui"

type MainView struct {
	wv *WorkspacesView
	tv *TopicsView
}

var _ View = &MainView{}

func newMainView(wv *WorkspacesView, tv *TopicsView) *MainView {
	return &MainView{
		wv: wv,
		tv: tv,
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

func (ui *UI) focusMainView(window string) {
	FocusView(window)

	wv := GetInternalView(WorkspacesViewName)
	tv := GetInternalView(TopicViewName)

	off := gocui.ColorBlue
	on := gocui.ColorGreen

	switch window {
	case WorkspacesViewName:
		wv.FrameColor = on
		tv.FrameColor = off
	case TopicViewName:
		wv.FrameColor = off
		tv.FrameColor = on
	}
}

func (mv *MainView) Init(ui *UI) {
	mv.tv.Init(ui)
	mv.wv.Init(ui)
}

func (mv *MainView) Render(ui *UI) error {
	if !Api().IsConfigInitialized {
		GetDialog[*ConfirmationDialog](ui).Open(func(b bool) {
			if !b {
				return
			}

			Api().InitConfiguration()

			mv.Init(ui)
			ui.FocusTopicsView()
		}, "No configuration found. Would you like to initialize this directory?")
		return nil
	}
	mv.tv.Render(ui)
	mv.wv.Render(ui)

	if ui.action.Command != nil {
		return gocui.ErrQuit
	}

	return nil
}

func (ui *UI) RefreshMainView() {
	tv := GetView[*TopicsView](ui)
	wv := GetView[*WorkspacesView](ui)
	tv.refreshTopics()
	wv.refreshWorkspaces()
}

func (ui *UI) RefreshWorkspaces() {
	wv := GetView[*WorkspacesView](ui)
	wv.refreshWorkspaces()
}
