package ui

import (
	"fmt"
	"log"
	"mynav/pkg/github"
	"mynav/pkg/system"

	"github.com/awesome-gocui/gocui"
)

type GithubPrView struct {
	tableRenderer *TableRenderer
	prs           github.GithubPullRequests
}

const GithubPrViewName = "GithubPrView"

var _ View = &GithubPrView{}

func NewGithubPrView() *GithubPrView {
	return &GithubPrView{}
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

func (g *GithubPrView) Init(ui *UI) {
	view := SetViewLayout(g.Name())

	view.Title = "Pull Requests"
	view.TitleColor = gocui.ColorBlue
	view.FrameColor = gocui.ColorGreen

	sizeX, sizeY := view.Size()
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

	go func() {
		g.refreshPrs()
		UpdateGui(func(_ *gocui.Gui) error {
			g.Render(ui)
			return nil
		})
	}()

	moveUp := func() {
		ui.FocusWorkspacesView()
	}

	moveLeft := func() {
		ui.FocusTmuxView()
	}

	KeyBinding(g.Name()).
		set('j', func() {
			g.tableRenderer.Down()
		}).
		set('k', func() {
			g.tableRenderer.Up()
		}).
		set('L', func() {
			if Api().Github.IsAuthenticated() {
				return
			}

			deviceAuth := Api().Github.AuthenticateWithDeviceAuth(func() {
				g.refreshPrs()
				UpdateGui(func(_ *gocui.Gui) error {
					g.Render(ui)
					return nil
				})
			})

			if deviceAuth != nil {
				GetDialog[*ToastDialog](ui).Open(deviceAuth.UserCode, false, "User device code - automatically copied to clipboard", func() {})
				system.CopyToClip(deviceAuth.UserCode)
				deviceAuth.OpenBrowser()
			}
		}).
		set('P', func() {
			if Api().Github.IsAuthenticated() {
				return
			}

			GetDialog[*EditorDialog](ui).Open(func(s string) {
				if err := Api().Github.AuthenticateWithPersonalAccessToken(s); err != nil {
					GetDialog[*ToastDialog](ui).OpenError(err.Error())
					return
				}

				g.refreshPrs()
			}, func() {}, "Personal Access Token", Small)
		}).
		set('O', func() {
			Api().Github.LogoutUser()
		}).
		set(gocui.KeyEsc, moveUp).
		set(gocui.KeyArrowUp, moveUp).
		set(gocui.KeyCtrlK, moveUp).
		set(gocui.KeyArrowLeft, moveLeft).
		set(gocui.KeyCtrlH, moveLeft).
		set('?', func() {
			GetDialog[*HelpView](ui).Open(githubPrViewKeyBindings, func() {})
		})
}

func (g *GithubPrView) Name() string {
	return GithubPrViewName
}

func (g *GithubPrView) Render(ui *UI) error {
	view := GetInternalView(g.Name())
	if view == nil {
		g.Init(ui)
		view = GetInternalView(g.Name())
	}

	if !Api().Github.IsAuthenticated() {
		view.Clear()
		fmt.Fprintln(view, "Not authenticated")
		fmt.Fprintln(view, "Press:")
		fmt.Fprintln(view, "'L' - to login with device code using a browser")
		fmt.Fprintln(view, "'P' - to login in with Personal access token")
		return nil
	}

	view.Clear()
	view.Subtitle = "Login: " + Api().Github.GetPrincipalLogin()

	isFocused := false
	if v := GetFocusedView(); v != nil && v.Name() == g.Name() {
		isFocused = true
	}

	g.tableRenderer.RenderWithSelectCallBack(view, func(_ int, _ *TableRow) bool {
		return isFocused
	})

	if g.prs == nil {
		fmt.Fprintln(view, "Loading...")
	}

	return nil
}

func (g *GithubPrView) RequiresManager() bool {
	return false
}
