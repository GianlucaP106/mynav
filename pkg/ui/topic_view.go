package ui

import (
	"mynav/pkg/core"
	"strings"

	"github.com/awesome-gocui/gocui"
)

type TopicsView struct {
	tableRenderer *TableRenderer
	search        string
	topics        core.Topics
}

const TopicViewName = "TopicsView"

var _ View = &TopicsView{}

func newTopicsView() *TopicsView {
	tvm := &TopicsView{}
	return tvm
}

func (tv *TopicsView) RequiresManager() bool {
	return false
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

func (tv *TopicsView) Name() string {
	return TopicViewName
}

func (tv *TopicsView) Init(ui *UI) {
	if GetInternalView(tv.Name()) != nil {
		return
	}

	view := SetViewLayout(tv.Name())

	view.FrameColor = gocui.ColorBlue
	view.Title = withSurroundingSpaces("Topics")
	view.TitleColor = gocui.ColorBlue

	sizeX, sizeY := view.Size()
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
		topicName := strings.Split(selectedWorkspace.ShortPath(), "/")[0]
		tv.selectTopicByName(topicName)
	}

	moveRight := func() {
		if Api().Core.GetTopicCount() > 0 {
			ui.FocusWorkspacesView()
		}
	}

	moveDown := func() {
		ui.FocusPortView()
	}

	KeyBinding(tv.Name()).
		set('j', func() {
			tv.tableRenderer.Down()
			wv := GetView[*WorkspacesView](ui)
			wv.refreshWorkspaces()
		}).
		set('k', func() {
			tv.tableRenderer.Up()
			wv := GetView[*WorkspacesView](ui)
			wv.refreshWorkspaces()
		}).
		set(gocui.KeyEnter, moveRight).
		set(gocui.KeyArrowRight, moveRight).
		set(gocui.KeyCtrlL, moveRight).
		set(gocui.KeyArrowDown, moveDown).
		set(gocui.KeyCtrlJ, moveDown).
		set('/', func() {
			GetDialog[*EditorDialog](ui).Open(func(s string) {
				tv.search = s
				ui.RefreshMainView()
			}, func() {}, "Search", Small)
		}).
		set(gocui.KeyEsc, func() {
			if tv.search != "" {
				tv.search = ""
				ui.RefreshMainView()
			}
		}).
		set('a', func() {
			GetDialog[*EditorDialog](ui).Open(func(s string) {
				if err := Api().Core.CreateTopic(s); err != nil {
					GetDialog[*ToastDialog](ui).OpenError(err.Error())
					return
				}

				// HACK: when there a is a new topic
				// This will result in the corresponding topic going to the top
				// because we are sorting by modifed time
				tv.tableRenderer.SetSelectedRow(0)
				ui.RefreshMainView()
			}, func() {}, "Topic name", Small)
		}).
		set('r', func() {
			t := tv.getSelectedTopic()
			if t == nil {
				return
			}

			GetDialog[*EditorDialog](ui).Open(func(s string) {
				if err := Api().Core.RenameTopic(t, s); err != nil {
					GetDialog[*ToastDialog](ui).OpenError(err.Error())
					return
				}
				ui.RefreshMainView()
			}, func() {}, "New topic name", Small)
		}).
		set('D', func() {
			if Api().Core.GetTopicCount() <= 0 {
				return
			}
			GetDialog[*ConfirmationDialog](ui).Open(func(b bool) {
				if b {
					Api().Core.DeleteTopic(tv.getSelectedTopic())
					ui.RefreshMainView()
				}
			}, "Are you sure you want to delete this topic? All its content will be deleted.")
		}).
		set('?', func() {
			GetDialog[*HelpView](ui).Open(topicKeyBindings, func() {})
		})
}

func (tv *TopicsView) Render(ui *UI) error {
	view := GetInternalView(tv.Name())
	if view == nil {
		tv.Init(ui)
		view = GetInternalView(tv.Name())
	}

	if tv.search != "" {
		view.Subtitle = withSurroundingSpaces("Searching: " + tv.search)
	} else {
		view.Subtitle = ""
	}

	view.Clear()
	currentViewSelected := false
	if v := GetFocusedView(); v != nil && v.Name() == tv.Name() {
		currentViewSelected = true
	}

	tv.tableRenderer.RenderWithSelectCallBack(view, func(_ int, _ *TableRow) bool {
		return currentViewSelected
	})
	return nil
}
