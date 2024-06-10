package ui

import (
	"fmt"
	"mynav/pkg/api"
	"mynav/pkg/utils"

	"github.com/awesome-gocui/gocui"
	"github.com/gookit/color"
)

const PortViewName = "PortView"

type PortView struct {
	listRenderer *ListRenderer
	ports        []*api.Port
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
	ports := Api().GetPorts().ToList().Sorted()
	pv.ports = ports

	newListSize := len(pv.ports)
	if pv.listRenderer != nil && newListSize != pv.listRenderer.listSize {
		pv.listRenderer.setListSize(newListSize)
	}
}

func (p *PortView) getSelectedPort() *api.Port {
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

	KeyBinding(p.Name()).
		set('j', func() {
			p.listRenderer.increment()
		}).
		set('k', func() {
			p.listRenderer.decrement()
		}).
		set(gocui.KeyArrowUp, func() {
			ui.FocusTopicsView()
		}).
		set(gocui.KeyEsc, func() {
			ui.FocusTopicsView()
		}).
		set(gocui.KeyEnter, func() {
			port := p.getSelectedPort()
			if port.TmuxSession != nil {
				ui.setAction(utils.AttachTmuxSessionCmd(port.TmuxSession.Name))
			}
		}).
		set('D', func() {
			port := p.getSelectedPort()
			if port != nil {
				GetDialog[*ConfirmationDialog](ui).Open(func(b bool) {
					if b {
						if err := Api().KillPort(port); err != nil {
							GetDialog[*ToastDialog](ui).Open(err.Error(), func() {
								ui.FocusPortView()
							})
						}
						p.refreshPorts()
					}
					ui.FocusPortView()
				}, "Are you sure you want to kill this port?")
			}
		}).
		set('?', func() {
			GetDialog[*HelpView](ui).Open(getKeyBindings(p.Name()), func() {
				ui.FocusPortView()
			})
		})
}

func (p *PortView) formatPort(port *api.Port, selected bool) string {
	view := GetInternalView(p.Name())
	sizeX, _ := view.Size()

	fifth := sizeX / 5

	exeLine := withSpacePadding(port.GetExeName(), fifth)

	portNumber := port.GetPortStr()
	portLine := withSpacePadding(portNumber, fifth)

	tmuxSessionLine := ""
	tmuxSize := fifth*3 + 5
	tmuxContent := ""
	if port.TmuxSession != nil {
		workspace := Api().GetWorkspaceByTmuxSession(port.TmuxSession)
		if workspace != nil {
			tmuxContent = "workspace: " + workspace.Name
		} else {
			tmuxContent = "tmux: " + port.TmuxSession.Name
		}
	} else {
		tmuxContent = "external"
	}
	tmuxSessionLine = withSpacePadding(tmuxContent, tmuxSize)

	line := portLine + exeLine + tmuxSessionLine
	if selected {
		line = color.New(color.BgCyan, color.Black).Sprint(line)
	} else {
		line = color.New(color.White).Sprint(line)
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
	return color.Blue.Sprint(portTitle + pidTitle + tmuxTitle)
}

func (p *PortView) Render(ui *UI) error {
	p.Init(ui)
	go func() {
		gui.Update(func(g *gocui.Gui) error {
			view := GetInternalView(p.Name())
			if p.ports == nil {
				p.refreshPorts()
			}

			currentViewSelected := GetFocusedView().Name() == p.Name()
			view.Clear()
			content := make([]string, 0)
			p.listRenderer.forEach(func(idx int) {
				port := p.ports[idx]
				line := p.formatPort(port, (idx == p.listRenderer.selected) && currentViewSelected)
				content = append(content, line)
			})

			fmt.Fprintln(view, p.formatPortListTitle())
			for _, line := range content {
				fmt.Fprintln(view, line)
			}

			return nil
		})
	}()

	if ui.action.Command != nil {
		return gocui.ErrQuit
	}

	return nil
}
