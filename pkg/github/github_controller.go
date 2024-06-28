package github

import (
	"context"
	"sync"

	gh "github.com/google/go-github/v62/github"
)

type GithubController struct {
	client        GithubClient
	authenticator *GithubAuthenticator
	login         string
	PullRequests  GithubPullRequests
	Repos         []*gh.Repository
	clientMutex   sync.Mutex
}

type GithubClient = *gh.Client

const CLIENT_ID = "Ov23lirJDAVBmN4oRLY0"

const FETCH_LIMIT = 1000

func NewGithubController(token *GithubAuthenticationToken, onLogin func(*GithubAuthenticationToken), onLogout func()) *GithubController {
	ga := NewGithubAuthenticator(CLIENT_ID, onLogin, onLogout, "repo", "read:org")
	gs := &GithubController{
		authenticator: ga,
	}

	if token != nil {
		gs.client = gs.authenticator.InitClient(token)
	}

	return gs
}

func (gs *GithubController) AuthenticateWithDevice(callback func()) *GithubDevicePreAuthentication {
	gda, f := gs.authenticator.AuthenticateWithDevice()

	go func() {
		gs.client = f()
		callback()
	}()

	return gda
}

func (gs *GithubController) AuthenticateWithPersonalAccessToken(token string) error {
	client, err := gs.authenticator.AuthenticateWithPersonalAccessToken(token)
	if err != nil {
		return err
	}
	gs.client = client
	return nil
}

func (gs *GithubController) LogoutUser() {
	gs.clientMutex.Lock()
	defer gs.clientMutex.Unlock()
	gs.client = nil
	gs.login = ""
	gs.authenticator.onLogout()
}

func (gs *GithubController) IsAuthenticated() bool {
	return gs.client != nil
}

func (gs *GithubController) Principal() *gh.User {
	principal, _, err := gs.client.Users.Get(context.TODO(), "")
	if err != nil {
		return nil
	}

	return principal
}

func (gs *GithubController) GetPrincipalLogin() string {
	if gs.login == "" {
		p := gs.Principal()
		gs.login = *p.Login
	}

	return gs.login
}

func (gc *GithubController) GetUserPullRequests() (GithubPullRequests, error) {
	if gc.PullRequests != nil {
		return gc.PullRequests, nil
	}

	// TODO: mutex
	gc.clientMutex.Lock()
	defer gc.clientMutex.Unlock()

	allRepos, err := gc.GetUserRepos()
	if err != nil {
		return nil, err
	}

	allPrs := NewGithubPrContainer()
	for _, repo := range allRepos {
		prs, _, err := gc.client.PullRequests.List(context.Background(), *repo.GetOwner().Login, *repo.Name, &gh.PullRequestListOptions{
			ListOptions: gh.ListOptions{
				PerPage: FETCH_LIMIT,
			},
		})
		if err != nil {
			return nil, err
		}

		allPrs.AddPrs(repo, prs...)
	}

	login := gc.GetPrincipalLogin()
	out := NewGithubPrContainer()
	for _, pr := range allPrs {
		belongsToUser, relation := func() (bool, string) {
			if pr.GetAssignee().GetLogin() == login {
				return true, "Assignee"
			}

			if pr.GetUser().GetLogin() == login {
				return true, "User"
			}

			if func() bool {
				for _, reviewer := range pr.RequestedReviewers {
					if reviewer.GetLogin() == login {
						return true
					}
				}
				return false
			}() {
				return true, "Requested Reviewer"
			}

			if func() bool {
				for _, u := range pr.Assignees {
					if u.GetLogin() == login {
						return true
					}
				}
				return false
			}() {
				return true, "Assignee"
			}

			return false, ""
		}()

		if belongsToUser {
			pr.Relation = relation
			out.AddFromPr(pr)
		}
	}

	gc.PullRequests = out

	return out, nil
}

func (gc *GithubController) GetUserRepos() ([]*gh.Repository, error) {
	if gc.Repos != nil {
		return gc.Repos, nil
	}

	allRepos := make(map[string]*gh.Repository)
	userRepos, _, err := gc.client.Repositories.ListByAuthenticatedUser(context.Background(), &gh.RepositoryListByAuthenticatedUserOptions{
		ListOptions: gh.ListOptions{
			PerPage: FETCH_LIMIT,
		},
	})
	if err != nil {
		return nil, err
	}

	for _, repo := range userRepos {
		allRepos[repo.GetFullName()] = repo
	}

	orgs, _, err := gc.client.Organizations.List(context.Background(), "", &gh.ListOptions{
		PerPage: FETCH_LIMIT,
	})
	if err != nil {
		return nil, err
	}

	for _, org := range orgs {
		orgRepos, _, err := gc.client.Repositories.ListByOrg(context.Background(), *org.Login, &gh.RepositoryListByOrgOptions{
			ListOptions: gh.ListOptions{
				PerPage: FETCH_LIMIT,
			},
		})
		if err != nil {
			return nil, err
		}

		for _, repo := range orgRepos {
			allRepos[repo.GetFullName()] = repo
		}
	}

	out := make([]*gh.Repository, 0)
	for _, r := range allRepos {
		out = append(out, r)
	}

	gc.Repos = out

	return out, nil
}
