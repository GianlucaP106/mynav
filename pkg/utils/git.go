package utils

import (
	"encoding/json"
	"io"
	"net/http"
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
	exec.Command("git", "clone", url, path).Run()
	return nil
}

func TrimGithubUrl(url string) string {
	items := strings.Split(url, "/")
	return strings.Join(items[len(items)-2:], "/")
}

type Release struct {
	TagName string `json:"tag_name"`
}

func GetLatestReleaseTag() (string, error) {
	url := "https://api.github.com/repos/GianlucaP106/mynav/releases/latest"

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var release Release
	err = json.Unmarshal(body, &release)
	if err != nil {
		return "", err
	}

	return release.TagName, nil
}
