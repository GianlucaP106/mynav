package ui

import (
	"fmt"
	"strconv"

	"github.com/awesome-gocui/gocui"
	"github.com/gookit/color"
)

type HeaderState struct {
	viewName string
}

func newHeaderState() *HeaderState {
	titleView := &HeaderState{
		viewName: "HeaderView",
	}

	return titleView
}

func (ui *UI) initHeaderView() *gocui.View {
	view := ui.setView(ui.header.viewName)
	return view
}

func (ui *UI) renderHeaderView() {
	currentView := ui.getView(ui.header.viewName)
	if currentView == nil {
		currentView = ui.initHeaderView()
	}

	sizeX, _ := currentView.Size()
	currentView.Clear()
	fmt.Fprintln(currentView, blankLine(sizeX))
	if !ui.api.IsConfigInitialized {
		fmt.Fprintln(currentView, displayWhiteText("Welcome to mynav, a workspace manager", Center, sizeX))
		return
	}

	line := withSpacePadding("", 10)

	if w := ui.api.GetSelectedWorkspace(); w != nil {
		selected := withSpacePadding("", 5)
		selected += withSurroundingSpaces("Last seen: ")
		selected += color.New(color.Blue).Sprint(w.ShortPath())
		line += selected
	}

	sessionCount, windowCount := ui.api.GetTmuxStats()
	tmux := withSpacePadding("", 5)
	tmux += strconv.Itoa(sessionCount) + withSurroundingSpaces("tmux sessions |")
	tmux += strconv.Itoa(windowCount) + withSurroundingSpaces("windows open")
	tmux = color.New(color.Green).Sprint(tmux)
	line += tmux

	numTopics, numWorkspaces := ui.api.GetSystemStats()
	generalStats := withSpacePadding("", 5)
	generalStats += strconv.Itoa(numTopics) + withSurroundingSpaces("topics |")
	generalStats += strconv.Itoa(numWorkspaces) + withSurroundingSpaces("workspaces")
	generalStats = color.New(color.Red).Sprint(generalStats)
	line += generalStats

	line = display(line, Center, sizeX)
	fmt.Fprintln(currentView, line)
}
