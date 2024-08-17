package ui

import (
	"fmt"
	"mynav/pkg/events"
	"mynav/pkg/github"
	"mynav/pkg/tui"

	"github.com/awesome-gocui/gocui"
)

type githubRepoView struct {
	view          *tui.View
	tableRenderer *tui.TableRenderer[*github.GithubRepository]
}

var _ viewable = new(githubRepoView)

func newGithubRepoView() *githubRepoView {
	return &githubRepoView{}
}

func getGithubRepoView() *githubRepoView {
	return getViewable[*githubRepoView]()
}

func (g *githubRepoView) getView() *tui.View {
	return g.view
}

func (g *githubRepoView) Focus() {
	focusView(g.getView().Name())
}

func (g *githubRepoView) init() {
	g.view = getViewPosition(GithubRepoView).Set()

	g.view.Title = tui.WithSurroundingSpaces("Repositories")

	styleView(g.view)

	g.tableRenderer = tui.NewTableRenderer[*github.GithubRepository]()
	sizeX, sizeY := g.view.Size()

	g.tableRenderer.InitTable(
		sizeX,
		sizeY,
		[]string{
			"Repo name",
			"Owner",
		},
		[]float64{
			0.5,
			0.5,
		},
	)

	events.AddEventListener(events.GithubReposChangesEvent, func(s string) {
		g.refresh()
		renderView(g)
	})

	moveRight := func() {
		getGithubPrView().Focus()
	}

	g.view.KeyBinding().
		Set('j', "Move down", func() {
			g.tableRenderer.Down()
		}).
		Set('k', "Move up", func() {
			g.tableRenderer.Up()
		}).
		Set(gocui.KeyArrowRight, "Focus PR View", moveRight).
		Set(gocui.KeyCtrlL, "Focus PR View", moveRight).
		Set('?', "Toggle cheatsheet", func() {
			OpenHelpDialog(g.view.GetKeybindings(), func() {})
		})
}

func (g *githubRepoView) refresh() {
	if !getApi().Github.IsAuthenticated() {
		return
	}

	repos := getApi().Github.GetUserRepos()

	rows := make([][]string, 0)
	rowValues := make([]*github.GithubRepository, 0)
	for _, repo := range repos {
		rowValues = append(rowValues, repo)
		rows = append(rows, []string{
			repo.GetName(),
			repo.GetOwner().GetLogin(),
		})
	}
	g.tableRenderer.FillTable(rows, rowValues)
}

func (g *githubRepoView) render() error {
	if !getApi().Github.IsAuthenticated() {
		g.view.Clear()
		fmt.Fprintln(g.view, "Not authenticated")
		return nil
	}

	g.view.Clear()
	isFocused := g.view.IsFocused()
	g.view = getViewPosition(g.view.Name()).Set()

	g.tableRenderer.RenderWithSelectCallBack(g.view, func(_ int, _ *tui.TableRow[*github.GithubRepository]) bool {
		return isFocused
	})

	if getApi().Github.IsLoading() {
		fmt.Fprintln(g.view, "Loading...")
	} else if g.tableRenderer.GetTableSize() == 0 {
		fmt.Fprintln(g.view, "No repos to display")
	}

	return nil
}
