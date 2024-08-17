package ui

import (
	"fmt"
	"mynav/pkg/constants"
	"mynav/pkg/events"
	"mynav/pkg/github"
	"mynav/pkg/system"
	"mynav/pkg/tui"

	"github.com/awesome-gocui/gocui"
)

type githubPrView struct {
	view          *tui.View
	tableRenderer *tui.TableRenderer[*github.GithubPullRequest]
}

var _ viewable = new(githubPrView)

func newGithubPrView() *githubPrView {
	return &githubPrView{}
}

func getGithubPrView() *githubPrView {
	return getViewable[*githubPrView]()
}

func (g *githubPrView) getView() *tui.View {
	return g.view
}

func (g *githubPrView) Focus() {
	focusView(g.getView().Name())
}

func (g *githubPrView) init() {
	g.view = GetViewPosition(constants.GithubPrViewName).Set()

	g.view.Title = tui.WithSurroundingSpaces("Pull Requests")

	tui.StyleView(g.view)

	sizeX, sizeY := g.view.Size()
	g.tableRenderer = tui.NewTableRenderer[*github.GithubPullRequest]()
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
		renderView(g)
	})

	g.view.KeyBinding().
		Set('j', "Move down", func() {
			g.tableRenderer.Down()
		}).
		Set('k', "Move up", func() {
			g.tableRenderer.Up()
		}).
		Set('o', "Open browser to PR", func() {
			pr := g.getSelectedPr()
			if pr == nil {
				return
			}

			system.OpenBrowser(pr.GetHTMLURL())
		}).
		Set('?', "Toggle cheatsheet", func() {
			OpenHelpDialog(g.view.GetKeybindings(), func() {})
		}).
		Set(gocui.KeyCtrlL, "Focus "+constants.GithubRepoViewName, func() {
			getGithubRepoView().Focus()
		}).
		Set(gocui.KeyArrowRight, "Focus "+constants.GithubRepoViewName, func() {
			getGithubRepoView().Focus()
		})
}

func (g *githubPrView) getSelectedPr() *github.GithubPullRequest {
	_, pr := g.tableRenderer.GetSelectedRow()
	if pr != nil {
		return *pr
	}

	return nil
}

func (g *githubPrView) refresh() {
	if !getApi().Github.IsAuthenticated() {
		return
	}

	gpr := getApi().Github.GetUserPullRequests()

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

func (g *githubPrView) render() error {
	g.view.Clear()
	if !getApi().Github.IsAuthenticated() {
		fmt.Fprintln(g.view, "Not authenticated")
		return nil
	}

	isFocused := g.view.IsFocused()
	g.tableRenderer.RenderWithSelectCallBack(g.view, func(_ int, _ *tui.TableRow[*github.GithubPullRequest]) bool {
		return isFocused
	})

	if getApi().Github.IsLoading() {
		fmt.Fprintln(g.view, "Loading...")
	} else if g.tableRenderer.GetTableSize() == 0 {
		fmt.Fprintln(g.view, "Nothing to show")
	}

	return nil
}
