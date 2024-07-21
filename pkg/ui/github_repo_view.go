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

	g.view.Title = "Repositories"
	g.view.TitleColor = gocui.ColorBlue
	g.view.FrameColor = gocui.ColorGreen

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
		g.refreshRepos()
		RenderView(g)
	})

	moveRight := func() {
		GetGithubPrView().Focus()
	}

	g.view.KeyBinding().
		set('j', func() {
			g.tableRenderer.Down()
		}, "Move down").
		set('k', func() {
			g.tableRenderer.Up()
		}, "Move up").
		set(gocui.KeyArrowRight, moveRight, "Focus PR View").
		set(gocui.KeyCtrlL, moveRight, "Focus PR View").
		set('?', func() {
			OpenHelpView(g.view.keybindingInfo.toList(), func() {})
		}, "Toggle cheatsheet")
}

func (g *GithubRepoView) refreshRepos() {
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
