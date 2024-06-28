package ui

import (
	"fmt"
	"mynav/pkg/system"
	"mynav/pkg/tmux"

	"github.com/awesome-gocui/gocui"
)

type PortView struct {
	view          *View
	tableRenderer *TableRenderer
	ports         []*Port
}

type Port struct {
	tmux *tmux.TmuxSession
	*system.Port
}

var _ Viewable = new(PortView)

const PortViewName = "PortView"

func NewPortView() *PortView {
	return &PortView{}
}

func GetPortView() *PortView {
	return GetViewable[*PortView]()
}

func (pv *PortView) Focus() {
	FocusView(pv.View().Name())
}

func (p *PortView) View() *View {
	return p.view
}

func (p *PortView) Init() {
	p.view = GetViewPosition(PortViewName).Set()

	p.view.FrameColor = gocui.ColorBlue
	p.view.Title = withSurroundingSpaces("Open Ports")
	p.view.TitleColor = gocui.ColorBlue

	sizeX, sizeY := p.view.Size()
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
		UpdateGui(func(_ *Gui) error {
			p.Render()
			return nil
		})
	}()

	moveUp := func() {
		GetTopicsView().Focus()
	}

	moveRight := func() {
		GetTmuxSessionView().Focus()
	}

	KeyBinding(p.view.Name()).
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
		setWithQuit(gocui.KeyEnter, func() bool {
			port := p.getSelectedPort()
			if port == nil {
				return false
			}

			if port.tmux == nil {
				return false
			}

			SetAction(tmux.GetAttachTmuxSessionCmd(port.tmux.Name))
			return true
		}).
		set('D', func() {
			port := p.getSelectedPort()
			if port == nil {
				return
			}

			if port.tmux == nil {
				OpenToastDialogError("Operation not allowed on external port")
				return
			}

			OpenConfirmationDialog(func(b bool) {
				if b {
					if err := Api().Port.KillPort(port.Port); err != nil {
						OpenToastDialogError(err.Error())
					}
					Api().Tmux.SyncPorts()
					p.refreshPorts()
				}
			}, "Are you sure you want to kill this port?")
		}).
		set('?', func() {
			OpenHelpView(portKeyBindings, func() {})
		})
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
			p.GetExeName(),
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

func (p *PortView) Render() error {
	p.view.Clear()

	currentViewSelected := IsViewFocused(p.view)

	p.tableRenderer.RenderWithSelectCallBack(p.view, func(_ int, _ *TableRow) bool {
		return currentViewSelected
	})

	if p.ports == nil {
		fmt.Fprintln(p.view, "Loading...")
	}

	return nil
}
