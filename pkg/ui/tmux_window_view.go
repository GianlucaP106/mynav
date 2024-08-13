package ui

import (
	"mynav/pkg/constants"
	"mynav/pkg/events"
	"mynav/pkg/system"
	"strconv"
	"time"

	"github.com/GianlucaP106/gotmux/gotmux"
	"github.com/awesome-gocui/gocui"
)

type TmuxWindowView struct {
	view          *View
	tableRenderer *TableRenderer[*gotmux.Window]
}

var _ Viewable = new(TmuxWindowView)

func NewTmuxWindowView() *TmuxWindowView {
	return &TmuxWindowView{}
}

func GetTmuxWindowView() *TmuxWindowView {
	return GetViewable[*TmuxWindowView]()
}

func (t *TmuxWindowView) View() *View {
	return t.view
}

func (t *TmuxWindowView) Focus() {
	FocusView(t.View().Name())
}

func (t *TmuxWindowView) getSelectedWindow() *gotmux.Window {
	_, value := t.tableRenderer.GetSelectedRow()
	if value != nil {
		return *value
	}

	return nil
}

func (t *TmuxWindowView) Init() {
	t.view = GetViewPosition(constants.TmuxWindowViewName).Set()

	t.view.Title = withSurroundingSpaces("Tmux Windows")
	t.view.TitleColor = gocui.ColorBlue
	t.view.FrameColor = gocui.ColorGreen

	sizeX, sizeY := t.view.Size()
	t.tableRenderer = NewTableRenderer[*gotmux.Window]()
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

	events.AddEventListener(constants.TmuxWindowChangeEventName, func(s string) {
		t.refresh()
		RenderView(t)
		events.Emit(constants.TmuxPaneChangeEventName)
	})

	t.refresh()

	t.view.KeyBinding().
		set('j', "Move down", func() {
			t.tableRenderer.Down()
			events.Emit(constants.TmuxPaneChangeEventName)
		}).
		set('k', "Move up", func() {
			t.tableRenderer.Up()
			events.Emit(constants.TmuxPaneChangeEventName)
		}).
		set(gocui.KeyEnter, "Focus Pane view", func() {
			GetTmuxPaneView().Focus()
		}).
		set(gocui.KeyEsc, "Focus session view", func() {
			GetTmuxSessionView().Focus()
		}).
		set(gocui.KeyCtrlL, "Focus Pane view", func() {
			GetTmuxPaneView().Focus()
		}).
		set(gocui.KeyArrowRight, "Focus Pane view", func() {
			GetTmuxPaneView().Focus()
		})
}

func (t *TmuxWindowView) refresh() {
	selectedSession := GetTmuxSessionView().getSelectedSession()
	if selectedSession == nil {
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

func (t *TmuxWindowView) Render() error {
	isFocused := t.view.IsFocused()
	t.view.Clear()
	t.tableRenderer.RenderWithSelectCallBack(t.view, func(i int, tr *TableRow[*gotmux.Window]) bool {
		return isFocused
	})

	return nil
}
