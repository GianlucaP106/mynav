package ui

import (
	"fmt"
	"mynav/pkg/constants"
	"strconv"

	"github.com/gookit/color"
)

type HeaderView struct {
	view *View
}

var _ Viewable = new(HeaderView)

func NewHeaderView() *HeaderView {
	return &HeaderView{}
}

func GetHeaderView() *HeaderView {
	return GetViewable[*HeaderView]()
}

func (hv *HeaderView) View() *View {
	return hv.view
}

func (hv *HeaderView) Init() {
	hv.view = GetViewPosition(constants.HeaderViewName).Set()
	hv.view.Frame = false
}

func (hv *HeaderView) Render() error {
	sizeX, _ := hv.view.Size()
	hv.view.Clear()
	if !Api().Configuration.IsConfigInitialized {
		fmt.Fprintln(hv.view, displayWhiteText("Welcome to mynav, a workspace manager", Center, sizeX))
		return nil
	}

	sep := withSurroundingSpaces("- ")
	line := ""
	line += "Tab: " + GetMainTabGroup().GetSelectedTab().Frame.Name() + " " + sep

	if w := Api().Core.GetSelectedWorkspace(); w != nil {
		line += "Last: " + w.ShortPath() + " " + sep
	}

	sessionCount, windowCount := Api().Tmux.GetTmuxStats()
	line += strconv.Itoa(sessionCount) + withSurroundingSpaces("sessions") + sep
	line += strconv.Itoa(windowCount) + withSurroundingSpaces("windows") + sep

	numTopics, numWorkspaces := Api().GetSystemStats()
	line += strconv.Itoa(numTopics) + withSurroundingSpaces("topics") + sep
	line += strconv.Itoa(numWorkspaces) + withSurroundingSpaces("workspaces")

	line = color.Blue.Sprint(line)
	line = display(line, Center, sizeX)
	fmt.Fprintln(hv.view, line)

	return nil
}
