package ui

import (
	"fmt"
	"mynav/pkg/api"
	"strings"

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
	if !exists {
		view = ui.setView(ui.topics.viewName)
	}

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
	ui.refreshTopics()

	if selectedWorkspace := ui.api.GetSelectedWorkspace(); selectedWorkspace != nil {
		topicName := strings.Split(selectedWorkspace.ShortPath(), "/")[0]
		ui.selectTopicByName(topicName)
	}

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
				ui.refreshTopics()
				ui.refreshWorkspaces()
			}, func() {}, "Search", Small)
		}).
		set(gocui.KeyEsc, func() {
			if ui.topics.search != "" {
				ui.topics.search = ""
				ui.refreshTopics()
				ui.refreshWorkspaces()
			}
		}).
		set('a', func() {
			ui.openEditorDialog(func(s string) {
				if err := ui.api.CreateTopic(s); err != nil {
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
			if ui.api.GetTopicCount() <= 0 {
				return
			}
			ui.openConfirmationDialog(func(b bool) {
				if b {
					ui.api.DeleteTopic(ui.getSelectedTopic())
					ui.refreshTopics()
					ui.refreshWorkspaces()
				}
			}, "Are you sure you want to delete this topic? All its content will be deleted.")
		}).
		set(gocui.KeyEnter, func() {
			if ui.api.GetTopicCount() > 0 {
				ui.setFocusedFsView(ui.workspaces.viewName)
			}
		}).
		set('?', func() {
			ui.openHelpView(ui.getKeyBindings(ui.topics.viewName))
		})
	return view
}

func (ui *UI) refreshTopics() {
	topics := ui.api.GetTopics().Sorted()

	if ui.topics.search != "" {
		topics = topics.FilterByNameContaining(ui.topics.search)
	}

	ui.topics.topics = topics

	newListSize := ui.topics.topics.Len()
	if ui.topics.listRenderer != nil && newListSize != ui.topics.listRenderer.listSize {
		ui.topics.listRenderer.setListSize(newListSize)
	}
}

func (ui *UI) getSelectedTopic() *api.Topic {
	if ui.topics.topics.Len() <= 0 {
		return nil
	}

	return ui.topics.topics[ui.topics.listRenderer.selected]
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

func (ui *UI) selectTopicByName(name string) {
	for idx, t := range ui.topics.topics {
		if t.Name == name {
			ui.topics.listRenderer.setSelected(idx)
		}
	}
}

func (ui *UI) renderTopicsView() {
	view := ui.initTopicsView()

	view.Clear()
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
