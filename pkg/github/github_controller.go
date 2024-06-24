package github

import (
	"context"
	"sync"

	gh "github.com/google/go-github/v62/github"
)

type GithubController struct {
	client        GithubClient
	authenticator *GithubAuthenticator
	onLogin       func(*GithubAuthenticationToken)
	onLogout      func()
	login         string
	PullRequests  GithubPullRequests
	clientMutex   sync.Mutex
}

type GithubClient = *gh.Client

const CLIENT_ID = "Ov23lirJDAVBmN4oRLY0"

func NewGithubController(token *GithubAuthenticationToken, onLogin func(*GithubAuthenticationToken), onLogout func()) *GithubController {
	ga := NewGithubAuthenticator(CLIENT_ID, "repo", "read:org")
	gs := &GithubController{
		authenticator: ga,
		onLogin:       onLogin,
		onLogout:      onLogout,
	}

	if token != nil {
		gs.client = gs.authenticator.InitClient(token)
		gs.onLogin(token)
	}

	return gs
}

func (gs *GithubController) AuthenticateWithDeviceAuth(callback func()) *GithubDevicePreAuthentication {
	gda := gs.authenticator.InitDeviceAuth()
	go func() {
		auth := gs.authenticator.AuthenticateDevice(gda)
		gs.client = gs.authenticator.InitClient(auth)
		gs.onLogin(auth)
		callback()
	}()

	return gda
}

func (gs *GithubController) AuthenticateWithPersonalAccessToken(token string) error {
	client, auth, err := gs.authenticator.AuthenticateWithPersonalAccessToken(token)
	if err != nil {
		return err
	}

	gs.client = client
	gs.onLogin(auth)

	return nil
}

func (gs *GithubController) LogoutUser() {
	gs.clientMutex.Lock()
	defer gs.clientMutex.Unlock()
	gs.client = nil
	gs.login = ""
	gs.onLogout()
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

func (gs *GithubController) GetUserPullRequests() (GithubPullRequests, error) {
	if gs.PullRequests == nil {
		prs, err := gs.FetchUserPullRequests()
		if err != nil {
			return nil, err
		}
		gs.PullRequests = prs
	}
	return gs.PullRequests, nil
}

func (gs *GithubController) FetchUserPullRequests() (GithubPullRequests, error) {
	gs.clientMutex.Lock()
	defer gs.clientMutex.Unlock()
	allRepos := make(map[string]*gh.Repository)
	const limit = 1000

	userRepos, _, err := gs.client.Repositories.ListByAuthenticatedUser(context.Background(), &gh.RepositoryListByAuthenticatedUserOptions{
		ListOptions: gh.ListOptions{
			PerPage: limit,
		},
	})
	if err != nil {
		return nil, err
	}

	for _, repo := range userRepos {
		allRepos[repo.GetFullName()] = repo
	}

	orgs, _, err := gs.client.Organizations.List(context.Background(), "", &gh.ListOptions{
		PerPage: limit,
	})
	if err != nil {
		return nil, err
	}

	for _, org := range orgs {
		orgRepos, _, err := gs.client.Repositories.ListByOrg(context.Background(), *org.Login, &gh.RepositoryListByOrgOptions{
			ListOptions: gh.ListOptions{
				PerPage: limit,
			},
		})
		if err != nil {
			return nil, err
		}

		for _, repo := range orgRepos {
			allRepos[repo.GetFullName()] = repo
		}
	}

	allPrs := NewGithubPrContainer()
	for _, repo := range allRepos {
		prs, _, err := gs.client.PullRequests.List(context.Background(), *repo.GetOwner().Login, *repo.Name, &gh.PullRequestListOptions{
			ListOptions: gh.ListOptions{
				PerPage: limit,
			},
		})
		if err != nil {
			return nil, err
		}

		allPrs.AddPrs(repo, prs...)
	}

	login := gs.GetPrincipalLogin()
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

	return out, nil
}
