package ui

import (
	"log"
	"mynav/pkg/tui"
)

const (
	WorkspacesView        = "WorkspacesView"
	TmuxSessionView       = "TmuxSessionView"
	TmuxWindowView        = "TmuxWindowView"
	TmuxPreviewView       = "TmuxPreviewView"
	TmuxPaneView          = "TmuxPaneView"
	WorkspaceInfoDialog   = "WorkspaceInfoDialog"
	TopicView             = "TopicsView"
	PortView              = "PortView"
	PsView                = "PsView"
	HelpDialog            = "HelpDialog"
	HeaderView            = "HeaderView"
	GithubRepoView        = "GithubRepoView"
	GithubPrView          = "GithubPrView"
	GithubProfileView     = "GithubProfileView"
	EditorDialog          = "EditorDialog"
	ConfirmationDialog    = "ConfirmationDialog"
	ToastDialog           = "ToastDialogView"
	SearchListDialog1View = "SearchListDialog1"
	SearchListDialog2View = "SearchListDialog2"
)

func getViewPosition(viewName string) *tui.ViewPosition {
	maxX, maxY := tui.ScreenSize()
	positionMap := map[string]*tui.ViewPosition{}

	top := maxY / 16
	bottom := ((maxY * 92) / 100)

	positionMap[WorkspacesView] = tui.NewViewPosition(
		WorkspacesView,
		(maxX/3)+1,
		top,
		maxX-2,
		bottom,
		0,
	)

	positionMap[TopicView] = tui.NewViewPosition(
		TopicView,
		2,
		top,
		maxX/3-1,
		bottom,
		0,
	)

	positionMap[TmuxSessionView] = tui.NewViewPosition(
		TmuxSessionView,
		2,
		top,
		maxX/3-1,
		maxY/2-1, 0,
	)

	positionMap[TmuxWindowView] = tui.NewViewPosition(
		TmuxWindowView,
		(maxX/3)+1,
		top,
		((maxX*2)/3)-1,
		maxY/2-1,
		0,
	)

	positionMap[TmuxPaneView] = tui.NewViewPosition(
		TmuxPaneView,
		((maxX*2)/3)+1,
		top,
		maxX-2,
		maxY/2-1,
		0,
	)

	positionMap[TmuxPreviewView] = tui.NewViewPosition(
		TmuxPreviewView,
		2,
		maxY/2+1,
		maxX-2,
		bottom, 0,
	)

	positionMap[PortView] = tui.NewViewPosition(
		PortView,
		maxX/2+1,
		top,
		maxX-2,
		bottom,
		0,
	)

	positionMap[PsView] = tui.NewViewPosition(
		PsView,
		2,
		top,
		maxX/2-1,
		bottom,
		0,
	)

	positionMap[GithubRepoView] = tui.NewViewPosition(
		GithubRepoView,
		maxX/2+1,
		top,
		maxX-4,
		bottom,
		0,
	)

	positionMap[GithubPrView] = tui.NewViewPosition(
		GithubPrView,
		2,
		maxY/2+1,
		maxX/2-1,
		bottom,
		0,
	)

	positionMap[GithubProfileView] = tui.NewViewPosition(
		GithubProfileView,
		2,
		top,
		maxX/2-1,
		maxY/2-1,
		0,
	)

	positionMap[HeaderView] = tui.NewViewPosition(
		HeaderView,
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

	return p
}
