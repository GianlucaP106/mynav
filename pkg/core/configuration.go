package core

import (
	"errors"
	"log"
	"mynav/pkg/filesystem"
	"mynav/pkg/git"
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
	c.DetectConfig()
	c.ConfigurationDatasource = NewConfigurationDatasource(c.GetConfigStorePath())
	return c
}

func (c *Configuration) InitConfig(dir string) (string, error) {
	path := filepath.Join(dir, ".mynav")
	if err := filesystem.CreateDir(path); err != nil {
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

func (c *Configuration) GetConfigDir() string {
	return c.path
}

func (c *Configuration) DetectConfig() bool {
	cwd, err := os.Getwd()
	if err != nil {
		log.Panicln(err)
	}
	dirEntries := filesystem.GetDirEntries(cwd)
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
			dirEntries = filesystem.GetDirEntries(cwd)
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
	tag, err := git.GetLatestReleaseTag()
	if err != nil {
		return false, ""
	}

	// TODO:
	res := semver.Compare(tag, "TODO")
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

	return isBeforeOneHourAgo(*time)
}

func isBeforeOneHourAgo(timestamp time.Time) bool {
	now := time.Now()
	oneHourAgo := now.Add(-1 * time.Hour)
	return timestamp.Before(oneHourAgo)
}
