package ui

import (
	"errors"
	"log"
	"mynav/pkg/core"
	"mynav/pkg/tui"

	"github.com/awesome-gocui/gocui"
)

var _ui *UI

func Start(api *core.Api) {
	g := tui.NewTui()
	defer g.Close()

	_ui = &UI{
		views: make([]viewable, 0),
		api:   api,
	}

	if getApi().LocalConfiguration.IsConfigInitialized {
		_ui.InitUI()
	} else if getApi().GlobalConfiguration.Standalone {
		_ui.initStandaloneUI()
	} else {
		_ui.askConfig()
	}

	err := g.MainLoop()
	if err != nil {
		if !errors.Is(err, gocui.ErrQuit) {
			log.Panicln(err)
		}
	}
}

func getViewable[T viewable]() T {
	for _, v := range _ui.views {
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
	for _, v := range _ui.views {
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

func refreshAsync(v viewable) {
	go func() {
		v.refresh()
		renderView(v)
	}()
}

func getMainTabGroup() *tui.TabGroup {
	return _ui.mainTabGroup
}

func getApi() *core.Api {
	return _ui.api
}
