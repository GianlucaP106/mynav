package ui

import (
	"mynav/pkg/api"
	"mynav/pkg/tui"

	"github.com/awesome-gocui/gocui"
)

type UI struct {
	mainTabGroup *tui.TabGroup
	api          *api.Api
	views        []viewable
}

type viewable interface {
	init()
	getView() *tui.View
	render() error
}

func (ui *UI) InitUI() *UI {
	ui.views = []viewable{
		newHeaderView(),
		newTopicsView(),
		newWorkspcacesView(),
		newPortView(),
		newPsView(),
		newTmuxSessionView(),
		newTmuxWindowView(),
		newTmuxPaneView(),
		newTmuxPreviewView(),
		newGithubProfileView(),
		newGithubPrView(),
		newGithubRepoView(),
	}

	ui.sealViews()

	ui.mainTabGroup = tui.NewTabGroup(focusView)
	ui.buildMainTab()
	ui.buildTmuxTab()
	ui.buildSystemTab()
	ui.buildGithubTab()

	ui.mainTabGroup.FocusTabByIndex(0)

	systemUpdate()

	if getApi().Configuration.GetLastTab() != "" {
		ui.mainTabGroup.FocusTab(getApi().Configuration.GetLastTab())
	}

	if getApi().Configuration.GetLastTab() == "main" {
		if getApi().Core.GetSelectedWorkspace() != nil {
			getWorkspacesView().Focus()
		} else {
			getTopicsView().Focus()
		}
	}

	ui.initGlobalKeybindings()
	return ui
}

func (ui *UI) initStandaloneUI() {
	ui.views = []viewable{
		newHeaderView(),
		newTmuxSessionView(),
		newTmuxWindowView(),
		newTmuxPaneView(),
		newTmuxPreviewView(),
		newPortView(),
		newPsView(),
		newGithubPrView(),
		newGithubProfileView(),
		newGithubRepoView(),
	}

	ui.sealViews()

	ui.mainTabGroup = tui.NewTabGroup(focusView)
	ui.buildTmuxTab()
	ui.buildSystemTab()
	ui.buildGithubTab()
	ui.mainTabGroup.FocusTabByIndex(0)

	systemUpdate()
	ui.initGlobalKeybindings()
}

func (ui *UI) askConfig() {
	openConfirmationDialog(func(b bool) {
		if !b {
			getApi().Configuration.SetStandalone(true)
			ui.initStandaloneUI()
			return
		}

		getApi().InitConfiguration()
		ui.InitUI()
	}, "No configuration found. Would you like to initialize this directory?")
}

func (ui *UI) initGlobalKeybindings() {
	quit := func() bool {
		return true
	}
	tui.NewKeybindingBuilder(nil).
		SetWithQuit(gocui.KeyCtrlC, quit, "Quit").
		SetWithQuit('q', quit, "Quit").
		SetWithQuit('q', quit, "Quit").
		Set(']', "Cycle tab right", func() {
			ui.mainTabGroup.IncrementSelectedTab(func(tab *tui.Tab) {
				getApi().Configuration.SetLastTab(tab.Frame.Name())
			})
		}).
		Set('[', "Cycle tab left", func() {
			ui.mainTabGroup.DecrementSelectedTab(func(tab *tui.Tab) {
				getApi().Configuration.SetLastTab(tab.Frame.Name())
			})
		}).
		Set('?', "Toggle cheatsheet", func() {
			OpenHelpDialog(nil, func() {})
		})
}

func (ui *UI) sealViews() {
	managers := make([]gocui.Manager, 0)
	for _, view := range ui.views {
		managers = append(managers, gocui.ManagerFunc(func(_ *gocui.Gui) error {
			return view.render()
		}))
	}

	tui.SetManagerFunctions(managers...)
	for _, v := range ui.views {
		v.init()
	}
}

func (ui *UI) buildMainTab() *tui.Tab {
	tab := ui.mainTabGroup.NewTab("main", getTopicsView().getView().Name())
	tab.AddView(getHeaderView().view, tui.None)
	tab.AddView(getTopicsView().view, tui.TopLeft)
	tab.AddView(getWorkspacesView().view, tui.TopRight)
	tab.GenerateNavigationKeyBindings()
	return tab
}

func (ui *UI) buildTmuxTab() *tui.Tab {
	tab := ui.mainTabGroup.NewTab("tmux", getTmuxSessionView().getView().Name())
	tab.AddView(getTmuxSessionView().view, tui.TopLeft)
	tab.AddView(getTmuxWindowView().view, tui.TopRight)
	tab.AddView(getTmuxPreviewView().view, tui.None)
	tab.AddView(getTmuxPaneView().view, tui.None)
	tab.AddView(getHeaderView().view, tui.None)
	tab.GenerateNavigationKeyBindings()
	return tab
}

func (ui *UI) buildSystemTab() *tui.Tab {
	tab := ui.mainTabGroup.NewTab("system", getPortView().getView().Name())
	tab.AddView(getPortView().view, tui.TopRight)
	tab.AddView(getPsView().view, tui.TopLeft)
	tab.AddView(getHeaderView().view, tui.None)
	tab.GenerateNavigationKeyBindings()
	return tab
}

func (ui *UI) buildGithubTab() *tui.Tab {
	tab := ui.mainTabGroup.NewTab("github", getGithubProfileView().getView().Name())
	tab.AddView(getGithubProfileView().view, tui.TopLeft)
	tab.AddView(getGithubRepoView().view, tui.TopRight)
	tab.AddView(getGithubPrView().view, tui.BottomLeft)
	tab.AddView(getHeaderView().view, tui.None)
	tab.GenerateNavigationKeyBindings()
	return tab
}
