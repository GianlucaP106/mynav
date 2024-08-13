package ui

import (
	"fmt"
	"mynav/pkg/constants"
	"mynav/pkg/events"
	"mynav/pkg/system"

	"github.com/GianlucaP106/gotmux/gotmux"
	"github.com/awesome-gocui/gocui"
)

type PortView struct {
	view          *View
	tableRenderer *TableRenderer[*Port]
}

type Port struct {
	tmux *gotmux.Session
	*system.Port
}

var _ Viewable = new(PortView)

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
	p.view = GetViewPosition(constants.PortViewName).Set()

	p.view.FrameColor = gocui.ColorBlue
	p.view.Title = withSurroundingSpaces("Open Ports")
	p.view.TitleColor = gocui.ColorBlue

	sizeX, sizeY := p.view.Size()
	p.tableRenderer = NewTableRenderer[*Port]()
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

	events.AddEventListener(constants.PortChangeEventName, func(_ string) {
		p.refresh()
		RenderView(p)
	})

	events.Emit(constants.PortChangeEventName)

	p.view.KeyBinding().
		set('j', "Move down", func() {
			p.tableRenderer.Down()
		}).
		set('k', "Move up", func() {
			p.tableRenderer.Up()
		}).
		set(gocui.KeyEnter, "Open associated tmux session (if it exists)", func() {
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
		}).
		set('D', "Kill port", func() {
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
				}
			}, "Are you sure you want to kill this port?")
		}).
		set('?', "Toggle cheatsheet", func() {
			OpenHelpView(p.view.keybindingInfo.toList(), func() {})
		})
}

func (pv *PortView) refresh() {
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

	rows := make([][]string, 0)
	rowValues := make([]*Port, 0)
	for _, p := range ports {
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

		rowValues = append(rowValues, p)
		rows = append(rows, []string{
			p.GetPortStr(),
			p.GetExeName(),
			linkedTo,
		})
	}

	pv.tableRenderer.FillTable(rows, rowValues)
}

func (p *PortView) getSelectedPort() *Port {
	_, port := p.tableRenderer.GetSelectedRow()
	if port != nil {
		return *port
	}

	return nil
}

func (p *PortView) Render() error {
	p.view.Clear()

	currentViewSelected := p.view.IsFocused()
	p.tableRenderer.RenderWithSelectCallBack(p.view, func(_ int, _ *TableRow[*Port]) bool {
		return currentViewSelected
	})

	if p.tableRenderer.GetTableSize() == 0 {
		fmt.Fprintln(p.view, "Nothing to show")
	}

	return nil
}
