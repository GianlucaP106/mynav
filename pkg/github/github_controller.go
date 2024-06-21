package github

import (
	"context"
	"errors"

	gh "github.com/google/go-github/v62/github"
)

type GithubController struct {
	client        GithubClient
	authenticator *GithubAuthenticator
}

type GithubClient = *gh.Client

const CLIENT_ID = "Ov23lirJDAVBmN4oRLY0"

func NewGithubController() *GithubController {
	ga := NewGithubAuthenticator(CLIENT_ID, "repo")
	gs := &GithubController{
		authenticator: ga,
	}
	return gs
}

func (gs *GithubController) InitGithubClient(auth *GithubAuthenticationToken) {
	http := gs.authenticator.HttpClient(auth)
	client := gh.NewClient(http)
	gs.client = client
}

func (gs *GithubController) AuthenticateWithDeviceAuth() *GithubDevicePreAuthentication {
	gda := gs.authenticator.InitAuth()

	go func() {
		auth := gs.authenticator.Authenticate(gda)
		gs.InitGithubClient(auth)
	}()

	return gda
}

func (gs *GithubController) IsAuthenticatedToGithub() bool {
	return gs.client != nil
}

func (gs *GithubController) GithubPrincipal() *gh.User {
	principal, _, err := gs.client.Users.Get(context.TODO(), "")
	if err != nil {
		return nil
	}

	return principal
}

func (gs *GithubController) AuthenticateWithPersonalAccessToken(token string) error {
	client := gh.NewClient(nil).WithAuthToken(token)
	gs.client = client

	if gs.GithubPrincipal() == nil {
		gs.client = nil
		return errors.New("invalid token")
	}

	return nil
}
