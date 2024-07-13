package ui

import "log"

type ViewPosition struct {
	viewName string
	x0       int
	x1       int
	y0       int
	y1       int
	overlaps byte
}

func GetViewPosition(viewName string) *ViewPosition {
	maxX, maxY := ScreenSize()
	positionMap := map[string]*ViewPosition{}
	positionMap[WorkspacesViewName] = &ViewPosition{
		x0: (maxX / 3) + 1,
		y0: maxY/2 - 15,
		x1: maxX - 2,
		y1: maxY/2 + 15,
	}

	positionMap[TopicViewName] = &ViewPosition{
		x0: 2,
		y0: maxY/2 - 15,
		x1: maxX/3 - 1,
		y1: maxY/2 + 15,
	}

	positionMap[TmuxSessionViewName] = &ViewPosition{
		x0: (maxX / 3) + 1,
		y0: (maxY / 2) - 1,
		x1: maxX - 2,
		y1: maxY - 4,
	}

	positionMap[PortViewName] = &ViewPosition{
		x0: 2,
		y0: (maxY / 2) - 1,
		x1: maxX/3 - 1,
		y1: maxY - 4,
	}

	positionMap[GithubPrViewName] = &ViewPosition{
		x0: maxX/2 + 1,
		y0: maxY/2 - 10,
		x1: maxX - 4,
		y1: maxY/2 + 10,
	}

	positionMap[GithubRepoViewName] = &ViewPosition{
		x0: 4,
		y0: maxY/2 - 10,
		x1: maxX/2 - 1,
		y1: maxY/2 + 10,
	}
	positionMap[HeaderViewName] = &ViewPosition{
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

func (p *ViewPosition) Set() *View {
	return SetView(p.viewName, p.x0, p.y0, p.x1, p.y1, p.overlaps)
}
