package app

import (
	"fmt"
	"mynav/pkg/core"
	"mynav/pkg/system"
	"mynav/pkg/tui"
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

func (s *Sessions) getLoading() bool {
	return s.loading
}

func (s *Sessions) setLoading(b bool) {
	s.loading = b
}

func (s *Sessions) show() {
	def := "No preview"
	session := s.selected()
	if session == nil {
		a.preview.show(def)
		return
	}

	p := a.api.WorkspacePreview(session.Workspace)
	if p == "" {
		a.preview.show(def)
		return
	}

	a.preview.show(p)
	a.info.show(session.Workspace)
}

func (s *Sessions) focus() {
	a.focusView(s.view)
	s.show()
}

func (s *Sessions) refreshAll() {
	a.refresh(nil, nil, false, true)
}

func (s *Sessions) refresh() {
	rows := make([][]string, 0)
	sessions := a.api.AllSessions()
	for _, s := range sessions {
		rows = append(rows, []string{
			s.Workspace.Name,
			strconv.Itoa(s.Windows),
			system.UnixTime(s.LastAttached).Format(system.TimeFormat()),
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

	isFocused := a.ui.IsFocused(s.view)
	s.table.RenderSelect(s.view, func(i int, tr *tui.TableRow[*core.Session]) bool {
		return isFocused
	})
}

func (s *Sessions) init() {
	s.view = a.ui.SetView(getViewPosition(SessionsView))
	s.view.Title = " Active Sessions "
	a.styleView(s.view)

	sizeX, sizeY := s.view.Size()
	titles := []string{
		"Workspace Name",
		"Windows",
		"Last Attached",
	}
	proportions := []float64{
		0.40,
		0.20,
		0.40,
	}
	styles := []color.Style{
		color.New(color.FgBlue, color.Bold),
		color.New(color.Magenta, color.Bold),
		color.New(color.FgDarkGray, color.OpItalic),
	}
	s.table = tui.NewTableRenderer[*core.Session]()
	s.table.Init(sizeX, sizeY, titles, proportions)
	s.table.SetStyles(styles)

	a.ui.KeyBinding(s.view).
		Set('j', "Move down", func() {
			s.table.Down()
			s.show()
		}).
		Set('k', "Move up", func() {
			s.table.Up()
			s.show()
		}).
		Set(gocui.KeyEnter, "Open Session", func() {
			session := s.selected()
			if session == nil {
				return
			}
			var err error
			a.runAction(func() {
				err = a.api.OpenWorkspace(session.Workspace)
			})
			if err != nil {
				toast(err.Error(), toastError)
			}
			s.refreshAll()
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

				s.refreshAll()
				toast("Killed session "+session.Workspace.Name, toastInfo)
			}, "Are you sure you want to delete the session?")
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
