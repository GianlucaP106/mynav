package github

import "github.com/google/go-github/v62/github"

type GithubPullRequests []*GithubPullRequest

type GithubPullRequest struct {
	Repo *GithubRepository
	*github.PullRequest
	Relation string
}

func NewGithubPrContainer() GithubPullRequests {
	return make(GithubPullRequests, 0)
}

func (g *GithubPullRequests) AddFromPr(pr *GithubPullRequest) {
	*g = append(*g, &GithubPullRequest{
		Repo:        pr.Repo,
		PullRequest: pr.PullRequest,
		Relation:    pr.Relation,
	})
}

func (g *GithubPullRequests) AddPrs(repo *GithubRepository, prs ...*github.PullRequest) {
	for _, pr := range prs {
		*g = append(*g, &GithubPullRequest{
			Repo:        repo,
			PullRequest: pr,
		})
	}
}
