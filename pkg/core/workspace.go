package core

import (
	"mynav/pkg/utils"
	"path/filepath"
	"sort"
	"time"
)

type Workspace struct {
	Topic     *Topic
	gitRemote *string
	Name      string
	Path      string
}

func newWorkspace(name string, topic *Topic) *Workspace {
	wsPath := filepath.Join(topic.Filesystem.path, filepath.Join(topic.Name, name))
	ws := &Workspace{
		Name:  name,
		Topic: topic,
		Path:  wsPath,
	}
	return ws
}

func (ws *Workspace) GetGitRemote() string {
	if ws.gitRemote == nil {
		ws.detectGitRemote()
	}
	return *(ws.gitRemote)
}

func (ws *Workspace) detectGitRemote() {
	gitPath := filepath.Join(ws.Path, ".git")
	if _, err := filepath.Abs(gitPath); err != nil {
		return
	}

	gitRemote, _ := utils.GitRemote(gitPath)
	ws.gitRemote = &gitRemote
}

func (ws *Workspace) OpenWorkspace() error {
	if err := utils.OpenTerminal(ws.Path); err != nil {
		return err
	}

	return nil
}

type Workspaces []*Workspace

func (w Workspaces) Len() int      { return len(w) }
func (w Workspaces) Swap(i, j int) { w[i], w[j] = w[j], w[i] }
func (w Workspaces) Less(i, j int) bool {
	return w[i].GetLastModifiedTime().After(w[j].GetLastModifiedTime())
}

func (w Workspaces) Sorted() Workspaces {
	sort.Sort(w)
	return w
}

func (w *Workspace) GetLastModifiedTime() time.Time {
	time, _ := utils.GetLastModifiedTime(w.Path)
	return time
}

func (w *Workspace) GetLastModifiedTimeFormatted() string {
	time := w.GetLastModifiedTime().Format(w.Topic.Filesystem.TimeFormat())
	return time
}
