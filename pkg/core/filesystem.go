package core

import (
	"errors"
	"mynav/pkg/utils"
	"os"
	"path/filepath"
)

type Filesystem struct {
	path        string
	Workspaces  Workspaces
	Topics      Topics
	Initialized bool
}

func NewFilesystem() *Filesystem {
	fs := &Filesystem{}
	fs.detectConfig()
	fs.InitFilesystem()
	return fs
}

func (fs *Filesystem) InitFilesystem() {
	if !fs.Initialized {
		return
	}
	wsp := make(Workspaces, 0)
	tps := make(Topics, 0)
	configDir := fs.path
	topics := utils.GetDirEntries(configDir)
	for _, entry := range topics {
		if !entry.IsDir() {
			continue
		}
		topic := newTopic(entry.Name(), fs)
		tps = append(tps, topic)
		workspaces := utils.GetDirEntries(filepath.Join(configDir, entry.Name()))
		for _, workspace := range workspaces {
			if !workspace.IsDir() {
				continue
			}
			w := newWorkspace(workspace.Name(), topic, fs)
			wsp = append(wsp, w)
		}
	}
	fs.Workspaces = wsp
	fs.Topics = tps
}

func (fs *Filesystem) CreateTopic(name string) error {
	newTopicPath := filepath.Join(fs.path, name)
	return os.Mkdir(newTopicPath, 0755)
}

func (fs *Filesystem) CreateWorkspace(name string, repoUrl string, topic *Topic) error {
	newWorkspacePath := filepath.Join(fs.path, topic.Name, name)
	if repoUrl != "" {
		if err := utils.GitClone(repoUrl, newWorkspacePath); err != nil {
			return errors.New("failed to clone repository")
		}
	} else {
		return os.Mkdir(newWorkspacePath, 0755)
	}

	return nil
}

func (fs *Filesystem) DeleteTopic(topic *Topic) error {
	topicPath := filepath.Join(fs.path, topic.Name)
	return os.RemoveAll(topicPath)
}

func (fs *Filesystem) DeleteWorkspace(workspace *Workspace) error {
	return os.RemoveAll(workspace.Path)
}
