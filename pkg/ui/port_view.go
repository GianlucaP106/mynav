package ui

import (
	"fmt"
	"mynav/pkg/constants"
	"mynav/pkg/events"
	"mynav/pkg/system"
	"mynav/pkg/tmux"

	"github.com/awesome-gocui/gocui"
)

type PortView struct {
	view          *View
	tableRenderer *TableRenderer[*Port]
}

type Port struct {
	tmux *tmux.TmuxSession
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
	screenX, screenY := ScreenSize()
	p.view = SetCenteredView(constants.PortViewName, screenX/2, screenY/3, 0)

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
		p.refreshPorts()
		RenderView(p)
	})

	events.Emit(constants.PortSyncNeededEventName)

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
