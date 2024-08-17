package ui

import (
	"mynav/pkg/core"
	"mynav/pkg/events"
	"mynav/pkg/system"
	"mynav/pkg/tui"
	"strconv"

	"github.com/GianlucaP106/gotmux/gotmux"
	"github.com/awesome-gocui/gocui"
)

type tmuxSessionView struct {
	view          *tui.View
	tableRenderer *tui.TableRenderer[*gotmux.Session]
}

var _ viewable = new(tmuxSessionView)

func newTmuxSessionView() *tmuxSessionView {
	return &tmuxSessionView{}
}

func getTmuxSessionView() *tmuxSessionView {
	return getViewable[*tmuxSessionView]()
}

func (tv *tmuxSessionView) getView() *tui.View {
	return tv.view
}

func (tv *tmuxSessionView) Focus() {
	focusView(tv.getView().Name())
}

func (tv *tmuxSessionView) init() {
	tv.view = getViewPosition(TmuxSessionView).Set()

	tv.view.Title = tui.WithSurroundingSpaces("Tmux Sessions")
	styleView(tv.view)

	sizeX, sizeY := tv.view.Size()
	tv.tableRenderer = tui.NewTableRenderer[*gotmux.Session]()
	titles := []string{
		"Windows",
		"Workspace",
		"Session name",
	}
	proportions := []float64{
		0.2,
		0.3,
		0.5,
	}
	tv.tableRenderer.InitTable(sizeX, sizeY, titles, proportions)

	events.AddEventListener(events.TmuxSessionChangeEvent, func(_ string) {
		tv.refresh()
		renderView(tv)
		events.Emit(events.TmuxWindowChangeEvent)
	})

	tv.refresh()

	tv.view.KeyBinding().
		Set('o', "Attach to session", func() {
			if core.IsTmuxSession() {
				openToastDialogError("You are already in a tmux session. Nested tmux sessions are not supported yet.")
				return
			}

			session := tv.getSelectedSession()
			if session == nil {
				return
			}

			tui.RunAction(func() {
				getApi().Tmux.AttachTmuxSession(session)
			})
		}).
		Set(gocui.KeyEnter, "Focus window view", func() {
			getTmuxWindowView().focus()
		}).
		Set('D', "Delete session", func() {
			session := tv.getSelectedSession()
			if session == nil {
				return
			}

			openConfirmationDialog(func(b bool) {
				if b {
					if err := getApi().Tmux.DeleteTmuxSession(session); err != nil {
						openToastDialogError(err.Error())
						return
					}
					events.Emit(events.WorkspaceChangeEvent)
				}
			}, "Are you sure you want to delete this session?")
		}).
		Set('X', "Kill tmux server (kill all sessions)", func() {
			if tv.getSelectedSession() == nil {
				return
			}

			openConfirmationDialog(func(b bool) {
				if b {
					if err := getApi().Tmux.KillTmuxServer(); err != nil {
						openToastDialogError(err.Error())
						return
					}
				}
			}, "Are you sure you want to delete ALL tmux sessions?")
		}).
		Set('W', "Kill ALL non-external tmux sessions (has a workspace)", func() {
			if getApi().Configuration.Standalone || getApi().Core.GetWorkspaceTmuxSessionCount() == 0 {
				return
			}

			openConfirmationDialog(func(b bool) {
				if b {
					if err := getApi().Core.DeleteAllWorkspaceTmuxSessions(); err != nil {
						openToastDialogError(err.Error())
						return
					}
				}
			}, "Are you sure you want to delete ALL non-external tmux sessions?")
		}).
		Set('j', "Move down", func() {
			tv.tableRenderer.Down()
			events.Emit(events.TmuxWindowChangeEvent)
		}).
		Set('k', "Move up", func() {
			tv.tableRenderer.Up()
			events.Emit(events.TmuxWindowChangeEvent)
		}).
		Set('c', "Open choose tree in session", func() {
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

			tui.RunAction(func() {
				session.Attach()
			})
		}).
		Set('a', "New external session (not associated to a workspace)", func() {
			if core.IsTmuxSession() {
				return
			}
			openEditorDialog(func(s string) {
				tui.RunAction(func() {
					getApi().Tmux.CreateAndAttachTmuxSession(s, "~")
				})
			}, func() {}, "New session name", smallEditorSize)
		}).
		Set('?', "Toggle cheatsheet", func() {
			OpenHelpDialog(tv.view.GetKeybindings(), func() {})
		})
}

func (tv *tmuxSessionView) getSelectedSession() *gotmux.Session {
	_, ts := tv.tableRenderer.GetSelectedRow()
	if ts != nil {
		return *ts
	}

	return nil
}

func (ts *tmuxSessionView) refresh() {
	sessions := make([]*gotmux.Session, 0)
	sessions = append(sessions, getApi().Tmux.GetTmuxSessions()...)

	rows := make([][]string, 0)
	for _, session := range sessions {
		workspace := "external"
		sessionName := session.Name
		if !getApi().Configuration.Standalone {
			w := getApi().Core.GetWorkspaceByTmuxSession(session)
			if w != nil {
				workspace = w.ShortPath()
				sessionName = system.ShortenPath(sessionName, 20)
			}
		}

		windows := strconv.Itoa(session.Windows)
		rows = append(rows, []string{
			windows,
			workspace,
			sessionName,
		})
	}

	ts.tableRenderer.FillTable(rows, sessions)
}

func (tv *tmuxSessionView) render() error {
	isViewFocused := tv.view.IsFocused()
	tv.view.Clear()
	tv.view.Resize(getViewPosition(tv.view.Name()))
	tv.tableRenderer.RenderWithSelectCallBack(tv.view, func(_ int, _ *tui.TableRow[*gotmux.Session]) bool {
		return isViewFocused
	})

	return nil
}
