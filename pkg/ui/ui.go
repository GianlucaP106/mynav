package ui

import (
	"mynav/pkg/api"
	"reflect"

	"github.com/awesome-gocui/gocui"
)

type UI struct {
	mainTabGroup *TabGroup
	api          *api.Api
	views        []Viewable
}

type Viewable interface {
	Init()
	View() *View
	Render() error
}

func (ui *UI) InitUI() *UI {
	ui.views = []Viewable{
		NewHeaderView(),
		NewTopicsView(),
		NewWorkspcacesView(),
		NewPortView(),
		NewPsView(),
		NewTmuxSessionView(),
		NewTmuxWindowView(),
		NewTmuxPaneView(),
		NewTmuxPreviewView(),
		NewGithubProfileView(),
		NewGithubPrView(),
		NewGithubRepoView(),
	}

	ui.sealViews()

	ui.mainTabGroup = NewTabGroup(
		ui.buildMainTab(),
		ui.buildTmuxTab(),
		ui.buildSystemTab(),
		ui.buildGithubTab(),
	)

	ui.mainTabGroup.FocusTabByIndex(0)

	SystemUpdate()

	if Api().Configuration.GetLastTab() != "" {
		ui.mainTabGroup.FocusTab(Api().Configuration.GetLastTab())
	}

	if Api().Configuration.GetLastTab() == "main" {
		if Api().Core.GetSelectedWorkspace() != nil {
			GetWorkspacesView().Focus()
		} else {
			GetTopicsView().Focus()
		}
	}

	ui.InitGlobalKeybindings()
	return ui
}

func (ui *UI) InitStandaloneUI() {
	ui.views = []Viewable{
		NewHeaderView(),
		NewTmuxSessionView(),
		NewTmuxWindowView(),
		NewTmuxPaneView(),
		NewTmuxPreviewView(),
		NewPortView(),
		NewPsView(),
		NewGithubPrView(),
		NewGithubProfileView(),
		NewGithubRepoView(),
	}

	ui.sealViews()

	ui.mainTabGroup = NewTabGroup(
		ui.buildTmuxTab(),
		ui.buildSystemTab(),
		ui.buildGithubTab(),
	)
	ui.mainTabGroup.FocusTabByIndex(0)

	SystemUpdate()
	ui.InitGlobalKeybindings()
}

func (ui *UI) AskConfig() {
	OpenConfirmationDialog(func(b bool) {
		if !b {
			Api().Configuration.SetStandalone(true)
			ui.InitStandaloneUI()
			return
		}

		Api().InitConfiguration()
		ui.InitUI()
	}, "No configuration found. Would you like to initialize this directory?")
}

func (ui *UI) addView(v Viewable) {
	t := reflect.TypeOf(v)
	for _, view := range ui.views {
		t2 := reflect.TypeOf(view)
		if t == t2 {
			return
		}
	}
	ui.views = append(ui.views, v)
}

func (ui *UI) InitGlobalKeybindings() {
	quit := func() bool {
		return true
	}
	NewKeybindingBuilder(nil).
		setWithQuit(gocui.KeyCtrlC, quit, "Quit").
		setWithQuit('q', quit, "Quit").
		setWithQuit('q', quit, "Quit").
		set(']', "Cycle tab right", func() {
			ui.mainTabGroup.IncrementSelectedTab(func(tab *Tab) {
				Api().Configuration.SetLastTab(tab.Frame.Name())
			})
		}).
		set('[', "Cycle tab left", func() {
			ui.mainTabGroup.DecrementSelectedTab(func(tab *Tab) {
				Api().Configuration.SetLastTab(tab.Frame.Name())
			})
		}).
		set('?', "Toggle cheatsheet", func() {
			OpenHelpView(nil, func() {})
		})
}

func (ui *UI) sealViews() {
	managers := make([]gocui.Manager, 0)
	for _, view := range ui.views {
		managers = append(managers, gocui.ManagerFunc(func(_ *gocui.Gui) error {
			return view.Render()
		}))
	}

	SetManagerFunctions(managers...)
	for _, v := range ui.views {
		v.Init()
	}
}

func (ui *UI) buildMainTab() *Tab {
	tab := NewTab("main", GetTopicsView().View().Name())
	tab.AddView(GetHeaderView(), None)
	tab.AddView(GetTopicsView(), TopLeft)
	tab.AddView(GetWorkspacesView(), TopRight)
	tab.GenerateNavigationKeyBindings()
	return tab
}

func (ui *UI) buildTmuxTab() *Tab {
	tab := NewTab("tmux", GetTmuxSessionView().View().Name())
	tab.AddView(GetTmuxSessionView(), TopLeft)
	tab.AddView(GetTmuxWindowView(), TopRight)
	tab.AddView(GetTmuxPreviewView(), None)
	tab.AddView(GetTmuxPaneView(), None)
	tab.AddView(GetHeaderView(), None)
	tab.GenerateNavigationKeyBindings()
	return tab
}

func (ui *UI) buildSystemTab() *Tab {
	tab := NewTab("system", GetPortView().View().Name())
	tab.AddView(GetPortView(), TopRight)
	tab.AddView(GetPsView(), TopLeft)
	tab.AddView(GetHeaderView(), None)
	tab.GenerateNavigationKeyBindings()
	return tab
}

func (ui *UI) buildGithubTab() *Tab {
	tab := NewTab("github", GetGithubProfileView().View().Name())
	tab.AddView(GetGithubProfileView(), TopLeft)
	tab.AddView(GetGithubRepoView(), TopRight)
	tab.AddView(GetGithubPrView(), BottomLeft)
	tab.AddView(GetHeaderView(), None)
	tab.GenerateNavigationKeyBindings()
	return tab
}
