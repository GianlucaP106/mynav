package api

import (
	"errors"
	"log"
	"mynav/pkg/utils"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/mod/semver"
)

type Configuration struct {
	ConfigurationDatasource *ConfigurationDatasource
	path                    string
	IsConfigInitialized     bool
}

func NewConfiguration() *Configuration {
	c := &Configuration{}
	initialized := c.DetectConfig()
	if initialized {
		cwd, _ := os.Getwd()
		c.InitConfig(cwd)
	}
	return c
}

func (c *Configuration) InitConfig(dir string) (string, error) {
	path := filepath.Join(dir, ".mynav")
	if err := utils.CreateDir(path); err != nil {
		return "", err
	}

	c.path = dir
	c.IsConfigInitialized = true
	c.ConfigurationDatasource = NewConfigurationDatasource(c.GetConfigStorePath())
	return c.path, nil
}

func (c *Configuration) GetConfigPath() string {
	return filepath.Join(c.path, ".mynav")
}

func (c *Configuration) DetectConfig() bool {
	cwd, err := os.Getwd()
	if err != nil {
		log.Panicln(err)
	}
	dirEntries := utils.GetDirEntries(cwd)
	homeDir, _ := os.UserHomeDir()

	configPath, err := func() (string, error) {
		for {
			for _, entry := range dirEntries {
				if cwd == "/" {
					return "", errors.New("no config present")
				}
				if entry.Name() == ".mynav" {
					if cwd == homeDir {
						break
					}

					if !entry.IsDir() {
						os.Remove(filepath.Join(cwd, entry.Name()))
						c.InitConfig(cwd)
						return cwd, nil
					}

					return cwd, nil
				}
			}
			cwd = filepath.Dir(cwd)
			dirEntries = utils.GetDirEntries(cwd)
		}
	}()
	if err != nil {
		return false
	}
	c.path = configPath
	c.IsConfigInitialized = true
	return true
}

func (c *Configuration) GetWorkspaceStorePath() string {
	return filepath.Join(c.GetConfigPath(), "workspaces.json")
}

func (c *Configuration) GetConfigStorePath() string {
	return filepath.Join(c.GetConfigPath(), "config.json")
}

func TimeFormat() string {
	return "Mon, 02 Jan 15:04:05"
}

func (c *Configuration) DetectUpdate() (update bool, newTag string) {
	tag, err := utils.GetLatestReleaseTag()
	if err != nil {
		return false, ""
	}

	res := semver.Compare(tag, VERSION)
	return res == 1, tag
}

func (c *Configuration) SetUpdateAsked() {
	now := time.Now()
	c.ConfigurationDatasource.Data.UpdateAsked = &now
	c.ConfigurationDatasource.SaveStore()
}

func (c *Configuration) IsUpdateAsked() bool {
	time := c.ConfigurationDatasource.Data.UpdateAsked
	if time == nil {
		return false
	}

	return !utils.IsBeforeOneHourAgo(*time)
}

func (c *Configuration) GetUpdateSystemCmd() []string {
	return []string{"sh", "-c", "curl -fsSL https://raw.githubusercontent.com/GianlucaP106/mynav/main/install.sh | bash"}
}
