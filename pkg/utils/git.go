package utils

import (
	"context"
	"log"
	"os/exec"
	"strings"

	"github.com/google/go-github/v62/github"
)

func GitRemote(path string) (string, error) {
	out, err := exec.Command("git", "--git-dir="+path, "config", "--get", "remote.origin.url").Output()
	if err != nil {
		return "", err
	}

	gitRemote := strings.ReplaceAll(string(out), "\n", "")
	return gitRemote, nil
}

func GitClone(url string, path string) error {
	return exec.Command("git", "clone", url, path).Run()
}

func TrimGithubUrl(url string) string {
	items := strings.Split(url, "/")
	return strings.Join(items[len(items)-2:], "/")
}

func GetLatestReleaseTag(repo string, owner string) string {
	ctx := context.Background()

	client := github.NewClient(nil)

	release, _, err := client.Repositories.GetLatestRelease(ctx, owner, repo)
	if err != nil {
		log.Panicln(err)
	}

	return *release.TagName
}
