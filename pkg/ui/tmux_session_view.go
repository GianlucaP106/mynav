package ui

import (
	"fmt"
	"mynav/pkg/api"
	"mynav/pkg/utils"
	"strconv"

	"github.com/awesome-gocui/gocui"
	"github.com/gookit/color"
)

const TmuxSessionViewName = "TmuxSessionView"

type TmuxSessionView struct {
	editor       Editor
	listRenderer *ListRenderer
	sessions     []*api.TmuxSession
	standalone   bool
}

var _ Dialog = &TmuxSessionView{}

func newTmuxSessionView() *TmuxSessionView {
	ts := &TmuxSessionView{}

	return ts
}

func (tv *TmuxSessionView) refresh(ui *UI) {
	tv.refreshTmuxSessions()
	if !tv.standalone {
		ui.RefreshWorkspaces()
	}
}

func (tv *TmuxSessionView) Open(ui *UI, standalone bool) {
	tv.standalone = standalone
	view := SetViewLayout(tv.Name())

	view.Title = withSurroundingSpaces("TMUX Sessions")
	view.TitleColor = gocui.ColorBlue
	view.FrameColor = gocui.ColorGreen

	_, sizeY := view.Size()
	tv.listRenderer = newListRenderer(0, sizeY, 0)
	tv.refreshTmuxSessions()

	tv.editor = gocui.EditorFunc(func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
		switch {
		case key == gocui.KeyEnter:
			if utils.IsTmuxSession() {
				GetDialog[*ToastDialog](ui).Open("You are already in a tmux session. Nested tmux sessions are not supported yet.", func() {
					FocusView(tv.Name())
				})
				return
			}

			session := tv.getSelectedSession()
			ui.setAction(utils.AttachTmuxSessionCmd(session.Name))
		case ch == 'd':
			if Api().GetTmuxSessionCount() == 0 {
				return
			}

			GetDialog[*ConfirmationDialog](ui).Open(func(b bool) {
				if b {
					session := tv.getSelectedSession()
					if err := Api().DeleteTmuxSession(session); err != nil {
						GetDialog[*ToastDialog](ui).Open(err.Error(), func() {
							FocusView(tv.Name())
						})
						return
					}
					tv.refresh(ui)
				}
				FocusView(tv.Name())
			}, "Are you sure you want to delete this session?")
		case ch == 'x':
			if Api().GetTmuxSessionCount() == 0 {
				return
			}

			GetDialog[*ConfirmationDialog](ui).Open(func(b bool) {
				if b {
					if err := Api().DeleteAllTmuxSessions(); err != nil {
						GetDialog[*ToastDialog](ui).Open(err.Error(), func() {
							FocusView(tv.Name())
						})
						return
					}
					tv.refresh(ui)
				}
				FocusView(tv.Name())
			}, "Are you sure you want to delete ALL tmux sessions?")
		case ch == 'w':
			if tv.standalone || Api().GetWorkspaceTmuxSessionCount() == 0 {
				return
			}

			GetDialog[*ConfirmationDialog](ui).Open(func(b bool) {
				if b {
					if err := Api().DeleteAllWorkspaceTmuxSessions(); err != nil {
						GetDialog[*ToastDialog](ui).Open(err.Error(), func() {
							FocusView(tv.Name())
						})
						return
					}
					tv.refresh(ui)
				}
				FocusView(tv.Name())
			}, "Are you sure you want to delete ALL non-external tmux sessions?")
		case key == gocui.KeyEsc:
			if !tv.standalone {
				tv.Close()
				ui.FocusTopicsView()
			}
		case ch == 'j':
			tv.listRenderer.increment()
		case ch == 'k':
			tv.listRenderer.decrement()
		case ch == 'a':
			if utils.IsTmuxSession() {
				return
			}
			GetDialog[*EditorDialog](ui).Open(func(s string) {
				ui.setAction(utils.NewTmuxSessionCmd(s, "~"))
			}, func() {
				FocusView(tv.Name())
			}, "New session name", Small)
		case ch == '?':
			GetDialog[*HelpView](ui).Open(getKeyBindings(tv.Name()), func() {
				FocusView(tv.Name())
			})
		}
	})

	view.Editor = tv.editor
	view.Editable = true
	FocusView(tv.Name())
}

func (tv *TmuxSessionView) getSelectedSession() *api.TmuxSession {
	return tv.sessions[tv.listRenderer.selected]
}

func (tv *TmuxSessionView) Close() {
	DeleteView(tv.Name())
}

func (ts *TmuxSessionView) refreshTmuxSessions() {
	out := make([]*api.TmuxSession, 0)
	for _, session := range Api().GetTmuxSessions() {
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

func (tv *TmuxSessionView) format(session *api.TmuxSession, selected bool, w *api.Workspace) string {
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
		return nil
	}

	if ui.action.Command != nil {
		return gocui.ErrQuit
	}

	view.Clear()
	fmt.Fprintln(view, tv.formatTitles())
	tv.listRenderer.forEach(func(idx int) {
		session := tv.sessions[idx]
		var potentialWorkspace *api.Workspace
		if !tv.standalone {
			potentialWorkspace = Api().GetWorkspaceByTmuxSession(session)
		}
		line := tv.format(session, idx == tv.listRenderer.selected, potentialWorkspace)
		fmt.Fprintln(view, line)
	})

	return nil
}
