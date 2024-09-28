package ui

import (
	"mynav/pkg/core"
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
		"Session name",
		"Workspace",
	}
	proportions := []float64{
		0.2,
		0.5,
		0.3,
	}
	tv.tableRenderer.InitTable(sizeX, sizeY, titles, proportions)

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

			var error error = nil
			runAction(func() {
				err := api().Tmux.AttachTmuxSession(session)
				if err != nil {
					error = err
				}
			})

			if error != nil {
				openToastDialogError(error.Error())
			}
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
					if err := api().Tmux.DeleteTmuxSession(session); err != nil {
						openToastDialogError(err.Error())
						return
					}

					refreshTmuxViews()
					refreshMainViews()
				}
			}, "Are you sure you want to delete this session?")
		}).
		Set('X', "Kill tmux server (kill all sessions)", func() {
			if tv.getSelectedSession() == nil {
				return
			}

			openConfirmationDialog(func(b bool) {
				if b {
					if err := api().Tmux.KillTmuxServer(); err != nil {
						openToastDialogError(err.Error())
						return
					}

					refreshTmuxViews()
					refreshMainViews()
				}
			}, "Are you sure you want to delete ALL tmux sessions?")
		}).
		Set('W', "Kill ALL non-external tmux sessions (has a workspace)", func() {
			if api().GlobalConfiguration.Standalone || tv.getSelectedSession() == nil {
				return
			}

			openConfirmationDialog(func(b bool) {
				if b {
					if err := api().Workspaces.DeleteAllWorkspaceTmuxSessions(); err != nil {
						openToastDialogError(err.Error())
						return
					}

					refreshTmuxViews()
					refreshMainViews()
				}
			}, "Are you sure you want to delete ALL non-external tmux sessions?")
		}).
		Set('j', "Move down", func() {
			tv.tableRenderer.Down()
			tv.refreshDown()
		}).
		Set('k', "Move up", func() {
			tv.tableRenderer.Up()
			tv.refreshDown()
		}).
		Set('c', "Open choose tree in session", func() {
			session := tv.getSelectedSession()
			if session == nil {
				return
			}

			var err error = nil
			runAction(func() {
				err = api().Tmux.OpenTmuxSessionChooseTree(session)
			})
			if err != nil {
				openToastDialogError(err.Error())
			}
		}).
		Set('a', "New external session (not associated to a workspace)", func() {
			if core.IsTmuxSession() {
				return
			}

			openEditorDialog(func(s string) {
				runAction(func() {
					api().Tmux.CreateAndAttachTmuxSession(s, "~")
				})
			}, func() {}, "New session name", smallEditorSize)
		}).
		Set('?', "Toggle cheatsheet", func() {
			openHelpDialog(tv.view.GetKeybindings(), func() {})
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
	sessions = append(sessions, api().Tmux.GetTmuxSessions()...)

	rows := make([][]string, 0)
	for _, session := range sessions {
		workspace := "external"
		sessionName := session.Name
		if !api().GlobalConfiguration.Standalone {
			w := api().Workspaces.GetWorkspaceByTmuxSession(session)
			if w != nil {
				workspace = w.ShortPath()
				sessionName = system.ShortenPath(sessionName, 20)
			}
		}

		windows := strconv.Itoa(session.Windows)
		rows = append(rows, []string{
			windows,
			sessionName,
			workspace,
		})
	}

	ts.tableRenderer.FillTable(rows, sessions)
}

func refreshTmuxViews() {
	ui.queueRefresh(func() {
		ts := getTmuxSessionView()
		ts.refresh()
		renderView(ts)
		ts.refreshDown()
	})
}

func (ts *tmuxSessionView) refreshDown() {
	go func() {
		twv := getTmuxWindowView()
		twv.refresh()
		renderView(twv)

		tpv := getTmuxPaneView()
		tpv.refresh()
		renderView(tpv)

		tpvv := getTmuxPreviewView()
		tpvv.refresh()
		renderView(tpvv)
	}()
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
