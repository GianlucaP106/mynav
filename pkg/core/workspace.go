package core

import (
	"mynav/pkg/system"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Workspace
type Workspace struct {
	Topic     *Topic
	gitRemote *string
	Name      string
	path      string
}

func newWorkspace(topic *Topic, name string) *Workspace {
	ws := &Workspace{
		Name:  name,
		Topic: topic,
		path:  filepath.Join(topic.path, name),
	}

	return ws
}

func (w *Workspace) ShortPath() string {
	return filepath.Join(w.Topic.Name, w.Name)
}

func (w *Workspace) Path() string {
	return w.path
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

func (w Workspaces) SortedByTopic() Workspaces {
	sort.Slice(w, func(i, j int) bool {
		return w[i].Topic.LastModifiedTime().After(w[j].Topic.LastModifiedTime())
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

func (w Workspaces) ByName(s string) Workspaces {
	if s == "" {
		return w
	}

	filtered := Workspaces{}
	for _, workspace := range w {
		if workspace.Name == s {
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

func (w Workspaces) RemoveDuplicates() Workspaces {
	set := map[*Workspace]struct{}{}
	for _, workspace := range w {
		set[workspace] = struct{}{}
	}
	out := make(Workspaces, 0)
	for w := range set {
		out = append(out, w)
	}
	return out
}

// WorkspaceRepository exposes crud on workspaces.
type WorkspaceRepository struct {
	container *Container[Workspace]
}

func newWorkspaceRepository(topics Topics) *WorkspaceRepository {
	w := &WorkspaceRepository{}
	w.load(topics)
	return w
}

func (w *WorkspaceRepository) load(topics Topics) {
	wc := newContainer[Workspace]()
	w.container = wc
	for _, topic := range topics {
		workspaceDirEntries := system.GetDirEntries(topic.path)
		for _, dirEntry := range workspaceDirEntries {
			if !dirEntry.IsDir() {
				continue
			}

			workspace := newWorkspace(topic, dirEntry.Name())
			wc.Add(workspace)
		}
	}
}

func (w *WorkspaceRepository) Save(workspace *Workspace) error {
	// if this workspace doesnt exist, we create a dir
	if !w.container.Contains(workspace) {
		if err := system.CreateDir(workspace.Path()); err != nil {
			return err
		}
	}

	// if the name is not the same as the end of its path it means the name changed
	// and if the workspace topic path is not the same as the dir of the workspace path, it means the topic changed
	// in both cases we rename the dir
	if workspace.Name != filepath.Base(workspace.path) || workspace.Topic.path != filepath.Dir(workspace.path) {
		newPath := filepath.Join(workspace.Topic.path, workspace.Name)
		if err := os.Rename(workspace.path, newPath); err != nil {
			return err
		}
		workspace.path = newPath
	}

	// save it to the container
	w.container.Add(workspace)
	return nil
}

func (w *WorkspaceRepository) Delete(workspace *Workspace) error {
	if err := os.RemoveAll(workspace.Path()); err != nil {
		return err
	}

	w.container.Remove(workspace)
	return nil
}

func (w *WorkspaceRepository) FindByShortPath(shortPath string) *Workspace {
	for _, w := range w.container.All() {
		if w.ShortPath() == shortPath {
			return w
		}
	}
	return nil
}

func (w *WorkspaceRepository) FindByPath(path string) *Workspace {
	for _, w := range w.container.All() {
		if w.Path() == path {
			return w
		}
	}
	return nil
}

func (w *WorkspaceRepository) All() Workspaces {
	return w.container.All()
}

func (w *WorkspaceRepository) Count() int {
	return w.container.Size()
}
