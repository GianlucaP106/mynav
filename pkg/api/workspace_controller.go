package api

import "mynav/pkg/utils"

type WorkspaceController struct {
	WorkspaceRepository   *WorkspaceRepository
	TmuxSessionController *TmuxSessionController
}

func NewWorkspaceController(topics Topics, storePath string, tr *TmuxSessionController) *WorkspaceController {
	wc := &WorkspaceController{}
	wc.TmuxSessionController = tr
	wc.WorkspaceRepository = NewWorkspaceRepository(topics, storePath)
	return wc
}

func (wc *WorkspaceController) CreateWorkspace(name string, topic *Topic) (*Workspace, error) {
	workspace := NewWorkspace(name, topic)

	if err := wc.WorkspaceRepository.SetSelectedWorkspace(workspace); err != nil {
		return nil, err
	}

	return workspace, nil
}

func (wc *WorkspaceController) DeleteWorkspace(w *Workspace) error {
	wc.DeleteWorkspaceTmuxSession(w)

	if err := wc.WorkspaceRepository.Delete(w); err != nil {
		return err
	}

	return nil
}

func (wc *WorkspaceController) RenameWorkspace(w *Workspace, newName string) error {
	s := wc.TmuxSessionController.GetTmuxSessionByWorkspace(w)

	if err := wc.WorkspaceRepository.Rename(w, newName); err != nil {
		return err
	}

	if s != nil {
		wc.TmuxSessionController.RenameTmuxSession(s, w.Path)
	}

	return nil
}

func (wc *WorkspaceController) SetDescription(description string, w *Workspace) {
	w.Metadata.Description = description
	wc.WorkspaceRepository.SetSelectedWorkspace(w)
}

func (wc *WorkspaceController) GetCreateOrAttachTmuxSessionCmd(w *Workspace) []string {
	if ts := wc.TmuxSessionController.GetTmuxSessionByWorkspace(w); ts != nil {
		wc.WorkspaceRepository.SetSelectedWorkspace(w)
		return utils.AttachTmuxSessionCmd(ts.Name)
	}

	wc.WorkspaceRepository.SetSelectedWorkspace(w)
	return utils.NewTmuxSessionCmd(w.Path, w.Path)
}

func (wc *WorkspaceController) DeleteWorkspaceTmuxSession(w *Workspace) {
	if ts := wc.TmuxSessionController.GetTmuxSessionByWorkspace(w); ts != nil {
		wc.TmuxSessionController.DeleteTmuxSession(ts)
	}
}

func (wc *WorkspaceController) GetSelectedWorkspace() *Workspace {
	return wc.WorkspaceRepository.GetSelectedWorkspace()
}

func (wc *WorkspaceController) SetSelectedWorkspace(w *Workspace) {
	wc.WorkspaceRepository.SetSelectedWorkspace(w)
}

func (wc *WorkspaceController) GetWorkspaceByTmuxSession(s *TmuxSession) *Workspace {
	for _, w := range wc.GetWorkspaces() {
		if ts := wc.TmuxSessionController.GetTmuxSessionByWorkspace(w); ts != nil && ts.Name == s.Name {
			return w
		}
	}

	return nil
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

func (wc *WorkspaceController) CloneRepo(repoUrl string, w *Workspace) error {
	err := w.CloneRepo(repoUrl)
	if err != nil {
		return err
	}
	w.GitRemote = &repoUrl
	wc.WorkspaceRepository.SetSelectedWorkspace(w)
	return nil
}
