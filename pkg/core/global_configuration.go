package core

import (
	"encoding/json"
	"io"
	"log"
	"mynav/pkg"
	"mynav/pkg/system"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/mod/semver"
)

type GlobalConfigurationDataSchema struct {
	UpdateAsked             *time.Time                 `json:"update-asked"`
	GithubToken             *GithubAuthenticationToken `json:"github-token"`
	LastTab                 string                     `json:"last-tab"`
	CustomWorspaceOpenerCmd []string                   `json:"custom-workspace-openner"`
	TerminalOpenerCmd       []string                   `json:"terminal-opener-command"`
	EnableGithubTab         bool                       `json:"enable-github-tab"`
}

type GlobalConfiguration struct {
	Datasource *Datasource[GlobalConfigurationDataSchema]
	Standalone bool
}

type Configuration struct {
	*LocalConfiguration
	*GlobalConfiguration
}

func NewGlobalConfiguration() (*GlobalConfiguration, error) {
	gc := &GlobalConfiguration{}
	gc.Standalone = system.IsCurrentProcessHomeDir()
	ds, err := NewDatasource(gc.GetConfigFile(), &GlobalConfigurationDataSchema{})
	if err != nil {
		return nil, err
	}

	gc.Datasource = ds
	return gc, nil
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
	tag, err := getLatestReleaseTag()
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

func (gc *GlobalConfiguration) GetGithubToken() *GithubAuthenticationToken {
	return gc.Datasource.GetData().GithubToken
}

func (gc *GlobalConfiguration) SetGithubToken(token *GithubAuthenticationToken) {
	data := gc.Datasource.GetData()
	data.GithubToken = token
	gc.Datasource.SaveData(data)
}

func (gc *GlobalConfiguration) UpdateMynav() error {
	return exec.Command("sh", "-c", "curl -fsSL https://raw.githubusercontent.com/GianlucaP106/mynav/main/install.sh | bash").Run()
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

func (c *GlobalConfiguration) GetCustomWorkspaceOpenerCmd() []string {
	data := c.Datasource.GetData()
	return data.CustomWorspaceOpenerCmd
}

func (c *GlobalConfiguration) SetCustomWorkspaceOpenerCmd(cmd string) {
	data := c.Datasource.GetData()
	command := []string{}
	if cmd != "" {
		command = append(command, strings.Split(cmd, " ")...)
	}

	data.CustomWorspaceOpenerCmd = command
	c.Datasource.SaveData(data)
}

func (c *GlobalConfiguration) GetTerminalOpenerCmd() []string {
	data := c.Datasource.GetData()
	return data.TerminalOpenerCmd
}

func (c *GlobalConfiguration) SetTerminalOpenerCmd(cmd string) {
	data := c.Datasource.GetData()
	command := []string{}
	if cmd != "" {
		command = append(command, strings.Split(cmd, " ")...)
	}

	data.TerminalOpenerCmd = command
	c.Datasource.SaveData(data)
}

func (c *GlobalConfiguration) GetGithubTabEnabled() bool {
	return c.Datasource.GetData().EnableGithubTab
}

func (c *GlobalConfiguration) SetGithubTabEnabled(enabled bool) {
	data := c.Datasource.GetData()
	data.EnableGithubTab = enabled
	c.Datasource.SaveData(data)
}

type Release struct {
	TagName string `json:"tag_name"`
}

func getLatestReleaseTag() (string, error) {
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
