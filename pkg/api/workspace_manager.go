package api

import (
	"errors"
	"mynav/pkg/utils"
	"os"
	"path/filepath"
)

//////// WorkspaceManager /////////////

type WorkspaceManager struct {
	Controller     *Controller
	WorkspaceStore *WorkspaceStore
	Workspaces     Workspaces
}

func newWorkspaceManager(c *Controller) *WorkspaceManager {
	wm := &WorkspaceManager{
		Controller: c,
	}

	wm.WorkspaceStore = newWorkspaceStore(filepath.Join(c.Configuration.GetConfigPath(), "workspaces.json"))

	workspaces := make(Workspaces, 0)
	for _, topic := range c.TopicManager.Topics {
		workspaceDirs := utils.GetDirEntries(topic.Path)
		for _, w := range workspaceDirs {
			if !w.IsDir() {
				continue
			}

			workspace := NewWorkspace(w.Name(), topic, wm.WorkspacePath(topic, w.Name()))
			metadata := wm.WorkspaceStore.Workspaces[workspace.GetShortPath()]
			workspace.Metadata = metadata
			workspaces = append(workspaces, workspace)
		}
	}

	wm.Workspaces = workspaces

	return wm
}

func (wm *WorkspaceManager) WorkspacePath(topic *Topic, name string) string {
	return filepath.Join(filepath.Join(wm.Controller.Configuration.path, topic.Name), name)
}

func (wm *WorkspaceManager) detectGitRemote(w *Workspace) {
	gitPath := filepath.Join(w.Path, ".git")
	if _, err := filepath.Abs(gitPath); err != nil {
		return
	}

	gitRemote, _ := utils.GitRemote(gitPath)
	w.GitRemote = &gitRemote
}

func (ws *WorkspaceManager) GetGitRemote(w *Workspace) string {
	if w.GitRemote == nil {
		ws.detectGitRemote(w)
	}
	return *(w.GitRemote)
}

func (wm *WorkspaceManager) CreateWorkspace(name string, repoUrl string, topic *Topic) (*Workspace, error) {
	newWorkspacePath := filepath.Join(wm.Controller.Configuration.path, topic.Name, name)
	if repoUrl != "" {
		if err := utils.GitClone(repoUrl, newWorkspacePath); err != nil {
			return nil, errors.New("failed to clone repository")
		}
	} else {
		if err := os.Mkdir(newWorkspacePath, 0755); err != nil {
			return nil, err
		}
	}

	workspace := NewWorkspace(name, topic, wm.WorkspacePath(topic, name))
	wm.Workspaces = append(wm.Workspaces, workspace)

	return workspace, nil
}

func (wm *WorkspaceManager) DeleteWorkspace(workspace *Workspace) error {
	if err := os.RemoveAll(workspace.Path); err != nil {
		return err
	}

	idx := 0
	for i, t := range wm.Workspaces {
		if t == workspace {
			idx = i
		}
	}

	wm.Workspaces = append(wm.Workspaces[:idx], wm.Workspaces[idx+1:]...)
	wm.WorkspaceStore.DeleteWorkspaceMetadata(workspace.GetShortPath())

	return nil
}

func (wm *WorkspaceManager) SetDescription(workspace *Workspace, description string) {
	m := &WorkspaceMetadata{
		Description: description,
	}
	wm.WorkspaceStore.SetWorkspaceMetadata(workspace.GetShortPath(), m)
	workspace.Metadata = m
}

type WorkspaceStore struct {
	Workspaces map[string]*WorkspaceMetadata `json:"workspaces"`
	Store      string                        `json:"-"`
}

func newWorkspaceStore(store string) *WorkspaceStore {
	w := utils.Load[WorkspaceStore](store)
	if w != nil {
		w.Store = store
		return w
	}

	w = &WorkspaceStore{
		Workspaces: map[string]*WorkspaceMetadata{},
		Store:      store,
	}
	w.Save()

	return w
}

func (ws *WorkspaceStore) SetWorkspaceMetadata(id string, m *WorkspaceMetadata) {
	ws.Workspaces[id] = m
	ws.Save()
}

func (ws *WorkspaceStore) DeleteWorkspaceMetadata(ids ...string) {
	for _, id := range ids {
		delete(ws.Workspaces, id)
	}
	ws.Save()
}

func (ws *WorkspaceStore) Save() {
	utils.Save(ws, ws.Store)
}
