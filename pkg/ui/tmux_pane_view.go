package ui

import (
	"mynav/pkg/system"
	"mynav/pkg/tui"
	"strconv"

	"github.com/GianlucaP106/gotmux/gotmux"
	"github.com/awesome-gocui/gocui"
)

type tmuxPaneView struct {
	view          *tui.View
	tableRenderer *tui.TableRenderer[*gotmux.Pane]
}

var _ viewable = new(tmuxPaneView)

func newTmuxPaneView() *tmuxPaneView {
	return &tmuxPaneView{}
}

func getTmuxPaneView() *tmuxPaneView {
	return getViewable[*tmuxPaneView]()
}

func (t *tmuxPaneView) getView() *tui.View {
	return t.view
}

func (t *tmuxPaneView) focus() {
	focusView(t.view.Name())
}

func (t *tmuxPaneView) init() {
	t.view = getViewPosition(TmuxPaneView).Set()

	t.view.Title = tui.WithSurroundingSpaces("Tmux Panes")
	styleView(t.view)

	sizeX, sizeY := t.view.Size()
	t.tableRenderer = tui.NewTableRenderer[*gotmux.Pane]()
	t.tableRenderer.InitTable(sizeX, sizeY, []string{
		"Current command",
		"Pid",
		"Path",
	}, []float64{
		0.3,
		0.2,
		0.50,
	})

	t.refresh()

	tpv := getTmuxPreviewView()
	twv := getTmuxWindowView()
	t.view.KeyBinding().
		Set('j', "Move down", func() {
			t.tableRenderer.Down()
			refreshAsync(tpv)
		}).
		Set('k', "Move up", func() {
			t.tableRenderer.Up()
			refreshAsync(tpv)
		}).
		Set('X', "Kill this pane", func() {
			pane := t.getSelectedPane()
			if pane == nil {
				return
			}

			openConfirmationDialog(func(b bool) {
				if !b {
					return
				}

				err := getApi().Tmux.KillTmuxPane(pane)
				if err != nil {
					openToastDialogError(err.Error())
				}

				refreshTmuxViewsAsync()
			}, "Are you sure you want to kill this pane?")
		}).
		Set(gocui.KeyEsc, "Focus window view", func() {
			twv.focus()
		}).
		Set(gocui.KeyCtrlH, "Focus window view", func() {
			twv.focus()
		}).
		Set(gocui.KeyArrowLeft, "Focus window view", func() {
			twv.focus()
		}).
		Set('?', "Toggle cheatsheet", func() {
			openHelpDialog(twv.view.GetKeybindings(), func() {})
		})
}

func (t *tmuxPaneView) refresh() {
	window := getTmuxWindowView().getSelectedWindow()
	if window == nil {
		t.tableRenderer.ClearTable()
		return
	}

	panes, err := window.ListPanes()
	if err != nil {
		return
	}

	rows := make([][]string, 0)
	for _, p := range panes {
		rows = append(rows, []string{
			p.CurrentCommand,
			strconv.Itoa(int(p.Pid)),
			system.ShortenPath(p.CurrentPath, 32),
		})
	}

	t.tableRenderer.FillTable(rows, panes)
}

func (t *tmuxPaneView) getSelectedPane() *gotmux.Pane {
	_, value := t.tableRenderer.GetSelectedRow()
	if value != nil {
		return *value
	}

	return nil
}

func (t *tmuxPaneView) render() error {
	isFocused := t.view.IsFocused()
	t.view.Clear()
	t.view = getViewPosition(t.view.Name()).Set()
	t.view.Resize(getViewPosition(t.view.Name()))
	t.tableRenderer.RenderWithSelectCallBack(t.view, func(i int, tr *tui.TableRow[*gotmux.Pane]) bool {
		return isFocused
	})

	return nil
}
