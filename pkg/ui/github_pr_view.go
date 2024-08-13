package ui

import (
	"fmt"
	"mynav/pkg/constants"
	"mynav/pkg/events"
	"mynav/pkg/github"
	"mynav/pkg/system"

	"github.com/awesome-gocui/gocui"
)

type GithubPrView struct {
	view          *View
	tableRenderer *TableRenderer[*github.GithubPullRequest]
}

var _ Viewable = new(GithubPrView)

func NewGithubPrView() *GithubPrView {
	return &GithubPrView{}
}

func GetGithubPrView() *GithubPrView {
	return GetViewable[*GithubPrView]()
}

func (g *GithubPrView) View() *View {
	return g.view
}

func (g *GithubPrView) Focus() {
	FocusView(g.View().Name())
}

func (g *GithubPrView) Init() {
	g.view = GetViewPosition(constants.GithubPrViewName).Set()

	g.view.Title = "Pull Requests"
	g.view.TitleColor = gocui.ColorBlue
	g.view.FrameColor = gocui.ColorGreen

	sizeX, sizeY := g.view.Size()
	g.tableRenderer = NewTableRenderer[*github.GithubPullRequest]()
	g.tableRenderer.InitTable(
		sizeX,
		sizeY,
		[]string{
			"Repo name",
			"Pr title",
			"Relation",
		},
		[]float64{
			0.30,
			0.30,
			0.40,
		})

	events.AddEventListener(constants.GithubPrsChangesEventName, func(_ string) {
		g.refresh()
		RenderView(g)
	})

	g.view.KeyBinding().
		set('j', "Move down", func() {
			g.tableRenderer.Down()
		}).
		set('k', "Move up", func() {
			g.tableRenderer.Up()
		}).
		set('o', "Open browser to PR", func() {
			pr := g.getSelectedPr()
			if pr == nil {
				return
			}

			system.OpenBrowser(pr.GetHTMLURL())
		}).
		set('?', "Toggle cheatsheet", func() {
			OpenHelpView(g.view.keybindingInfo.toList(), func() {})
		}).
		set(gocui.KeyCtrlL, "Focus "+constants.GithubRepoViewName, func() {
			GetGithubRepoView().Focus()
		}).
		set(gocui.KeyArrowRight, "Focus "+constants.GithubRepoViewName, func() {
			GetGithubRepoView().Focus()
		})
}

func (g *GithubPrView) getSelectedPr() *github.GithubPullRequest {
	_, pr := g.tableRenderer.GetSelectedRow()
	if pr != nil {
		return *pr
	}

	return nil
}

func (g *GithubPrView) refresh() {
	if !Api().Github.IsAuthenticated() {
		return
	}

	gpr := Api().Github.GetUserPullRequests()

	rows := make([][]string, 0)
	rowValues := make([]*github.GithubPullRequest, 0)
	for _, pr := range gpr {
		rowValues = append(rowValues, pr)
		rows = append(rows, []string{
			pr.Repo.GetName(),
			pr.GetTitle(),
			pr.Relation,
		})
	}

	g.tableRenderer.FillTable(rows, rowValues)
}

func (g *GithubPrView) Render() error {
	g.view.Clear()
	if !Api().Github.IsAuthenticated() {
		fmt.Fprintln(g.view, "Not authenticated")
		return nil
	}

	isFocused := g.view.IsFocused()
	g.tableRenderer.RenderWithSelectCallBack(g.view, func(_ int, _ *TableRow[*github.GithubPullRequest]) bool {
		return isFocused
	})

	if Api().Github.IsLoading() {
		fmt.Fprintln(g.view, "Loading...")
	} else if g.tableRenderer.GetTableSize() == 0 {
		fmt.Fprintln(g.view, "Nothing to show")
	}

	return nil
}
