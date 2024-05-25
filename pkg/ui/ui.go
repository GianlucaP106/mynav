package ui

import (
	"errors"
	"log"

	"github.com/awesome-gocui/gocui"
)

type UI struct {
	views   map[string]View
	dialogs map[string]Dialog

	action *Action
}

func Start() *Action {
	InitApi()

	g := NewGui()
	defer g.Close()

	ui := &UI{
		views: map[string]View{},
	}

	managers := ui.InitViews()
	managers = append(managers, ui.InitDialogs()...)
	SetScreenManagers(managers...)
	for _, v := range ui.views {
		if v.RequiresManager() {
			v.Init(ui)
		}
	}

	quit := func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	}
	KeyBinding("").
		setKeybinding("", gocui.KeyCtrlC, quit).
		setKeybinding("", 'q', quit).
		set('?', func() {
			GetDialog[*HelpView](ui).Open(nil, func() {
				ui.FocusTopicsView()
			})
		})

	ui.handleUpdate()

	err := g.MainLoop()
	if err != nil {
		if !errors.Is(err, gocui.ErrQuit) {
			log.Panicln(err)
		}
	}

	return ui.action
}

func (ui *UI) handleUpdate() {
	if Api().IsConfigInitialized && !Api().IsUpdateAsked() {
		Api().SetUpdateAsked()
		update, newTag := Api().DetectUpdate()
		if update {
			GetDialog[*ConfirmationDialog](ui).Open(func(b bool) {
				if b {
					ui.setActionEnd(Api().GetUpdateSystemCmd())
				}
			}, "A new update of mynav is available! Would you like to update to version "+newTag+"?")
		}
	}
}

func SetViewLayout(viewName string) *gocui.View {
	maxX, maxY := ScreenSize()
	views := map[string]func() *gocui.View{}
	views[WorkspacesViewName] = func() *gocui.View {
		view, _ := SetView(WorkspacesViewName, (maxX)/3+1, 8, maxX-2, maxY-4, 0)
		return view
	}

	views[TopicViewName] = func() *gocui.View {
		view, _ := SetView(TopicViewName, 2, 8, maxX/3-1, maxY-4, 0)
		return view
	}

	views[HeaderStateName] = func() *gocui.View {
		view, _ := SetView(HeaderStateName, 2, 1, maxX-2, 5, 0)
		return view
	}

	views[TmuxSessionViewName] = func() *gocui.View {
		return SetCenteredView(TmuxSessionViewName, maxX/2, maxY/3, 0)
	}

	f := views[viewName]
	if f == nil {
		log.Panicln("invalid view")
	}

	return f()
}
