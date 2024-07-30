package core

import (
	"errors"
	"mynav/pkg/constants"
	"mynav/pkg/events"
	"mynav/pkg/system"
	"mynav/pkg/tmux"
	"strings"
)

type WorkspaceController struct {
	WorkspaceRepository *WorkspaceRepository
	TmuxController      *tmux.TmuxController
}

func NewWorkspaceController(topics Topics, storePath string, tr *tmux.TmuxController) *WorkspaceController {
	wc := &WorkspaceController{}
	wc.TmuxController = tr
	wc.WorkspaceRepository = NewWorkspaceRepository(topics, storePath)
	return wc
}

func (wc *WorkspaceController) CreateWorkspace(name string, topic *Topic) (*Workspace, error) {
	if err := wc.periodValidation(name); err != nil {
		return nil, err
	}

	workspace := NewWorkspace(name, topic)

	if err := wc.WorkspaceRepository.SetSelectedWorkspace(workspace); err != nil {
		return nil, err
	}

	events.Emit(constants.TopicChangeEventName)

	return workspace, nil
}

func (wc *WorkspaceController) DeleteWorkspace(w *Workspace) error {
	wc.DeleteWorkspaceTmuxSession(w)

	if err := wc.WorkspaceRepository.Delete(w); err != nil {
		return err
	}

	events.Emit(constants.TopicChangeEventName)
	return nil
}

func (wc *WorkspaceController) MoveWorkspace(w *Workspace, newTopic *Topic) error {
	if w.Topic.Name == newTopic.Name {
		return errors.New("workspace is already in this topic")
	}

	if wc.GetWorkspaces().FilterByTopic(newTopic).GetWorkspaceByName(w.Name) != nil {
		return errors.New("workspace with this name already exists")
	}

	s := wc.TmuxController.GetTmuxSessionByName(w.Path)

	if err := wc.WorkspaceRepository.Move(w, newTopic); err != nil {
		return err
	}

	if s != nil {
		wc.TmuxController.RenameTmuxSession(s, w.Path)
	}

	events.Emit(constants.TopicChangeEventName)
	events.Emit(constants.PortSyncNeededEventName)
	return nil
}

func (wc *WorkspaceController) RenameWorkspace(w *Workspace, newName string) error {
	if err := wc.periodValidation(newName); err != nil {
		return err
	}

	s := wc.TmuxController.GetTmuxSessionByName(w.Path)

	if err := wc.WorkspaceRepository.Rename(w, newName); err != nil {
		return err
	}

	if s != nil {
		wc.TmuxController.RenameTmuxSession(s, w.Path)
	}

	events.Emit(constants.WorkspaceChangeEventName)
	events.Emit(constants.PortSyncNeededEventName)
	return nil
}

func (wc *WorkspaceController) GetWorkspaceTmuxSessionCount() int {
	out := 0
	for _, w := range wc.WorkspaceRepository.Container.All() {
		if wc.TmuxController.GetTmuxSessionByName(w.Path) != nil {
			out++
		}
	}
	return out
}

func (wc *WorkspaceController) SetDescription(description string, w *Workspace) {
	w.Metadata.Description = description
	wc.WorkspaceRepository.SetSelectedWorkspace(w)
	events.Emit(constants.WorkspaceChangeEventName)
}

func (wc *WorkspaceController) OpenNeovimInWorkspace(w *Workspace) error {
	wc.WorkspaceRepository.SetSelectedWorkspace(w)
	return system.CommandWithRedirect("nvim", w.Path).Run()
}

func (wc *WorkspaceController) OpenTerminalInWorkspace(w *Workspace) error {
	cmd, err := system.GetOpenTerminalCmd(w.Path)
	if err != nil {
		return err
	}

	wc.WorkspaceRepository.SetSelectedWorkspace(w)
	return system.CommandWithRedirect(cmd...).Run()
}

