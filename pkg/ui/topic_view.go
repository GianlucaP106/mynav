package ui

import (
	"mynav/pkg/core"
	"mynav/pkg/events"
	"mynav/pkg/persistence"
	"mynav/pkg/tui"

	"github.com/awesome-gocui/gocui"
)

type topicsView struct {
	view          *tui.View
	tableRenderer *tui.TableRenderer[*core.Topic]
	search        *persistence.Value[string]
	globalSearch  *persistence.Value[string]
}

var _ viewable = new(topicsView)

func newTopicsView() *topicsView {
	return &topicsView{
		search: persistence.NewValue(""),
	}
}

func getTopicsView() *topicsView {
	return getViewable[*topicsView]()
}

func (tv *topicsView) getView() *tui.View {
	return tv.view
}

func (tv *topicsView) Focus() {
	focusView(tv.getView().Name())
}

func (tv *topicsView) init() {
	tv.view = getViewPosition(TopicView).Set()

	tv.view.Title = tui.WithSurroundingSpaces("Topics")
	styleView(tv.view)

	sizeX, sizeY := tv.view.Size()
	tv.tableRenderer = tui.NewTableRenderer[*core.Topic]()
	titles := []string{
		"Name",
		"Last Modified",
	}
	colProportions := []float64{
		0.5,
		0.5,
	}
	tv.tableRenderer.InitTable(sizeX, sizeY, titles, colProportions)

	events.AddEventListener(events.TopicChangeEvent, func(_ string) {
		tv.refresh()
		wv := getWorkspacesView()
		wv.refresh()
		renderView(tv)
		renderView(wv)
	})

	tv.refresh()

	if selectedWorkspace := getApi().Core.GetSelectedWorkspace(); selectedWorkspace != nil {
		tv.selectTopicByName(selectedWorkspace.Topic.Name)
	}

	moveRight := func() {
		if getApi().Core.GetTopicCount() > 0 {
			getWorkspacesView().Focus()
		}
	}

	tv.view.KeyBinding().
		Set('j', "Move down", func() {
			tv.tableRenderer.Down()
			events.Emit(events.WorkspaceChangeEvent)
		}).
		Set('k', "Move up", func() {
			tv.tableRenderer.Up()
			events.Emit(events.WorkspaceChangeEvent)
		}).
		Set(gocui.KeyEnter, "Open topic", moveRight).
		Set('/', "Search by name", func() {
			openEditorDialog(func(s string) {
				tv.search.Set(s)
				tv.view.Subtitle = tui.WithSurroundingSpaces("Searching: " + s)
				tv.refresh()
				getWorkspacesView().refresh()
			}, func() {}, "Search", smallEditorSize)
		}).
		Set(gocui.KeyEsc, "Escape search", func() {
			if tv.search.Get() != "" {
				tv.search.Set("")
				tv.view.Subtitle = ""
				tv.refresh()
				getWorkspacesView().refresh()
			}
		}).
		Set('a', "Create a topic", func() {
			openEditorDialog(func(s string) {
				if err := getApi().Core.CreateTopic(s); err != nil {
					openToastDialogError(err.Error())
					return
				}

				// HACK: when there a is a new topic
				// This will result in the corresponding topic going to the top
				// because we are sorting by modifed time
				tv.tableRenderer.SelectRow(0)
			}, func() {}, "Topic name", smallEditorSize)
		}).
		Set('r', "Rename topic", func() {
			t := tv.getSelectedTopic()
			if t == nil {
				return
			}

			openEditorDialogWithDefaultValue(func(s string) {
				if err := getApi().Core.RenameTopic(t, s); err != nil {
					openToastDialogError(err.Error())
					return
				}
			}, func() {}, "New topic name", smallEditorSize, t.Name)
		}).
		Set('s', "Search for a workspace", func() {
			sd := new(*searchListDialog[*core.Workspace])
			*sd = openSearchListDialog(searchDialogConfig[*core.Workspace]{
				onSearch: func(s string) ([][]string, []*core.Workspace) {
					workspaces := getApi().Core.GetWorkspaces().Sorted().FilterByNameContaining(s)
					rows := make([][]string, 0)
					for _, w := range workspaces {
						rows = append(rows, []string{
							w.Topic.Name,
							w.Name,
						})
					}

					return rows, workspaces
				},
				onSelect: func(w *core.Workspace) {
					tv.selectTopicByName(w.Topic.Name)
					wv := getWorkspacesView()
					wv.refresh()
					wv.selectWorkspaceByShortPath(w.ShortPath())

					if *sd != nil {
						(*sd).close()
					}

					wv.Focus()
				},
				onSelectDescription: "Go to workspace",
				searchViewTitle:     "Search a workspace",
				tableViewTitle:      "Result",
				tableTitles: []string{
					"Topic",
					"Name",
				}, tableProportions: []float64{
					0.5,
					0.5,
				},
			})
		}).
		Set('D', "Delete topic", func() {
			if getApi().Core.GetTopicCount() <= 0 {
				return
			}
			openConfirmationDialog(func(b bool) {
				if !b {
					return
				}

				if err := getApi().Core.DeleteTopic(tv.getSelectedTopic()); err != nil {
					openToastDialogError(err.Error())
				}
			}, "Are you sure you want to delete this topic? All its content will be deleted.")
		}).
		Set('?', "Toggle cheatsheet", func() {
			OpenHelpDialog(tv.view.GetKeybindings(), func() {})
		})
}

func (tv *topicsView) refresh() {
	topics := getApi().Core.GetTopics().Sorted()

	search := tv.search.Get()
	if search != "" {
		topics = topics.FilterByNameContaining(search)
	}

	rows := make([][]string, 0)
	rowValues := make([]*core.Topic, 0)
	for _, topic := range topics {
		rowValues = append(rowValues, topic)
		rows = append(rows, []string{
			topic.Name,
			topic.GetLastModifiedTimeFormatted(),
		})
	}

	tv.tableRenderer.FillTable(rows, rowValues)
}

func (tv *topicsView) getSelectedTopic() *core.Topic {
	_, t := tv.tableRenderer.GetSelectedRow()
	if t != nil {
		return *t
	}
	return nil
}

func (tv *topicsView) selectTopicByName(name string) {
	tv.tableRenderer.SelectRowByValue(func(t *core.Topic) bool {
		return t.Name == name
	})
}

func (tv *topicsView) render() error {
	currentViewSelected := tv.view.IsFocused()
	tv.view.Clear()
	tv.view = getViewPosition(tv.view.Name()).Set()

	tv.tableRenderer.RenderWithSelectCallBack(tv.view, func(_ int, _ *tui.TableRow[*core.Topic]) bool {
		return currentViewSelected
	})

	return nil
}
