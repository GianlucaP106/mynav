package api

import (
	"mynav/pkg/utils"
	"path/filepath"
	"sort"
	"strings"
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

func (w Workspaces) FilterByNameContaining(s string) Workspaces {
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
	return out.Sorted()
}

func (w Workspaces) GetWorkspace(idx int) *Workspace {
	if idx >= len(w) || idx < 0 {
		return nil
	}
	return w[idx]
}

func (w Workspaces) Sorted() Workspaces {
	sort.Sort(w)
	return w
}
