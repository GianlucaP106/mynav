package ui

import (
	"errors"
	"fmt"
	"log"
	"mynav/pkg/system"

	"github.com/awesome-gocui/gocui"
)

type UI struct {
	views   map[string]View
	dialogs map[string]Dialog

	action *Action
}

func Start() *Action {
	if err := InitApi(); err != nil {
		fmt.Println(err.Error())
		return nil
	}

	g := NewGui()
	defer g.Close()
	ui := &UI{
		views:  map[string]View{},
		action: &Action{},
	}

	managers := ui.InitViews()
	managers = append(managers, ui.InitDialogs()...)
	SetScreenManagers(managers...)

	quit := func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	}
	KeyBinding("").
		setKeybinding("", gocui.KeyCtrlC, quit).
		setKeybinding("", 'q', quit).
		setKeybinding("", 'q', quit).
		set('t', func() {
			ui.FocusTmuxView()
		}).
		set('?', func() {
			GetDialog[*HelpView](ui).Open(nil, func() {})
		})

	err := g.MainLoop()
	if err != nil {
		if !errors.Is(err, gocui.ErrQuit) {
			log.Panicln(err)
		}
	}

	return ui.action
}

func (ui *UI) handleUpdate() bool {
	if Api().Core.IsConfigInitialized && !Api().Core.IsUpdateAsked() {
		Api().Core.SetUpdateAsked()
		update, newTag := Api().Core.DetectUpdate()
		if update {
			GetDialog[*ConfirmationDialog](ui).Open(func(b bool) {
				if b {
					ui.setActionEnd(system.GetUpdateSystemCmd())
				}
			}, "A new update of mynav is available! Would you like to update to version "+newTag+"?")
			return true
		}
	}
	return false
}

func SetViewLayout(viewName string) *gocui.View {
	maxX, maxY := ScreenSize()
	views := map[string]func() *gocui.View{}
	views[WorkspacesViewName] = func() *gocui.View {
		view, _ := SetView(WorkspacesViewName, (maxX/3)+1, 8, maxX-2, (maxY / 2), 0)
		return view
	}

	views[TmuxSessionViewName] = func() *gocui.View {
		view, _ := SetView(TmuxSessionViewName, (maxX/3)+1, (maxY/2)+1, maxX-2, maxY-4, 0)
		return view
	}

	views[TopicViewName] = func() *gocui.View {
		view, _ := SetView(TopicViewName, 2, 8, maxX/3-1, (maxY / 2), 0)
		return view
	}

	views[PortViewName] = func() *gocui.View {
		view, _ := SetView(PortViewName, 2, (maxY/2)+1, maxX/3-1, maxY-4, 0)
		return view
	}

	views[HeaderStateName] = func() *gocui.View {
		view, _ := SetView(HeaderStateName, 2, 1, maxX-2, 5, 0)
		return view
	}

	f := views[viewName]
	if f == nil {
		log.Panicln("invalid view")
	}

	return f()
}
