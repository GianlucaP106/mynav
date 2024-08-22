package core

import (
	"context"
	"errors"
	"log"
	"mynav/pkg/system"
	"net/http"
	"time"

	"github.com/google/go-github/v62/github"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

type GithubAuthenticator struct {
	oauthConfig *oauth2.Config
}

type GithubDevicePreAuthentication struct {
	*oauth2.DeviceAuthResponse
}

type GithubAuthenticationToken struct {
	PersonalAccessToken *string
	DeviceToken         *oauth2.Token
}

const CLIENT_ID = "Ov23lirJDAVBmN4oRLY0"

func NewGithubAuthenticator(scopes ...string) *GithubAuthenticator {
	c := &oauth2.Config{
		ClientID: CLIENT_ID,
		Scopes:   scopes,
		Endpoint: githuboauth.Endpoint,
	}

	ga := &GithubAuthenticator{
		oauthConfig: c,
	}

	return ga
}

func (ga *GithubAuthenticator) AuthenticateWithDevice() (*GithubDevicePreAuthentication, func() (*github.Client, *GithubAuthenticationToken)) {
	deviceAuth, err := ga.oauthConfig.DeviceAuth(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	gda := &GithubDevicePreAuthentication{
		DeviceAuthResponse: deviceAuth,
	}

	return gda, func() (*github.Client, *GithubAuthenticationToken) {
		auth := ga.pollAuthenticateDevice(gda)
		client := ga.initClient(auth)
		return client, auth
	}
}

func (ga *GithubAuthenticator) pollAuthenticateDevice(da *GithubDevicePreAuthentication) *GithubAuthenticationToken {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	token, err := ga.oauthConfig.DeviceAccessToken(ctx, da.DeviceAuthResponse)
	if err != nil {
		log.Fatal(err)
	}

	return &GithubAuthenticationToken{
		DeviceToken: token,
	}
}

func (ga *GithubAuthenticator) AuthenticateWithPersonalAccessToken(token string) (*github.Client, *GithubAuthenticationToken, error) {
	gt := &GithubAuthenticationToken{
		PersonalAccessToken: &token,
	}

	client := ga.initClient(gt)
	if client == nil {
		return nil, nil, nil
	}

	if _, _, err := client.Users.Get(context.Background(), ""); err != nil {
		return nil, nil, errors.New("invalid token")
	}

	return client, gt, nil
}

func (ga *GithubAuthenticator) httpClient(auth *GithubAuthenticationToken) *http.Client {
	return ga.oauthConfig.Client(context.Background(), auth.DeviceToken)
}

func (gda *GithubDevicePreAuthentication) OpenBrowser() {
	system.OpenBrowser(gda.VerificationURI)
}

func (gda *GithubAuthenticator) initClient(auth *GithubAuthenticationToken) *github.Client {
	var client *github.Client
	if auth.PersonalAccessToken != nil {
		client = github.NewClient(nil).WithAuthToken(*auth.PersonalAccessToken)
	} else {
		http := gda.httpClient(auth)
		client = github.NewClient(http)
	}

	return client
}
