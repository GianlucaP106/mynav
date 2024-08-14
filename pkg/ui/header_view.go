package ui

import (
	"fmt"
	"mynav/pkg/constants"

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
	screenX, screenY := ScreenSize()
	if screenY < 30 || screenX < 50 {
		hv.view.Clear()
		return nil
	}

	sizeX, _ := hv.view.Size()
	hv.view.Clear()
	if !Api().Configuration.IsConfigInitialized {
		fmt.Fprintln(hv.view, displayWhiteText("Welcome to mynav, a workspace manager", Center, sizeX))
		return nil
	}

	selectedTabName := GetMainTabGroup().GetSelectedTab().Frame.Name()
	isMainTab := selectedTabName == "main"
	sep := withSurroundingSpaces("- ")
	line := ""
	line += "Tab: " + selectedTabName
	if isMainTab {
		line += " " + sep
	}

	if isMainTab {
		if w := Api().Core.GetSelectedWorkspace(); w != nil {
			line += "Last workspace: " + w.ShortPath()
		}
	}

	s := color.New(color.Yellow, color.Bold)
	line = s.Sprint(line)
	line = display(line, Left, sizeX)
	fmt.Fprintln(hv.view, line)

	return nil
}
