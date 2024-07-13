package ui

import (
	"mynav/pkg/core"

	"github.com/awesome-gocui/gocui"
)

type TopicsView struct {
	view          *View
	tableRenderer *TableRenderer
	search        string
	topics        core.Topics
}

var _ Viewable = new(TopicsView)

const TopicViewName = "TopicsView"

func NewTopicsView() *TopicsView {
	return &TopicsView{}
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
	tv.view = GetViewPosition(TopicViewName).Set()

	tv.view.FrameColor = gocui.ColorBlue
	tv.view.Title = withSurroundingSpaces("Topics")
	tv.view.TitleColor = gocui.ColorBlue

	sizeX, sizeY := tv.view.Size()
	tv.tableRenderer = NewTableRenderer()
	titles := []string{
		"Name",
		"Last Modified",
	}
	colProportions := []float64{
		0.5,
		0.5,
	}
	tv.tableRenderer.InitTable(sizeX, sizeY, titles, colProportions)
	tv.refreshTopics()

	if selectedWorkspace := Api().Core.GetSelectedWorkspace(); selectedWorkspace != nil {
		tv.selectTopicByName(selectedWorkspace.Topic.Name)
	}

	moveRight := func() {
		if Api().Core.GetTopicCount() > 0 {
			GetWorkspacesView().Focus()
		}
	}

	KeyBinding(tv.view.Name()).
		set('j', func() {
			tv.tableRenderer.Down()
			GetWorkspacesView().refreshWorkspaces()
		}).
		set('k', func() {
			tv.tableRenderer.Up()
			GetWorkspacesView().refreshWorkspaces()
		}).
		set(gocui.KeyEnter, moveRight).
		set(gocui.KeyArrowRight, moveRight).
		set(gocui.KeyCtrlL, moveRight).
		set('/', func() {
			OpenEditorDialog(func(s string) {
				tv.search = s
				tv.view.Subtitle = withSurroundingSpaces("Searching: " + tv.search)
				RefreshAllData()
			}, func() {}, "Search", Small)
		}).
		set(gocui.KeyEsc, func() {
			if tv.search != "" {
				tv.search = ""
				tv.view.Subtitle = ""
				RefreshAllData()
			}
		}).
		set('a', func() {
			OpenEditorDialog(func(s string) {
				if err := Api().Core.CreateTopic(s); err != nil {
					OpenToastDialogError(err.Error())
					return
				}

				// HACK: when there a is a new topic
				// This will result in the corresponding topic going to the top
				// because we are sorting by modifed time
				tv.tableRenderer.SetSelectedRow(0)
				RefreshAllData()
			}, func() {}, "Topic name", Small)
		}).
		set('r', func() {
			t := tv.getSelectedTopic()
			if t == nil {
				return
			}

			OpenEditorDialogWithDefaultValue(func(s string) {
				if err := Api().Core.RenameTopic(t, s); err != nil {
					OpenToastDialogError(err.Error())
					return
				}

				RefreshAllData()
			}, func() {}, "New topic name", Small, t.Name)
		}).
		set('D', func() {
			if Api().Core.GetTopicCount() <= 0 {
				return
			}
			OpenConfirmationDialog(func(b bool) {
				if b {
					Api().Core.DeleteTopic(tv.getSelectedTopic())
					RefreshAllData()
				}
			}, "Are you sure you want to delete this topic? All its content will be deleted.")
		}).
		set('?', func() {
			OpenHelpView(topicKeyBindings, func() {})
		})
}

func (tv *TopicsView) refreshTopics() {
	topics := Api().Core.GetTopics().Sorted()

	if tv.search != "" {
		topics = topics.FilterByNameContaining(tv.search)
	}

	tv.topics = topics
	tv.syncTopicsToTable()
}

func (tv *TopicsView) syncTopicsToTable() {
	rows := make([][]string, 0)
	for _, topic := range tv.topics {
		rows = append(rows, []string{
			topic.Name,
			topic.GetLastModifiedTimeFormatted(),
		})
	}

	tv.tableRenderer.FillTable(rows)
}

func (tv *TopicsView) getSelectedTopic() *core.Topic {
	if tv.topics.Len() <= 0 {
		return nil
	}

	return tv.topics[tv.tableRenderer.GetSelectedRowIndex()]
}

func (tv *TopicsView) selectTopicByName(name string) {
	for idx, t := range tv.topics {
		if t.Name == name {
			tv.tableRenderer.SetSelectedRow(idx)
		}
	}
}

func (tv *TopicsView) Render() error {
	tv.view.Clear()
	currentViewSelected := IsViewFocused(tv.view)

	tv.tableRenderer.RenderWithSelectCallBack(tv.view, func(_ int, _ *TableRow) bool {
		return currentViewSelected
	})

	return nil
}
