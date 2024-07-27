package ui

import (
	"fmt"
	"mynav/pkg/constants"
	"mynav/pkg/events"
	"mynav/pkg/system"
	"sync"

	"github.com/awesome-gocui/gocui"
)

type GithubProfileView struct {
	view *View
}

var _ Viewable = new(GithubProfileView)

func NewGithubProfileView() *GithubProfileView {
	return &GithubProfileView{}
}

func GetGithubProfileView() *GithubProfileView {
	return GetViewable[*GithubProfileView]()
}

func (g *GithubProfileView) Init() {
	g.view = GetViewPosition(constants.GithubProfileViewName).Set()

	g.view.Title = "Profile"
	g.view.TitleColor = gocui.ColorBlue
	g.view.FrameColor = gocui.ColorGreen

	g.view.KeyBinding().
		set('L', func() {
			if Api().Github.IsAuthenticated() {
				return
			}

			tdMu := &sync.Mutex{}
			td := new(*ToastDialog)
			events.AddEventListener(constants.GithubDeviceAuthenticatedEventName, func(listenerId string) {
				tdMu.Lock()
				defer tdMu.Unlock()

				if td != nil && *td != nil {
					UpdateGui(func(g *Gui) error {
						(*td).Close()
						return nil
					})
				}

				Api().Github.LoadData()
				events.RemoveEventListener(constants.GithubDeviceAuthenticatedEventName, listenerId)
			})

			deviceAuth := Api().Github.InitWithDeviceAuth()

			if deviceAuth != nil {
				tdMu.Lock()
				(*td) = OpenToastDialog(deviceAuth.UserCode, false, "User device code - automatically copied to clipboard", func() {})
				tdMu.Unlock()
				system.CopyToClip(deviceAuth.UserCode)
				deviceAuth.OpenBrowser()
			}
		}, "Login with device code and browser").
		set('o', func() {
			profile := Api().Github.GetProfile()
			if profile.IsLoaded() {
				profile.OpenBrowser()
			}
		}, "Open in browser").
		set('P', func() {
			if Api().Github.IsAuthenticated() {
				return
			}

			OpenEditorDialog(func(s string) {
				if err := Api().Github.InitWithPAT(s); err != nil {
					OpenToastDialogError(err.Error())
					return
				}
			}, func() {}, "Personal Access Token", Small)
		}, "Login with personal access token").
		set('O', func() {
			Api().Github.LogoutUser()
		}, "Logout").
		set('?', func() {
			OpenHelpView(g.view.keybindingInfo.toList(), func() {})
		}, "Toggle cheatsheet")
}

func (g *GithubProfileView) Render() error {
	g.view.Clear()
	if !Api().Github.IsAuthenticated() {
		fmt.Fprintln(g.view, "Not authenticated")
		fmt.Fprintln(g.view, "Press:")
		fmt.Fprintln(g.view, "'L' - to login with device code using a browser")
		fmt.Fprintln(g.view, "'P' - to login in with Personal access token")
		return nil
	}

	profile := Api().Github.GetProfile()
	fmt.Fprintln(g.view, "Login: ", profile.Login)
	fmt.Fprintln(g.view, "Email: ", profile.Email)
	fmt.Fprintln(g.view, "Name: ", profile.Name)
	fmt.Fprintln(g.view, "Url: ", profile.Url)

	return nil
}

func (g *GithubProfileView) View() *View {
	return g.view
}

func (g *GithubProfileView) Focus() {
	FocusView(g.View().Name())
}
