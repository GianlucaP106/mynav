package app

import (
	"github.com/GianlucaP106/mynav/pkg/tui"
)

// Views.
const (
	HeaderView        = "HeaderView"
	TopicView         = "TopicsView"
	WorkspacesView    = "WorkspacesView"
	SessionsView      = "SessionsView"
	WorkspaceInfoView = "WorkspaceInfoView"
	PreviewView       = "TmuxPreviewView"
	Header1View       = "Header1View"
	Header2View       = "Header2View"
	Header3View       = "Header3View"
	Header4View       = "Header4View"
)

// Dialogs.
const (
	EditorDialog           = "EditorDialog"
	ConfirmationDialog     = "ConfirmationDialog"
	ToastDialog            = "ToastDialogView"
	HelpDialog             = "HelpDialog"
	SearchListDialog1View  = "SearchListDialog1"
	SearchListDialog2View  = "SearchListDialog2"
	SearchListDialog3View  = "SearchListDialog3"
	SearchListDialogBgView = "SearchListDialogBg"
)

func getViewPosition(viewName string) *tui.ViewPosition {
	maxX, maxY := a.ui.Size()
	positionMap := map[string]*tui.ViewPosition{}

	thirdX := maxX / 3

	// header
	positionMap[HeaderView] = tui.NewViewPosition(
		HeaderView,
		0, 0,
		thirdX/3-1, 2,
		0,
	)

	positionMap[Header2View] = tui.NewViewPosition(
		Header2View,
		thirdX/3, 0,
		(11*thirdX/20)-1, 2,
		0,
	)

	positionMap[Header3View] = tui.NewViewPosition(
		Header3View,
		(11 * thirdX / 20), 0,
		(77*thirdX/100)-1, 2,
		0,
	)

	positionMap[Header4View] = tui.NewViewPosition(
		Header4View,
		(77 * thirdX / 100), 0,
		maxX/3-1, 2,
		0,
	)

	// topics
	positionMap[TopicView] = tui.NewViewPosition(
		TopicView,
		0, 3,
		maxX/3-1, maxY/3-1,
		0,
	)

	// workspaces
	positionMap[WorkspacesView] = tui.NewViewPosition(
		WorkspacesView,
		0, maxY/3,
		maxX/3-1, 2*maxY/3-1,
		0,
	)

	// sessions
	positionMap[SessionsView] = tui.NewViewPosition(
		SessionsView,
		0, 2*maxY/3,
		maxX/3-1, maxY-1,
		0,
	)

	// workspace info
	positionMap[WorkspaceInfoView] = tui.NewViewPosition(
		WorkspaceInfoView,
		maxX/3, 0,
		maxX-1, 6,
		0,
	)

	// preview
	positionMap[PreviewView] = tui.NewViewPosition(
		PreviewView,
		maxX/3, 7,
		maxX-1, maxY-1,
		0,
	)

	p := positionMap[viewName]
	// if p == nil {
	// 	log.Panic("invalid view")
	// }

	return p
}
