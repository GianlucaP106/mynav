package ui

import (
	"fmt"
	"mynav/pkg/constants"
	"mynav/pkg/events"
	"mynav/pkg/system"
	"mynav/pkg/tui"

	"github.com/GianlucaP106/gotmux/gotmux"
	"github.com/awesome-gocui/gocui"
)

type portView struct {
	view          *tui.View
	tableRenderer *tui.TableRenderer[*port]
}

type port struct {
	tmux *gotmux.Session
	*system.Port
}

var _ viewable = new(portView)

func newPortView() *portView {
	return &portView{}
}

func getPortView() *portView {
	return getViewable[*portView]()
}

func (pv *portView) Focus() {
	focusView(pv.getView().Name())
}

func (p *portView) getView() *tui.View {
	return p.view
}

func (p *portView) init() {
	p.view = GetViewPosition(constants.PortViewName).Set()

	p.view.Title = tui.WithSurroundingSpaces("Open Ports")

	tui.StyleView(p.view)

	sizeX, sizeY := p.view.Size()
	p.tableRenderer = tui.NewTableRenderer[*port]()
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
		renderView(p)
	})

	events.Emit(constants.PortChangeEventName)

	p.view.KeyBinding().
		Set('j', "Move down", func() {
			p.tableRenderer.Down()
		}).
		Set('k', "Move up", func() {
			p.tableRenderer.Up()
		}).
		Set(gocui.KeyEnter, "Open associated tmux session (if it exists)", func() {
			port := p.getSelectedPort()
			if port == nil {
				return
			}

			if port.tmux == nil {
				return
			}

			tui.RunAction(func() {
				getApi().Tmux.AttachTmuxSession(port.tmux)
			})
		}).
		Set('D', "Kill port", func() {
			port := p.getSelectedPort()
			if port == nil {
				return
			}

			if port.tmux == nil {
				openToastDialogError("Operation not allowed on external port")
				return
			}

			openConfirmationDialog(func(b bool) {
				if b {
					if err := getApi().Port.KillPort(port.Port); err != nil {
						openToastDialogError(err.Error())
					}
				}
			}, "Are you sure you want to kill this port?")
		}).
		Set('?', "Toggle cheatsheet", func() {
			OpenHelpDialog(p.view.GetKeybindings(), func() {})
		})
}

func (pv *portView) refresh() {
	ports := make([]*port, 0)
	for _, p := range getApi().Port.GetPorts().ToList().Sorted() {
		if t := getApi().Tmux.GetTmuxSessionByPort(p); t != nil {
			ports = append(ports, &port{
				tmux: t,
				Port: p,
			})
		} else {
			ports = append(ports, &port{
				tmux: nil,
				Port: p,
			})
		}
	}

	rows := make([][]string, 0)
	rowValues := make([]*port, 0)
	for _, p := range ports {
		linkedTo := func() string {
			if p.tmux == nil {
				return "external"
			}

			workspace := getApi().Core.GetWorkspaceByTmuxSession(p.tmux)
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

func (p *portView) getSelectedPort() *port {
	_, port := p.tableRenderer.GetSelectedRow()
	if port != nil {
		return *port
	}

	return nil
}

func (p *portView) render() error {
	p.view.Clear()

	currentViewSelected := p.view.IsFocused()
	p.tableRenderer.RenderWithSelectCallBack(p.view, func(_ int, _ *tui.TableRow[*port]) bool {
		return currentViewSelected
	})

	if p.tableRenderer.GetTableSize() == 0 {
		fmt.Fprintln(p.view, "Nothing to show")
	}

	return nil
}
