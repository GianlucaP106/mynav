package core

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type Filesystem struct {
	path string
}

func newFilesystem(path string) *Filesystem {
	c := &Filesystem{
		path: path,
	}

	return c
}

func (c *Filesystem) CreateTopic(name string) (*Topic, error) {
	t := newTopic(c.path, name)
	if err := CreateDir(t.Path()); err != nil {
		return nil, err
	}
	return t, nil
}

func (c *Filesystem) RenameTopic(t *Topic, name string) error {
	if name == "" {
		return errors.New("name must not be empty")
	}

	oldPath := t.Path()
	t.Name = name
	newPath := t.Path()

	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}

	return nil
}

func (c *Filesystem) DeleteTopic(t *Topic) error {
	if err := os.RemoveAll(t.Path()); err != nil {
		return err
	}

	return nil
}

func (c *Filesystem) CreateWorkspace(t *Topic, name string) (*Workspace, error) {
	if name == "" {
		return nil, errors.New("name must not be empty")
	}

	name = strings.ReplaceAll(name, ".", "_")
	w := newWorkspace(t, name)
	if err := CreateDir(w.Path()); err != nil {
		return nil, err
	}
	return w, nil
}

func (c *Filesystem) MoveWorkspace(w *Workspace, topic *Topic) error {
	if w.Topic.Name == topic.Name {
		return errors.New("workspace is already in this topic")
	}

	if err := os.Rename(w.Path(), filepath.Join(c.path, topic.Name, w.Name)); err != nil {
		return err
	}

	w.Topic = topic
	return nil
}

func (c *Filesystem) RenameWorkspace(w *Workspace, name string) error {
	if name == "" {
		return errors.New("name must not be empty")
	}

	name = strings.ReplaceAll(name, ".", "_")

	oldPath := w.Path()
	w.Name = name
	newPath := w.Path()

	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}

	return nil
}

func (c *Filesystem) DeleteWorkspace(w *Workspace) error {
	if err := os.RemoveAll(w.Path()); err != nil {
		return err
	}
	return nil
}

func (c *Filesystem) Topics() Topics {
	topics := make(Topics, 0)
	for _, topicDir := range GetDirEntries(c.path) {
		if !topicDir.IsDir() || topicDir.Name() == ".mynav" {
			continue
		}

		topicName := topicDir.Name()
		topic := newTopic(c.path, topicName)
		topics = append(topics, topic)
	}
	return topics
}

func (c *Filesystem) Workspaces(t *Topic) Workspaces {
	workspaces := make(Workspaces, 0)
	for _, dirEntry := range GetDirEntries(t.Path()) {
		if !dirEntry.IsDir() {
			continue
		}

		workspace := newWorkspace(t, dirEntry.Name())
		workspaces = append(workspaces, workspace)
	}
	return workspaces
}

func (f *Filesystem) Workspace(shortPath string) *Workspace {
	topicName, workspaceName := filepath.Dir(shortPath), filepath.Base(shortPath)
	workspacePath := filepath.Join(f.path, topicName, workspaceName)
	if !Exists(workspacePath) {
		return nil
	}

	return newWorkspace(newTopic(f.path, topicName), workspaceName)
}

func (c *Filesystem) TopicsCount() int {
	dirs := GetDirEntries(c.path)
	count := 0
	for _, fi := range dirs {
		if fi.IsDir() {
			count++
		}
	}
	return count
}

func (c *Filesystem) AllWorkspaces() Workspaces {
	out := make(Workspaces, 0)
	for _, t := range c.Topics() {
		for _, w := range c.Workspaces(t) {
			out = append(out, w)
		}
	}
	return out
}

func (c *Filesystem) WorkspacesCount() int {
	dirs := GetDirEntries(c.path)
	count := 0
	for _, topic := range dirs {
		if topic.IsDir() {
			for _, workspace := range GetDirEntries(filepath.Join(c.path, topic.Name())) {
				if workspace.IsDir() {
					count++
				}
			}
		}
	}
	return count
}
