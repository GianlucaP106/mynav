package ui

import (
	"fmt"
	"mynav/pkg/persistence"
	"mynav/pkg/system"
	"mynav/pkg/tui"
	"sync"
)

type githubProfileView struct {
	view           *tui.View
	isFetchingData *persistence.Value[bool]
}

var _ viewable = new(githubProfileView)

func newGithubProfileView() *githubProfileView {
	return &githubProfileView{
		isFetchingData: persistence.NewValue(false),
	}
}

func getGithubProfileView() *githubProfileView {
	return getViewable[*githubProfileView]()
}

func (g *githubProfileView) refresh() {}

func (g *githubProfileView) init() {
	g.view = getViewPosition(GithubProfileView).Set()

	g.view.Title = tui.WithSurroundingSpaces("Profile")

	styleView(g.view)

	isAuthenticated := getApi().Github.IsAuthenticated()
	if isAuthenticated {
		g.fetchData()
	}

	g.view.KeyBinding().
		Set('L', "Login with device code and browser", func() {
			if getApi().Github.IsAuthenticated() {
				return
			}

			tdMu := &sync.Mutex{}
			td := new(*toastDialog)

			deviceAuth, poll := getApi().Github.InitWithDeviceAuth()
			go func() {
				poll()

				tdMu.Lock()
				defer tdMu.Unlock()

				if td != nil && *td != nil {
					tui.UpdateTui(func(t *tui.Tui) error {
						(*td).close()
						g.view.Focus()
						return nil
					})
				}

				g.fetchData()
			}()

			if deviceAuth != nil {
				tdMu.Lock()
				(*td) = openToastDialog(fmt.Sprintf("Code: %s - %s", deviceAuth.UserCode, deviceAuth.VerificationURI), toastDialogNeutralType, "User device code - automatically copied to clipboard", func() {})
				tdMu.Unlock()
				system.CopyToClip(deviceAuth.UserCode)
				deviceAuth.OpenBrowser()
			}
		}).
		Set('o', "Open in browser", func() {
			profile := getApi().Github.GetPrincipal()
			if profile != nil {
				system.OpenBrowser(profile.GetHTMLURL())
			}
		}).
		Set('u', "Copy profile url to cliboard", func() {
			user := getApi().Github.GetPrincipal()
			if user == nil {
				return
			}

			url := user.GetHTMLURL()
			system.CopyToClip(url)
			openToastDialog(url, toastDialogNeutralType, "Profile URL copied", func() {})
		}).
		Set('P', "Login with personal access token", func() {
			if getApi().Github.IsAuthenticated() {
				return
			}

			openEditorDialog(func(s string) {
				if err := getApi().Github.InitWithPAT(s); err != nil {
					openToastDialogError(err.Error())
					return
				}
			}, func() {}, "Personal Access Token", smallEditorSize)
		}).
		Set('O', "Logout", func() {
			getApi().Github.LogoutUser()
			openToastDialog("Successfully logged out - restart mynav to clear the github views", toastDialogSuccessType, "Note", func() {})
		}).
		Set('R', "Refetch all github data", func() {
			g.refetchData()
		}).
		Set('?', "Toggle cheatsheet", func() {
			openHelpDialog(g.view.GetKeybindings(), func() {})
		})
}

func (g *githubProfileView) fetchData() {
	go func() {
		g.isFetchingData.Set(true)
		getApi().Github.LoadProfile()
		refreshAsync(g)

		getApi().Github.LoadUserRepos()
		refreshAsync(getGithubRepoView())

		getApi().Github.LoadUserPullRequests()
		refreshAsync(getGithubPrView())
		g.isFetchingData.Set(false)
	}()
}

func (g *githubProfileView) refetchData() {
	repoView := getGithubRepoView()
	prView := getGithubPrView()
	repoView.tableRenderer.ClearTable()
	prView.tableRenderer.ClearTable()
	g.fetchData()
}

func (g *githubProfileView) render() error {
	g.view.Clear()
	g.view.Resize(getViewPosition(g.view.Name()))
	if !getApi().Github.IsAuthenticated() {
		fmt.Fprintln(g.view, "Not authenticated")
		fmt.Fprintln(g.view, "Press:")
		fmt.Fprintln(g.view, "'L' - to login with device code using a browser")
		fmt.Fprintln(g.view, "'P' - to login in with Personal access token")
		return nil
	}

	profile := getApi().Github.GetPrincipal()
	fmt.Fprintln(g.view, "Login: ", profile.GetLogin())
	fmt.Fprintln(g.view, "Email: ", profile.GetEmail())
	fmt.Fprintln(g.view, "Name: ", profile.GetName())
	fmt.Fprintln(g.view, "Url: ", profile.GetHTMLURL())

	return nil
}

func (g *githubProfileView) getView() *tui.View {
	return g.view
}

func (g *githubProfileView) Focus() {
	focusView(g.getView().Name())
}
