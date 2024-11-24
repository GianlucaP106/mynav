package app

import (
	"fmt"
	"mynav/pkg/core"
	"mynav/pkg/system"
	"mynav/pkg/tui"
	"strconv"
	"strings"

	"github.com/gookit/color"
)

type WorkspaceInfo struct {
	view          *tui.View
	workspaceInfo *tui.TableRenderer[*core.Workspace]
	sessionInfo   *tui.TableRenderer[*core.Session]
}

func newWorkspaceInfo() *WorkspaceInfo {
	w := &WorkspaceInfo{}
	return w
}

func (w *WorkspaceInfo) init() {
	w.view = a.ui.SetView(getViewPosition(WorkspaceInfoView))
	w.view.Title = " Workspace "
	a.styleView(w.view)

	// workspace info table
	sizeX, sizeY := w.view.Size()
	w.workspaceInfo = tui.NewTableRenderer[*core.Workspace]()
	w.workspaceInfo.Init(sizeX, sizeY, []string{
		"Name",
		"Topic",
		"Last Modified",
		"Git Remote",
	}, []float64{
		0.20,
		0.20,
		0.20,
		0.40,
	})
	w.workspaceInfo.SetStyles([]color.Style{
		color.Primary.Style,
		color.Secondary.Style,
		color.Comment.Style,
		color.Question.Style,
	})

	// session info table
	w.sessionInfo = tui.NewTableRenderer[*core.Session]()
	w.sessionInfo.Init(sizeX, sizeY, []string{
		"Active Session",
		"Windows",
		"Panes",
		"Last Attached",
		"Created",
		"Running",
	}, []float64{
		0.20,
		0.10,
		0.10,
		0.20,
		0.20,
		0.20,
	})
	w.sessionInfo.SetStyles([]color.Style{
		color.Success.Style,
		color.New(color.Magenta, color.Bold),
		color.Success.Style,
		color.Comment.Style,
		color.Comment.Style,
		color.Secondary.Style,
	})
}

func (w *WorkspaceInfo) show(workspace *core.Workspace) {
	// clear and resize
	w.view.Clear()
	w.view = a.ui.SetView(getViewPosition(w.view.Name()))

	// workspace info
	remote, _ := workspace.GitRemote()
	if remote == "" {
		remote = "None"
	}
	row := [][]string{{
		workspace.Name,
		workspace.Topic.Name,
		workspace.LastModifiedTimeFormatted(),
		remote,
	}}
	w.workspaceInfo.Fill(row, []*core.Workspace{workspace})
	w.workspaceInfo.RenderSelect(w.view, func(i int, tr *tui.TableRow[*core.Workspace]) bool {
		return false
	})

	session := a.api.Session(workspace)
	if session == nil {
		return
	}

	// seperate with newline
	fmt.Fprintln(w.view)

	panes, _ := session.ListPanes()
	commands := []string{}
	for _, p := range panes {
		commands = append(commands, p.CurrentCommand)
	}

	// session info
	row2 := [][]string{{
		"Yes",
		strconv.Itoa(session.Windows),
		strconv.Itoa(len(panes)),
		system.UnixTime(session.LastAttached).Format(system.TimeFormat()),
		system.UnixTime(session.Created).Format(system.TimeFormat()),
		strings.Join(commands, ","),
	}}
	w.sessionInfo.Fill(row2, []*core.Session{session})
	w.sessionInfo.RenderSelect(w.view, func(i int, tr *tui.TableRow[*core.Session]) bool {
		return false
	})
}
