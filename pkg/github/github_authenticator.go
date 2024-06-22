package github

import (
	"context"
	"log"
	"mynav/pkg/system"
	"net/http"
	"time"

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

func (ga *GithubAuthenticator) InitAuth() *GithubDevicePreAuthentication {
	deviceAuth, err := ga.oauthConfig.DeviceAuth(context.TODO())
	if err != nil {
		// TODO:
		log.Fatal(err)
	}

	return &GithubDevicePreAuthentication{
		DeviceAuthResponse: deviceAuth,
	}
}

func (ga *GithubAuthenticator) Authenticate(da *GithubDevicePreAuthentication) *GithubAuthenticationToken {
	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
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

func (ga *GithubAuthenticator) HttpClient(auth *GithubAuthenticationToken) *http.Client {
	return ga.oauthConfig.Client(context.Background(), auth.DeviceToken)
}

func (gda *GithubDevicePreAuthentication) OpenBrowser() {
	system.OpenBrowser(gda.VerificationURI)
}
