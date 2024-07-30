package ui

import (
	"fmt"
	"mynav/pkg/constants"
	"mynav/pkg/persistence"
	"mynav/pkg/tasks"
	"mynav/pkg/tmux"
	"strconv"

	"github.com/awesome-gocui/gocui"
)

type PsView struct {
	view          *View
	tableRenderer *TableRenderer[*tmux.TmuxPaneProcess]

	// tmp
	isLoading *persistence.Value[bool]
}

var _ Viewable = new(PsView)

func NewPsView() *PsView {
	return &PsView{
		isLoading: persistence.NewValue(false),
	}
}

func GetPsView() *PsView {
	return GetViewable[*PsView]()
}

func (p *PsView) Init() {
	p.view = GetViewPosition(constants.PsViewName).Set()

	p.view.FrameColor = gocui.ColorBlue
	p.view.Title = withSurroundingSpaces("Processes")
	p.view.TitleColor = gocui.ColorBlue

	p.tableRenderer = NewTableRenderer[*tmux.TmuxPaneProcess]()

	sizeX, sizeY := p.view.Size()
	p.tableRenderer.InitTable(sizeX, sizeY, []string{
		"Exe",
		"Session",
		"Pid",
	}, []float64{
		0.2,
		0.6,
		0.2,
	})

	tasks.QueueTask(func() {
		// TODO: move this to core, system or tmux
		p.isLoading.Set(true)
		processes := make([]*tmux.TmuxPaneProcess, 0)
		for _, ts := range Api().Tmux.GetTmuxSessions() {
			ps := Api().Tmux.GetTmuxSessionChildProcesses(ts)
			processes = append(processes, ps...)
		}

		rows := make([][]string, 0)
		for _, tpp := range processes {
			name, err := tpp.Process.Name()
			if err != nil {
				continue
			}

			pid := strconv.Itoa(int(tpp.Process.Pid))
			rows = append(rows, []string{
				name,
				tpp.Pane.Session.Name,
				pid,
			})
		}

		p.tableRenderer.FillTable(rows, processes)
		p.isLoading.Set(false)
		RenderView(p)
	})

	p.view.KeyBinding().
		set('?', "Toggle cheatsheet", func() {
			OpenHelpView(p.view.keybindingInfo.toList(), func() {})
		})
}

func (p *PsView) View() *View {
	return p.view
}

func (p *PsView) Render() error {
	p.view.Clear()

	isFocused := p.view.IsFocused()
	p.tableRenderer.RenderWithSelectCallBack(p.view, func(i int, tr *TableRow[*tmux.TmuxPaneProcess]) bool {
		return isFocused
	})

	if p.isLoading.Get() {
		fmt.Fprintln(p.view, "Loading...")
	} else if p.tableRenderer.GetTableSize() == 0 {
		fmt.Fprintln(p.view, "Nothing to show")
	}

	return nil
}
