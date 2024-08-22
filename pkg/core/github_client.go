package core

import (
	"context"

	"github.com/google/go-github/v62/github"
)

type GithubClient struct {
	client        *github.Client
	authenticator *GithubAuthenticator
}

func NewGithubClient(token *GithubAuthenticationToken) GithubClient {
	ga := NewGithubAuthenticator("repo", "read:org")
	client := GithubClient{
		client:        nil,
		authenticator: ga,
	}

	if token != nil {
		client.client = ga.initClient(token)
	}

	return client
}

func (g *GithubClient) AuthenticateWithDevice() (*GithubDevicePreAuthentication, func() *GithubAuthenticationToken) {
	gda, f := g.authenticator.AuthenticateWithDevice()

	poll := func() *GithubAuthenticationToken {
		client, token := f()
		g.client = client
		return token
	}

	return gda, poll
}

func (g *GithubClient) AuthenticateWithPersonalAccessToken(token string) (*GithubAuthenticationToken, error) {
	client, gtoken, err := g.authenticator.AuthenticateWithPersonalAccessToken(token)
	if err != nil {
		return nil, err
	}

	g.client = client
	return gtoken, nil
}

func (gs *GithubClient) IsAuthenticated() bool {
	return gs.client != nil
}

func (g *GithubClient) PullRequests(owner string, repoName string, fetchLimit int) ([]*github.PullRequest, error) {
	prs, _, err := g.client.PullRequests.List(context.Background(), owner, repoName, &github.PullRequestListOptions{
		ListOptions: github.ListOptions{
			PerPage: fetchLimit,
		},
	})

	return prs, err
}

func (g *GithubClient) PrincipalRepositories(fetchLimit int) ([]*github.Repository, error) {
	userRepos, _, err := g.client.Repositories.ListByAuthenticatedUser(context.Background(), &github.RepositoryListByAuthenticatedUserOptions{
		ListOptions: github.ListOptions{
			PerPage: fetchLimit,
		},
	})

	return userRepos, err
}

func (g *GithubClient) PrincipalOrganizations(fetchLimit int) ([]*github.Organization, error) {
	orgs, _, err := g.client.Organizations.List(context.Background(), "", &github.ListOptions{
		PerPage: fetchLimit,
	})

	return orgs, err
}

func (g *GithubClient) RepositoriesByOrg(org string, fetchLimit int) ([]*github.Repository, error) {
	orgRepos, _, err := g.client.Repositories.ListByOrg(context.Background(), org, &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{
			PerPage: fetchLimit,
		},
	})

	return orgRepos, err
}

func (g *GithubClient) Principal() (*github.User, error) {
	principal, _, err := g.client.Users.Get(context.Background(), "")
	return principal, err
}
