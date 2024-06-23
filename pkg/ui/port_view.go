package ui

import (
	"fmt"
	"mynav/pkg/system"
	"mynav/pkg/tmux"

	"github.com/awesome-gocui/gocui"
)

const PortViewName = "PortView"

type PortView struct {
	tableRenderer *TableRenderer
	ports         []*Port
}

type Port struct {
	tmux *tmux.TmuxSession
	*system.Port
}

var _ View = &PortView{}

func newPortView() *PortView {
	pv := &PortView{
		ports: nil,
	}
	return pv
}

func (p *PortView) Name() string {
	return PortViewName
}

func (p *PortView) RequiresManager() bool {
	return false
}

func (pv *PortView) refreshPorts() {
	ports := make([]*Port, 0)

	if len(Api().Port.GetPorts()) == 0 {
		Api().Tmux.SyncPorts()
	}

	for _, p := range Api().Port.GetPorts().ToList().Sorted() {
		if t := Api().Tmux.GetTmuxSessionByPort(p); t != nil {
			ports = append(ports, &Port{
				tmux: t,
				Port: p,
			})
		} else {
			ports = append(ports, &Port{
				tmux: nil,
				Port: p,
			})
		}
	}

	pv.ports = ports
	pv.syncPortsToTable()
}

func (pv *PortView) syncPortsToTable() {
	rows := make([][]string, 0)
	for _, p := range pv.ports {
		linkedTo := func() string {
			if p.tmux == nil {
				return "external"
			}

			workspace := Api().Core.GetWorkspaceByTmuxSession(p.tmux)
			if workspace != nil {
				return "workspace: " + workspace.ShortPath()
			} else {
				return "tmux: " + p.tmux.Name
			}
		}()

		rows = append(rows, []string{
			p.GetPortStr(),
			p.Exe,
			linkedTo,
		})
	}

	pv.tableRenderer.FillTable(rows)
}

func (p *PortView) getSelectedPort() *Port {
	if len(p.ports) <= 0 {
		return nil
	}

	return p.ports[p.tableRenderer.GetSelectedRowIndex()]
}

func (p *PortView) Init(ui *UI) {
	if GetInternalView(p.Name()) != nil {
		return
	}

	view := SetViewLayout(p.Name())

	view.FrameColor = gocui.ColorBlue
	view.Title = withSurroundingSpaces("Open Ports")
	view.TitleColor = gocui.ColorBlue

	sizeX, sizeY := view.Size()
	p.tableRenderer = NewTableRenderer()
	p.tableRenderer.InitTable(
		sizeX,
		sizeY,
		[]string{
			"Port",
			"Exe",
			"Linked to",
		},
		[]float64{
			0.25,
			0.25,
			0.5,
		})

	go func() {
		p.refreshPorts()
		UpdateGui(func(_ *gocui.Gui) error {
			p.Render(ui)
			return nil
		})
	}()

	moveUp := func() {
		ui.FocusTopicsView()
	}

	moveRight := func() {
		ui.FocusTmuxView()
	}

	KeyBinding(p.Name()).
		set('j', func() {
			p.tableRenderer.Down()
		}).
		set('k', func() {
			p.tableRenderer.Up()
		}).
		set(gocui.KeyArrowUp, moveUp).
		set(gocui.KeyCtrlK, moveUp).
		set(gocui.KeyEsc, moveUp).
		set(gocui.KeyArrowRight, moveRight).
		set(gocui.KeyCtrlL, moveRight).
		set(gocui.KeyEnter, func() {
			port := p.getSelectedPort()
			if port == nil {
				return
			}

			if port.tmux == nil {
				return
			}

			ui.setAction(system.GetAttachTmuxSessionCmd(port.tmux.Name))
		}).
		set('D', func() {
			port := p.getSelectedPort()
			if port == nil {
				return
			}

			if port.tmux == nil {
				GetDialog[*ToastDialog](ui).OpenError("Operation not allowed on external port")
				return
			}

			GetDialog[*ConfirmationDialog](ui).Open(func(b bool) {
				if b {
					if err := Api().Port.KillPort(port.Port); err != nil {
						GetDialog[*ToastDialog](ui).OpenError(err.Error())
					}
					Api().Tmux.SyncPorts()
					p.refreshPorts()
				}
			}, "Are you sure you want to kill this port?")
		}).
		set('?', func() {
			GetDialog[*HelpView](ui).Open(portKeyBindings, func() {})
		})
}

func (p *PortView) Render(ui *UI) error {
	p.Init(ui)
	view := GetInternalView(p.Name())
	view.Clear()

	currentViewSelected := false
	if v := GetFocusedView(); v != nil && v.Name() == p.Name() {
		currentViewSelected = true
	}

	p.tableRenderer.RenderWithSelectCallBack(view, func(_ int, _ *TableRow) bool {
		return currentViewSelected
	})

	if p.ports == nil {
		fmt.Fprintln(view, "Loading...")
	}

	if ui.action.Command != nil {
		return gocui.ErrQuit
	}

	return nil
}
