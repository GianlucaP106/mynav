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
}

func newTopicsState() *TopicsState {
	tvm := &TopicsState{
		viewName: "TopicsView",
	}
	return tvm
}

func (ui *UI) initTopicsView() *gocui.View {
	view := ui.getView(ui.topics.viewName)
	if view != nil {
		return view
	}

	view = ui.setView(ui.topics.viewName)
	_, sizeY := view.Size()
	ui.topics.listRenderer = newListRenderer(0, sizeY/3, 0)

	ui.refreshWorkspaces()

	view.FrameColor = gocui.ColorBlue
	view.Title = withSurroundingSpaces("Topics")
	view.TitleColor = gocui.ColorBlue

	ui.keyBinding(ui.topics.viewName).
		set('j', func() {
			ui.topics.listRenderer.increment()
			ui.refreshWorkspaces()
		}).
		set('k', func() {
			ui.topics.listRenderer.decrement()
			ui.refreshWorkspaces()
		}).
		set('a', func() {
			ui.openEditorDialog(func(s string) {
				if err := ui.controller.CreateTopic(s); err != nil {
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
			if ui.controller.GetTopicCount() <= 0 {
				return
			}
			ui.openConfirmationDialog(func(b bool) {
				if b {
					ui.controller.DeleteTopic(ui.getSelectedTopic())
					ui.refreshWorkspaces()
				}
			}, "Are you sure you want to delete this topic? All its content will be deleted.")
		}).
		set(gocui.KeyEnter, func() {
			if ui.controller.GetTopicCount() > 0 {
				ui.setFocusedFsView(ui.workspaces.viewName)
			}
		})

	return view
}

func (ui *UI) refreshTopics() {
	newListSize := ui.controller.GetTopicCount()
	if newListSize != ui.topics.listRenderer.listSize {
		ui.topics.listRenderer.setListSize(newListSize)
	}
}

func (ui *UI) getSelectedTopic() *api.Topic {
	return ui.controller.GetTopic(ui.topics.listRenderer.selected)
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
	topics := ui.controller.GetTopics()
	if topics.Len() <= 0 {
		return
	}

	ui.refreshTopics()

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
