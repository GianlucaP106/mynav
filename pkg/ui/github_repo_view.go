package ui

import (
	"fmt"
	"mynav/pkg/github"

	"github.com/awesome-gocui/gocui"
)

type GithubRepoView struct {
	view          *View
	tableRenderer *TableRenderer
	repos         []*github.GithubRepository
}

var _ Viewable = new(GithubRepoView)

const GithubRepoViewName = "GithubRepoView"

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
	g.view = GetViewPosition(GithubRepoViewName).Set()

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

	KeyBinding(g.view.Name()).
		set('j', func() {
			g.tableRenderer.Down()
		}).
		set('k', func() {
			g.tableRenderer.Up()
		}).
		set(gocui.KeyArrowRight, moveRight).
		set(gocui.KeyCtrlL, moveRight)
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

	isFocused := IsViewFocused(g.view)
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
