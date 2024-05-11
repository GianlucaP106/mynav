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
	Metadata  *WorkspaceMetadata
	Name      string
	Path      string
}

func NewWorkspace(name string, topic *Topic) *Workspace {
	shortWsPath := filepath.Join(topic.Name, name)
	wsPath := filepath.Join(topic.Filesystem.path, shortWsPath)

	ws := &Workspace{
		Name:  name,
		Topic: topic,
		Path:  wsPath,
	}

	// TODO: https://github.com/GianlucaP106/mynav/issues/34
	store := LoadMetadataStore(ws.GetWorkspaceStorePath())
	for id, w := range store.Workspaces {
		if id == shortWsPath {
			ws.Metadata = w
			break
		}
	}

	return ws
}

func (ws *Workspace) GetGitRemote() string {
	if ws.gitRemote == nil {
		ws.DetectGitRemote()
	}
	return *(ws.gitRemote)
}

func (ws *Workspace) DetectGitRemote() {
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

func (w *Workspace) GetLastModifiedTime() time.Time {
	time, _ := utils.GetLastModifiedTime(w.Path)
	return time
}

func (w *Workspace) GetLastModifiedTimeFormatted() string {
	time := w.GetLastModifiedTime().Format(w.Topic.Filesystem.TimeFormat())
	return time
}

func (w *Workspace) GetDescription() string {
	if w.Metadata == nil {
		return ""
	}
	return w.Metadata.Description
}

type WorkspaceStore struct {
	Workspaces map[string]*WorkspaceMetadata `json:"workspaces"`
}

type WorkspaceMetadata struct {
	Description string `json:"description"`
}

func (w *Workspace) GetWorkspaceStorePath() string {
	return filepath.Join(w.Topic.Filesystem.GetConfigPath(), "workspaces.json")
}

// TODO: https://github.com/GianlucaP106/mynav/issues/34
func (w *Workspace) SaveDescription(description string) {
	if w.Metadata == nil {
		w.Metadata = &WorkspaceMetadata{}
	}
	w.Metadata.Description = description
	id := filepath.Join(w.Topic.Name, w.Name)

	storePath := w.GetWorkspaceStorePath()
	store := LoadMetadataStore(storePath)
	store.Workspaces[id] = w.Metadata
	SaveMetadataStore(store, storePath)
}

func SaveMetadataStore(data *WorkspaceStore, store string) {
	utils.Save(data, store)
}

func LoadMetadataStore(store string) *WorkspaceStore {
	s := utils.Load[WorkspaceStore](store)
	if s == nil {
		s = &WorkspaceStore{
			Workspaces: map[string]*WorkspaceMetadata{},
		}
	}
	return s
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
