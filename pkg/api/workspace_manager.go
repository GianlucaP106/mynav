package api

import (
	"errors"
	"mynav/pkg/utils"
	"os"
	"path/filepath"
)

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
	wm.loadWorkspaces()
	wm.SyncTmuxSessions()
	return wm
}

func (wm *WorkspaceManager) loadWorkspaces() {
	workspaces := make(Workspaces, 0)
	for _, topic := range wm.Controller.TopicManager.Topics {
		workspaceDirs := utils.GetDirEntries(topic.Path)
		for _, w := range workspaceDirs {
			if !w.IsDir() {
				continue
			}

			workspace := NewWorkspace(w.Name(), topic, wm.WorkspacePath(topic, w.Name()))
			metadata := wm.WorkspaceStore.Workspaces[workspace.ShortPath()]
			if metadata == nil {
				metadata = &WorkspaceMetadata{}
			}
			workspace.Metadata = metadata
			workspaces = append(workspaces, workspace)
		}
	}

	wm.Workspaces = workspaces
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

func (wm *WorkspaceManager) GetWorkspaceByPath(path string) *Workspace {
	for _, workspace := range wm.Workspaces {
		if workspace.Path == path {
			return workspace
		}
	}
	return nil
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
	wm.WorkspaceStore.DeleteWorkspaceMetadata(workspace.ShortPath())

	return nil
}

func (wm *WorkspaceManager) GetOrCreateTmuxSession(workspace *Workspace) (foundExisting bool, sessionName string) {
	m := wm.WorkspaceStore.Workspaces[workspace.ShortPath()]
	if m == nil {
		m = &WorkspaceMetadata{}
	}

	if m.TmuxSession != nil {
		return true, m.TmuxSession.Name
	}

	m.TmuxSession = &utils.TmuxSession{
		Name:       workspace.Path,
		NumWindows: 0,
	}
	wm.WorkspaceStore.SetWorkspaceMetadata(workspace.ShortPath(), m)
	return false, workspace.Path
}

func (wm *WorkspaceManager) SyncTmuxSessions() {
	sessions := utils.GetTmuxSessions()

	for _, metadata := range wm.WorkspaceStore.Workspaces {
		if metadata.TmuxSession != nil && sessions[metadata.TmuxSession.Name] == nil {
			metadata.TmuxSession = nil
		}
	}

	for _, session := range sessions {
		workspace := wm.GetWorkspaceByPath(session.Name)
		if workspace == nil {
			continue
		}
		workspace.Metadata.TmuxSession = &utils.TmuxSession{
			Name:       session.Name,
			NumWindows: session.NumWindows,
		}

		wm.WorkspaceStore.SetWorkspaceMetadata(workspace.ShortPath(), workspace.Metadata)
	}

	wm.WorkspaceStore.Save()
}

func (wm *WorkspaceManager) SetDescription(workspace *Workspace, description string) {
	m := &WorkspaceMetadata{
		Description: description,
	}
	wm.WorkspaceStore.SetWorkspaceMetadata(workspace.ShortPath(), m)
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
