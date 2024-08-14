package ui

import (
	"mynav/pkg/constants"
	"mynav/pkg/events"
	"mynav/pkg/system"
	"strconv"

	"github.com/GianlucaP106/gotmux/gotmux"
	"github.com/awesome-gocui/gocui"
)

type TmuxPaneView struct {
	view          *View
	tableRenderer *TableRenderer[*gotmux.Pane]
}

var _ Viewable = new(TmuxPaneView)

func NewTmuxPaneView() *TmuxPaneView {
	return &TmuxPaneView{}
}

func GetTmuxPaneView() *TmuxPaneView {
	return GetViewable[*TmuxPaneView]()
}

func (t *TmuxPaneView) View() *View {
	return t.view
}

func (t *TmuxPaneView) Focus() {
	FocusView(t.view.Name())
}

func (t *TmuxPaneView) Init() {
	t.view = GetViewPosition(constants.TmuxPaneViewName).Set()

	t.view.Title = withSurroundingSpaces("Tmux Panes")
	StyleView(t.view)

	sizeX, sizeY := t.view.Size()
	t.tableRenderer = NewTableRenderer[*gotmux.Pane]()
	t.tableRenderer.InitTable(sizeX, sizeY, []string{
		"Current command",
		"Pid",
		"Path",
	}, []float64{
		0.25,
		0.25,
		0.50,
	})

	events.AddEventListener(constants.TmuxPaneChangeEventName, func(s string) {
		t.refresh()
		RenderView(t)
		events.Emit(constants.TmuxPreviewChangeEventName)
	})

	t.refresh()

	t.view.KeyBinding().
		set('j', "Move down", func() {
			t.tableRenderer.Down()
			events.Emit(constants.TmuxPreviewChangeEventName)
		}).
		set('k', "Move up", func() {
			t.tableRenderer.Up()
			events.Emit(constants.TmuxPreviewChangeEventName)
		}).
		set(gocui.KeyEsc, "Focus window view", func() {
			GetTmuxWindowView().Focus()
		}).
		set(gocui.KeyCtrlH, "Focus window view", func() {
			GetTmuxWindowView().Focus()
		}).
		set(gocui.KeyArrowLeft, "Focus window view", func() {
			GetTmuxWindowView().Focus()
		})
}

func (t *TmuxPaneView) refresh() {
	window := GetTmuxWindowView().getSelectedWindow()
	if window == nil {
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

func (t *TmuxPaneView) getSelectedPane() *gotmux.Pane {
	_, value := t.tableRenderer.GetSelectedRow()
	if value != nil {
		return *value
	}

	return nil
}

func (t *TmuxPaneView) Render() error {
	isFocused := t.view.IsFocused()
	t.view.Clear()
	t.tableRenderer.RenderWithSelectCallBack(t.view, func(i int, tr *TableRow[*gotmux.Pane]) bool {
		return isFocused
	})

	return nil
}
