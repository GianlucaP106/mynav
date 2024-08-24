package ui

import (
	"mynav/pkg/core"
	"mynav/pkg/tui"

	"github.com/awesome-gocui/gocui"
)

type UI struct {
	mainTabGroup *tui.TabGroup
	api          *core.Api
	views        []viewable
}

type viewable interface {
	init()
	getView() *tui.View
	refresh()
	render() error
}

func (ui *UI) InitUI() *UI {
	ui.views = []viewable{
		newHeaderView(),
		newTopicsView(),
		newWorkspcacesView(),
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
	ui.buildGithubTab()

	ui.mainTabGroup.FocusTabByIndex(0)

	systemUpdate()

	if getApi().GlobalConfiguration.GetLastTab() != "" {
		ui.mainTabGroup.FocusTab(getApi().GlobalConfiguration.GetLastTab())
	}

	if getApi().GlobalConfiguration.GetLastTab() == "main" {
		if getApi().Core.GetSelectedWorkspace() != nil {
			getWorkspacesView().focus()
		} else {
			getTopicsView().focus()
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

		// NOTE: not using these views for the moment - https://github.com/GianlucaP106/mynav/issues/287
		// newPortView(),
		// newPsView(),

		newGithubPrView(),
		newGithubProfileView(),
		newGithubRepoView(),
	}

	ui.sealViews()

	ui.mainTabGroup = tui.NewTabGroup(focusView)
	ui.buildTmuxTab()
	ui.buildGithubTab()
	ui.mainTabGroup.FocusTabByIndex(0)

	systemUpdate()
	ui.initGlobalKeybindings()
}

func (ui *UI) askConfig() {
	openConfirmationDialog(func(b bool) {
		if !b {
			getApi().GlobalConfiguration.SetStandalone(true)
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
				getApi().GlobalConfiguration.SetLastTab(tab.Frame.Name())
			})
		}).
		Set('[', "Cycle tab left", func() {
			ui.mainTabGroup.DecrementSelectedTab(func(tab *tui.Tab) {
				getApi().GlobalConfiguration.SetLastTab(tab.Frame.Name())
			})
		}).
		Set('?', "Toggle cheatsheet", func() {
			openHelpDialog(nil, func() {})
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
	tab.AddView(getHeaderView().view, tui.NoPosition)
	tab.AddView(getTopicsView().view, tui.TopLeftPosition)
	tab.AddView(getWorkspacesView().view, tui.TopRightPosition)
	tab.GenerateNavigationKeyBindings()
	return tab
}

func (ui *UI) buildTmuxTab() *tui.Tab {
	tab := ui.mainTabGroup.NewTab("tmux", getTmuxSessionView().getView().Name())
	tab.AddView(getTmuxSessionView().view, tui.TopLeftPosition)
	tab.AddView(getTmuxWindowView().view, tui.TopRightPosition)
	tab.AddView(getTmuxPreviewView().view, tui.NoPosition)
	tab.AddView(getTmuxPaneView().view, tui.NoPosition)
	tab.AddView(getHeaderView().view, tui.NoPosition)
	tab.GenerateNavigationKeyBindings()
	return tab
}

func (ui *UI) buildGithubTab() *tui.Tab {
	tab := ui.mainTabGroup.NewTab("github", getGithubProfileView().getView().Name())
	tab.AddView(getGithubProfileView().view, tui.TopLeftPosition)
	tab.AddView(getGithubRepoView().view, tui.TopRightPosition)
	tab.AddView(getGithubPrView().view, tui.BottomLeftPosition)
	tab.AddView(getHeaderView().view, tui.NoPosition)
	tab.GenerateNavigationKeyBindings()
	return tab
}

func runAction(f func()) {
	tui.Suspend()
	f()
	tui.Resume()
	refreshFsAsync()
	refreshTmuxViewsAsync()
}
