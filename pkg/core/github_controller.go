package core

import (
	"mynav/pkg/persistence"
	"strings"

	"github.com/google/go-github/v62/github"
)

type GithubController struct {
	client        GithubClient
	config        *GlobalConfiguration
	profile       *persistence.Value[*github.User]
	prContainer   *persistence.Container[github.PullRequest]
	repoContainer *persistence.Container[github.Repository]
}

const FETCH_LIMIT = 1000

func NewGithubController(config *GlobalConfiguration) *GithubController {
	token := config.GetGithubToken()
	g := &GithubController{
		config:        config,
		client:        NewGithubClient(token),
		prContainer:   persistence.NewContainer[github.PullRequest](),
		repoContainer: persistence.NewContainer[github.Repository](),
		profile:       persistence.NewValue[*github.User](nil),
	}

	return g
}

func (g *GithubController) IsAuthenticated() bool {
	return g.client.IsAuthenticated()
}

func (g *GithubController) InitWithDeviceAuth() (*GithubDevicePreAuthentication, func()) {
	pa, poll := g.client.AuthenticateWithDevice()
	return pa, func() {
		token := poll()
		g.config.SetGithubToken(token)
	}
}

func (g *GithubController) InitWithPAT(token string) error {
	gtoken, err := g.client.AuthenticateWithPersonalAccessToken(token)
	if err != nil {
		return err
	}

	g.config.SetGithubToken(gtoken)
	return nil
}

func (g *GithubController) LogoutUser() {
	g.config.SetGithubToken(nil)
}

func (g *GithubController) LoadUserPullRequests() {
	prs, err := g.fetchUserPullRequests()
	if err != nil {
		return
	}

	g.prContainer.SetAll(prs, func(gpr *github.PullRequest) string {
		return gpr.GetURL()
	})
}

func (g *GithubController) LoadUserRepos() {
	repos, err := g.fetchUserRepos()
	if err != nil {
		return
	}

	g.repoContainer.SetAll(repos, func(gr *github.Repository) string {
		return gr.GetURL()
	})
}

func (g *GithubController) LoadProfile() {
	profile, err := g.fetchProfile()
	if err != nil {
		return
	}

	g.profile.Set(profile)
}

func (g *GithubController) GetUserPullRequests() []*github.PullRequest {
	return g.prContainer.All()
}

func (g *GithubController) GetUserRepos() []*github.Repository {
	return g.repoContainer.All()
}

func (g *GithubController) GetPrincipal() *github.User {
	return g.profile.Get()
}

func (g *GithubController) GetPrRelation(pr *github.PullRequest, principal *github.User) (bool, string) {
	login := principal.GetLogin()
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
}

func (g *GithubController) fetchProfile() (*github.User, error) {
	principal, err := g.client.Principal()
	if err != nil {
		return nil, err
	}

	return principal, nil
}

func (gc *GithubController) fetchUserPullRequests() ([]*github.PullRequest, error) {
	allRepos, err := gc.fetchUserRepos()
	if err != nil {
		return nil, err
	}

	allPrs := make([]*github.PullRequest, 0)
	for _, repo := range allRepos {
		prs, err := gc.client.PullRequests(*repo.GetOwner().Login, *repo.Name, FETCH_LIMIT)
		if err != nil {
			return nil, err
		}

		allPrs = append(allPrs, prs...)
	}

	principal, err := gc.client.Principal()
	if err != nil {
		return nil, err
	}

	out := make([]*github.PullRequest, 0)
	for _, pr := range allPrs {
		belongsToUser, _ := gc.GetPrRelation(pr, principal)
		if belongsToUser {
			out = append(out, pr)
		}
	}

	return out, nil
}

func (gc *GithubController) fetchUserRepos() ([]*github.Repository, error) {
	allRepos := make(map[string]*github.Repository)
	userRepos, err := gc.client.PrincipalRepositories(FETCH_LIMIT)
	if err != nil {
		return nil, err
	}

	for _, repo := range userRepos {
		allRepos[repo.GetFullName()] = repo
	}

	orgs, err := gc.client.PrincipalOrganizations(FETCH_LIMIT)
	if err != nil {
		return nil, err
	}

	for _, org := range orgs {
		orgRepos, err := gc.client.RepositoriesByOrg(org.GetLogin(), FETCH_LIMIT)
		if err != nil {
			return nil, err
		}

		for _, repo := range orgRepos {
			allRepos[repo.GetFullName()] = repo
		}
	}

	out := make([]*github.Repository, 0)
	for _, r := range allRepos {
		out = append(out, r)
	}

	return out, nil
}

func TrimGithubUrl(url string) string {
	items := strings.Split(url, "/")
	return strings.Join(items[len(items)-2:], "/")
}
