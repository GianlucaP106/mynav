package core

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
	exec.Command("git", "clone", url, path).Run()
	return nil
}
