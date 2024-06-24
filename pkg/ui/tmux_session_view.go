package ui

import (
	"mynav/pkg/system"
	"mynav/pkg/tmux"
	"strconv"

	"github.com/awesome-gocui/gocui"
)

type TmuxSessionView struct {
	view          *View
	ui            *UI
	tableRenderer *TableRenderer
	sessions      []*tmux.TmuxSession
}

const TmuxSessionViewName = "TmuxSessionView"

func NewTmuxSessionView(ui *UI) *TmuxSessionView {
	return &TmuxSessionView{
		ui: ui,
	}
}

func (tv *TmuxSessionView) Init() {
	if tv.ui.Standalone {
		screenX, screenY := ScreenSize()
		tv.view = SetCenteredView(TmuxSessionViewName, screenX/2, screenY/3, 0)
	} else {
		tv.view = SetViewLayout(TmuxSessionViewName)
	}

	tv.view.Title = withSurroundingSpaces("TMUX Sessions")
	tv.view.TitleColor = gocui.ColorBlue
	tv.view.FrameColor = gocui.ColorGreen

	sizeX, sizeY := tv.view.Size()
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
		if !tv.ui.Standalone {
			FocusWorkspacesView()
		}
	}

	moveLeft := func() {
		if !tv.ui.Standalone {
			FocusPortView()
		}
	}

	moveRight := func() {
		if !tv.ui.Standalone {
			FocusPrView()
		}
	}

	KeyBinding(TmuxSessionViewName).
		setWithQuit(gocui.KeyEnter, func() bool {
			if system.IsTmuxSession() {
				OpenToastDialogError("You are already in a tmux session. Nested tmux sessions are not supported yet.")
				return false
			}

			session := tv.getSelectedSession()
			SetAction(system.GetAttachTmuxSessionCmd(session.Name))
			return true
		}).
		set('D', func() {
			if Api().Tmux.GetTmuxSessionCount() == 0 {
				return
			}

			OpenConfirmationDialog(func(b bool) {
				if b {
					session := tv.getSelectedSession()
					if err := Api().Tmux.DeleteTmuxSession(session); err != nil {
						OpenToastDialogError(err.Error())
						return
					}
					tv.ui.RefreshAllViews()
				}
			}, "Are you sure you want to delete this session?")
		}).
		set('X', func() {
			if Api().Tmux.GetTmuxSessionCount() == 0 {
				return
			}

			OpenConfirmationDialog(func(b bool) {
				if b {
					if err := Api().Tmux.DeleteAllTmuxSessions(); err != nil {
						OpenToastDialogError(err.Error())
						return
					}
					tv.ui.RefreshAllViews()
				}
			}, "Are you sure you want to delete ALL tmux sessions?")
		}).
		set('W', func() {
			if tv.ui.Standalone || Api().Core.GetWorkspaceTmuxSessionCount() == 0 {
				return
			}

			OpenConfirmationDialog(func(b bool) {
				if b {
					if err := Api().Core.DeleteAllWorkspaceTmuxSessions(); err != nil {
						OpenToastDialogError(err.Error())
						return
					}
					tv.ui.RefreshAllViews()
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
			OpenEditorDialog(func(s string) {
				SetAction(system.GetNewTmuxSessionCmd(s, "~"))
			}, func() {}, "New session name", Small)
		}).
		set('?', func() {
			OpenHelpView(getTmuxKeyBindings(tv.ui.Standalone), func() {})
		}).
		set(gocui.KeyEsc, moveUp).
		set(gocui.KeyArrowUp, moveUp).
		set(gocui.KeyCtrlK, moveUp).
		set(gocui.KeyArrowLeft, moveLeft).
		set(gocui.KeyCtrlH, moveLeft).
		set(gocui.KeyArrowRight, moveRight).
		set(gocui.KeyCtrlL, moveRight)
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

func (tv *TmuxSessionView) syncSessionsToTable() {
	rows := make([][]string, 0)
	for _, session := range tv.sessions {
		workspace := "external"
		if !tv.ui.Standalone {
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
	tv.tableRenderer.FillTable(rows)
}

func (tv *TmuxSessionView) Render() error {
	isViewFocused := false
	if fv := GetFocusedView(); fv != nil {
		isViewFocused = tv.view.Name() == GetFocusedView().Name()
	}

	tv.view.Clear()
	tv.tableRenderer.RenderWithSelectCallBack(tv.view, func(_ int, _ *TableRow) bool {
		return isViewFocused
	})

	return nil
}
