package app

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/GianlucaP106/mynav/pkg/core"
	"github.com/GianlucaP106/mynav/pkg/tui"
	"github.com/gookit/color"
)

type Info struct {
	view          *tui.View
	workspaceInfo *tui.TableRenderer[*core.Workspace]
	sessionInfo   *tui.TableRenderer[*core.Session]
}

func newInfo() *Info {
	w := &Info{}
	return w
}

func (w *Info) init() {
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
		"Last Active",
		"Session Started",
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

func (i *Info) showSession(session *core.Session) {
	if session == nil {
		i.sessionInfo.Clear()
		return
	}

	panes, _ := session.ListPanes()
	// name -> count
	commands := map[string]int{}
	for _, p := range panes {
		commands[p.CurrentCommand]++
	}
	commandsList := []string{}
	for c := range commands {
		commandsList = append(commandsList, c)
	}
	sort.Slice(commandsList, func(i, j int) bool {
		return commands[commandsList[i]] < commands[commandsList[j]]
	})
	commandStrs := []string{}
	for _, c := range commandsList {
		count := commands[c]
		commandStrs = append(commandStrs, fmt.Sprintf("%dx %s", count, c))
	}

	// session info
	activity := core.UnixTime(session.Activity)
	created := core.UnixTime(session.Created)
	row2 := &tui.TableRow[*core.Session]{
		Cols: []string{
			"Yes",
			strconv.Itoa(session.Windows),
			strconv.Itoa(len(panes)),
			core.TimeAgo(activity),
			core.TimeAgo(created),
			strings.Join(commandStrs, ", "),
		},
		Value: session,
	}
	i.sessionInfo.Fill([]*tui.TableRow[*core.Session]{row2})
}

func (i *Info) show(workspace *core.Workspace) {
	if workspace == nil {
		i.workspaceInfo.Clear()
		i.sessionInfo.Clear()
		return
	}

	// workspace info
	remote, _ := workspace.GitRemote()
	if remote == "" {
		remote = "None"
	}

	timeStr := fmt.Sprintf("%s (%s)", workspace.LastModified().Format(core.TimeFormat()), core.TimeAgo(workspace.LastModified()))
	row := &tui.TableRow[*core.Workspace]{
		Cols: []string{
			workspace.Name,
			workspace.Topic.Name,
			timeStr,
			remote,
		},
		Value: workspace,
	}
	i.workspaceInfo.Fill([]*tui.TableRow[*core.Workspace]{row})

	session := a.api.Session(workspace)
	i.showSession(session)
}

func (w *Info) render() {
	w.view.Clear()
	a.ui.Resize(w.view, getViewPosition(w.view.Name()))

	if w.workspaceInfo.Size() > 0 {
		w.workspaceInfo.RenderTable(w.view, func(i int, tr *tui.TableRow[*core.Workspace]) bool {
			return false
		}, nil)
		fmt.Fprintln(w.view)
	}

	if w.sessionInfo.Size() > 0 {
		w.sessionInfo.RenderTable(w.view, func(i int, tr *tui.TableRow[*core.Session]) bool {
			return false
		}, nil)
	}
}
