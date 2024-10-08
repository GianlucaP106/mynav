package ui

import (
	"fmt"
	"mynav/pkg/core"
	"mynav/pkg/system"
	"mynav/pkg/tui"
	"sort"
	"strings"

	"github.com/awesome-gocui/gocui"
	"github.com/google/go-github/v62/github"
)

type githubRepoView struct {
	view          *tui.View
	tableRenderer *tui.TableRenderer[*github.Repository]
	search        string
}

var _ viewable = new(githubRepoView)

func newGithubRepoView() *githubRepoView {
	return &githubRepoView{}
}

func (g *githubRepoView) getView() *tui.View {
	return g.view
}

func (g *githubRepoView) focus() {
	ui.focusView(g.getView().Name())
}

func (g *githubRepoView) init() {
	g.view = getViewPosition(GithubRepoView).Set()

	g.view.Title = tui.WithSurroundingSpaces("Repositories")

	ui.styleView(g.view)

	g.tableRenderer = tui.NewTableRenderer[*github.Repository]()
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

	g.refresh()

	moveRight := func() {
		ui.getGithubPrView().Focus()
	}

	g.view.KeyBinding().
		Set('j', "Move down", func() {
			g.tableRenderer.Down()
		}).
		Set('k', "Move up", func() {
			g.tableRenderer.Up()
		}).
		Set('c', "Clone repo to a workspace", func() {
			repo := g.getSelectedRepo()
			if repo == nil {
				return
			}

			if api().GlobalConfiguration.Standalone {
				openToastDialog("Cannot clone to worksapce without a configured directory.", toastDialogNeutralType, "Note", func() {})
				return
			}

			sd := new(*searchListDialog[*core.Workspace])
			*sd = openSearchListDialog(searchDialogConfig[*core.Workspace]{
				tableViewTitle:  "workspaces",
				searchViewTitle: "Filter",
				tableTitles: []string{
					"Topic/Name",
				},
				tableProportions: []float64{
					1.0,
				},
				focusList: true,
				initial: func() ([][]string, []*core.Workspace) {
					workspaces := api().Workspaces.GetWorkspaces()
					rows := make([][]string, 0)
					for _, w := range workspaces {
						rows = append(rows, []string{
							w.ShortPath(),
						})
					}
					return rows, workspaces
				},
				onSearch: func(s string) ([][]string, []*core.Workspace) {
					workspaces := api().Workspaces.GetWorkspaces().FilterByNameContaining(s)
					rows := make([][]string, 0)
					for _, w := range workspaces {
						rows = append(rows, []string{
							w.ShortPath(),
						})
					}
					return rows, workspaces
				},
				onSelect: func(a *core.Workspace) {
					if *sd != nil {
						(*sd).close()
					}

					g.focus()
					err := a.CloneRepo(repo.GetHTMLURL())
					if err != nil {
						openToastDialogError(err.Error())
						return
					}

					go func() {
						wv := ui.getWorkspacesView()
						tv := ui.getTopicsView()
						refreshMainViews()
						ui.mainTabGroup.FocusTabByIndex(0)
						wv.focus()
						tv.selectTopicByName(a.Topic.Name)
						wv.selectWorkspaceByShortPath(a.ShortPath())
					}()
				},
			})
		}).
		Set('o', "Open repo in browser", func() {
			repo := g.getSelectedRepo()
			if repo == nil {
				return
			}

			system.OpenBrowser(repo.GetHTMLURL())
		}).
		Set('R', "Refetch all github data", func() {
			ui.getGithubProfileView().refetchData()
		}).
		Set('u', "Copy repo url to clipboard", func() {
			repo := g.getSelectedRepo()
			if repo == nil {
				return
			}

			url := repo.GetHTMLURL()
			system.CopyToClip(url)
			openToastDialog(url, toastDialogNeutralType, "Repo URL copied to clipboard", func() {})
		}).
		Set('/', "", func() {
			openEditorDialog(func(s string) {
				if s == "" {
					return
				}

				g.search = s
				g.view.Subtitle = fmt.Sprintf("Searching: %s", g.search)
				g.refresh()
			}, func() {}, "Search by name", smallEditorSize)
		}).
		Set(gocui.KeyArrowRight, "Focus PR View", moveRight).
		Set(gocui.KeyCtrlL, "Focus PR View", moveRight).
		Set('?', "Toggle cheatsheet", func() {
			openHelpDialog(g.view.GetKeybindings(), func() {})
		})
}

func (g *githubRepoView) getSelectedRepo() *github.Repository {
	_, value := g.tableRenderer.GetSelectedRow()
	if value != nil {
		return *value
	}

	return nil
}

func (g *githubRepoView) refresh() {
	if !api().Github.IsAuthenticated() {
		return
	}

	repos := api().Github.GetUserRepos()

	if g.search != "" {
		o := make([]*github.Repository, 0)
		for _, r := range repos {
			if strings.Contains(r.GetName(), g.search) {
				o = append(o, r)
			}
		}
		repos = o
	}

	sort.Slice(repos, func(i, j int) bool {
		return repos[i].GetUpdatedAt().After(repos[j].GetUpdatedAt().Time)
	})

	rows := make([][]string, 0)
	rowValues := make([]*github.Repository, 0)
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
	g.view.Clear()
	g.view.Resize(getViewPosition(g.view.Name()))
	if !api().Github.IsAuthenticated() {
		fmt.Fprintln(g.view, "Not authenticated")
		return nil
	}

	isFocused := g.view.IsFocused()
	g.tableRenderer.RenderWithSelectCallBack(g.view, func(_ int, _ *tui.TableRow[*github.Repository]) bool {
		return isFocused
	})

	if ui.getGithubProfileView().isFetchingData.Get() {
		fmt.Fprintln(g.view, "Loading...")
	} else if g.tableRenderer.GetTableSize() == 0 {
		fmt.Fprintln(g.view, "No repos to display")
	}

	return nil
}
