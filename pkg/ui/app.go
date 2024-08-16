package ui

import (
	"errors"
	"log"
	"mynav/pkg/api"

	"github.com/awesome-gocui/gocui"
)

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

func RenderView(v Viewable) {
	UpdateGui(func(g *Gui) error {
		v.Render()
		return nil
	})
}

func FocusView(viewName string) {
	SetFocusView(viewName)
	views := make([]*View, 0)
	for _, v := range _ui.views {
		views = append(views, v.View())
	}

	off := gocui.AttrDim | gocui.ColorWhite
	on := gocui.ColorWhite

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
