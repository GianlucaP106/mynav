package ui

import (
	"fmt"
	"mynav/pkg/events"
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

func (g *githubProfileView) init() {
	g.view = getViewPosition(GithubProfileView).Set()

	g.view.Title = tui.WithSurroundingSpaces("Profile")

	styleView(g.view)

	g.view.KeyBinding().
		Set('L', "Login with device code and browser", func() {
			if getApi().Github.IsAuthenticated() {
				return
			}

			tdMu := &sync.Mutex{}
			td := new(*toastDialog)
			events.AddEventListener(events.GithubDeviceAuthenticatedEvent, func(listenerId string) {
				tdMu.Lock()
				defer tdMu.Unlock()

				if td != nil && *td != nil {
					tui.UpdateTui(func(g *tui.Tui) error {
						(*td).close()
						return nil
					})
				}

				getApi().Github.LoadData()
				events.RemoveEventListener(events.GithubDeviceAuthenticatedEvent, listenerId)
			})

			deviceAuth := getApi().Github.InitWithDeviceAuth()

			if deviceAuth != nil {
				tdMu.Lock()
				(*td) = openToastDialog(deviceAuth.UserCode, false, "User device code - automatically copied to clipboard", func() {})
				tdMu.Unlock()
				system.CopyToClip(deviceAuth.UserCode)
				deviceAuth.OpenBrowser()
			}
		}).
		Set('o', "Open in browser", func() {
			profile := getApi().Github.GetProfile()
			if profile.IsLoaded() {
				profile.OpenBrowser()
			}
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
		}).
		Set('?', "Toggle cheatsheet", func() {
			OpenHelpDialog(g.view.GetKeybindings(), func() {})
		})
}

func (g *githubProfileView) render() error {
	g.view.Clear()
	g.view = getViewPosition(g.view.Name()).Set()
	if !getApi().Github.IsAuthenticated() {
		fmt.Fprintln(g.view, "Not authenticated")
		fmt.Fprintln(g.view, "Press:")
		fmt.Fprintln(g.view, "'L' - to login with device code using a browser")
		fmt.Fprintln(g.view, "'P' - to login in with Personal access token")
		return nil
	}

	profile := getApi().Github.GetProfile()
	fmt.Fprintln(g.view, "Login: ", profile.Login)
	fmt.Fprintln(g.view, "Email: ", profile.Email)
	fmt.Fprintln(g.view, "Name: ", profile.Name)
	fmt.Fprintln(g.view, "Url: ", profile.Url)

	return nil
}

func (g *githubProfileView) getView() *tui.View {
	return g.view
}

func (g *githubProfileView) Focus() {
	focusView(g.getView().Name())
}
