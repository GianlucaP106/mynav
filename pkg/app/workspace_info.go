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
		0.40,
		0.20,
	})
	w.workspaceInfo.SetStyles([]color.Style{
		workspaceNameColor,
		topicNameColor,
		timestampColor,
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
		sessionMarkerColor,
		alternateSessionMarkerColor,
		sessionMarkerColor,
		timestampColor,
		timestampColor,
		topicNameColor,
	})
}

func (w *WorkspaceInfo) show(workspace *core.Workspace) {
	if workspace == nil {
		w.workspaceInfo.Clear()
		w.sessionInfo.Clear()
		return
	}

	// workspace info
	remote, _ := workspace.GitRemote()
	if remote == "" {
		remote = "None"
	}

	timeStr := fmt.Sprintf("%s (%s)", workspace.LastModifiedTimeFormatted(), system.TimeAgo(workspace.LastModifiedTime()))
	row := [][]string{{
		workspace.Name,
		workspace.Topic.Name,
		timeStr,
		remote,
	}}
	w.workspaceInfo.Fill(row, []*core.Workspace{workspace})

	session := a.api.Session(workspace)
	if session == nil {
		w.sessionInfo.Clear()
		return
	}

	panes, _ := session.ListPanes()
	commands := []string{}
	for _, p := range panes {
		commands = append(commands, p.CurrentCommand)
	}

	// session info
	lastAttached := system.UnixTime(session.LastAttached)
	created := system.UnixTime(session.Created)
	row2 := [][]string{{
		"Yes",
		strconv.Itoa(session.Windows),
		strconv.Itoa(len(panes)),
		system.TimeAgo(lastAttached),
		created.Format(system.TimeFormat()),
		strings.Join(commands, ","),
	}}
	w.sessionInfo.Fill(row2, []*core.Session{session})
}

func (w *WorkspaceInfo) render() {
	w.view.Clear()
	a.ui.Resize(w.view, getViewPosition(w.view.Name()))

	if w.workspaceInfo.Size() == 0 {
		return
	}

	w.workspaceInfo.RenderTable(w.view, func(i int, tr *tui.TableRow[*core.Workspace]) bool {
		return false
	}, nil)

	if w.sessionInfo.Size() == 0 {
		return
	}

	fmt.Fprintln(w.view)
	w.sessionInfo.RenderTable(w.view, func(i int, tr *tui.TableRow[*core.Session]) bool {
		return false
	}, nil)
}
