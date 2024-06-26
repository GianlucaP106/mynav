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

const GithubPrViewName = "GithubPrView"

var _ Viewable = new(GithubPrView)

func NewGithubPrView() *GithubPrView {
	return &GithubPrView{}
}

func GetGithubPrView() *GithubPrView {
	return GetViewable[*GithubPrView]()
}

func FocusGithubPrView() {
	FocusView(GithubPrViewName)
}

func (g *GithubPrView) View() *View {
	return g.view
}

func (g *GithubPrView) Init() {
	g.view = SetViewLayout(g.Name())

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

	// moveUp := func() {
	// 	FocusWorkspacesView()
	// }
	//
	// moveLeft := func() {
	// 	FocusTmuxView()
	// }

	KeyBinding(g.Name()).
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

			deviceAuth := Api().Github.AuthenticateWithDeviceAuth(func() {
				g.refreshPrs()
				UpdateGui(func(_ *Gui) error {
					g.Render()
					return nil
				})
			})

			if deviceAuth != nil {
				OpenToastDialog(deviceAuth.UserCode, false, "User device code - automatically copied to clipboard", func() {})
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
		// set(gocui.KeyEsc, moveUp).
		// set(gocui.KeyArrowUp, moveUp).
		// set(gocui.KeyCtrlK, moveUp).
		// set(gocui.KeyArrowLeft, moveLeft).
		// set(gocui.KeyCtrlH, moveLeft).
		set('?', func() {
			OpenHelpView(githubPrViewKeyBindings, func() {})
		})
}

func (g *GithubPrView) refreshPrs() {
	if !Api().Github.IsAuthenticated() {
		return
	}

	gpr, err := Api().Github.GetUserPullRequests()
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

func (g *GithubPrView) Name() string {
	return GithubPrViewName
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

	isFocused := false
	if v := GetFocusedView(); v != nil && v.Name() == g.Name() {
		isFocused = true
	}

	g.tableRenderer.RenderWithSelectCallBack(g.view, func(_ int, _ *TableRow) bool {
		return isFocused
	})

	if g.prs == nil {
		fmt.Fprintln(g.view, "Loading...")
	}

	return nil
}
