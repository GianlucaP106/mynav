package ui

import (
	"mynav/pkg/constants"
	"mynav/pkg/core"
	"mynav/pkg/events"
	"strconv"

	"github.com/GianlucaP106/gotmux/gotmux"
	"github.com/awesome-gocui/gocui"
)

type TmuxSessionView struct {
	view          *View
	tableRenderer *TableRenderer[*gotmux.Session]
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
	tv.view = GetViewPosition(constants.TmuxSessionViewName).Set()

	tv.view.Title = withSurroundingSpaces("Tmux Sessions")
	StyleView(tv.view)

	sizeX, sizeY := tv.view.Size()
	tv.tableRenderer = NewTableRenderer[*gotmux.Session]()
	titles := []string{
		"Workspace",
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
		tv.refresh()
		RenderView(tv)
		events.Emit(constants.TmuxWindowChangeEventName)
	})

	tv.refresh()

	tv.view.KeyBinding().
		set('o', "Attach to session", func() {
			if core.IsTmuxSession() {
				OpenToastDialogError("You are already in a tmux session. Nested tmux sessions are not supported yet.")
				return
			}

			session := tv.getSelectedSession()
			RunAction(func() {
				Api().Tmux.AttachTmuxSession(session)
			})
		}).
		set(gocui.KeyEnter, "Focus window view", func() {
			GetTmuxWindowView().Focus()
		}).
		set('D', "Delete session", func() {
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
		}).
		set('X', "Kill ALL tmux sessions", func() {
			if Api().Tmux.GetTmuxSessionCount() == 0 {
				return
			}

			OpenConfirmationDialog(func(b bool) {
				if b {
					if err := Api().Tmux.KillTmuxServer(); err != nil {
						OpenToastDialogError(err.Error())
						return
					}
				}
			}, "Are you sure you want to delete ALL tmux sessions?")
		}).
		set('W', "Kill ALL non-externalt mux sessions (has a workspace)", func() {
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
		}).
		set('j', "Move down", func() {
			tv.tableRenderer.Down()
			events.Emit(constants.TmuxWindowChangeEventName)
		}).
		set('k', "Move up", func() {
			tv.tableRenderer.Up()
			events.Emit(constants.TmuxWindowChangeEventName)
		}).
		set('c', "Open choose tree in session", func() {
			// TODO: move this flow in core
			session := tv.getSelectedSession()
			if session == nil {
				return
			}

			windows, err := session.ListWindows()
			if err != nil {
				return
			}

			var window *gotmux.Window
			for _, w := range windows {
				if w != nil {
					window = w
					break
				}
			}

			if window == nil {
				window, err = session.New()
				if err != nil {
					return
				}
			}

			var pane *gotmux.Pane
			pane, err = window.GetPaneByIndex(0)
			if err != nil {
				// TODO: create pane - blocked by https://github.com/GianlucaP106/gotmux/issues/12
				return
			}

			if pane == nil {
				// TODO: create pane - blocked by https://github.com/GianlucaP106/gotmux/issues/12
				return
			}

			err = pane.ChooseTree(&gotmux.ChooseTreeOptions{
				SessionsCollapsed: true,
			})
			if err != nil {
				return
			}

			RunAction(func() {
				session.Attach()
			})
		}).
		set('a', "New external session (not associated to a workspace)", func() {
			if core.IsTmuxSession() {
				return
			}
			OpenEditorDialog(func(s string) {
				RunAction(func() {
					Api().Tmux.CreateAndAttachTmuxSession(s, "~")
				})
			}, func() {}, "New session name", Small)
		}).
		set('?', "Toggle cheatsheet", func() {
			OpenHelpView(tv.view.keybindingInfo.toList(), func() {})
		})
}

func (tv *TmuxSessionView) getSelectedSession() *gotmux.Session {
	_, ts := tv.tableRenderer.GetSelectedRow()
	if ts != nil {
		return *ts
	}

	return nil
}

func (ts *TmuxSessionView) refresh() {
	sessions := make([]*gotmux.Session, 0)
	sessions = append(sessions, Api().Tmux.GetTmuxSessions()...)

	rows := make([][]string, 0)
	for _, session := range sessions {
		workspace := "external"
		if !Api().Configuration.Standalone {
			w := Api().Core.GetWorkspaceByTmuxSession(session)
			if w != nil {
				workspace = w.ShortPath()
			}
		}

		windows := strconv.Itoa(session.Windows)
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
	tv.tableRenderer.RenderWithSelectCallBack(tv.view, func(_ int, _ *TableRow[*gotmux.Session]) bool {
		return isViewFocused
	})

	return nil
}
