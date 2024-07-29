package github

import (
	"mynav/pkg/constants"
	"mynav/pkg/events"
	"mynav/pkg/persistence"
	"mynav/pkg/tasks"
	"strings"
)

type GithubController struct {
	client GithubClient

	isLoading     *persistence.Value[bool]
	profile       *persistence.Value[GithubProfile]
	prContainer   *persistence.Container[GithubPullRequest]
	repoContainer *persistence.Container[GithubRepository]
}

const FETCH_LIMIT = 1000

func NewGithubController(token *GithubAuthenticationToken, onLogin func(*GithubAuthenticationToken), onLogout func()) *GithubController {
	g := &GithubController{
		client:        NewGithubClient(token, onLogin, onLogout),
		prContainer:   persistence.NewContainer[GithubPullRequest](),
		repoContainer: persistence.NewContainer[GithubRepository](),
		profile:       persistence.NewValue(GithubProfile{}),
		isLoading:     persistence.NewValue(false),
	}

	g.LoadData()

	return g
}

func (g *GithubController) LoadData() {
	if g.IsAuthenticated() {
		tasks.AddTask(func() {
			g.isLoading.Set(true)
			defer g.isLoading.Set(false)
			g.LoadProfile()
			g.LoadUserRepos()
			g.LoadUserPullRequests()
		})
	}
}

func (g *GithubController) IsAuthenticated() bool {
	return g.client.IsAuthenticated()
}

func (g *GithubController) IsLoading() bool {
	return g.isLoading.Get()
}

func (g *GithubController) InitWithDeviceAuth() *GithubDevicePreAuthentication {
	return g.client.AuthenticateWithDevice()
}

func (g *GithubController) InitWithPAT(token string) error {
	return g.client.AuthenticateWithPersonalAccessToken(token)
}

func (g *GithubController) LogoutUser() {
	if g.isLoading.Get() {
		return
	}
	g.client.Logout()
}

func (g *GithubController) LoadUserPullRequests() {
	prs, err := g.fetchUserPullRequests()
	if err != nil {
		return
	}

	g.prContainer.SetAll(prs, func(gpr *GithubPullRequest) string {
		return gpr.GetURL()
	})

	events.Emit(constants.GithubPrsChangesEventName)
}

func (g *GithubController) LoadUserRepos() {
	repos, err := g.fetchUserRepos()
	if err != nil {
		return
	}

	g.repoContainer.SetAll(repos, func(gr *GithubRepository) string {
		return gr.GetURL()
	})

	events.Emit(constants.GithubReposChangesEventName)
}

func (g *GithubController) LoadProfile() {
	profile, err := g.fetchProfile()
	if err != nil {
		return
	}

	g.profile.Set(*profile)
}

func (g *GithubController) GetUserPullRequests() GithubPullRequests {
	return g.prContainer.All()
}

func (g *GithubController) GetUserRepos() []*GithubRepository {
	return g.repoContainer.All()
}

func (g *GithubController) GetProfile() GithubProfile {
	return g.profile.Get()
}

func (g *GithubController) fetchProfile() (*GithubProfile, error) {
	principal, err := g.client.Principal()
	if err != nil {
		return nil, err
	}

	pr := &GithubProfile{
		Login: principal.GetLogin(),
		Name:  principal.GetName(),
		Email: principal.GetEmail(),
		Url:   principal.GetHTMLURL(),
	}

	return pr, nil
}

func (gc *GithubController) fetchUserPullRequests() (GithubPullRequests, error) {
	allRepos, err := gc.fetchUserRepos()
	if err != nil {
		return nil, err
	}

	allPrs := NewGithubPrContainer()
	for _, repo := range allRepos {
		prs, err := gc.client.PullRequests(*repo.GetOwner().Login, *repo.Name, FETCH_LIMIT)
		if err != nil {
			return nil, err
		}

		allPrs.AddPrs(repo, prs...)
	}

	principal, err := gc.client.Principal()
	if err != nil {
		return nil, err
	}

	login := principal.GetLogin()
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

func (gc *GithubController) fetchUserRepos() ([]*GithubRepository, error) {
	allRepos := make(map[string]*GithubRepository)
	userRepos, err := gc.client.PrincipalRepositories(FETCH_LIMIT)
	if err != nil {
		return nil, err
	}

	for _, repo := range userRepos {
		allRepos[repo.GetFullName()] = &GithubRepository{
			Repository: repo,
		}
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
			allRepos[repo.GetFullName()] = &GithubRepository{
				Repository: repo,
			}
		}
	}

	out := make([]*GithubRepository, 0)
	for _, r := range allRepos {
		out = append(out, r)
	}

	return out, nil
}

func TrimGithubUrl(url string) string {
	items := strings.Split(url, "/")
	return strings.Join(items[len(items)-2:], "/")
}
