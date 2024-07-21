package configuration

import (
	"log"
	"mynav/pkg"
	"mynav/pkg/git"
	"mynav/pkg/github"
	"mynav/pkg/persistence"
	"mynav/pkg/system"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/mod/semver"
)

type GlobalConfigurationDataSchema struct {
	UpdateAsked *time.Time                        `json:"update-asked"`
	GithubToken *github.GithubAuthenticationToken `json:"github-token"`
	LastTab     string                            `json:"last-tab"`
}

type GlobalConfiguration struct {
	Datasource *persistence.Datasource[GlobalConfigurationDataSchema]
	Standalone bool
}

type Configuration struct {
	*LocalConfiguration
	*GlobalConfiguration
}

func NewGlobalConfiguration() *GlobalConfiguration {
	gc := &GlobalConfiguration{}
	gc.Standalone = system.IsCurrentProcessHomeDir()
	gc.Datasource = persistence.NewDatasource[GlobalConfigurationDataSchema](gc.GetConfigFile())
	gc.Datasource.LoadData()
	if gc.Datasource.GetData() == nil {
		gc.Datasource.SaveData(&GlobalConfigurationDataSchema{})
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
	data := gc.Datasource.GetData()
	data.UpdateAsked = &now
	gc.Datasource.SaveData(data)
}

func (gc *GlobalConfiguration) IsUpdateAsked() bool {
	time := gc.Datasource.GetData().UpdateAsked
	if time == nil {
		return false
	}

	return isBeforeOneHourAgo(*time)
}

func (gc *GlobalConfiguration) GetGithubToken() *github.GithubAuthenticationToken {
	return gc.Datasource.GetData().GithubToken
}

func (gc *GlobalConfiguration) SetGithubToken(token *github.GithubAuthenticationToken) {
	data := gc.Datasource.GetData()
	data.GithubToken = token
	gc.Datasource.SaveData(data)
}

func (gc *GlobalConfiguration) UpdateMynav() error {
	return system.Command("sh", "-c", "curl -fsSL https://raw.githubusercontent.com/GianlucaP106/mynav/main/install.sh | bash").Run()
}

func (c *GlobalConfiguration) SetLastTab(lastTab string) {
	data := c.Datasource.GetData()
	data.LastTab = lastTab
	c.Datasource.SaveData(data)
}

func (c *GlobalConfiguration) GetLastTab() string {
	return c.Datasource.GetData().LastTab
}

func (c *GlobalConfiguration) SetStandalone(s bool) {
	c.Standalone = s
}
