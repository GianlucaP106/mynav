package ui

import (
	"mynav/pkg/constants"
	"mynav/pkg/core"
	"mynav/pkg/events"
	"mynav/pkg/persistence"

	"github.com/awesome-gocui/gocui"
)

type TopicsView struct {
	view          *View
	tableRenderer *TableRenderer[*core.Topic]
	search        *persistence.Value[string]
	globalSearch  *persistence.Value[string]
}

var _ Viewable = new(TopicsView)

func NewTopicsView() *TopicsView {
	return &TopicsView{
		search: persistence.NewValue(""),
	}
}

func GetTopicsView() *TopicsView {
	return GetViewable[*TopicsView]()
}

func (tv *TopicsView) View() *View {
	return tv.view
}

func (tv *TopicsView) Focus() {
	FocusView(tv.View().Name())
}

func (tv *TopicsView) Init() {
	tv.view = GetViewPosition(constants.TopicViewName).Set()

	tv.view.FrameColor = gocui.ColorBlue
	tv.view.Title = withSurroundingSpaces("Topics")
	tv.view.TitleColor = gocui.ColorBlue

	sizeX, sizeY := tv.view.Size()
	tv.tableRenderer = NewTableRenderer[*core.Topic]()
	titles := []string{
		"Name",
		"Last Modified",
	}
	colProportions := []float64{
		0.5,
		0.5,
	}
	tv.tableRenderer.InitTable(sizeX, sizeY, titles, colProportions)

	events.AddEventListener(constants.TopicChangeEventName, func(_ string) {
		tv.refreshTopics()
		wv := GetWorkspacesView()
		wv.refreshWorkspaces()
		RenderView(tv)
		RenderView(wv)
	})

	tv.refreshTopics()

	if selectedWorkspace := Api().Core.GetSelectedWorkspace(); selectedWorkspace != nil {
		tv.selectTopicByName(selectedWorkspace.Topic.Name)
	}

	moveRight := func() {
		if Api().Core.GetTopicCount() > 0 {
			GetWorkspacesView().Focus()
		}
	}

	tv.view.KeyBinding().
		set('j', "Move down", func() {
			tv.tableRenderer.Down()
			events.Emit(constants.WorkspaceChangeEventName)
		}).
		set('k', "Move up", func() {
			tv.tableRenderer.Up()
			events.Emit(constants.WorkspaceChangeEventName)
		}).
		set(gocui.KeyEnter, "Open topic", moveRight).
		set('/', "Search by name", func() {
			OpenEditorDialog(func(s string) {
				tv.search.Set(s)
				tv.view.Subtitle = withSurroundingSpaces("Searching: " + s)
				tv.refreshTopics()
				GetWorkspacesView().refreshWorkspaces()
			}, func() {}, "Search", Small)
		}).
		set(gocui.KeyEsc, "Escape search", func() {
			if tv.search.Get() != "" {
				tv.search.Set("")
				tv.view.Subtitle = ""
				tv.refreshTopics()
				GetWorkspacesView().refreshWorkspaces()
			}
		}).
		set('a', "Create a topic", func() {
			OpenEditorDialog(func(s string) {
				if err := Api().Core.CreateTopic(s); err != nil {
					OpenToastDialogError(err.Error())
					return
				}

				// HACK: when there a is a new topic
				// This will result in the corresponding topic going to the top
				// because we are sorting by modifed time
				tv.tableRenderer.SelectRow(0)
			}, func() {}, "Topic name", Small)
		}).
		set('r', "Rename topic", func() {
			t := tv.getSelectedTopic()
			if t == nil {
				return
			}

			OpenEditorDialogWithDefaultValue(func(s string) {
				if err := Api().Core.RenameTopic(t, s); err != nil {
					OpenToastDialogError(err.Error())
					return
				}
			}, func() {}, "New topic name", Small, t.Name)
		}).
		set('s', "Search for a workspace", func() {
			sd := new(*SearchListDialog[*core.Workspace])
			*sd = OpenSearchListDialog(SearchDialogConfig[*core.Workspace]{
				onSearch: func(s string) ([][]string, []*core.Workspace) {
					workspaces := Api().Core.GetWorkspaces().Sorted().FilterByNameContaining(s)
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
					wv := GetWorkspacesView()
					wv.refreshWorkspaces()
					wv.selectWorkspaceByShortPath(w.ShortPath())

					if *sd != nil {
						(*sd).Close()
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
		set('D', "Delete topic", func() {
			if Api().Core.GetTopicCount() <= 0 {
				return
			}
			OpenConfirmationDialog(func(b bool) {
				if !b {
					return
				}

				if err := Api().Core.DeleteTopic(tv.getSelectedTopic()); err != nil {
					OpenToastDialogError(err.Error())
				}
			}, "Are you sure you want to delete this topic? All its content will be deleted.")
		}).
		set('?', "Toggle cheatsheet", func() {
			OpenHelpView(tv.view.keybindingInfo.toList(), func() {})
		})
}

func (tv *TopicsView) refreshTopics() {
	topics := Api().Core.GetTopics().Sorted()

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

func (tv *TopicsView) getSelectedTopic() *core.Topic {
	_, t := tv.tableRenderer.GetSelectedRow()
	if t != nil {
		return *t
	}
	return nil
}

func (tv *TopicsView) selectTopicByName(name string) {
	tv.tableRenderer.SelectRowByValue(func(t *core.Topic) bool {
		return t.Name == name
	})
}

func (tv *TopicsView) Render() error {
	tv.view.Clear()
	currentViewSelected := tv.view.IsFocused()

	tv.tableRenderer.RenderWithSelectCallBack(tv.view, func(_ int, _ *TableRow[*core.Topic]) bool {
		return currentViewSelected
	})

	return nil
}
