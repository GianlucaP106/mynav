package core

import (
	"os"
	"path/filepath"
	"sort"
	"time"
)

type Workspace struct {
	Topic     *Topic
	gitRemote *string
	Name      string
}

func newWorkspace(topic *Topic, name string) *Workspace {
	ws := &Workspace{
		Name:  name,
		Topic: topic,
	}

	return ws
}

func (w *Workspace) Path() string {
	return filepath.Join(w.Topic.Path(), w.Name)
}

func (w *Workspace) ShortPath() string {
	return filepath.Join(w.Topic.Name, w.Name)
}

func (w *Workspace) LastModified() time.Time {
	fi, _ := os.Stat(w.Path())
	return fi.ModTime()
}

func (w *Workspace) GitRemote() (string, error) {
	if w.gitRemote != nil {
		return *(w.gitRemote), nil
	}

	gitPath := filepath.Join(w.Path(), ".git")
	if !Exists(gitPath) {
		return "", nil
	}

	if _, err := filepath.Abs(gitPath); err != nil {
		return "", err
	}

	gitRemote, _ := GitRemote(gitPath)
	w.gitRemote = &gitRemote
	return *(w.gitRemote), nil
}

func (w *Workspace) CloneRepo(url string) error {
	err := GitClone(url, w.Path())
	if err != nil {
		return err
	}
	return nil
}

type Workspaces []*Workspace

func (w Workspaces) Sorted() Workspaces {
	sort.Slice(w, func(i, j int) bool {
		return w[i].LastModified().After(w[j].LastModified())
	})
	return w
}
