package ui

import (
	"mynav/pkg/core"
	"mynav/pkg/system"
	"mynav/pkg/tasks"
	"mynav/pkg/tui"
	"strconv"
	"time"

	"github.com/GianlucaP106/gotmux/gotmux"
	"github.com/awesome-gocui/gocui"
)

type tmuxWindowView struct {
	view          *tui.View
	tableRenderer *tui.TableRenderer[*gotmux.Window]
}

var _ viewable = new(tmuxWindowView)

func newTmuxWindowView() *tmuxWindowView {
	return &tmuxWindowView{}
}

func getTmuxWindowView() *tmuxWindowView {
	return getViewable[*tmuxWindowView]()
}

func (t *tmuxWindowView) getView() *tui.View {
	return t.view
}

func (t *tmuxWindowView) focus() {
	focusView(t.getView().Name())
}

func (t *tmuxWindowView) getSelectedWindow() *gotmux.Window {
	_, value := t.tableRenderer.GetSelectedRow()
	if value != nil {
		return *value
	}

	return nil
}

func (t *tmuxWindowView) init() {
	t.view = getViewPosition(TmuxWindowView).Set()

	t.view.Title = tui.WithSurroundingSpaces("Tmux Windows")
	styleView(t.view)

	sizeX, sizeY := t.view.Size()
	t.tableRenderer = tui.NewTableRenderer[*gotmux.Window]()
	t.tableRenderer.InitTable(sizeX, sizeY, []string{
		"Name",
		"# Panes",
		"Activity",
		"Active",
	}, []float64{
		0.25,
		0.25,
		0.25,
		0.25,
	})

	t.refresh()

	tv := getTmuxSessionView()
	tpv := getTmuxPaneView()
	t.view.KeyBinding().
		Set('j', "Move down", func() {
			t.tableRenderer.Down()
			t.refreshDown()
		}).
		Set('k', "Move up", func() {
			t.tableRenderer.Up()
			t.refreshDown()
		}).
		Set('o', "Open tmux session", func() {
			if core.IsTmuxSession() {
				openToastDialogError("You are already in a tmux session. Nested tmux sessions are not supported yet.")
				return
			}

			session := getTmuxSessionView().getSelectedSession()
			if session == nil {
				return
			}

			runAction(func() {
				getApi().Tmux.AttachTmuxSession(session)
			})
		}).
		Set('X', "Kill this window", func() {
			w := t.getSelectedWindow()
			if w == nil {
				return
			}

			openConfirmationDialog(func(b bool) {
				if !b {
					return
				}

				err := getApi().Tmux.KillTmuxWindow(w)
				if err != nil {
					openToastDialogError(err.Error())
				}

				tv.refreshTmuxViewsAsync()
			}, "Are you sure you want to kill this window?")
		}).
		Set(gocui.KeyEsc, "Focus session view", func() {
			tv.Focus()
		}).
		Set(gocui.KeyEnter, "Focus Pane view", func() {
			tpv.focus()
		}).
		Set(gocui.KeyCtrlL, "Focus Pane view", func() {
			tpv.focus()
		}).
		Set(gocui.KeyArrowRight, "Focus Pane view", func() {
			tpv.focus()
		}).
		Set('?', "Toggle cheatsheet", func() {
			openHelpDialog(t.view.GetKeybindings(), func() {})
		})
}

func (t *tmuxWindowView) refresh() {
	selectedSession := getTmuxSessionView().getSelectedSession()
	if selectedSession == nil {
		t.tableRenderer.ClearTable()
		return
	}

	rows := make([][]string, 0)
	windows, _ := selectedSession.ListWindows()
	for _, w := range windows {
		active := "No"
		if w.Active {
			active = "Yes"
		}

		activityInt, err := strconv.Atoi(w.Activity)
		if err != nil {
			continue
		}

		time := time.Unix(int64(activityInt), 0)

		rows = append(rows, []string{
			w.Name,
			strconv.Itoa(w.Panes),
			time.Format(system.TimeFormat()),
			active,
		})
	}

	t.tableRenderer.FillTable(rows, windows)
}

func (t *tmuxWindowView) refreshDown() {
	tasks.QueueTask(func() {
		tpv := getTmuxPaneView()
		tpv.refresh()
		renderView(tpv)
		tpv2 := getTmuxPreviewView()
		tpv2.refresh()
		renderView(tpv2)
	})
}

func (t *tmuxWindowView) render() error {
	isFocused := t.view.IsFocused()
	t.view.Clear()
	t.view = getViewPosition(t.view.Name()).Set()
	t.view.Resize(getViewPosition(t.view.Name()))
	t.tableRenderer.RenderWithSelectCallBack(t.view, func(i int, tr *tui.TableRow[*gotmux.Window]) bool {
		return isFocused
	})

	return nil
}
