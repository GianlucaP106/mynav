package ui

import (
	"fmt"
	"mynav/pkg/constants"
	"mynav/pkg/github"

	"github.com/awesome-gocui/gocui"
)

type GithubRepoView struct {
	view          *View
	tableRenderer *TableRenderer
	repos         []*github.GithubRepository
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

	g.view.Title = "Repository"
	g.view.TitleColor = gocui.ColorBlue
	g.view.FrameColor = gocui.ColorGreen

	g.tableRenderer = NewTableRenderer()
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

	go func() {
		g.refreshRepos()
		UpdateGui(func(_ *Gui) error {
			g.Render()
			return nil
		})
	}()

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

	repos, _ := Api().Github.GetUserReposLocked()
	g.repos = repos

	g.syncReposToTable()
}

func (g *GithubRepoView) syncReposToTable() {
	rows := make([][]string, 0)
	for _, repo := range g.repos {
		rows = append(rows, []string{
			repo.GetName(),
			repo.GetOwner().GetLogin(),
		})
	}
	g.tableRenderer.FillTable(rows)
}

func (g *GithubRepoView) Render() error {
	if !Api().Github.IsAuthenticated() {
		g.view.Clear()
		fmt.Fprintln(g.view, "Not authenticated")
		return nil
	}

	g.view.Clear()
	g.view.Subtitle = "Login: " + Api().Github.GetPrincipalLogin()

	isFocused := g.view.IsFocused()
	g.tableRenderer.render(g.view, func(_ int, _ *TableRow) bool {
		return isFocused
	})

	if g.repos == nil {
		fmt.Fprintln(g.view, "Loading...")
	} else if len(g.repos) == 0 {
		fmt.Fprintln(g.view, "No repos to display")
	}

	return nil
}
