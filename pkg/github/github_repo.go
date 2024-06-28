package github

import gh "github.com/google/go-github/v62/github"

type GithubRepository struct {
	*gh.Repository
}
