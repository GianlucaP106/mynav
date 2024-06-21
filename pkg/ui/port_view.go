package ui

import (
	"fmt"
	"mynav/pkg/system"

	"github.com/awesome-gocui/gocui"
	"github.com/gookit/color"
)

const PortViewName = "PortView"

type PortView struct {
	listRenderer  *ListRenderer
	ports         system.PortList
	externalPorts system.PortList
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
	externalPorts := make(system.PortList, 0)
	ports := make(system.PortList, 0)

	if len(Api().Port.GetPorts()) == 0 {
		Api().Tmux.SyncPorts()
	}

	for _, p := range Api().Port.GetPorts().ToList().Sorted() {
		if Api().Tmux.GetTmuxSessionByPort(p) != nil {
			ports = append(ports, p)
		} else {
			externalPorts = append(externalPorts, p)
		}
	}

	pv.ports = ports
	pv.externalPorts = externalPorts

	newListSize := len(pv.ports)
	if pv.listRenderer != nil && newListSize != pv.listRenderer.listSize {
		pv.listRenderer.setListSize(newListSize)
	}
}

func (p *PortView) getSelectedPort() *system.Port {
	if len(p.ports) <= 0 {
		return nil
	}

	return p.ports[p.listRenderer.selected]
}

func (p *PortView) Init(ui *UI) {
	if GetInternalView(p.Name()) != nil {
		return
	}

	view := SetViewLayout(p.Name())

	view.FrameColor = gocui.ColorBlue
	view.Title = withSurroundingSpaces("Open Ports")
	view.TitleColor = gocui.ColorBlue

	_, sizeY := view.Size()
	p.listRenderer = newListRenderer(0, sizeY, 0)

	moveUp := func() {
		ui.FocusTopicsView()
	}

	moveRight := func() {
		ui.FocusTmuxView()
	}

	KeyBinding(p.Name()).
		set('j', func() {
			p.listRenderer.increment()
		}).
		set('k', func() {
			p.listRenderer.decrement()
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

			if ts := Api().Tmux.GetTmuxSessionByPort(port); ts != nil {
				ui.setAction(system.GetAttachTmuxSessionCmd(ts.Name))
			}
		}).
		set('D', func() {
			port := p.getSelectedPort()
			if port != nil {
				GetDialog[*ConfirmationDialog](ui).Open(func(b bool) {
					if b {
						if err := Api().Port.KillPort(port); err != nil {
							GetDialog[*ToastDialog](ui).Open(err.Error(), func() {})
						}
						Api().Tmux.SyncPorts()
						p.refreshPorts()
					}
				}, "Are you sure you want to kill this port?")
			}
		}).
		set('?', func() {
			GetDialog[*HelpView](ui).Open(portKeyBindings, func() {})
		})
}

func (p *PortView) formatPort(port *system.Port, selected bool) string {
	view := GetInternalView(p.Name())
	sizeX, _ := view.Size()

	fifth := sizeX / 5

	exeLine := withSpacePadding(port.GetExeName(), fifth)

	portNumber := port.GetPortStr()
	portLine := withSpacePadding(portNumber, fifth)

	tmuxSessionLine := ""
	tmuxSize := fifth*3 + 5
	tmuxContent := ""
	if ts := Api().Tmux.GetTmuxSessionByPort(port); ts != nil {
		workspace := Api().Core.GetWorkspaceByTmuxSession(ts)
		if workspace != nil {
			tmuxContent = "workspace: " + workspace.ShortPath()
		} else {
			tmuxContent = "tmux: " + ts.Name
		}
	} else {
		tmuxContent = "external"
	}
	tmuxSessionLine = withSpacePadding(tmuxContent, tmuxSize)

	line := portLine + exeLine + tmuxSessionLine
	if selected {
		line = color.New(color.BgCyan, color.Black).Sprint(line)
	} else {
		line = color.New(color.Blue).Sprint(line)
	}

	return line
}

func (p *PortView) formatPortListTitle() string {
	view := GetInternalView(p.Name())
	sizeX, _ := view.Size()
	fifth := sizeX / 5
	portTitle := withSpacePadding("port", fifth)
	pidTitle := withSpacePadding("exe", fifth)
	tmuxTitle := withSpacePadding("linked to", (fifth*3)+5)
	return portTitle + pidTitle + tmuxTitle
}

func (p *PortView) Render(ui *UI) error {
	p.Init(ui)
	go func() {
		gui.Update(func(g *gocui.Gui) error {
			view := GetInternalView(p.Name())
			if p.ports == nil {
				p.refreshPorts()
			}

			sizeX, _ := view.Size()
			currentViewSelected := GetFocusedView().Name() == p.Name()
			view.Clear()
			content := make([]string, 0)
			p.listRenderer.forEach(func(idx int) {
				port := p.ports[idx]
				line := p.formatPort(port, (idx == p.listRenderer.selected) && currentViewSelected)
				content = append(content, line)
			})

			fmt.Fprintln(view, p.formatPortListTitle())
			fmt.Fprintln(view, withCharPadding("", sizeX, "-"))

			if len(content) > 0 {
				for _, line := range content {
					fmt.Fprintln(view, line)
				}
			} else {
				fmt.Fprintln(view, display("No workspace ports", Left, sizeX))
			}

			fmt.Fprintln(view, withCharPadding("", sizeX, "-"))

			for _, port := range p.externalPorts {
				fmt.Fprintln(view, p.formatPort(port, false))
			}

			return nil
		})
	}()

	if ui.action.Command != nil {
		return gocui.ErrQuit
	}

	return nil
}
