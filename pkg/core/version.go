package core

import (
	"encoding/json"
	"io"
	"net/http"
	"os/exec"

	"golang.org/x/mod/semver"
)

type updater struct{}

const Version = "v2.1.1"

// Updates mynav by running update script.
func (u *updater) UpdateMynav() error {
	return exec.Command("sh", "-c", "curl -fsSL https://raw.githubusercontent.com/GianlucaP106/mynav/main/install.sh | bash").Run()
}

// Returns if a mynav update is available.
func (u *updater) UpdateAvailable() (bool, string) {
	tag, err := u.getLatestReleaseTag()
	if err != nil {
		return false, ""
	}

	res := semver.Compare(tag, Version)
	return res == 1, tag
}

// Returns the latest tag
func (u *updater) getLatestReleaseTag() (string, error) {
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

	var release struct {
		TagName string `json:"tag_name"`
	}
	err = json.Unmarshal(body, &release)
	if err != nil {
		return "", err
	}

	return release.TagName, nil
}
