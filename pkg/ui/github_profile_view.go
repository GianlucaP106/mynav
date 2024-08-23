package ui

import (
	"fmt"
	"mynav/pkg/system"
	"mynav/pkg/tui"
	"sync"
)

type githubProfileView struct {
	view *tui.View
}

var _ viewable = new(githubProfileView)

func newGithubProfileView() *githubProfileView {
	return &githubProfileView{}
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
		g.loadData()
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
					tui.UpdateTui(func(g *tui.Tui) error {
						(*td).close()
						return nil
					})
				}

				g.loadData()
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
			system.OpenBrowser(profile.GetHTMLURL())
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
		Set('?', "Toggle cheatsheet", func() {
			openHelpDialog(g.view.GetKeybindings(), func() {})
		})
}

func (g *githubProfileView) loadData() {
	go func() {
		getApi().Github.LoadProfile()
		refreshAsync(g)

		getApi().Github.LoadUserRepos()
		refreshAsync(getGithubRepoView())

		getApi().Github.LoadUserPullRequests()
		refreshAsync(getGithubPrView())
	}()
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
	fmt.Fprintln(g.view, "Url: ", profile.GetURL())

	return nil
}

func (g *githubProfileView) getView() *tui.View {
	return g.view
}

func (g *githubProfileView) Focus() {
	focusView(g.getView().Name())
}
