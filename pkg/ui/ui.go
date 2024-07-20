package ui

import (
	"errors"
	"log"
	"mynav/pkg/api"
	"mynav/pkg/constants"

	"github.com/awesome-gocui/gocui"
)

type Viewable interface {
	Init()
	View() *View
	Render() error
}

type UI struct {
	mainTabGroup *TabGroup
	api          *api.Api
	views        []Viewable
}

var _ui *UI

func Start(api *api.Api) {
	g := NewGui()
	defer g.Close()

	_ui = &UI{
		views: make([]Viewable, 0),
		api:   api,
	}

	if Api().Configuration.IsConfigInitialized {
		_ui.InitUI()
	} else if Api().Configuration.Standalone {
		_ui.InitStandaloneUI()
	} else {
		_ui.AskConfig()
	}

	err := g.MainLoop()
	if err != nil {
		if !errors.Is(err, gocui.ErrQuit) {
			log.Panicln(err)
		}
	}
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

func (ui *UI) InitUI() *UI {
	ui.views = []Viewable{
		NewHeaderView(),
		NewTopicsView(),
		NewWorkspcacesView(),
		NewPortView(),
		NewTmuxSessionView(),
		NewGithubPrView(),
		NewGithubRepoView(),
	}

	SetViewManagers(ui.views)
	InitViewables(ui.views)

	tab1 := NewTab("main", GetTopicsView().View().Name())
	tab1.AddView(GetHeaderView(), None)
	tab1.AddView(GetTopicsView(), TopLeft)
	tab1.AddView(GetWorkspacesView(), TopRight)
	tab1.GenerateNavigationKeyBindings()

	tab2 := NewTab("tmux", GetTmuxSessionView().View().Name())
	tab2.AddView(GetTmuxSessionView(), None)
	tab2.AddView(GetHeaderView(), None)

	tab3 := NewTab("system", GetPortView().View().Name())
	tab3.AddView(GetPortView(), None)
	tab3.AddView(GetHeaderView(), None)

	tab4 := NewTab("github", GetGithubPrView().View().Name())
	tab4.AddView(GetGithubRepoView(), TopLeft)
	tab4.AddView(GetGithubPrView(), TopRight)
	tab4.AddView(GetHeaderView(), None)

	ui.mainTabGroup = NewTabGroup([]*Tab{
		tab1,
		tab2,
		tab3,
		tab4,
	})

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
	tmv := NewTmuxSessionView()
	ui.views = make([]Viewable, 0)
	ui.views = append(ui.views, tmv)

	SetViewManagers(ui.views)
	InitViewables(ui.views)

	tab := NewTab("tab1", constants.TmuxSessionViewName)
	tab.AddView(tmv, None)
	ui.mainTabGroup = NewTabGroup([]*Tab{tab})
	ui.mainTabGroup.FocusTabByIndex(0)

	SystemUpdate()
	ui.InitGlobalKeybindings()
}

func (ui *UI) InitGlobalKeybindings() {
	quit := func() bool {
		return true
	}
	NewKeybindingBuilder(nil).
		setWithQuit(gocui.KeyCtrlC, quit, "Quit").
		setWithQuit('q', quit, "Quit").
		setWithQuit('q', quit, "Quit").
		set(']', func() {
			ui.mainTabGroup.IncrementSelectedTab(func(tab *Tab) {
				Api().Configuration.SetLastTab(tab.Frame.Name())
			})
		}, "Cycle tab right").
		set('[', func() {
			ui.mainTabGroup.DecrementSelectedTab(func(tab *Tab) {
				Api().Configuration.SetLastTab(tab.Frame.Name())
			})
		}, "Cycle tab left").
		set('?', func() {
			OpenHelpView(nil, func() {})
		}, "Toggle cheatsheet")
}

func SetViewManagers(vs []Viewable) {
	managers := make([]gocui.Manager, 0)
	for _, view := range vs {
		managers = append(managers, gocui.ManagerFunc(func(_ *gocui.Gui) error {
			return view.Render()
		}))
	}

	SetManagerFunctions(managers...)
}

func InitViewables(vs []Viewable) {
	for _, v := range vs {
		v.Init()
	}
}

func GetViewable[T Viewable]() T {
	for _, v := range _ui.views {
		if v, ok := v.(T); ok {
			return v
		}
	}

	panic("invalid view")
}

func FocusView(viewName string) {
	SetFocusView(viewName)
	views := make([]*View, 0)
	for _, v := range _ui.views {
		views = append(views, v.View())
	}

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

func GetMainTabGroup() *TabGroup {
	return _ui.mainTabGroup
}

func Api() *api.Api {
	return _ui.api
}

func RunAction(action func()) {
	gocui.Suspend()
	action()
	gocui.Resume()
}
