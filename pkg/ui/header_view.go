package ui

import (
	"fmt"
	"mynav/pkg/tui"

	"github.com/gookit/color"
)

type headerView struct {
	view *tui.View
}

var _ viewable = new(headerView)

func newHeaderView() *headerView {
	return &headerView{}
}

func getHeaderView() *headerView {
	return getViewable[*headerView]()
}

func (hv *headerView) getView() *tui.View {
	return hv.view
}

func (h *headerView) refresh() {}

func (hv *headerView) init() {
	hv.view = getViewPosition(HeaderView).Set()
	hv.view.Frame = false
}

func (hv *headerView) render() error {
	hv.view.Clear()
	hv.view.Resize(getViewPosition(hv.view.Name()))
	screenX, screenY := tui.ScreenSize()
	if screenY < 30 || screenX < 50 {
		hv.view.Clear()
		return nil
	}

	selectedTabName := getMainTabGroup().GetSelectedTab().Frame.Name()
	isMainTab := selectedTabName == "main"
	sep := tui.WithSurroundingSpaces("- ")
	line := ""
	line += "Tab: " + selectedTabName
	if isMainTab {
		line += " " + sep
	}

	if isMainTab {
		if w := getApi().Core.GetSelectedWorkspace(); w != nil {
			line += "Last workspace: " + w.ShortPath()
		}
	}

	sizeX, _ := hv.view.Size()
	s := color.New(color.Yellow, color.Bold)
	line = s.Sprint(line)
	line = tui.Display(line, tui.LeftAlign, sizeX)
	fmt.Fprintln(hv.view, line)

	return nil
}
