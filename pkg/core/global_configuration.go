package core

import (
	"log"
	"mynav/pkg"
	"mynav/pkg/git"
	"mynav/pkg/github"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/mod/semver"
)

type GlobalConfigurationDataSchema struct {
	UpdateAsked *time.Time                        `json:"update-asked"`
	GithubToken *github.GithubAuthenticationToken `json:"github-token"`
}

type GlobalConfiguration struct {
	Datasource *Datasource[GlobalConfigurationDataSchema]
}

func NewGlobalConfiguration() *GlobalConfiguration {
	gc := &GlobalConfiguration{}
	gc.Datasource = NewDatasource[GlobalConfigurationDataSchema](gc.GetConfigFile())
	gc.Datasource.LoadData()
	if gc.Datasource.Data == nil {
		gc.Datasource.Data = &GlobalConfigurationDataSchema{}
	}

	return gc
}

func (gc *GlobalConfiguration) GetGlobalConfigDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Panicln(err)
	}

	return filepath.Join(dir, ".mynav")
}

func (gc *GlobalConfiguration) GetConfigFile() string {
	dir := gc.GetGlobalConfigDir()
	return filepath.Join(dir, "config.json")
}

func (gc *GlobalConfiguration) DetectUpdate() (update bool, newTag string) {
	tag, err := git.GetLatestReleaseTag()
	if err != nil {
		return false, ""
	}

	res := semver.Compare(tag, pkg.VERSION)
	return res == 1, tag
}

func (gc *GlobalConfiguration) SetUpdateAsked() {
	now := time.Now()
	gc.Datasource.Data.UpdateAsked = &now
	gc.Datasource.SaveData()
}

func (gc *GlobalConfiguration) IsUpdateAsked() bool {
	time := gc.Datasource.Data.UpdateAsked
	if time == nil {
		return false
	}

	return isBeforeOneHourAgo(*time)
}

func (gc *GlobalConfiguration) GetGithubToken() *github.GithubAuthenticationToken {
	return gc.Datasource.Data.GithubToken
}

func (gc *GlobalConfiguration) SetGithubToken(token *github.GithubAuthenticationToken) {
	gc.Datasource.Data.GithubToken = token
	gc.Datasource.SaveData()
}

func GetUpdateSystemCmd() []string {
	return []string{"sh", "-c", "curl -fsSL https://raw.githubusercontent.com/GianlucaP106/mynav/main/install.sh | bash"}
}
