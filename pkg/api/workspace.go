package api

import (
	"mynav/pkg/utils"
	"path/filepath"
	"sort"
	"time"
)

type WorkspaceMetadata struct {
	TmuxSession *utils.TmuxSession `json:"tmux-session"`
	Description string             `json:"description"`
}

type Workspace struct {
	Topic     *Topic
	GitRemote *string
	Metadata  *WorkspaceMetadata
	Name      string
	Path      string
}

func NewWorkspace(name string, topic *Topic, path string) *Workspace {
	ws := &Workspace{
		Name:     name,
		Topic:    topic,
		Path:     path,
		Metadata: &WorkspaceMetadata{},
	}

	return ws
}

func (w *Workspace) ShortPath() string {
	return filepath.Join(w.Topic.Name, w.Name)
}

func (w *Workspace) GetLastModifiedTime() time.Time {
	time, _ := utils.GetLastModifiedTime(w.Path)
	return time
}

func (w *Workspace) GetLastModifiedTimeFormatted() string {
	time := w.GetLastModifiedTime().Format(TimeFormat())
	return time
}

func (w *Workspace) GetDescription() string {
	if w.Metadata == nil {
		return ""
	}
	return w.Metadata.Description
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
