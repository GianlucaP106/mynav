package ui

import (
	"log"
	"mynav/pkg/constants"
	"mynav/pkg/tui"
)

func GetViewPosition(viewName string) *tui.ViewPosition {
	maxX, maxY := tui.ScreenSize()
	positionMap := map[string]*tui.ViewPosition{}

	top := maxY / 20
	bottom := ((maxY * 92) / 100)

	positionMap[constants.WorkspacesViewName] = tui.NewViewPosition(
		constants.WorkspacesViewName,
		(maxX/3)+1,
		top,
		maxX-2,
		bottom,
		0,
	)

	positionMap[constants.TopicViewName] = tui.NewViewPosition(
		constants.TopicViewName,
		2,
		top,
		maxX/3-1,
		bottom,
		0,
	)

	positionMap[constants.TmuxSessionViewName] = tui.NewViewPosition(
		constants.TmuxSessionViewName,
		2,
		top,
		maxX/3-1,
		maxY/2-1, 0,
	)

	positionMap[constants.TmuxWindowViewName] = tui.NewViewPosition(
		constants.TmuxWindowViewName,
		(maxX/3)+1,
		top,
		((maxX*2)/3)-1,
		maxY/2-1,
		0,
	)

	positionMap[constants.TmuxPaneViewName] = tui.NewViewPosition(
		constants.TmuxPaneViewName,
		((maxX*2)/3)+1,
		top,
		maxX-2,
		maxY/2-1,
		0,
	)

	positionMap[constants.TmuxPreviewViewName] = tui.NewViewPosition(
		constants.TmuxPreviewViewName,
		2,
		maxY/2+1,
		maxX-2,
		bottom, 0,
	)

	positionMap[constants.PortViewName] = tui.NewViewPosition(
		constants.PortViewName,
		maxX/2+1,
		top,
		maxX-2,
		bottom,
		0,
	)

	positionMap[constants.PsViewName] = tui.NewViewPosition(
		constants.PsViewName,
		2,
		top,
		maxX/2-1,
		bottom,
		0,
	)

	positionMap[constants.GithubRepoViewName] = tui.NewViewPosition(
		constants.GithubRepoViewName,
		maxX/2+1,
		top,
		maxX-4,
		bottom,
		0,
	)

	positionMap[constants.GithubPrViewName] = tui.NewViewPosition(
		constants.GithubPrViewName,
		2,
		maxY/2+1,
		maxX/2-1,
		bottom,
		0,
	)

	positionMap[constants.GithubProfileViewName] = tui.NewViewPosition(
		constants.GithubProfileViewName,
		2,
		top,
		maxX/2-1,
		maxY/2-1,
		0,
	)

	positionMap[constants.HeaderViewName] = tui.NewViewPosition(
		constants.HeaderViewName,
		2,
		0,
		maxX-2,
		2,
		0,
	)

	p := positionMap[viewName]
	if p == nil {
		log.Panicln("invalid view")
	}

	p.SetName(viewName)

	return p
}
