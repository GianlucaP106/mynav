package ui

import (
	"fmt"
	"mynav/pkg/api"

	"github.com/awesome-gocui/gocui"
	"github.com/gookit/color"
)

type TopicsState struct {
	listRenderer *ListRenderer
	viewName     string
	search       string
	topics       api.Topics
}

func newTopicsState() *TopicsState {
	tvm := &TopicsState{
		viewName: "TopicsView",
	}
	return tvm
}

func (ui *UI) initTopicsView() *gocui.View {
	exists := false
	view := ui.getView(ui.topics.viewName)
	exists = view != nil
	view = ui.setView(ui.topics.viewName)

	if ui.topics.search != "" {
		view.Subtitle = withSurroundingSpaces("Searching: " + ui.topics.search)
	} else {
		view.Subtitle = ""
	}

	view.FrameColor = gocui.ColorBlue
	view.Title = withSurroundingSpaces("Topics")
	view.TitleColor = gocui.ColorBlue

	if exists {
		return view
	}
	_, sizeY := view.Size()
	ui.topics.listRenderer = newListRenderer(0, sizeY/3, 0)

	ui.refreshWorkspaces()

	ui.keyBinding(ui.topics.viewName).
		set('j', func() {
			ui.topics.listRenderer.increment()
			ui.refreshWorkspaces()
		}).
		set('k', func() {
			ui.topics.listRenderer.decrement()
			ui.refreshWorkspaces()
		}).
		set('/', func() {
			ui.openEditorDialog(func(s string) {
				ui.topics.search = s
			}, func() {}, "Search", Small)
		}).
		set(gocui.KeyEsc, func() {
			if ui.topics.search != "" {
				ui.topics.search = ""
			}
		}).
		set('a', func() {
			ui.openEditorDialog(func(s string) {
				if _, err := ui.controller.TopicManager.CreateTopic(s); err != nil {
					ui.openToastDialog(err.Error())
					return
				}

				// HACK: when there a is a new topic
				// This will result in the corresponding topic going to the top
				// because we are sorting by modifed time
				ui.topics.listRenderer.setSelected(0)
				ui.refreshTopics()
				ui.refreshWorkspaces()
			}, func() {
			}, "Topic name", Small)
		}).
		set('d', func() {
			if ui.controller.TopicManager.Topics.Len() <= 0 {
				return
			}
			ui.openConfirmationDialog(func(b bool) {
				if b {
					ui.controller.TopicManager.DeleteTopic(ui.getSelectedTopic())
					ui.refreshWorkspaces()
				}
			}, "Are you sure you want to delete this topic? All its content will be deleted.")
		}).
		set(gocui.KeyEnter, func() {
			if ui.controller.TopicManager.Topics.Len() > 0 {
				ui.setFocusedFsView(ui.workspaces.viewName)
			}
		})

	return view
}

func (ui *UI) refreshTopics() {
	topics := ui.controller.TopicManager.Topics.Sorted()

	if ui.topics.search != "" {
		topics = topics.FilterByNameContaining(ui.topics.search)
	}

	ui.topics.topics = topics

	newListSize := ui.topics.topics.Len()
	if newListSize != ui.topics.listRenderer.listSize {
		ui.topics.listRenderer.setListSize(newListSize)
	}
}

func (ui *UI) getSelectedTopic() *api.Topic {
	return ui.controller.TopicManager.Topics.GetTopic(ui.topics.listRenderer.selected)
}

func (ui *UI) formatTopic(topic *api.Topic, selected bool) []string {
	sizeX, _ := ui.getView(ui.topics.viewName).Size()
	style, blankLine := func() (color.Style, string) {
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
		blankLine,
		line,
		blankLine,
	}
	return out
}

func (ui *UI) renderTopicsView() {
	view := ui.initTopicsView()

	view.Clear()
	ui.refreshTopics()
	topics := ui.topics.topics
	content := make([]string, 0)
	ui.topics.listRenderer.forEach(func(idx int) {
		topic := topics[idx]
		selected := idx == ui.topics.listRenderer.selected
		content = append(content, ui.formatTopic(topic, selected)...)
	})
	for _, line := range content {
		fmt.Fprintln(view, line)
	}
}
