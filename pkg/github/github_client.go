package github

import (
	"context"
	"mynav/pkg/constants"
	"mynav/pkg/events"
	"mynav/pkg/tasks"
	"sync"

	gh "github.com/google/go-github/v62/github"
)

type GithubClient struct {
	client        *gh.Client
	authenticator *GithubAuthenticator
	mu            *sync.RWMutex
}

func NewGithubClient(token *GithubAuthenticationToken, onLogin func(*GithubAuthenticationToken), onLogout func()) GithubClient {
	ga := NewGithubAuthenticator(onLogin, onLogout, "repo", "read:org")
	client := GithubClient{
		client:        nil,
		mu:            &sync.RWMutex{},
		authenticator: ga,
	}

	if token != nil {
		client.client = ga.initClient(token)
	}

	return client
}

func (g *GithubClient) AuthenticateWithDevice() *GithubDevicePreAuthentication {
	gda, f := g.authenticator.AuthenticateWithDevice()

	tasks.AddTask(func() {
		client := f()
		g.mu.Lock()
		defer g.mu.Unlock()
		g.client = client
		events.Emit(constants.GithubDeviceAuthenticatedEventName)
	})
	// go func() {
	// }()

	return gda
}

func (g *GithubClient) AuthenticateWithPersonalAccessToken(token string) error {
	client, err := g.authenticator.AuthenticateWithPersonalAccessToken(token)
	if err != nil {
		return err
	}

	g.mu.Lock()
	defer g.mu.Unlock()
	g.client = client
	return nil
}

func (gs *GithubClient) Logout() {
	gs.mu.Lock()
	defer gs.mu.Unlock()
	gs.authenticator.onLogout()
	gs.client = nil
}

func (gs *GithubClient) IsAuthenticated() bool {
	gs.mu.RLock()
	defer gs.mu.RUnlock()
	return gs.client != nil
}

func (g *GithubClient) PullRequests(owner string, repoName string, fetchLimit int) ([]*gh.PullRequest, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	prs, _, err := g.client.PullRequests.List(context.Background(), owner, repoName, &gh.PullRequestListOptions{
		ListOptions: gh.ListOptions{
			PerPage: fetchLimit,
		},
	})

	return prs, err
}

func (g *GithubClient) PrincipalRepositories(fetchLimit int) ([]*gh.Repository, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	userRepos, _, err := g.client.Repositories.ListByAuthenticatedUser(context.Background(), &gh.RepositoryListByAuthenticatedUserOptions{
		ListOptions: gh.ListOptions{
			PerPage: fetchLimit,
		},
	})

	return userRepos, err
}

func (g *GithubClient) PrincipalOrganizations(fetchLimit int) ([]*gh.Organization, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	orgs, _, err := g.client.Organizations.List(context.Background(), "", &gh.ListOptions{
		PerPage: fetchLimit,
	})

	return orgs, err
}

func (g *GithubClient) RepositoriesByOrg(org string, fetchLimit int) ([]*gh.Repository, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	orgRepos, _, err := g.client.Repositories.ListByOrg(context.Background(), org, &gh.RepositoryListByOrgOptions{
		ListOptions: gh.ListOptions{
			PerPage: fetchLimit,
		},
	})

	return orgRepos, err
}

func (g *GithubClient) Principal() (*gh.User, error) {
	principal, _, err := g.client.Users.Get(context.TODO(), "")
	return principal, err
}
