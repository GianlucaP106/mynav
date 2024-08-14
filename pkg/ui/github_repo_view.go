package ui

import (
	"fmt"
	"mynav/pkg/constants"
	"mynav/pkg/events"
	"mynav/pkg/github"

	"github.com/awesome-gocui/gocui"
)

type GithubRepoView struct {
	view          *View
	tableRenderer *TableRenderer[*github.GithubRepository]
}

var _ Viewable = new(GithubRepoView)

func NewGithubRepoView() *GithubRepoView {
	return &GithubRepoView{}
}

func GetGithubRepoView() *GithubRepoView {
	return GetViewable[*GithubRepoView]()
}

func (g *GithubRepoView) View() *View {
	return g.view
}

func (g *GithubRepoView) Focus() {
	FocusView(g.View().Name())
}

func (g *GithubRepoView) Init() {
	g.view = GetViewPosition(constants.GithubRepoViewName).Set()

	g.view.Title = withSurroundingSpaces("Repositories")

	StyleView(g.view)

	g.tableRenderer = NewTableRenderer[*github.GithubRepository]()
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

	events.AddEventListener(constants.GithubReposChangesEventName, func(s string) {
		g.refresh()
		RenderView(g)
	})

	moveRight := func() {
		GetGithubPrView().Focus()
	}

	g.view.KeyBinding().
		set('j', "Move down", func() {
			g.tableRenderer.Down()
		}).
		set('k', "Move up", func() {
			g.tableRenderer.Up()
		}).
		set(gocui.KeyArrowRight, "Focus PR View", moveRight).
		set(gocui.KeyCtrlL, "Focus PR View", moveRight).
		set('?', "Toggle cheatsheet", func() {
			OpenHelpView(g.view.keybindingInfo.toList(), func() {})
		})
}

func (g *GithubRepoView) refresh() {
	if !Api().Github.IsAuthenticated() {
		return
	}

	repos := Api().Github.GetUserRepos()

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

func (g *GithubRepoView) Render() error {
	if !Api().Github.IsAuthenticated() {
		g.view.Clear()
		fmt.Fprintln(g.view, "Not authenticated")
		return nil
	}

	g.view.Clear()

	isFocused := g.view.IsFocused()
	g.tableRenderer.render(g.view, func(_ int, _ *TableRow[*github.GithubRepository]) bool {
		return isFocused
	})

	if Api().Github.IsLoading() {
		fmt.Fprintln(g.view, "Loading...")
	} else if g.tableRenderer.GetTableSize() == 0 {
		fmt.Fprintln(g.view, "No repos to display")
	}

	return nil
}
