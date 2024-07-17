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
	screenX, screenY := ScreenSize()
	p.view = SetCenteredView(PortViewName, screenX/2, screenY/3, 0)

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
		Api().Tmux.SyncPorts()
		p.refreshPorts()
		UpdateGui(func(_ *Gui) error {
			p.Render()
			return nil
		})
	}()

	p.view.KeyBinding().
		set('j', func() {
			p.tableRenderer.Down()
		}, "Move down").
		set('k', func() {
			p.tableRenderer.Up()
		}, "Move up").
		set(gocui.KeyEnter, func() {
			port := p.getSelectedPort()
			if port == nil {
				return
			}

			if port.tmux == nil {
				return
			}

			RunAction(func() {
				Api().Tmux.AttachTmuxSession(port.tmux)
			})
		}, "Open associated tmux session (if it exists)").
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
		}, "Kill port").
		set('?', func() {
			OpenHelpView(p.view.keybindingInfo.toList(), func() {})
		}, "Toggle cheatsheet")
}

func (pv *PortView) refreshPorts() {
	ports := make([]*Port, 0)

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

	currentViewSelected := p.view.IsFocused()

	p.tableRenderer.RenderWithSelectCallBack(p.view, func(_ int, _ *TableRow) bool {
		return currentViewSelected
	})

	if p.ports == nil {
		fmt.Fprintln(p.view, "Loading...")
	}

	return nil
}
