package api

import "mynav/pkg/utils"

type WorkspaceController struct {
	WorkspaceRepository *WorkspaceRepository
	TmuxCommunicator    *TmuxCommunicator
}

func NewWorkspaceController(topics Topics, storePath string) *WorkspaceController {
	wc := &WorkspaceController{}
	wc.TmuxCommunicator = NewTmuxCommunicator()
	wc.WorkspaceRepository = NewWorkspaceRepository(topics, storePath)
	wc.syncTmuxSessions()
	return wc
}

func (wc *WorkspaceController) CreateWorkspace(name string, topic *Topic) (*Workspace, error) {
	workspace := NewWorkspace(name, topic)
	if err := wc.WorkspaceRepository.Save(workspace); err != nil {
		return nil, err
	}

	return workspace, nil
}

func (wc *WorkspaceController) DeleteWorkspace(w *Workspace) error {
	wc.DeleteTmuxSession(w)

	if err := wc.WorkspaceRepository.Delete(w); err != nil {
		return err
	}

	return nil
}

func (wc *WorkspaceController) SetDescription(description string, w *Workspace) {
	w.Metadata.Description = description
	wc.WorkspaceRepository.Save(w)
}

func (wc *WorkspaceController) CreateOrAttachTmuxSession(w *Workspace) []string {
	if w.Metadata.TmuxSession != nil {
		return utils.AttachTmuxSessionCmd(w.Metadata.TmuxSession.Name)
	}

	ts := NewTmuxSession(w.Path)
	w.Metadata.TmuxSession = ts
	wc.WorkspaceRepository.Save(w)
	return utils.NewTmuxSessionCmd(ts.Name, ts.Name)
}

func (wc *WorkspaceController) DeleteTmuxSession(w *Workspace) {
	if w.Metadata.TmuxSession != nil {
		wc.TmuxCommunicator.DeleteSession(w.Metadata.TmuxSession)
	}
	w.Metadata.TmuxSession = nil
	wc.WorkspaceRepository.Save(w)
}

func (wc *WorkspaceController) GetSelectedWorkspace() *Workspace {
	return wc.WorkspaceRepository.GetSelectedWorkspace()
}

func (wm *WorkspaceController) GetTmuxStats() (sessionCount int, windowCount int) {
	sessionCount = 0
	windowCount = 0
	for _, w := range wm.WorkspaceRepository.GetContainer() {
		if w.Metadata.TmuxSession != nil {
			sessionCount++
			windowCount += w.Metadata.TmuxSession.NumWindows
		}
	}
	return
}

func (wc *WorkspaceController) GetWorkspaceCount() int {
	return len(wc.WorkspaceRepository.GetContainer())
}

func (wc *WorkspaceController) GetWorkspacesByTopicCount(t *Topic) int {
	return wc.GetWorkspaces().ByTopic(t).Len()
}

func (wc *WorkspaceController) GetWorkspaces() Workspaces {
	return wc.WorkspaceRepository.GetContainer().ToList()
}

func (wc *WorkspaceController) DeleteWorkspacesByTopic(t *Topic) {
	for _, w := range wc.WorkspaceRepository.WorkspaceContainer.ToList().ByTopic(t) {
		wc.DeleteWorkspace(w)
	}
}

func (wc *WorkspaceController) syncTmuxSessions() {
	sessions := wc.TmuxCommunicator.GetSessions()
	for _, w := range wc.WorkspaceRepository.WorkspaceContainer {
		if w.Metadata.TmuxSession != nil && !sessions.Exists(w.Path) {
			w.Metadata.TmuxSession = nil
			wc.WorkspaceRepository.Save(w)
		}
	}

	for id, session := range sessions {
		if w := wc.WorkspaceRepository.FindByPath(id); w != nil {
			w.Metadata.TmuxSession = session
			wc.WorkspaceRepository.Save(w)
		}
	}
}