func (wc *WorkspaceController) CreateOrAttachTmuxSession(w *Workspace) error {
	if ts := wc.TmuxController.GetTmuxSessionByName(w.Path); ts != nil {
		wc.WorkspaceRepository.SetSelectedWorkspace(w)
		if err := wc.TmuxController.AttachTmuxSession(ts); err != nil {
			return err
		}
		events.Emit(constants.WorkspaceChangeEventName)
		return nil
	}

	wc.WorkspaceRepository.SetSelectedWorkspace(w)
	if err := wc.TmuxController.CreateAndAttachTmuxSession(w.Path, w.Path); err != nil {
		return err
	}

	events.Emit(constants.WorkspaceChangeEventName)
	return nil
}

func (wc *WorkspaceController) DeleteWorkspaceTmuxSession(w *Workspace) {
	if ts := wc.TmuxController.GetTmuxSessionByName(w.Path); ts != nil {
		wc.TmuxController.DeleteTmuxSession(ts)
		events.Emit(constants.WorkspaceChangeEventName)
	}
}

func (wc *WorkspaceController) DeleteAllWorkspaceTmuxSessions() error {
	for _, w := range wc.GetWorkspaces() {
		if s := wc.TmuxController.GetTmuxSessionByName(w.Path); s != nil {
			if err := wc.TmuxController.DeleteTmuxSession(s); err != nil {
				return err
			}
		}
	}

	events.Emit(constants.TmuxSessionChangeEventName)
	return nil
}

func (wc *WorkspaceController) GetSelectedWorkspace() *Workspace {
	return wc.WorkspaceRepository.GetSelectedWorkspace()
}

func (wc *WorkspaceController) SetSelectedWorkspace(w *Workspace) {
	wc.WorkspaceRepository.SetSelectedWorkspace(w)
}

func (wc *WorkspaceController) GetWorkspaceByTmuxSession(s *tmux.TmuxSession) *Workspace {
	for _, w := range wc.GetWorkspaces() {
		if ts := wc.TmuxController.GetTmuxSessionByName(w.Path); ts != nil && ts.Name == s.Name {
			return w
		}
	}

	return nil
}

func (wc *WorkspaceController) GetWorkspaceCount() int {
	return wc.WorkspaceRepository.Container.Size()
}

func (wc *WorkspaceController) GetWorkspacesByTopicCount(t *Topic) int {
	return wc.GetWorkspaces().FilterByTopic(t).Len()
}

func (wc *WorkspaceController) GetWorkspaces() Workspaces {
	return wc.WorkspaceRepository.Container.All()
}

func (wc *WorkspaceController) DeleteWorkspacesByTopic(t *Topic) error {
	var workspaces Workspaces = wc.WorkspaceRepository.Container.All()
	for _, w := range workspaces.FilterByTopic(t) {
		ts := wc.TmuxController.GetTmuxSessionByName(w.Path)
		if ts != nil {
			// TODO: deleteManytmuxSessions
			wc.TmuxController.TmuxRepository.DeleteSession(ts)
		}

		if err := wc.WorkspaceRepository.Delete(w); err != nil {
			return err
		}
	}

	events.Emit(constants.PortSyncNeededEventName)
	events.Emit(constants.TmuxSessionChangeEventName)
	events.Emit(constants.WorkspaceChangeEventName)
	return nil
}

func (wc *WorkspaceController) GetPortsByWorkspace(w *Workspace) system.PortList {
	out := make(system.PortList, 0)

	session := wc.TmuxController.GetTmuxSessionByName(w.Path)
	if session == nil {
		return out
	}

	return session.Ports.ToList().Sorted()
}

func (wc *WorkspaceController) CloneRepo(repoUrl string, w *Workspace) error {
	err := w.CloneRepo(repoUrl)
	if err != nil {
		return err
	}
	w.GitRemote = &repoUrl
	wc.WorkspaceRepository.SetSelectedWorkspace(w)
	events.Emit(constants.WorkspaceChangeEventName)
	return nil
}

func (wc *WorkspaceController) periodValidation(name string) error {
	if strings.ContainsRune(name, '.') {
		return errors.New("workspace name cannot contain '.'")
	}
	return nil
}
