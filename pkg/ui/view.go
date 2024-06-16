package ui

import (
	"github.com/awesome-gocui/gocui"
)

type View interface {
	Name() string
	Render(ui *UI) error
	Init(ui *UI)
	RequiresManager() bool
}

func (ui *UI) InitViews() []gocui.Manager {
	tv := newTopicsView()
	wv := newWorkspacesView(tv)
	pv := newPortView()
	tmsv := newTmuxSessionView()
	ui.SetViews(
		newMainView(wv, tv, pv, tmsv),
		tv,
		wv,
		pv,
		tmsv,
		newHeaderView(),
	)

	managers := []gocui.Manager{}
	for _, view := range ui.views {
		if view.RequiresManager() {
			manFunc := func(_ *gocui.Gui) error {
				return view.Render(ui)
			}
			managers = append(managers, gocui.ManagerFunc(manFunc))
		}
	}

	return managers
}

func (ui *UI) SetView(v View) {
	ui.views[v.Name()] = v
}

func (ui *UI) SetViews(views ...View) {
	for _, v := range views {
		ui.SetView(v)
	}
}

func GetView[T View](ui *UI) T {
	for _, view := range ui.views {
		v, ok := view.(T)
		if ok {
			return v
		}
	}
	panic("invalid view type")
}
