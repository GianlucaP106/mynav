package ui

import (
	"mynav/pkg/constants"
	"mynav/pkg/events"
	"mynav/pkg/tmux"
	"strconv"

	"github.com/awesome-gocui/gocui"
)

type TmuxSessionView struct {
	view          *View
	tableRenderer *TableRenderer[*tmux.TmuxSession]
}

var _ Viewable = new(TmuxSessionView)

func NewTmuxSessionView() *TmuxSessionView {
	return &TmuxSessionView{}
}

func GetTmuxSessionView() *TmuxSessionView {
	return GetViewable[*TmuxSessionView]()
}

func (tv *TmuxSessionView) View() *View {
	return tv.view
}

func (tv *TmuxSessionView) Focus() {
	FocusView(tv.View().Name())
}

func (tv *TmuxSessionView) Init() {
	screenX, screenY := ScreenSize()
	tv.view = SetCenteredView(constants.TmuxSessionViewName, screenX/2, screenY/3, 0)

	tv.view.Title = withSurroundingSpaces("TMUX Sessions")
	tv.view.TitleColor = gocui.ColorBlue
	tv.view.FrameColor = gocui.ColorGreen

	sizeX, sizeY := tv.view.Size()
	tv.tableRenderer = NewTableRenderer[*tmux.TmuxSession]()
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

	events.AddEventListener(constants.TmuxSessionChangeEventName, func(_ string) {
		tv.refreshTmuxSessions()
		RenderView(tv)
	})

	tv.refreshTmuxSessions()

	tv.view.KeyBinding().
		set(gocui.KeyEnter, func() {
			if tmux.IsTmuxSession() {
				OpenToastDialogError("You are already in a tmux session. Nested tmux sessions are not supported yet.")
				return
			}

			session := tv.getSelectedSession()
			RunAction(func() {
				Api().Tmux.AttachTmuxSession(session)
			})
		}, "Attach to session").
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
					events.Emit(constants.WorkspaceChangeEventName)
				}
			}, "Are you sure you want to delete this session?")
		}, "Delete session").
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
				}
			}, "Are you sure you want to delete ALL tmux sessions?")
		}, "Kill ALL tmux sessions").
		set('W', func() {
			if Api().Configuration.Standalone || Api().Core.GetWorkspaceTmuxSessionCount() == 0 {
				return
			}

			OpenConfirmationDialog(func(b bool) {
				if b {
					if err := Api().Core.DeleteAllWorkspaceTmuxSessions(); err != nil {
						OpenToastDialogError(err.Error())
						return
					}
				}
			}, "Are you sure you want to delete ALL non-external tmux sessions?")
		}, "Kill ALL non-external (has a workspace) tmux sessions").
		set('j', func() {
			tv.tableRenderer.Down()
		}, "Move down").
		set('k', func() {
			tv.tableRenderer.Up()
		}, "Move up").
		set('a', func() {
			if tmux.IsTmuxSession() {
				return
			}
			OpenEditorDialog(func(s string) {
				RunAction(func() {
					Api().Tmux.CreateAndAttachTmuxSession(s, "~")
				})
			}, func() {}, "New session name", Small)
		}, "New external session (not associated to a workspace)").
		set('?', func() {
			OpenHelpView(tv.view.keybindingInfo.toList(), func() {})
		}, "Toggle cheatsheet")
}

func (tv *TmuxSessionView) getSelectedSession() *tmux.TmuxSession {
	_, ts := tv.tableRenderer.GetSelectedRow()
	if ts != nil {
		return *ts
	}

	return nil
}

func (ts *TmuxSessionView) refreshTmuxSessions() {
	sessions := make([]*tmux.TmuxSession, 0)
	for _, session := range Api().Tmux.GetTmuxSessions() {
		sessions = append(sessions, session)
	}

	rows := make([][]string, 0)
	for _, session := range sessions {
		workspace := "external"
		if !Api().Configuration.Standalone {
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

	ts.tableRenderer.FillTable(rows, sessions)
}

func (tv *TmuxSessionView) Render() error {
	isViewFocused := tv.view.IsFocused()

	tv.view.Clear()
	tv.tableRenderer.RenderWithSelectCallBack(tv.view, func(_ int, _ *TableRow[*tmux.TmuxSession]) bool {
		return isViewFocused
	})

	return nil
}
