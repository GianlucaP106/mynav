package core

import (
	"errors"
	"io/fs"
	"mynav/pkg/utils"
	"os"
	"path/filepath"
	"sync"
)

type Filesystem struct {
	path              string
	Workspaces        Workspaces
	Topics            Topics
	ConfigInitialized bool
}

func NewFilesystem() *Filesystem {
	fs := &Filesystem{}
	fs.detectConfig()
	fs.InitFilesystem()
	return fs
}

func (fs *Filesystem) InitFilesystem() {
	if !fs.ConfigInitialized {
		return
	}

	topics := utils.GetDirEntries(fs.path)
	tps := make(Topics, 0)
	wsp := make(Workspaces, 0)
	for _, entry := range topics {
		if !entry.IsDir() {
			continue
		}
		topic := newTopic(entry.Name(), fs)
		tps = append(tps, topic)

		workspaces := topic.InitTopicWorkspaces()
		wsp = append(wsp, workspaces...)
	}

	fs.Workspaces = wsp
	fs.Topics = tps
}

func (f *Filesystem) InitFilesystemAsync() {
	if !f.ConfigInitialized {
		return
	}

	topics := make([]fs.FileInfo, 0)
	for _, fsItem := range utils.GetDirEntries(f.path) {
		if fsItem.IsDir() {
			topics = append(topics, fsItem)
		}
	}

	resultChan := make(chan Workspaces, len(topics))
	var wg sync.WaitGroup

	resultTopics := make(Topics, 0)
	for _, entry := range topics {
		if !entry.IsDir() {
			continue
		}
		topic := newTopic(entry.Name(), f)
		resultTopics = append(resultTopics, topic)

		wg.Add(1)
		go func() {
			defer wg.Done()
			workspaces := topic.InitTopicWorkspaces()
			resultChan <- workspaces
		}()
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	resultWorkspaces := make(Workspaces, 0)
	for w := range resultChan {
		resultWorkspaces = append(resultWorkspaces, w...)
	}

	f.Workspaces = resultWorkspaces
	f.Topics = resultTopics
}

func (fs *Filesystem) CreateTopic(name string) (*Topic, error) {
	newTopicPath := filepath.Join(fs.path, name)
	if err := os.Mkdir(newTopicPath, 0755); err != nil {
		return nil, err
	}

	topic := newTopic(name, fs)
	fs.Topics = append(fs.Topics, topic)

	return topic, nil
}

func (fs *Filesystem) CreateWorkspace(name string, repoUrl string, topic *Topic) (*Workspace, error) {
	newWorkspacePath := filepath.Join(fs.path, topic.Name, name)
	if repoUrl != "" {
		if err := utils.GitClone(repoUrl, newWorkspacePath); err != nil {
			return nil, errors.New("failed to clone repository")
		}
	}

	if err := os.Mkdir(newWorkspacePath, 0755); err != nil {
		return nil, err
	}

	workspace := newWorkspace(name, topic)
	fs.Workspaces = append(fs.Workspaces, workspace)

	return workspace, nil
}

func (fs *Filesystem) DeleteTopic(topic *Topic) error {
	topicPath := filepath.Join(fs.path, topic.Name)

	idx := 0
	for i, t := range fs.Topics {
		if t == topic {
			idx = i
		}
	}

	if err := os.RemoveAll(topicPath); err != nil {
		return err
	}

	fs.Topics = append(fs.Topics[:idx], fs.Topics[idx+1:]...)
	return nil
}

func (fs *Filesystem) DeleteWorkspace(workspace *Workspace) error {
	if err := os.RemoveAll(workspace.Path); err != nil {
		return err
	}
	idx := 0
	for i, t := range fs.Workspaces {
		if t == workspace {
			idx = i
		}
	}
	fs.Workspaces = append(fs.Workspaces[:idx], fs.Workspaces[idx+1:]...)
	return nil
}
