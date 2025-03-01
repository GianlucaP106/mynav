package core

import (
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/GianlucaP106/mynav/pkg/system"
)

// Workspace
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

func (w *Workspace) LastModifiedTime() time.Time {
	time, _ := system.GetLastModifiedTime(w.Path())
	return time
}

func (w *Workspace) LastModifiedTimeFormatted() string {
	time := w.LastModifiedTime().Format(system.TimeFormat())
	return time
}

func (w *Workspace) GitRemote() (string, error) {
	if w.gitRemote != nil {
		return *(w.gitRemote), nil
	}

	gitPath := filepath.Join(w.Path(), ".git")
	if !system.Exists(gitPath) {
		return "", nil
	}

	if _, err := filepath.Abs(gitPath); err != nil {
		return "", err
	}

	gitRemote, _ := system.GitRemote(gitPath)
	w.gitRemote = &gitRemote
	return *(w.gitRemote), nil
}

func (w *Workspace) CloneRepo(url string) error {
	err := system.GitClone(url, w.Path())
	if err != nil {
		return err
	}
	return nil
}

// Workspaces is a collection of Workspace.
type Workspaces []*Workspace

func (w Workspaces) Sorted() Workspaces {
	sort.Slice(w, func(i, j int) bool {
		return w[i].LastModifiedTime().After(w[j].LastModifiedTime())
	})
	return w
}

func (w Workspaces) ByNameContaining(s string) Workspaces {
	if s == "" {
		return w
	}

	filtered := Workspaces{}
	for _, workspace := range w {
		if strings.Contains(workspace.Name, s) {
			filtered = append(filtered, workspace)
		}
	}
	return filtered
}

func (w Workspaces) ByTopic(topic *Topic) Workspaces {
	out := make(Workspaces, 0)
	for _, workspace := range w {
		if workspace.Topic == topic {
			out = append(out, workspace)
		}
	}
	return out
}
