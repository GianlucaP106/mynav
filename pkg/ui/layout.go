package ui

import "log"

func SetViewLayout(viewName string) *View {
	maxX, maxY := ScreenSize()

	views := map[string]func() *View{}
	views[WorkspacesViewName] = func() *View {
		view, _ := SetView(WorkspacesViewName, (maxX/3)+1, 4, maxX-2, (maxY/2)-2, 0)
		return view
	}

	views[TmuxSessionViewName] = func() *View {
		view, _ := SetView(TmuxSessionViewName, (maxX/3)+1, (maxY/2)-1, maxX-2, maxY-4, 0)
		return view
	}

	views[TopicViewName] = func() *View {
		view, _ := SetView(TopicViewName, 2, 4, maxX/3-1, (maxY/2)-2, 0)
		return view
	}

	views[PortViewName] = func() *View {
		view, _ := SetView(PortViewName, 2, (maxY/2)-1, maxX/3-1, maxY-4, 0)
		return view
	}

	views[GithubPrViewName] = func() *View {
		view := SetCenteredView(GithubPrViewName, 75, 20, 0)
		return view
	}

	views[HeaderViewName] = func() *View {
		view, _ := SetView(HeaderViewName, 2, 1, maxX-2, 3, 0)
		return view
	}

	f := views[viewName]
	if f == nil {
		log.Panicln("invalid view")
	}

	return f()
}
