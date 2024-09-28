package ui

import (
	"fmt"
	"mynav/pkg/core"
	"mynav/pkg/system"
	"mynav/pkg/tui"

	"github.com/awesome-gocui/gocui"
	"github.com/google/go-github/v62/github"
)

type githubPrView struct {
	view          *tui.View
	tableRenderer *tui.TableRenderer[*github.PullRequest]
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
	g.view = getViewPosition(GithubPrView).Set()

	g.view.Title = tui.WithSurroundingSpaces("Pull Requests")

	styleView(g.view)

	sizeX, sizeY := g.view.Size()
	g.tableRenderer = tui.NewTableRenderer[*github.PullRequest]()
	g.tableRenderer.InitTable(
		sizeX,
		sizeY,
		[]string{
			"Repo Name",
			"Pr title",
			"Relation",
		},
		[]float64{
			0.30,
			0.30,
			0.40,
		})

	g.refresh()

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
		Set('u', "Copy PR URL to clipboard", func() {
			pr := g.getSelectedPr()
			if pr == nil {
				return
			}

			url := pr.GetHTMLURL()
			system.CopyToClip(url)
			openToastDialog(url, toastDialogNeutralType, "Copied PR URL to clipboard", func() {})
		}).
		Set('R', "Refetch all github data", func() {
			getGithubProfileView().refetchData()
		}).
		Set(gocui.KeyCtrlL, "Focus "+GithubRepoView, func() {
			getGithubRepoView().focus()
		}).
		Set(gocui.KeyArrowRight, "Focus "+GithubRepoView, func() {
			getGithubRepoView().focus()
		}).
		Set('?', "Toggle cheatsheet", func() {
			openHelpDialog(g.view.GetKeybindings(), func() {})
		})
}

func (g *githubPrView) getSelectedPr() *github.PullRequest {
	_, pr := g.tableRenderer.GetSelectedRow()
	if pr != nil {
		return *pr
	}

	return nil
}

func (g *githubPrView) refresh() {
	if !api().Github.IsAuthenticated() {
		return
	}

	gpr := api().Github.GetUserPullRequests()

	principal := api().Github.GetPrincipal()
	rows := make([][]string, 0)
	for _, pr := range gpr {
		_, relation := api().Github.GetPrRelation(pr, principal)
		rows = append(rows, []string{
			core.ExtractRepoFromPrUrl(pr.GetHTMLURL()),
			pr.GetTitle(),
			relation,
		})
	}

	g.tableRenderer.FillTable(rows, gpr)
}

func (g *githubPrView) render() error {
	g.view.Clear()
	g.view.Resize(getViewPosition(g.view.Name()))
	if !api().Github.IsAuthenticated() {
		fmt.Fprintln(g.view, "Not authenticated")
		return nil
	}

	isFocused := g.view.IsFocused()
	g.tableRenderer.RenderWithSelectCallBack(g.view, func(_ int, _ *tui.TableRow[*github.PullRequest]) bool {
		return isFocused
	})

	if getGithubProfileView().isFetchingData.Get() {
		fmt.Fprintln(g.view, "Loading...")
	} else if g.tableRenderer.GetTableSize() == 0 {
		fmt.Fprintln(g.view, "Nothing to show")
	}

	return nil
}
