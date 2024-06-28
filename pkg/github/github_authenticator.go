package github

import (
	"context"
	"errors"
	"log"
	"mynav/pkg/system"
	"net/http"
	"time"

	gh "github.com/google/go-github/v62/github"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

type GithubAuthenticator struct {
	oauthConfig *oauth2.Config
	onLogin     func(*GithubAuthenticationToken)
	onLogout    func()
}

type GithubDevicePreAuthentication struct {
	*oauth2.DeviceAuthResponse
}

type GithubAuthenticationToken struct {
	PersonalAccessToken *string
	DeviceToken         *oauth2.Token
}

func NewGithubAuthenticator(clientId string, onLogin func(*GithubAuthenticationToken), onLogout func(), scopes ...string) *GithubAuthenticator {
	c := &oauth2.Config{
		ClientID: clientId,
		Scopes:   scopes,
		Endpoint: githuboauth.Endpoint,
	}

	ga := &GithubAuthenticator{
		oauthConfig: c,
		onLogin:     onLogin,
		onLogout:    onLogout,
	}

	return ga
}

func (ga *GithubAuthenticator) AuthenticateWithDevice() (*GithubDevicePreAuthentication, func() *gh.Client) {
	deviceAuth, err := ga.oauthConfig.DeviceAuth(context.TODO())
	if err != nil {
		// TODO:
		log.Fatal(err)
	}

	gda := &GithubDevicePreAuthentication{
		DeviceAuthResponse: deviceAuth,
	}

	return gda, func() *gh.Client {
		auth := ga.PollAuthenticateDevice(gda)
		client := ga.InitClient(auth)
		ga.onLogin(auth)
		return client
	}
}

func (ga *GithubAuthenticator) PollAuthenticateDevice(da *GithubDevicePreAuthentication) *GithubAuthenticationToken {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	token, err := ga.oauthConfig.DeviceAccessToken(ctx, da.DeviceAuthResponse)
	if err != nil {
		// TODO:
		log.Fatal(err)
	}

	return &GithubAuthenticationToken{
		DeviceToken: token,
	}
}

func (ga *GithubAuthenticator) AuthenticateWithPersonalAccessToken(token string) (*gh.Client, error) {
	gt := &GithubAuthenticationToken{
		PersonalAccessToken: &token,
	}

	client := ga.InitClient(gt)
	if client == nil {
		return nil, nil
	}

	if _, _, err := client.Users.Get(context.Background(), ""); err != nil {
		return nil, errors.New("invalid token")
	}

	ga.onLogin(gt)

	return client, nil
}

func (ga *GithubAuthenticator) HttpClient(auth *GithubAuthenticationToken) *http.Client {
	return ga.oauthConfig.Client(context.Background(), auth.DeviceToken)
}

func (gda *GithubDevicePreAuthentication) OpenBrowser() {
	system.OpenBrowser(gda.VerificationURI)
}

func (gda *GithubAuthenticator) InitClient(auth *GithubAuthenticationToken) *gh.Client {
	var client *gh.Client
	if auth.PersonalAccessToken != nil {
		client = gh.NewClient(nil).WithAuthToken(*auth.PersonalAccessToken)
	} else {
		http := gda.HttpClient(auth)
		client = gh.NewClient(http)
	}

	return client
}
