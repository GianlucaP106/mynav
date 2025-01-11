package app

import (
	"fmt"
	"mynav/pkg/core"
	"mynav/pkg/system"
	"mynav/pkg/tui"
	"sort"
	"strconv"

	"github.com/awesome-gocui/gocui"
	"github.com/gookit/color"
)

// Sessions view displaying active workspace sessions.
type Sessions struct {
	view  *tui.View
	table *tui.TableRenderer[*core.Session]

	// loading flag to display loading (not atomic as it should only be touched in the mainloop)
	loading bool
}

func newSessionsView() *Sessions {
	s := &Sessions{}
	return s
}

func (s *Sessions) selected() *core.Session {
	_, session := s.table.SelectedRow()
	if session != nil {
		return *session
	}
	return nil
}

func (s *Sessions) selectSession(session *core.Session) {
	s.table.SelectRowByValue(func(session2 *core.Session) bool {
		return session2.Name == session.Name
	})
}

func (s *Sessions) getLoading() bool {
	return s.loading
}

func (s *Sessions) setLoading(b bool) {
	s.loading = b
}

func (s *Sessions) showInfo() {
	session := s.selected()
	if session == nil {
		a.info.show(nil)
		return
	}

	a.info.show(session.Workspace)
}

func (s *Sessions) refreshPreview() {
	session := s.selected()
	if session == nil {
		a.preview.refresh(nil)
		return
	}

	a.preview.refresh(session)
}

func (s *Sessions) focus() {
	a.focusView(s.view)
	s.refreshDown()
}

func (s *Sessions) refreshDown() {
	s.showInfo()
	a.worker.Queue(func() {
		s.refreshPreview()
		a.ui.Update(func() {
			a.preview.render()
		})
	})
}

func (s *Sessions) refresh() {
	sessions := a.api.AllSessions()

	// sort by last attached
	sort.Slice(sessions, func(i, j int) bool {
		t1 := system.UnixTime(sessions[i].LastAttached)
		t2 := system.UnixTime(sessions[j].LastAttached)
		return t1.After(t2)
	})

	// fill table
	rows := make([][]string, 0)
	for _, s := range sessions {
		timeStr := system.TimeAgo(system.UnixTime(s.LastAttached))
		rows = append(rows, []string{
			s.Workspace.ShortPath(),
			strconv.Itoa(s.Windows),
			timeStr,
		})
	}
	s.table.Fill(rows, sessions)
}

func (s *Sessions) render() {
	s.view.Clear()
	a.ui.Resize(s.view, getViewPosition(s.view.Name()))

	// update page row marker
	row, _ := s.table.SelectedRow()
	size := s.table.Size()
	s.view.Subtitle = fmt.Sprintf(" %d / %d ", min(row+1, size), size)

	if s.getLoading() {
		fmt.Fprintln(s.view, "Loading...")
		return
	}

	// renders table and updates the last modified time
	isFocused := a.ui.IsFocused(s.view)
	s.table.RenderTable(s.view, func(i int, tr *tui.TableRow[*core.Session]) bool {
		return isFocused
	}, func(i int, tr *tui.TableRow[*core.Session]) {
		newTime := system.TimeAgo(system.UnixTime(tr.Value.LastAttached))
		tr.Cols[len(tr.Cols)-1] = newTime
	})
}

func (s *Sessions) init() {
	s.view = a.ui.SetView(getViewPosition(SessionsView))
	s.view.Title = " Sessions "
	a.styleView(s.view)

	sizeX, sizeY := s.view.Size()
	titles := []string{
		"Workspace",
		"Windows",
		"Last Attached",
	}
	proportions := []float64{
		0.40,
		0.20,
		0.40,
	}
	styles := []color.Style{
		workspaceNameColor,
		sessionMarkerColor,
		timestampColor,
	}
	s.table = tui.NewTableRenderer[*core.Session]()
	s.table.Init(sizeX, sizeY, titles, proportions)
	s.table.SetStyles(styles)

	down := func() {
		s.table.Down()
		s.refreshDown()
	}
	up := func() {
		s.table.Up()
		s.refreshDown()
	}
	a.ui.KeyBinding(s.view).
		Set('j', "Move down", down).
		Set('k', "Move up", up).
		Set(gocui.KeyArrowDown, "Move down", down).
		Set(gocui.KeyArrowUp, "Move up", up).
		Set('g', "Go to top", func() {
			s.table.Top()
		}).
		Set('G', "Go to bottom", func() {
			s.table.Bottom()
		}).
		Set(gocui.KeyEnter, "Open Session", func() {
			session := s.selected()
			if session == nil {
				return
			}

			if core.IsTmuxSession() {
				toast("A tmux session is already active", toastWarn)
				return
			}

			var err error
			a.runAction(func() {
				err = a.api.OpenSession(session.Workspace)
			})
			if err != nil {
				toast(err.Error(), toastError)
			} else {
				toast("Detached from session "+session.Workspace.Name, toastInfo)
			}

			a.refresh(nil, nil, session)
		}).
		Set('D', "Kill session", func() {
			session := s.selected()
			if session == nil {
				return
			}
			alert(func(b bool) {
				if !b {
					return
				}

				if err := a.api.KillSession(session.Workspace); err != nil {
					toast(err.Error(), toastError)
					return
				}

				a.refresh(nil, nil, session)
				toast("Killed session "+session.Workspace.Name, toastInfo)
			}, fmt.Sprintf("Are you sure you want to delete session for %s?", session.Workspace.Name))
		}).
		Set('h', "Focus workspaces view", func() {
			a.workspaces.focus()
		}).
		Set(gocui.KeyArrowLeft, "Focus workspaces view", func() {
			a.workspaces.focus()
		}).
		Set('?', "Toggle cheatsheet", func() {
			help(s.view)
		})
}
