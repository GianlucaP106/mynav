package core

import (
	"errors"
	"mynav/pkg/system"
	"strings"

	"github.com/GianlucaP106/gotmux/gotmux"
)

type WorkspaceController struct {
	workspaceRepository *WorkspaceRepository
	tmuxController      *TmuxController
	globalConfiguration *GlobalConfiguration
	localConfiguration  *LocalConfiguration
}

func NewWorkspaceController(topics Topics, tr *TmuxController, gc *GlobalConfiguration, lc *LocalConfiguration) *WorkspaceController {
	wc := &WorkspaceController{}
	wc.tmuxController = tr
	wc.globalConfiguration = gc
	wc.localConfiguration = lc
	wc.workspaceRepository = NewWorkspaceRepository(topics, lc.GetWorkspaceStorePath())
	return wc
}

func (wc *WorkspaceController) CreateWorkspace(name string, topic *Topic) (*Workspace, error) {
	if err := wc.periodValidation(name); err != nil {
		return nil, err
	}

	workspace := NewWorkspace(name, topic)

	if err := wc.workspaceRepository.SetSelectedWorkspace(workspace); err != nil {
		return nil, err
	}

	return workspace, nil
}

func (wc *WorkspaceController) DeleteWorkspace(w *Workspace) error {
	wc.DeleteWorkspaceTmuxSession(w)

	if err := wc.workspaceRepository.Delete(w); err != nil {
		return err
	}

	return nil
}

func (wc *WorkspaceController) MoveWorkspace(w *Workspace, newTopic *Topic) error {
	if w.Topic.Name == newTopic.Name {
		return errors.New("workspace is already in this topic")
	}

	if wc.GetWorkspaces().FilterByTopic(newTopic).GetWorkspaceByName(w.Name) != nil {
		return errors.New("workspace with this name already exists")
	}

	s := wc.tmuxController.GetTmuxSessionByName(w.Path)

	if err := wc.workspaceRepository.Move(w, newTopic); err != nil {
		return err
	}

	if s != nil {
		wc.tmuxController.RenameTmuxSession(s, w.Path)
	}

	return nil
}

func (wc *WorkspaceController) RenameWorkspace(w *Workspace, newName string) error {
	if err := wc.periodValidation(newName); err != nil {
		return err
	}

	s := wc.tmuxController.GetTmuxSessionByName(w.Path)

	if err := wc.workspaceRepository.Rename(w, newName); err != nil {
		return err
	}

	if s != nil {
		wc.tmuxController.RenameTmuxSession(s, w.Path)
	}

	return nil
}

func (wc *WorkspaceController) GetWorkspaceTmuxSessionCount() int {
	out := 0
	for _, w := range wc.workspaceRepository.container.All() {
		if wc.tmuxController.GetTmuxSessionByName(w.Path) != nil {
			out++
		}
	}
	return out
}

func (wc *WorkspaceController) SetDescription(description string, w *Workspace) {
	w.Metadata.Description = description
	wc.workspaceRepository.SetSelectedWorkspace(w)
}

func (wc *WorkspaceController) OpenNeovimInWorkspace(w *Workspace) error {
	wc.workspaceRepository.SetSelectedWorkspace(w)
	return system.CommandWithRedirect("nvim", w.Path).Run()
}

func (wc *WorkspaceController) OpenTerminalInWorkspace(w *Workspace) error {
	cmd, err := system.GetOpenTerminalCmd()
	if err != nil {
		return err
	}

	cmd = append(cmd, w.Path)

	wc.workspaceRepository.SetSelectedWorkspace(w)
	return system.CommandWithRedirect(cmd...).Run()
}

func (wc *WorkspaceController) CreateOrAttachTmuxSession(w *Workspace) error {
	if ts := wc.tmuxController.GetTmuxSessionByName(w.Path); ts != nil {
		wc.workspaceRepository.SetSelectedWorkspace(w)
		if err := wc.tmuxController.AttachTmuxSession(ts); err != nil {
			return err
		}

		return nil
	}

	wc.workspaceRepository.SetSelectedWorkspace(w)
	if err := wc.tmuxController.CreateAndAttachTmuxSession(w.Path, w.Path); err != nil {
		return err
	}

	return nil
}

func (wc *WorkspaceController) OpenWorkspace(workspace *Workspace) error {
	cmd := wc.globalConfiguration.GetCustomWorkspaceOpenerCmd()
	if len(cmd) > 0 {
		cmd = append(cmd, workspace.Path)
		c := system.CommandWithRedirect(cmd...)
		err := c.Run()
		if err != nil {
			return err
		}
	} else if !IsParentAppInstance() {
		err := wc.OpenNeovimInWorkspace(workspace)
		if err != nil {
			return err
		}
	} else {
		err := wc.CreateOrAttachTmuxSession(workspace)
		if err != nil {
			return err
		}
	}

	return nil
}

func (wc *WorkspaceController) DeleteWorkspaceTmuxSession(w *Workspace) {
	if ts := wc.tmuxController.GetTmuxSessionByName(w.Path); ts != nil {
		wc.tmuxController.DeleteTmuxSession(ts)
	}
}

func (wc *WorkspaceController) DeleteAllWorkspaceTmuxSessions() error {
	for _, w := range wc.GetWorkspaces() {
		if s := wc.tmuxController.GetTmuxSessionByName(w.Path); s != nil {
			if err := wc.tmuxController.DeleteTmuxSession(s); err != nil {
				return err
			}
		}
	}

	return nil
}

func (wc *WorkspaceController) GetSelectedWorkspace() *Workspace {
	return wc.workspaceRepository.GetSelectedWorkspace()
}

func (wc *WorkspaceController) SetSelectedWorkspace(w *Workspace) {
	wc.workspaceRepository.SetSelectedWorkspace(w)
}

func (wc *WorkspaceController) GetWorkspaceByTmuxSession(s *gotmux.Session) *Workspace {
	for _, w := range wc.GetWorkspaces() {
		if w.Path == s.Name {
			return w
		}
	}

	return nil
}

func (wc *WorkspaceController) GetWorkspaceCount() int {
	return wc.workspaceRepository.container.Size()
}

func (wc *WorkspaceController) GetWorkspacesByTopicCount(t *Topic) int {
	return wc.GetWorkspaces().FilterByTopic(t).Len()
}

func (wc *WorkspaceController) GetWorkspaces() Workspaces {
	return wc.workspaceRepository.container.All()
}

func (wc *WorkspaceController) GetWorkspaceByShortPath(s string) *Workspace {
	return wc.workspaceRepository.Find(s)
}

func (wc *WorkspaceController) DeleteWorkspacesByTopic(t *Topic) error {
	var workspaces Workspaces = wc.workspaceRepository.container.All()
	for _, w := range workspaces.FilterByTopic(t) {
		ts := wc.tmuxController.GetTmuxSessionByName(w.Path)
		if ts != nil {
			ts.Kill()
		}

		if err := wc.workspaceRepository.Delete(w); err != nil {
			return err
		}
	}

	return nil
}

func (wc *WorkspaceController) CloneRepo(repoUrl string, w *Workspace) error {
	err := w.CloneRepo(repoUrl)
	if err != nil {
		return err
	}

	w.GitRemote = &repoUrl
	wc.workspaceRepository.SetSelectedWorkspace(w)
	return nil
}

func (wc *WorkspaceController) periodValidation(name string) error {
	if strings.ContainsRune(name, '.') {
		return errors.New("workspace name cannot contain '.'")
	}
	return nil
}
