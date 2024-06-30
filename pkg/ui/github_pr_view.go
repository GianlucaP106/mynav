package ui

import (
	"fmt"
	"log"
	"mynav/pkg/github"
	"mynav/pkg/system"

	"github.com/awesome-gocui/gocui"
)

type GithubPrView struct {
	view          *View
	tableRenderer *TableRenderer
	prs           github.GithubPullRequests
}

var _ Viewable = new(GithubPrView)

const GithubPrViewName = "GithubPrView"

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
	g.view = GetViewPosition(GithubPrViewName).Set()

	g.view.Title = "Pull Requests"
	g.view.TitleColor = gocui.ColorBlue
	g.view.FrameColor = gocui.ColorGreen

	sizeX, sizeY := g.view.Size()
	g.tableRenderer = NewTableRenderer()
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

	// TODO: make this more formal
	go func() {
		g.refreshPrs()
		UpdateGui(func(_ *Gui) error {
			g.Render()
			return nil
		})
	}()

	moveLeft := func() {
		GetGithubRepoView().Focus()
	}

	KeyBinding(g.view.Name()).
		set('j', func() {
			g.tableRenderer.Down()
		}).
		set('k', func() {
			g.tableRenderer.Up()
		}).
		set('o', func() {
			pr := g.getSelectedPr()
			if pr == nil {
				return
			}

			system.OpenBrowser(pr.GetHTMLURL())
		}).
		set('L', func() {
			if Api().Github.IsAuthenticated() {
				return
			}

			td := new(*ToastDialog)
			deviceAuth := Api().Github.AuthenticateWithDevice(func() {
				if *td != nil {
					(*td).Close()
				}
				g.refreshPrs()
				UpdateGui(func(_ *Gui) error {
					g.Render()
					return nil
				})
			})

			if deviceAuth != nil {
				(*td) = OpenToastDialog(deviceAuth.UserCode, false, "User device code - automatically copied to clipboard", func() {})
				system.CopyToClip(deviceAuth.UserCode)
				deviceAuth.OpenBrowser()
			}
		}).
		set('P', func() {
			if Api().Github.IsAuthenticated() {
				return
			}

			OpenEditorDialog(func(s string) {
				if err := Api().Github.AuthenticateWithPersonalAccessToken(s); err != nil {
					OpenToastDialogError(err.Error())
					return
				}

				g.refreshPrs()
			}, func() {}, "Personal Access Token", Small)
		}).
		set('O', func() {
			Api().Github.LogoutUser()
		}).
		set('?', func() {
			OpenHelpView(githubPrViewKeyBindings, func() {})
		}).
		set(gocui.KeyEsc, moveLeft).
		set(gocui.KeyArrowLeft, moveLeft).
		set(gocui.KeyCtrlH, moveLeft)
}

func (g *GithubPrView) refreshPrs() {
	if !Api().Github.IsAuthenticated() {
		return
	}

	gpr, err := Api().Github.GetUserPullRequestsLocked()
	if err != nil {
		log.Panicln(err)
	}

	g.prs = gpr
	g.syncPrsToTable()
}

func (g *GithubPrView) getSelectedPr() *github.GithubPullRequest {
	idx := g.tableRenderer.GetSelectedRowIndex()
	if idx < 0 || idx >= len(g.prs) {
		return nil
	}
	return g.prs[idx]
}

func (g *GithubPrView) syncPrsToTable() {
	rows := make([][]string, 0)
	for _, pr := range g.prs {
		rows = append(rows, []string{
			pr.Repo.GetName(),
			pr.GetTitle(),
			pr.Relation,
		})
	}
	g.tableRenderer.FillTable(rows)
}

func (g *GithubPrView) Render() error {
	if !Api().Github.IsAuthenticated() {
		g.view.Clear()
		fmt.Fprintln(g.view, "Not authenticated")
		fmt.Fprintln(g.view, "Press:")
		fmt.Fprintln(g.view, "'L' - to login with device code using a browser")
		fmt.Fprintln(g.view, "'P' - to login in with Personal access token")
		return nil
	}

	g.view.Clear()
	g.view.Subtitle = "Login: " + Api().Github.GetPrincipalLogin()

	isFocused := IsViewFocused(g.view)

	g.tableRenderer.RenderWithSelectCallBack(g.view, func(_ int, _ *TableRow) bool {
		return isFocused
	})

	if g.prs == nil {
		fmt.Fprintln(g.view, "Loading...")
	} else if len(g.prs) == 0 {
		fmt.Fprintln(g.view, "No PRs to display")
	}

	return nil
}
