package ui

import (
	"fmt"
	"mynav/pkg/tasks"
	"mynav/pkg/tui"
	"strconv"

	"github.com/shirou/gopsutil/process"
)

type psView struct {
	view          *tui.View
	tableRenderer *tui.TableRenderer[*process.Process]
	psProcessing  *tasks.Task
}

var _ viewable = new(psView)

func newPsView() *psView {
	return &psView{}
}

func getPsView() *psView {
	return getViewable[*psView]()
}

func (p *psView) init() {
	p.view = getViewPosition(PsView).Set()

	p.view.Title = tui.WithSurroundingSpaces("Processes")

	styleView(p.view)

	p.tableRenderer = tui.NewTableRenderer[*process.Process]()

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

	p.psProcessing = tasks.QueueTask(func() {
		rows := make([][]string, 0)
		processes := make([]*process.Process, 0)
		for _, ts := range getApi().Tmux.GetTmuxSessions() {
			ps := getApi().Tmux.GetTmuxSessionChildProcesses(ts)
			for _, proc := range ps {
				name, err := proc.Name()
				if err != nil {
					continue
				}

				pid := strconv.Itoa(int(proc.Pid))
				rows = append(rows, []string{
					name,
					ts.Name,
					pid,
				})
				processes = append(processes, proc)
			}
		}

		p.tableRenderer.FillTable(rows, processes)
		renderView(p)
	})

	p.view.KeyBinding().
		Set('j', "Move down", func() {
			p.tableRenderer.Down()
		}).
		Set('k', "Move up", func() {
			p.tableRenderer.Up()
		}).
		Set('?', "Toggle cheatsheet", func() {
			OpenHelpDialog(p.view.GetKeybindings(), func() {})
		})
}

func (p *psView) getView() *tui.View {
	return p.view
}

func (p *psView) render() error {
	p.view.Clear()
	isFocused := p.view.IsFocused()
	p.view.Resize(getViewPosition(p.view.Name()))
	p.tableRenderer.RenderWithSelectCallBack(p.view, func(i int, tr *tui.TableRow[*process.Process]) bool {
		return isFocused
	})

	if p.psProcessing.IsStarted() && !p.psProcessing.IsCompleted() {
		fmt.Fprintln(p.view, "Loading...")
	} else if p.tableRenderer.GetTableSize() == 0 {
		fmt.Fprintln(p.view, "Nothing to show")
	}

	return nil
}
