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
}

var _ Dialog = &TmuxSessionView{}

func newTmuxSessionView() *TmuxSessionView {
	ts := &TmuxSessionView{}

	return ts
}

func (tv *TmuxSessionView) Open(ui *UI) {
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
			GetDialog[*ConfirmationDialog](ui).Open(func(b bool) {
				if b {
					session := tv.getSelectedSession()
					Api().DeleteTmuxSession(session)
					tv.refreshTmuxSessions()
					ui.RefreshWorkspaces()
				}
				FocusView(tv.Name())
			}, "Are you sure you want to delete this session?")
		case key == gocui.KeyEsc:
			tv.Close()
			ui.FocusTopicsView()
		case ch == 'j':
			tv.listRenderer.increment()
		case ch == 'k':
			tv.listRenderer.decrement()
		case ch == 'a':
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

	// TODO: investigate the weird behaviour of len when there is color
	// if selected {
	// 	// sessionName = color.New(color.Blue).Sprint(sessionName)
	// 	// windows = color.Green.Sprint(windows)
	// 	workspace = color.Blue.Sprint(workspace)
	// } else {
	// 	// sessionName = color.White.Sprint(sessionName)
	// 	// windows = color.White.Sprint(windows)
	// 	workspace = color.White.Sprint(workspace)
	// }

	line = withSpacePadding(sessionName, 3*fifth)
	line += withSpacePadding(windows, fifth)
	line += withSpacePadding(workspace, fifth)

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
	tv.listRenderer.forEach(func(idx int) {
		session := tv.sessions[idx]
		potentialWorkspace := Api().GetWorkspaceByTmuxSession(session)
		line := tv.format(session, idx == tv.listRenderer.selected, potentialWorkspace)
		fmt.Fprintln(view, line)
	})

	return nil
}