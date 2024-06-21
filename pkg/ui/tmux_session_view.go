package ui

import (
	"fmt"
	"mynav/pkg/core"
	"mynav/pkg/system"
	"mynav/pkg/tmux"
	"strconv"

	"github.com/awesome-gocui/gocui"
	"github.com/gookit/color"
)

const TmuxSessionViewName = "TmuxSessionView"

type TmuxSessionView struct {
	listRenderer *ListRenderer
	sessions     []*tmux.TmuxSession
	standalone   bool
}

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

	_, sizeY := view.Size()
	tv.listRenderer = newListRenderer(0, sizeY, 0)
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

	KeyBinding(TmuxSessionViewName).
		set(gocui.KeyEnter, func() {
			if system.IsTmuxSession() {
				GetDialog[*ToastDialog](ui).Open("You are already in a tmux session. Nested tmux sessions are not supported yet.", func() {})
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
						GetDialog[*ToastDialog](ui).Open(err.Error(), func() {})
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
						GetDialog[*ToastDialog](ui).Open(err.Error(), func() {})
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
						GetDialog[*ToastDialog](ui).Open(err.Error(), func() {})
						return
					}
					ui.RefreshMainView()
				}
			}, "Are you sure you want to delete ALL non-external tmux sessions?")
		}).
		set('j', func() {
			tv.listRenderer.increment()
		}).
		set('k', func() {
			tv.listRenderer.decrement()
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
		set(gocui.KeyCtrlH, moveLeft)

	if tv.standalone {
		FocusView(tv.Name())
	}
}

func (tv *TmuxSessionView) getSelectedSession() *tmux.TmuxSession {
	return tv.sessions[tv.listRenderer.selected]
}

func (ts *TmuxSessionView) refreshTmuxSessions() {
	out := make([]*tmux.TmuxSession, 0)
	for _, session := range Api().Tmux.GetTmuxSessions() {
		out = append(out, session)
	}

	ts.sessions = out

	if ts.listRenderer != nil {
		newListSize := len(ts.sessions)
		if ts.listRenderer.listSize != newListSize {
			ts.listRenderer.setListSize(newListSize)
		}
	}
}

func (tv *TmuxSessionView) Name() string {
	return TmuxSessionViewName
}

func (tv *TmuxSessionView) formatTitles() string {
	view := GetInternalView(tv.Name())
	sizeX, _ := view.Size()

	fifth := (sizeX / 5) + 1
	line := ""

	line += withSpacePadding("Workspace | external", fifth)
	line += withSpacePadding("Windows Open", fifth)
	line += withSpacePadding("Session Name", 3*fifth)

	return line
}

func (tv *TmuxSessionView) format(session *tmux.TmuxSession, selected bool, w *core.Workspace) string {
	view := GetInternalView(tv.Name())
	sizeX, _ := view.Size()

	fifth := (sizeX / 5) + 1

	line := ""

	sessionName := session.Name + " "

	windows := strconv.Itoa(session.NumWindows) + " windows"

	workspace := ""
	if w != nil {
		workspace = w.ShortPath()
	} else {
		workspace = "external"
	}

	line += withSpacePadding(workspace, fifth)
	line += withSpacePadding(windows, fifth)
	line += withSpacePadding(sessionName, 3*fifth)

	if selected {
		line = color.New(color.BgCyan, color.Black).Sprint(line)
	} else {
		line = color.New(color.Blue).Sprint(line)
	}

	return line
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

	view.Clear()
	fmt.Fprintln(view, tv.formatTitles())
	tv.listRenderer.forEach(func(idx int) {
		session := tv.sessions[idx]
		var potentialWorkspace *core.Workspace
		if !tv.standalone {
			potentialWorkspace = Api().Core.GetWorkspaceByTmuxSession(session)
		}

		isViewFocused := false
		if fv := GetFocusedView(); fv != nil {
			isViewFocused = tv.Name() == GetFocusedView().Name()
		}

		line := tv.format(session, isViewFocused && idx == tv.listRenderer.selected, potentialWorkspace)
		fmt.Fprintln(view, line)
	})

	return nil
}
