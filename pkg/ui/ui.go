package ui

import (
	"errors"
	"log"
	"mynav/pkg/core"
	"mynav/pkg/tui"

	"github.com/awesome-gocui/gocui"
)

type UI struct {
	mainTabGroup *tui.TabGroup
	api          *core.Api
	refreshQueue chan func()
	views        []viewable
}

type viewable interface {
	init()
	getView() *tui.View
	refresh()
	render() error
}

const (
	offFrameColor = gocui.AttrDim | gocui.ColorWhite
	onFrameColor  = gocui.ColorWhite
)

var ui *UI

func Start(a *core.Api) {
	g := tui.NewTui()
	defer g.Close()

	ui = newUI(a)

	if api().LocalConfiguration.IsInitialized {
		ui.init()
	} else if api().GlobalConfiguration.Standalone {
		ui.initStanialone()
	} else {
		ui.askConfig()
	}

	err := g.MainLoop()
	if err != nil {
		if !errors.Is(err, gocui.ErrQuit) {
			log.Panicln(err)
		}
	}
}

func newUI(a *core.Api) *UI {
	return &UI{
		views: make([]viewable, 0),
		api:   a,
	}
}

func (ui *UI) init() *UI {
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

	if api().GlobalConfiguration.GetLastTab() != "" {
		ui.mainTabGroup.FocusTab(api().GlobalConfiguration.GetLastTab())
	}

	if api().GlobalConfiguration.GetLastTab() == "main" {
		if api().Workspaces.GetSelectedWorkspace() != nil {
			getWorkspacesView().focus()
		} else {
			getTopicsView().focus()
		}
	}

	ui.initGlobalKeybindings()
	ui.initRefreshExecutor()
	return ui
}

func (ui *UI) initStanialone() {
	ui.views = []viewable{
		newHeaderView(),
		newTmuxSessionView(),
		newTmuxWindowView(),
		newTmuxPaneView(),
		newTmuxPreviewView(),
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
	ui.initRefreshExecutor()
}

func (ui *UI) initRefreshExecutor() {
	ui.refreshQueue = make(chan func(), 10)
	go func() {
		for {
			task := <-ui.refreshQueue
			task()
		}
	}()
}

func (ui *UI) queueRefresh(f func()) {
	ui.refreshQueue <- f
}

func (ui *UI) askConfig() {
	openConfirmationDialog(func(b bool) {
		if !b {
			api().GlobalConfiguration.SetStandalone(true)
			ui.initStanialone()
			return
		}

		api().InitConfiguration()
		ui.init()
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
		Set('S', "Open settings", func() {
			openSettingsDialog()
		}).
		Set(']', "Cycle tab right", func() {
			ui.mainTabGroup.IncrementSelectedTab(func(tab *tui.Tab) {
				api().GlobalConfiguration.SetLastTab(tab.Frame.Name())
			})
		}).
		Set('[', "Cycle tab left", func() {
			ui.mainTabGroup.DecrementSelectedTab(func(tab *tui.Tab) {
				api().GlobalConfiguration.SetLastTab(tab.Frame.Name())
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

func getViewable[T viewable]() T {
	for _, v := range ui.views {
		if v, ok := v.(T); ok {
			return v
		}
	}

	panic("invalid view")
}

func renderView(v viewable) {
	tui.UpdateTui(func(g *tui.Tui) error {
		v.render()
		return nil
	})
}

func focusView(viewName string) {
	tui.GetView(viewName).Focus()

	views := make([]*tui.View, 0)
	for _, v := range ui.views {
		views = append(views, v.getView())
	}

	for _, v := range views {
		if v.Name() == viewName {
			v.FrameColor = onFrameColor
		} else {
			v.FrameColor = offFrameColor
		}
	}
}

func refresh(v viewable) {
	ui.queueRefresh(func() {
		v.refresh()
		renderView(v)
	})
}

func api() *core.Api {
	return ui.api
}

func runAction(f func()) {
	tui.Suspend()
	f()
	tui.Resume()
	refreshMainViews()
	refreshTmuxViews()
}

func styleView(v *tui.View) {
	v.FrameRunes = tui.ThickFrame
	v.TitleColor = gocui.AttrBold | gocui.ColorYellow
}

func systemUpdate() bool {
	if api().LocalConfiguration.IsInitialized && !api().GlobalConfiguration.IsUpdateAsked() {
		api().GlobalConfiguration.SetUpdateAsked()
		update, newTag := api().GlobalConfiguration.DetectUpdate()
		if update {
			openConfirmationDialog(func(b bool) {
				if b {
					runAction(func() {
						api().GlobalConfiguration.UpdateMynav()
					})
				}
			}, "A new update of mynav is available! Would you like to update to version "+newTag+"?")
			return true
		}
	}
	return false
}
