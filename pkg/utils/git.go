package utils

import (
	"os/exec"
	"strings"
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
