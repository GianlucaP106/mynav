package core

import (
	"errors"
	"log"
	"mynav/pkg/system"
	"os"
	"path/filepath"
	"time"
)

type LocalConfiguration struct {
	path          string
	IsInitialized bool
}

func NewLocalConfiguration() *LocalConfiguration {
	c := &LocalConfiguration{}
	c.DetectConfig()
	return c
}

func (c *LocalConfiguration) InitConfig(dir string) (string, error) {
	path := filepath.Join(dir, ".mynav")
	if err := system.CreateDir(path); err != nil {
		return "", err
	}

	c.path = dir
	c.IsInitialized = true
	return c.path, nil
}

func (c *LocalConfiguration) GetConfigPath() string {
	return filepath.Join(c.path, ".mynav")
}

func (c *LocalConfiguration) DetectConfig() bool {
	cwd, err := os.Getwd()
	if err != nil {
		log.Panicln(err)
	}
	dirEntries := system.GetDirEntries(cwd)
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
			dirEntries = system.GetDirEntries(cwd)
		}
	}()
	if err != nil {
		return false
	}
	c.path = configPath
	c.IsInitialized = true
	return true
}

func (c *LocalConfiguration) GetWorkspaceStorePath() string {
	return filepath.Join(c.GetConfigPath(), "workspaces.json")
}

func (c *LocalConfiguration) GetSocketPath() string {
	return filepath.Join(c.GetConfigPath(), "mynav.sock")
}

func IsParentAppInstance() bool {
	return !IsTmuxSession()
}

func isBeforeOneHourAgo(timestamp time.Time) bool {
	now := time.Now()
	oneHourAgo := now.Add(-1 * time.Hour)
	return timestamp.Before(oneHourAgo)
}
