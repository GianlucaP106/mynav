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
	if !ui.controller.Configuration.ConfigInitialized {
		fmt.Fprintln(currentView, blankLine(sizeX))
		fmt.Fprintln(currentView, displayWhiteText("Welcome to mynav, a workspace manager", Center, sizeX))
		return
	}

	line := withSpacePadding("", 10)
	sessionCount, windowCount := ui.controller.WorkspaceManager.GetTmuxStats()
	tmux := " "
	tmux += strconv.Itoa(sessionCount) + withSurroundingSpaces("tmux sessions :")
	tmux += strconv.Itoa(windowCount) + withSurroundingSpaces("windows open")
	tmux = color.New(color.Green).Sprint(tmux)
	line += tmux

	numTopics, numWorkspaces := ui.controller.GetSystemStats()
	generalStats := withSpacePadding("", 5)
	generalStats += strconv.Itoa(numTopics) + withSurroundingSpaces("topics :")
	generalStats += strconv.Itoa(numWorkspaces) + withSurroundingSpaces("workspaces")
	generalStats = color.New(color.Red).Sprint(generalStats)
	line += generalStats

	if w := ui.controller.WorkspaceManager.GetSelectedWorkspace(); w != nil {
		selected := withSpacePadding("", 5)
		selected += withSurroundingSpaces("Last: ")
		selected += color.New(color.Blue).Sprint(w.ShortPath())
		line += selected
	}

	line = display(line, Center, sizeX)
	fmt.Fprintln(currentView, line)
}
