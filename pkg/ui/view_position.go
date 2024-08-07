package ui

import (
	"log"
	"mynav/pkg/constants"
)

type ViewPosition struct {
	viewName string
	x0       int
	y0       int
	x1       int
	y1       int
	overlaps byte
}

func NewViewPosition(
	viewName string,
	x0 int,
	y0 int,
	x1 int,
	y1 int,
	overlaps byte,
) *ViewPosition {
	return &ViewPosition{
		viewName: viewName,
		x0:       x0,
		y0:       y0,
		x1:       x1,
		y1:       y1,
		overlaps: overlaps,
	}
}

func (p *ViewPosition) Set() *View {
	return SetView(p.viewName, p.x0, p.y0, p.x1, p.y1, p.overlaps)
}

func GetViewPosition(viewName string) *ViewPosition {
	maxX, maxY := ScreenSize()
	positionMap := map[string]*ViewPosition{}
	positionMap[constants.WorkspacesViewName] = &ViewPosition{
		x0: (maxX / 3) + 1,
		y0: maxY/2 - maxY/3,
		x1: maxX - 2,
		y1: maxY/2 + maxY/3,
	}

	positionMap[constants.TopicViewName] = &ViewPosition{
		x0: 2,
		y0: maxY/2 - maxY/3,
		x1: maxX/3 - 1,
		y1: maxY/2 + maxY/3,
	}

	positionMap[constants.TmuxSessionViewName] = &ViewPosition{
		x0: (maxX / 3) + 1,
		y0: (maxY / 2) - 1,
		x1: maxX - 2,
		y1: maxY - 4,
	}

	positionMap[constants.PortViewName] = &ViewPosition{
		x0: maxX/2 + 1,
		y0: maxY/2 - maxY/3,
		x1: maxX - 2,
		y1: maxY/2 + maxY/3,
	}

	positionMap[constants.PsViewName] = &ViewPosition{
		x0: 2,
		y0: maxY/2 - maxY/3,
		x1: maxX/2 - 1,
		y1: maxY/2 + maxY/3,
	}

	positionMap[constants.GithubRepoViewName] = &ViewPosition{
		x0: maxX/2 + 1,
		y0: maxY/2 - maxY/3,
		x1: maxX - 4,
		y1: maxY/2 + maxY/3,
	}

	positionMap[constants.GithubPrViewName] = &ViewPosition{
		x0: 4,
		y0: maxY/2 + 1,
		x1: maxX/2 - 1,
		y1: maxY/2 + maxY/3,
	}

	positionMap[constants.GithubProfileViewName] = &ViewPosition{
		x0: 4,
		y0: maxY/2 - maxY/3,
		x1: maxX/2 - 1,
		y1: maxY/2 - 1,
	}

	positionMap[constants.HeaderViewName] = &ViewPosition{
		x0: 2,
		y0: 1,
		x1: maxX - 2,
		y1: 3,
	}

	p := positionMap[viewName]
	if p == nil {
		log.Panicln("invalid view")
	}

	p.viewName = viewName

	return p
}
