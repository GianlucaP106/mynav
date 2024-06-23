package ui

import (
	"mynav/pkg/system"
	"mynav/pkg/tmux"
	"strconv"

	"github.com/awesome-gocui/gocui"
)

type TmuxSessionView struct {
	tableRenderer *TableRenderer
	sessions      []*tmux.TmuxSession
	standalone    bool
}

const TmuxSessionViewName = "TmuxSessionView"

var _ View = &TmuxSessionView{}

func newTmuxSessionView() *TmuxSessionView {
	ts := &TmuxSessionView{}

	return ts
}

func (tv *TmuxSessionView) RequiresManager() bool {
	return false
}

func (tv *TmuxSessionView) Init(ui *UI) {
	var view *gocui.View
	if tv.standalone {
		screenX, screenY := ScreenSize()
		view = SetCenteredView(TmuxSessionViewName, screenX/2, screenY/3, 0)
	} else {
		view = SetViewLayout(tv.Name())
	}

	view.Title = withSurroundingSpaces("TMUX Sessions")
	view.TitleColor = gocui.ColorBlue
	view.FrameColor = gocui.ColorGreen

	sizeX, sizeY := view.Size()
	tv.tableRenderer = NewTableRenderer()
	titles := []string{
		"workspace",
		"Windows",
		"Session name",
	}
	proportions := []float64{
		0.2,
		0.2,
		0.6,
	}
	tv.tableRenderer.InitTable(sizeX, sizeY, titles, proportions)
	tv.refreshTmuxSessions()

	moveUp := func() {
		if !tv.standalone {
			ui.FocusWorkspacesView()
		}
	}

	moveLeft := func() {
		if !tv.standalone {
			ui.FocusPortView()
		}
	}

	moveRight := func() {
		if !tv.standalone {
			ui.FocusPrView()
		}
	}

	KeyBinding(TmuxSessionViewName).
		set(gocui.KeyEnter, func() {
			if system.IsTmuxSession() {
				GetDialog[*ToastDialog](ui).OpenError("You are already in a tmux session. Nested tmux sessions are not supported yet.")
				return
			}

			session := tv.getSelectedSession()
			ui.setAction(system.GetAttachTmuxSessionCmd(session.Name))
		}).
		set('D', func() {
			if Api().Tmux.GetTmuxSessionCount() == 0 {
				return
			}

			GetDialog[*ConfirmationDialog](ui).Open(func(b bool) {
				if b {
					session := tv.getSelectedSession()
					if err := Api().Tmux.DeleteTmuxSession(session); err != nil {
						GetDialog[*ToastDialog](ui).OpenError(err.Error())
						return
					}
					ui.RefreshMainView()
				}
			}, "Are you sure you want to delete this session?")
		}).
		set('X', func() {
			if Api().Tmux.GetTmuxSessionCount() == 0 {
				return
			}

			GetDialog[*ConfirmationDialog](ui).Open(func(b bool) {
				if b {
					if err := Api().Tmux.DeleteAllTmuxSessions(); err != nil {
						GetDialog[*ToastDialog](ui).OpenError(err.Error())
						return
					}
					ui.RefreshMainView()
				}
			}, "Are you sure you want to delete ALL tmux sessions?")
		}).
		set('W', func() {
			if tv.standalone || Api().Core.GetWorkspaceTmuxSessionCount() == 0 {
				return
			}

			GetDialog[*ConfirmationDialog](ui).Open(func(b bool) {
				if b {
					if err := Api().Core.DeleteAllWorkspaceTmuxSessions(); err != nil {
						GetDialog[*ToastDialog](ui).OpenError(err.Error())
						return
					}
					ui.RefreshMainView()
				}
			}, "Are you sure you want to delete ALL non-external tmux sessions?")
		}).
		set('j', func() {
			tv.tableRenderer.Down()
		}).
		set('k', func() {
			tv.tableRenderer.Up()
		}).
		set('a', func() {
			if system.IsTmuxSession() {
				return
			}
			GetDialog[*EditorDialog](ui).Open(func(s string) {
				ui.setAction(system.GetNewTmuxSessionCmd(s, "~"))
			}, func() {}, "New session name", Small)
		}).
		set('?', func() {
			GetDialog[*HelpView](ui).Open(getTmuxKeyBindings(tv.standalone), func() {})
		}).
		set(gocui.KeyEsc, moveUp).
		set(gocui.KeyArrowUp, moveUp).
		set(gocui.KeyCtrlK, moveUp).
		set(gocui.KeyArrowLeft, moveLeft).
		set(gocui.KeyCtrlH, moveLeft).
		set(gocui.KeyArrowRight, moveRight).
		set(gocui.KeyCtrlL, moveRight)

	if tv.standalone {
		FocusView(tv.Name())
	}
}

func (tv *TmuxSessionView) getSelectedSession() *tmux.TmuxSession {
	return tv.sessions[tv.tableRenderer.GetSelectedRowIndex()]
}

func (ts *TmuxSessionView) refreshTmuxSessions() {
	out := make([]*tmux.TmuxSession, 0)
	for _, session := range Api().Tmux.GetTmuxSessions() {
		out = append(out, session)
	}

	ts.sessions = out
	ts.syncSessionsToTable()
}

func (ts *TmuxSessionView) syncSessionsToTable() {
	rows := make([][]string, 0)
	for _, session := range ts.sessions {
		workspace := "external"
		if !ts.standalone {
			w := Api().Core.GetWorkspaceByTmuxSession(session)
			if w != nil {
				workspace = w.ShortPath()
			}
		}

		windows := strconv.Itoa(session.NumWindows)
		rows = append(rows, []string{
			workspace,
			windows,
			session.Name,
		})
	}
	ts.tableRenderer.FillTable(rows)
}

func (tv *TmuxSessionView) Name() string {
	return TmuxSessionViewName
}

func (tv *TmuxSessionView) Render(ui *UI) error {
	view := GetInternalView(tv.Name())
	if view == nil {
		tv.Init(ui)
		view = GetInternalView(tv.Name())
	}

	if ui.action.Command != nil {
		return gocui.ErrQuit
	}

	isViewFocused := false
	if fv := GetFocusedView(); fv != nil {
		isViewFocused = tv.Name() == GetFocusedView().Name()
	}

	view.Clear()
	tv.tableRenderer.RenderWithSelectCallBack(view, func(_ int, _ *TableRow) bool {
		return isViewFocused
	})

	return nil
}
