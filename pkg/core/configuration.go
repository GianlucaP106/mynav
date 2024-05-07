package core

import (
	"errors"
	"log"
	"mynav/pkg/utils"
	"os"
	"path/filepath"
)

func (fs *Filesystem) CreateConfigurationFile() string {
	dir, _ := os.Getwd()

	path := dir + "/.mynav"
	if _, err := os.Create(path); err != nil {
		log.Panicln(err)
	}

	fs.path = dir
	fs.Initialized = true
	return path
}

func (fs *Filesystem) detectConfig() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Panicln(err)
	}
	dirEntries := utils.GetDirEntries(cwd)

	configPath, err := func() (string, error) {
		for {
			for _, entry := range dirEntries {
				if !entry.IsDir() && entry.Name() == ".mynav" {
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
		return
	}
	fs.path = configPath
	fs.Initialized = true
}

func (fs *Filesystem) getTimeFormat() string {
	return "Mon, 02 Jan 15:04:05"
}
