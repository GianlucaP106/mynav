package ui

import (
	"fmt"
	"mynav/pkg/api"
	"strings"

	"github.com/awesome-gocui/gocui"
	"github.com/gookit/color"
)

type TopicsView struct {
	listRenderer *ListRenderer
	search       string
	topics       api.Topics
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
	topics := Api().GetTopics().Sorted()

	if tv.search != "" {
		topics = topics.FilterByNameContaining(tv.search)
	}

	tv.topics = topics

	newListSize := tv.topics.Len()
	if tv.listRenderer != nil && newListSize != tv.listRenderer.listSize {
		tv.listRenderer.setListSize(newListSize)
	}
}

func (tv *TopicsView) getSelectedTopic() *api.Topic {
	if tv.topics.Len() <= 0 {
		return nil
	}

	return tv.topics[tv.listRenderer.selected]
}

func (tv *TopicsView) selectTopicByName(name string) {
	for idx, t := range tv.topics {
		if t.Name == name {
			tv.listRenderer.setSelected(idx)
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

	_, sizeY := view.Size()
	tv.listRenderer = newListRenderer(0, sizeY, 0)
	tv.refreshTopics()

	if selectedWorkspace := Api().GetSelectedWorkspace(); selectedWorkspace != nil {
		topicName := strings.Split(selectedWorkspace.ShortPath(), "/")[0]
		tv.selectTopicByName(topicName)
	}

	moveRight := func() {
		if Api().GetTopicCount() > 0 {
			ui.FocusWorkspacesView()
		}
	}

	moveDown := func() {
		ui.FocusPortView()
	}

	KeyBinding(tv.Name()).
		set('j', func() {
			tv.listRenderer.increment()
			wv := GetView[*WorkspacesView](ui)
			wv.refreshWorkspaces()
		}).
		set('k', func() {
			tv.listRenderer.decrement()
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
				ui.FocusTopicsView()
			}, func() {
				ui.FocusTopicsView()
			}, "Search", Small)
		}).
		set(gocui.KeyEsc, func() {
			if tv.search != "" {
				tv.search = ""
				ui.RefreshMainView()
			}
		}).
		set('a', func() {
			GetDialog[*EditorDialog](ui).Open(func(s string) {
				if err := Api().CreateTopic(s); err != nil {
					GetDialog[*ToastDialog](ui).Open(err.Error(), func() {
						ui.FocusTopicsView()
					})
					return
				}

				// HACK: when there a is a new topic
				// This will result in the corresponding topic going to the top
				// because we are sorting by modifed time
				tv.listRenderer.setSelected(0)
				ui.RefreshMainView()
				ui.FocusTopicsView()
			}, func() {
				ui.FocusTopicsView()
			}, "Topic name", Small)
		}).
		set('r', func() {
			t := tv.getSelectedTopic()
			if t == nil {
				return
			}

			GetDialog[*EditorDialog](ui).Open(func(s string) {
				if err := Api().RenameTopic(t, s); err != nil {
					GetDialog[*ToastDialog](ui).Open(err.Error(), func() {
						ui.FocusTopicsView()
					})
					return
				}
				ui.RefreshMainView()
				ui.FocusTopicsView()
			}, func() {
				ui.FocusTopicsView()
			}, "New topic name", Small)
		}).
		set('D', func() {
			if Api().GetTopicCount() <= 0 {
				return
			}
			GetDialog[*ConfirmationDialog](ui).Open(func(b bool) {
				if b {
					Api().DeleteTopic(tv.getSelectedTopic())
					ui.RefreshMainView()
				}
				ui.FocusTopicsView()
			}, "Are you sure you want to delete this topic? All its content will be deleted.")
		}).
		set('?', func() {
			GetDialog[*HelpView](ui).Open(topicKeyBindings, func() {
				ui.FocusTopicsView()
			})
		})
}

func (tv *TopicsView) formatTopic(topic *api.Topic, selected bool) []string {
	sizeX, _ := GetInternalView(tv.Name()).Size()
	style, _ := func() (color.Style, string) {
		if selected {
			return color.New(color.Black, color.BgCyan), highlightedBlankLine(sizeX + 5) // +5 for extra padding
		}
		return color.New(color.Blue), blankLine(sizeX)
	}()

	modTime := topic.GetLastModifiedTimeFormatted()
	name := withSpacePadding(withSurroundingSpaces(topic.Name), sizeX/3)
	modTime = withSpacePadding(modTime, ((sizeX*2)/3)+5) // +5 for extra padding

	line := style.Sprint(name + modTime)

	out := []string{
		line,
	}
	return out
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

	topics := tv.topics
	content := make([]string, 0)
	tv.listRenderer.forEach(func(idx int) {
		topic := topics[idx]
		selected := (idx == tv.listRenderer.selected) && currentViewSelected
		content = append(content, tv.formatTopic(topic, selected)...)
	})
	for _, line := range content {
		fmt.Fprintln(view, line)
	}
	return nil
}
