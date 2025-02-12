package core

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/GianlucaP106/mynav/pkg/system"
)

// Data for the local config store.
type LocalConfigData struct {
	SelectedWorkspace string `json:"selected-workspace"`
}

// LocalConfig is the LocalConfig configuration.
type LocalConfig struct {
	datasource *Datasource[LocalConfigData]
	path       string
}

func newLocalConfig(dir string) (*LocalConfig, error) {
	c := &LocalConfig{}
	// if dir is passed we initialize it and dont detect
	if dir != "" {
		// check if dir is home dir
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		if home == dir {
			return nil, errors.New("mynav cannot be initialized in the home directory")
		}

		// set up dir
		if err := c.setupDir(dir); err != nil {
			return nil, err
		}

		// set up datasource in the dir
		if err := c.setupDatasource(dir); err != nil {
			return nil, err
		}
		return c, nil
	}

	// if dir is not passed we detect
	path, err := c.detect()
	if err != nil {
		return nil, err
	}

	// return no error and nil if no config here
	if path == "" {
		return nil, nil
	}

	// if config exists set up datasource
	if err := c.setupDatasource(path); err != nil {
		return nil, err
	}

	return c, nil
}

func (l *LocalConfig) setupDatasource(rootdir string) error {
	ds, err := newDatasource(filepath.Join(rootdir, ".mynav", "config.json"), &LocalConfigData{})
	if err != nil {
		return err
	}

	l.path = rootdir
	l.datasource = ds
	return nil
}

func (c *LocalConfig) setupDir(rootdir string) error {
	path := filepath.Join(rootdir, ".mynav")
	return system.CreateDir(path)
}

func (c *LocalConfig) detect() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Panicln(err)
	}
	dirEntries := system.GetDirEntries(cwd)
	homeDir, _ := os.UserHomeDir()

	for {
		for _, entry := range dirEntries {
			if cwd == "/" {
				return "", nil
			}
			if entry.Name() == ".mynav" {
				if cwd == homeDir {
					break
				}

				if !entry.IsDir() {
					os.Remove(filepath.Join(cwd, entry.Name()))
					c.setupDir(cwd)
					return cwd, nil
				}

				return cwd, nil
			}
		}
		cwd = filepath.Dir(cwd)
		dirEntries = system.GetDirEntries(cwd)
	}
}

func (g *LocalConfig) SetSelectedWorkspace(s string) {
	data := g.datasource.Get()
	data.SelectedWorkspace = s
	g.datasource.Save(data)
}

func (l *LocalConfig) ConfigData() *LocalConfigData {
	return l.datasource.Get()
}

func isBeforeOneHourAgo(timestamp time.Time) bool {
	now := time.Now()
	oneHourAgo := now.Add(-1 * time.Hour)
	return timestamp.Before(oneHourAgo)
}
