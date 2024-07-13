package ui

import (
	"errors"
	"fmt"
	"log"
	"mynav/pkg/system"

	"github.com/awesome-gocui/gocui"
)

type Viewable interface {
	Init()
	View() *View
	Render() error
}

type UI struct {
	MainTabGroup *TabGroup
	Views        []Viewable

	// TODO: move to backend
	Standalone bool
}

var _ui *UI

func Start() *Action {
	if err := InitApi(); err != nil {
		fmt.Println(err.Error())
		return nil
	}

	g := NewGui()
	defer g.Close()

	_ui = &UI{
		Views: make([]Viewable, 0),
	}

	InitViews(_ui, false, true)

	err := g.MainLoop()
	if err != nil {
		if !errors.Is(err, gocui.ErrQuit) {
			log.Panicln(err)
		}
	}

	return action
}

func InitViews(ui *UI, standalone bool, askToInit bool) *UI {
	ui.Standalone = standalone

	setGlobalKeybindings := func() {
		quit := func() bool {
			return true
		}
		NewKeybindingBuilder("").
			setWithQuit(gocui.KeyCtrlC, quit).
			setWithQuit('q', quit).
			setWithQuit('q', quit).
			set(']', func() {
				ui.MainTabGroup.IncrementSelectedTab(func(tab *Tab) {
					Api().Core.SetLastTab(tab.Frame.Name())
				})
			}).
			set('[', func() {
				ui.MainTabGroup.DecrementSelectedTab(func(tab *Tab) {
					Api().Core.SetLastTab(tab.Frame.Name())
				})
			}).
			set('?', func() {
				OpenHelpView(nil, func() {})
			})
	}

	// TODO: refactor standlone configuration
	if ui.Standalone || system.IsCurrentProcessHomeDir() || (!Api().Core.IsConfigInitialized && !askToInit) {
		ui.Standalone = true
		tmv := NewTmuxSessionView()
		ui.Views = make([]Viewable, 0)
		ui.Views = append(ui.Views, tmv)

		SetViewManagers(ui.Views)
		InitViewables(ui.Views)

		tab := NewTab("tab1", TmuxSessionViewName)
		tab.AddView(tmv)
		ui.MainTabGroup = NewTabGroup([]*Tab{tab})
		ui.MainTabGroup.FocusTabByIndex(0)

		SystemUpdate()
		setGlobalKeybindings()
		return ui
	}

	if !Api().Core.IsConfigInitialized {
		OpenConfirmationDialog(func(b bool) {
			if !b {
				InitViews(ui, true, false)
				return
			}

			Api().InitConfiguration()
			InitViews(ui, false, false)
		}, "No configuration found. Would you like to initialize this directory?")
		return nil
	}

	ui.Views = []Viewable{
		NewHeaderView(),
		NewTopicsView(),
		NewWorkspcacesView(),
		NewPortView(),
		NewTmuxSessionView(),
		NewGithubPrView(),
		NewGithubRepoView(),
	}

	SetViewManagers(ui.Views)
	InitViewables(ui.Views)

	tab1 := NewTab("main", GetTopicsView().View().Name())
	tab1.AddView(GetHeaderView())
	tab1.AddView(GetTopicsView())
	tab1.AddView(GetWorkspacesView())

	tab2 := NewTab("tmux", GetTmuxSessionView().View().Name())
	tab2.AddView(GetTmuxSessionView())
	tab2.AddView(GetHeaderView())

	tab3 := NewTab("system", GetPortView().View().Name())
	tab3.AddView(GetPortView())
	tab3.AddView(GetHeaderView())

	tab4 := NewTab("github", GetGithubPrView().View().Name())
	tab4.AddView(GetGithubPrView())
	tab4.AddView(GetGithubRepoView())
	tab4.AddView(GetHeaderView())

	ui.MainTabGroup = NewTabGroup([]*Tab{
		tab1,
		tab2,
		tab3,
		tab4,
	})

	ui.MainTabGroup.FocusTabByIndex(0)

	SystemUpdate()

	if Api().Core.GetLastTab() != "" {
		ui.MainTabGroup.FocusTab(Api().Core.GetLastTab())
	} else if Api().Core.GetSelectedWorkspace() != nil {
		GetWorkspacesView().Focus()
	} else {
		GetTopicsView().Focus()
	}

	setGlobalKeybindings()

	return ui
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
	for _, v := range _ui.Views {
		if v, ok := v.(T); ok {
			return v
		}
	}

	panic("invalid view")
}

func FocusView(viewName string) {
	SetFocusView(viewName)
	views := make([]*View, 0)
	for _, v := range _ui.Views {
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
	return _ui.MainTabGroup
}

func IsStandlaone() bool {
	return _ui.Standalone
}

func RefreshAllData() {
	if !_ui.Standalone {
		GetTopicsView().refreshTopics()
		GetPortView().refreshPorts()
		GetWorkspacesView().refreshWorkspaces()
	}
	GetTmuxSessionView().refreshTmuxSessions()
}
