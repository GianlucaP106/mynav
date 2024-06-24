package ui

import (
	"fmt"
	"strconv"

	"github.com/gookit/color"
)

type HeaderView struct {
	view *View
}

const HeaderViewName = "HeaderView"

func NewHeaderView() *HeaderView {
	return &HeaderView{}
}

func (hv *HeaderView) Init() {
	hv.view = SetViewLayout(HeaderViewName)
}

func (hv *HeaderView) Render() error {
	sizeX, _ := hv.view.Size()
	hv.view.Clear()
	fmt.Fprintln(hv.view, blankLine(sizeX))
	if !Api().Core.IsConfigInitialized {
		fmt.Fprintln(hv.view, displayWhiteText("Welcome to mynav, a workspace manager", Center, sizeX))
		return nil
	}

	line := ""
	if w := Api().Core.GetSelectedWorkspace(); w != nil {
		selected := withSpacePadding("", 5)
		selected += withSurroundingSpaces("Last seen: ")
		selected += color.New(color.Blue).Sprint(w.ShortPath())
		line += selected
	}

	sessionCount, windowCount := Api().Tmux.GetTmuxStats()
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
	fmt.Fprintln(hv.view, line)

	return nil
}
