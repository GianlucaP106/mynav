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
}

type GithubDevicePreAuthentication struct {
	*oauth2.DeviceAuthResponse
}

type GithubAuthenticationToken struct {
	PersonalAccessToken *string
	DeviceToken         *oauth2.Token
}

func NewGithubAuthenticator(clientId string, scopes ...string) *GithubAuthenticator {
	c := &oauth2.Config{
		ClientID: clientId,
		Scopes:   scopes,
		Endpoint: githuboauth.Endpoint,
	}

	ga := &GithubAuthenticator{
		oauthConfig: c,
	}

	return ga
}

func (ga *GithubAuthenticator) InitDeviceAuth() *GithubDevicePreAuthentication {
	deviceAuth, err := ga.oauthConfig.DeviceAuth(context.TODO())
	if err != nil {
		// TODO:
		log.Fatal(err)
	}

	return &GithubDevicePreAuthentication{
		DeviceAuthResponse: deviceAuth,
	}
}

func (ga *GithubAuthenticator) AuthenticateDevice(da *GithubDevicePreAuthentication) *GithubAuthenticationToken {
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

func (ga *GithubAuthenticator) AuthenticateWithPersonalAccessToken(token string) (*gh.Client, *GithubAuthenticationToken, error) {
	outErr := errors.New("invalid token")
	gt := &GithubAuthenticationToken{
		PersonalAccessToken: &token,
	}

	client := ga.InitClient(gt)
	if client == nil {
		return nil, nil, outErr
	}

	if _, _, err := client.Users.Get(context.Background(), ""); err != nil {
		return nil, nil, errors.New("invalid token")
	}

	return client, gt, nil
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
