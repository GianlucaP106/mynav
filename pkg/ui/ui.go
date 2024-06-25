package ui

import (
	"errors"
	"fmt"
	"log"
	"mynav/pkg/system"

	"github.com/awesome-gocui/gocui"
)

type UI struct {
	WorkspacesView  *WorkspacesView
	TopicsView      *TopicsView
	PortsView       *PortView
	TmuxSessionView *TmuxSessionView
	GithubPrView    *GithubPrView
	HeaderView      *HeaderView

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

	_ui = InitViews(false, true)

	err := g.MainLoop()
	if err != nil {
		if !errors.Is(err, gocui.ErrQuit) {
			log.Panicln(err)
		}
	}

	return action
}

func InitViews(standalone bool, askToInit bool) *UI {
	ui := &UI{
		Standalone: standalone,
	}

	if ui.Standalone || system.IsCurrentProcessHomeDir() || (!Api().Core.IsConfigInitialized && !askToInit) {
		ui.Standalone = true
		ui.TmuxSessionView = NewTmuxSessionView(ui)
		SetScreenManagers([]gocui.Manager{
			gocui.ManagerFunc(func(g *gocui.Gui) error {
				return ui.TmuxSessionView.Render()
			}),
		}...)
		ui.TmuxSessionView.Init()
		FocusViewInternal(TmuxSessionViewName)
		SystemUpdate()
		setGlobalKeybindings()
		return ui
	}

	if !Api().Core.IsConfigInitialized {
		OpenConfirmationDialog(func(b bool) {
			if !b {
				InitViews(true, false)
				return
			}

			Api().InitConfiguration()
			InitViews(false, false)
		}, "No configuration found. Would you like to initialize this directory?")
		return nil
	}

	ui.TopicsView = NewTopicsView(ui)
	ui.WorkspacesView = NewWorkspcacesView(ui)
	ui.PortsView = NewPortView()
	ui.TmuxSessionView = NewTmuxSessionView(ui)
	ui.GithubPrView = NewGithubPrView()
	ui.HeaderView = NewHeaderView()

	SetScreenManagers([]gocui.Manager{
		gocui.ManagerFunc(func(g *gocui.Gui) error {
			return ui.TopicsView.Render()
		}),

		gocui.ManagerFunc(func(g *gocui.Gui) error {
			return ui.WorkspacesView.Render()
		}),

		gocui.ManagerFunc(func(g *gocui.Gui) error {
			return ui.PortsView.Render()
		}),

		gocui.ManagerFunc(func(g *gocui.Gui) error {
			return ui.TmuxSessionView.Render()
		}),

		gocui.ManagerFunc(func(g *gocui.Gui) error {
			return ui.GithubPrView.Render()
		}),

		gocui.ManagerFunc(func(g *gocui.Gui) error {
			return ui.HeaderView.Render()
		}),
	}...)

	ui.TopicsView.Init()
	ui.WorkspacesView.Init()
	ui.PortsView.Init()
	ui.TmuxSessionView.Init()
	ui.GithubPrView.Init()
	ui.HeaderView.Init()

	SystemUpdate()

	if Api().Core.GetSelectedWorkspace() != nil {
		FocusWorkspacesView()
	} else {
		FocusTopicsView()
	}

	setGlobalKeybindings()

	return ui
}

func setGlobalKeybindings() {
	quit := func() bool {
		return true
	}
	KeyBinding("").
		setWithQuit(gocui.KeyCtrlC, quit).
		setWithQuit('q', quit).
		setWithQuit('q', quit).
		set('t', func() {
			FocusTmuxView()
		}).
		set('?', func() {
			OpenHelpView(nil, func() {})
		})
}

func FocusTopicsView() {
	FocusView(TopicViewName)
}

func FocusWorkspacesView() {
	FocusView(WorkspacesViewName)
}

func FocusPortView() {
	FocusView(PortViewName)
}

func FocusTmuxView() {
	FocusView(TmuxSessionViewName)
}

func FocusPrView() {
	FocusView(GithubPrViewName)
}

func FocusView(viewName string) {
	FocusViewInternal(viewName)

	wv := GetView(WorkspacesViewName)
	tv := GetView(TopicViewName)
	pv := GetView(PortViewName)
	tmv := GetView(TmuxSessionViewName)
	gprv := GetView(GithubPrViewName)
	views := []*View{wv, tv, pv, tmv, gprv}

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

func (ui *UI) RefreshAllViews() {
	if !ui.Standalone {
		ui.TopicsView.refreshTopics()
		ui.PortsView.refreshPorts()
		ui.WorkspacesView.refreshWorkspaces()
	}
	ui.TmuxSessionView.refreshTmuxSessions()
}

func SystemUpdate() bool {
	if Api().Core.IsConfigInitialized && !Api().Core.IsUpdateAsked() {
		Api().Core.SetUpdateAsked()
		update, newTag := Api().Core.DetectUpdate()
		if update {
			OpenConfirmationDialog(func(b bool) {
				if b {
					SetActionEnd(system.GetUpdateSystemCmd())
				}
			}, "A new update of mynav is available! Would you like to update to version "+newTag+"?")
			return true
		}
	}
	return false
}

func SetViewLayout(viewName string) *View {
	maxX, maxY := ScreenSize()
	views := map[string]func() *View{}
	views[WorkspacesViewName] = func() *View {
		view, _ := SetView(WorkspacesViewName, (maxX/3)+1, 8, maxX-2, (maxY / 2), 0)
		return view
	}

	views[TmuxSessionViewName] = func() *View {
		view, _ := SetView(TmuxSessionViewName, (maxX/3)+1, (maxY/2)+1, ((2*maxX)/3)-1, maxY-4, 0)
		return view
	}

	views[TopicViewName] = func() *View {
		view, _ := SetView(TopicViewName, 2, 8, maxX/3-1, (maxY / 2), 0)
		return view
	}

	views[PortViewName] = func() *View {
		view, _ := SetView(PortViewName, 2, (maxY/2)+1, maxX/3-1, maxY-4, 0)
		return view
	}

	views[GithubPrViewName] = func() *View {
		view, _ := SetView(GithubPrViewName, ((2*maxX)/3)+1, (maxY/2)+1, maxX-2, maxY-4, 0)
		return view
	}

	views[HeaderViewName] = func() *View {
		view, _ := SetView(HeaderViewName, 2, 1, maxX-2, 5, 0)
		return view
	}

	f := views[viewName]
	if f == nil {
		log.Panicln("invalid view")
	}

	return f()
}
