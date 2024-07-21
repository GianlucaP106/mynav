package github

import "mynav/pkg/system"

type GithubProfile struct {
	Login string
	Name  string
	Email string
	Url   string
}

func (g GithubProfile) OpenBrowser() {
	system.OpenBrowser(g.Url)
}

func (g GithubProfile) IsLoaded() bool {
	return g.Login != ""
}
