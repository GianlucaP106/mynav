package core

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/GianlucaP106/mynav/pkg/system"
)

type Container struct {
	topics map[string]*Topic
	path   string
	mu     sync.RWMutex
}

func newContainer(path string) *Container {
	c := &Container{
		topics: make(map[string]*Topic),
		path:   path,
	}

	for _, topicDir := range system.GetDirEntries(c.path) {
		if !topicDir.IsDir() || topicDir.Name() == ".mynav" {
			continue
		}

		topicName := topicDir.Name()
		topic := newTopic(c.path, topicName)
		for _, dirEntry := range system.GetDirEntries(topic.Path()) {
			if !dirEntry.IsDir() {
				continue
			}

			workspace := newWorkspace(topic, dirEntry.Name())
			topic.workspaces[workspace.Name] = workspace
		}
		c.topics[topic.Name] = topic
	}
	return c
}

func (c *Container) CreateTopic(name string) (*Topic, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	t := newTopic(c.path, name)
	if err := system.CreateDir(t.Path()); err != nil {
		return nil, err
	}

	c.topics[name] = t
	return t, nil
}

func (c *Container) RenameTopic(t *Topic, name string) error {
	if name == "" {
		return errors.New("name must not be empty")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.topics[name] != nil {
		return errors.New("topic with this name already exists")
	}

	oldName := t.Name
	oldPath := t.Path()
	t.Name = name
	newPath := t.Path()

	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}

	delete(c.topics, oldName)
	c.topics[name] = t
	return nil
}

func (c *Container) DeleteTopic(t *Topic) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := os.RemoveAll(t.Path()); err != nil {
		return err
	}

	delete(c.topics, t.Name)
	return nil
}

func (c *Container) CreateWorkspace(t *Topic, name string) (*Workspace, error) {
	if name == "" {
		return nil, errors.New("name must not be empty")
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	name = strings.ReplaceAll(name, ".", "_")
	w := newWorkspace(t, name)
	if err := system.CreateDir(w.Path()); err != nil {
		return nil, err
	}
	t.workspaces[name] = w
	return w, nil
}

func (c *Container) MoveWorkspace(w *Workspace, topic *Topic) error {
	if w.Topic == topic {
		return errors.New("workspace is already in this topic")
	}

	topic.mu.Lock()
	defer topic.mu.Unlock()

	w.Topic.mu.Lock()
	defer w.Topic.mu.Unlock()

	// if there exists a workspace with this name in the same topic
	if topic.workspaces[w.Name] != nil {
		return errors.New("workspace with this name already exists")
	}

	if err := os.Rename(w.Path(), filepath.Join(c.path, topic.Name, w.Name)); err != nil {
		return err
	}

	delete(w.Topic.workspaces, w.Name)
	topic.workspaces[w.Name] = w
	w.Topic = topic
	return nil
}

func (c *Container) RenameWorkspace(w *Workspace, name string) error {
	if name == "" {
		return errors.New("name must not be empty")
	}

	t := w.Topic
	t.mu.Lock()
	defer t.mu.Unlock()

	name = strings.ReplaceAll(name, ".", "_")
	if t.workspaces[name] != nil {
		return errors.New("workspace with this name already exists")
	}

	oldName := w.Name
	oldPath := w.Path()
	w.Name = name
	newPath := w.Path()

	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}

	delete(t.workspaces, oldName)
	t.workspaces[name] = w
	return nil
}

func (c *Container) DeleteWorkspace(w *Workspace) error {
	t := w.Topic
	t.mu.Lock()
	defer t.mu.Unlock()
	if err := os.RemoveAll(w.Path()); err != nil {
		return err
	}

	delete(t.workspaces, w.Name)
	return nil
}

func (c *Container) AllTopics() Topics {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make(Topics, 0, len(c.topics))
	for _, t := range c.topics {
		out = append(out, t)
	}
	return out
}

func (c *Container) TopicsCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.topics)
}

func (c *Container) AllWorkspaces() Workspaces {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make(Workspaces, 0)
	for _, t := range c.topics {
		t.mu.RLock()
		for _, w := range t.workspaces {
			out = append(out, w)
		}
		t.mu.RUnlock()
	}
	return out
}

func (c *Container) WorkspacesCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	count := 0
	for _, t := range c.topics {
		t.mu.RLock()
		count += len(t.workspaces)
		t.mu.RUnlock()
	}
	return count
}
