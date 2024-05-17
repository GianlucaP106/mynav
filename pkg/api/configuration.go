package api

import (
	"errors"
	"log"
	"mynav/pkg/utils"
	"os"
	"path/filepath"
)

type Configuration struct {
	path                string
	IsConfigInitialized bool
}

func NewConfiguration() *Configuration {
	c := &Configuration{}
	c.DetectConfig()
	return c
}

func (c *Configuration) InitConfig(dir string) string {
	path := filepath.Join(dir, ".mynav")
	if err := utils.CreateDir(path); err != nil {
		log.Panicln(err)
	}

	c.path = dir
	c.IsConfigInitialized = true
	return c.path
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

	configPath, err := func() (string, error) {
		for {
			for _, entry := range dirEntries {
				if entry.Name() == ".mynav" {
					if !entry.IsDir() {
						os.Remove(filepath.Join(cwd, entry.Name()))
						c.InitConfig(cwd)
						return cwd, nil
					}
					homeDir, _ := os.UserHomeDir()
					if cwd == homeDir {
						break
					}

					return cwd, nil
				}
			}
			cwd = filepath.Dir(cwd)
			if cwd == "/" {
				return "", errors.New("no config present")
			}
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

func TimeFormat() string {
	return "Mon, 02 Jan 15:04:05"
}
