package app

import (
	"fmt"
	"strconv"

	"github.com/GianlucaP106/mynav/pkg/core"
	"github.com/GianlucaP106/mynav/pkg/tui"
	"github.com/awesome-gocui/gocui"
	"github.com/gookit/color"
)

type Topics struct {
	view  *tui.View
	table *tui.TableRenderer[*core.Topic]
}

func newTopicsView() *Topics {
	t := &Topics{}
	return t
}

func (tv *Topics) focus() {
	a.focusView(tv.view)
}

func (tv *Topics) refreshDown() {
	a.worker.DebounceLoad(func() {
		// run refresh continuously
		a.workspaces.refresh()
	}, func() {
		// once refresh is done and no more events are in we update the preview and info
		// show preview in the worker
		a.workspaces.refreshPreview()
		a.ui.Update(func() {
			// in the main routine we set loading to false
			a.workspaces.setLoading(false)

			// show info
			a.workspaces.showInfo()

			// render workspaces and preview
			a.workspaces.render()
			a.preview.render()
		})
	}, func() {
		a.ui.Update(func() {
			// if refresh takes long, we set loading to true
			a.workspaces.setLoading(true)
		})
	})
}

func (tv *Topics) refresh() {
	topics := a.api.Topics().Sorted()

	rows := make([][]string, 0)
	rowValues := make([]*core.Topic, 0)
	for _, topic := range topics {
		rowValues = append(rowValues, topic)
		topicWorkspaces := a.api.Workspaces(topic)
		timeStr := core.TimeAgo(topic.LastModified())
		rows = append(rows, []string{
			topic.Name,
			strconv.Itoa(len(topicWorkspaces)),
			timeStr,
		})
	}

	tv.table.Fill(rows, rowValues)
}

func (tv *Topics) selected() *core.Topic {
	_, t := tv.table.SelectedRow()
	if t != nil {
		return *t
	}
	return nil
}

func (tv *Topics) selectTopic(t *core.Topic) {
	tv.table.SelectRowByValue(func(t2 *core.Topic) bool {
		return t2.Name == t.Name
	})
}

func (tv *Topics) render() {
	currentViewSelected := a.ui.IsFocused(tv.view)
	tv.view.Clear()
	a.ui.Resize(tv.view, getViewPosition(tv.view.Name()))

	// update row marker
	row, _ := tv.table.SelectedRow()
	size := tv.table.Size()
	tv.view.Subtitle = fmt.Sprintf(" %d / %d ", min(row+1, size), size)

	// renders table and updates the last modified time
	tv.table.RenderTable(tv.view, func(_ int, _ *tui.TableRow[*core.Topic]) bool {
		return currentViewSelected
	}, func(i int, tr *tui.TableRow[*core.Topic]) {
		newTime := core.TimeAgo(tr.Value.LastModified())
		tr.Cols[len(tr.Cols)-1] = newTime
	})
}

func (tv *Topics) init() {
	tv.view = a.ui.SetView(getViewPosition(TopicView))
	tv.view.Title = " Topics "
	a.styleView(tv.view)

	sizeX, sizeY := tv.view.Size()
	tv.table = tui.NewTableRenderer[*core.Topic]()
	titles := []string{
		"Name",
		"Workspaces",
		"Last Modified",
	}
	colProportions := []float64{
		0.4,
		0.2,
		0.4,
	}
	styles := []color.Style{
		topicNameColor,
		alternateSessionMarkerColor,
		timestampColor,
	}
	tv.table.Init(sizeX, sizeY, titles, colProportions)
	tv.table.SetStyles(styles)

	moveRight := func() {
		if a.api.TopicCount() > 0 {
			a.workspaces.focus()
		}
	}

	down := func() {
		tv.table.Down()
		tv.refreshDown()
	}
	up := func() {
		tv.table.Up()
		tv.refreshDown()
	}
	a.ui.KeyBinding(tv.view).
		Set('j', "Move down", down).
		Set('k', "Move up", up).
		Set(gocui.KeyArrowDown, "Move down", down).
		Set(gocui.KeyArrowUp, "Move up", up).
		Set(gocui.KeyEnter, "Open topic", moveRight).
		Set('g', "Go to top", func() {
			tv.table.Top()
		}).
		Set('G', "Go to bottom", func() {
			tv.table.Bottom()
		}).
		Set('a', "Create a topic", func() {
			editor(func(s string) {
				topic, err := a.api.NewTopic(s)
				if err != nil {
					toast(err.Error(), toastError)
					return
				}

				a.refresh(topic, nil, nil)
				toast("Created topic "+topic.Name, toastInfo)
			}, func() {}, "Topic name", smallEditorSize, "")
		}).
		Set('r', "Rename topic", func() {
			t := tv.selected()
			if t == nil {
				return
			}

			editor(func(s string) {
				if err := a.api.RenameTopic(t, s); err != nil {
					toast(err.Error(), toastError)
					return
				}

				a.refresh(t, nil, nil)
				toast("Renamed topic "+t.Name, toastInfo)
			}, func() {}, "New topic name", smallEditorSize, t.Name)
		}).
		Set('D', "Delete topic", func() {
			t := tv.selected()
			if t == nil {
				return
			}
			alert(func(b bool) {
				if !b {
					return
				}
				if err := a.api.DeleteTopic(t); err != nil {
					toast(err.Error(), toastError)
				}

				a.refreshAll()
				toast("Deleted topic "+t.Name, toastInfo)
			}, fmt.Sprintf("Are you sure you want to delete topic %s? All its content will be deleted.", t.Name))
		}).
		Set('l', "Focus workspace view", func() {
			a.workspaces.focus()
		}).
		Set(gocui.KeyArrowRight, "Focus workspace view", func() {
			a.workspaces.focus()
		}).
		Set('?', "Toggle cheatsheet", func() {
			help(tv.view)
		})
}
