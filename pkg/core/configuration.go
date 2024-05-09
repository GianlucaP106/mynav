package core

import (
	"errors"
	"log"
	"mynav/pkg/utils"
	"os"
	"path/filepath"
)

func (fs *Filesystem) InitConfiguration(dir string) string {
	path := filepath.Join(dir, ".mynav")
	if err := utils.CreateDir(path); err != nil {
		log.Panicln(err)
	}

	fs.path = dir
	fs.ConfigInitialized = true
	return fs.path
}

func (fs *Filesystem) GetConfigPath() string {
	return filepath.Join(fs.path, ".mynav")
}

func (fs *Filesystem) DetectConfig() bool {
	cwd, err := os.Getwd()
	if err != nil {
		log.Panicln(err)
	}
	dirEntries := utils.GetDirEntries(cwd)

	configPath, err := func() (string, error) {
		for {
			for _, entry := range dirEntries {
				if entry.Name() == ".mynav" {
					// if .mynav as a file exists (prev version of mynav)
					if !entry.IsDir() {
						os.Remove(filepath.Join(cwd, entry.Name()))
						fs.InitConfiguration(cwd)
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
	fs.path = configPath
	fs.ConfigInitialized = true
	return true
}

func (fs *Filesystem) TimeFormat() string {
	return "Mon, 02 Jan 15:04:05"
}
