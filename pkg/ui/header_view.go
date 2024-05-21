package ui

import (
	"fmt"
	"strconv"

	"github.com/gookit/color"
)

const HeaderStateName = "HeaderView"

type HeaderView struct{}

var _ View = &HeaderView{}

func newHeaderState() *HeaderView {
	titleView := &HeaderView{}

	return titleView
}

func (hv *HeaderView) Name() string {
	return HeaderStateName
}

func (hv *HeaderView) RequiresManager() bool {
	return true
}

func (hv *HeaderView) Init(ui *UI) {
	if GetInternalView(hv.Name()) != nil {
		return
	}

	SetViewLayout(hv.Name())
}

func (hv *HeaderView) Render(ui *UI) error {
	currentView := GetInternalView(hv.Name())
	if currentView == nil {
		return nil
	}

	sizeX, _ := currentView.Size()
	currentView.Clear()
	fmt.Fprintln(currentView, blankLine(sizeX))
	if !Api().IsConfigInitialized {
		fmt.Fprintln(currentView, displayWhiteText("Welcome to mynav, a workspace manager", Center, sizeX))
		return nil
	}

	line := withSpacePadding("", 10)

	if w := Api().GetSelectedWorkspace(); w != nil {
		selected := withSpacePadding("", 5)
		selected += withSurroundingSpaces("Last seen: ")
		selected += color.New(color.Blue).Sprint(w.ShortPath())
		line += selected
	}

	sessionCount, windowCount := Api().GetTmuxStats()
	tmux := withSpacePadding("", 5)
	tmux += strconv.Itoa(sessionCount) + withSurroundingSpaces("tmux sessions |")
	tmux += strconv.Itoa(windowCount) + withSurroundingSpaces("windows open")
	tmux = color.New(color.Green).Sprint(tmux)
	line += tmux

	numTopics, numWorkspaces := Api().GetSystemStats()
	generalStats := withSpacePadding("", 5)
	generalStats += strconv.Itoa(numTopics) + withSurroundingSpaces("topics |")
	generalStats += strconv.Itoa(numWorkspaces) + withSurroundingSpaces("workspaces")
	generalStats = color.New(color.Blue).Sprint(generalStats)
	line += generalStats

	line = display(line, Center, sizeX)
	fmt.Fprintln(currentView, line)

	return nil
}
