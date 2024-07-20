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
	PortController      *system.PortController
}

func NewWorkspaceController(topics Topics, storePath string, tr *tmux.TmuxController, pc *system.PortController) *WorkspaceController {
	wc := &WorkspaceController{}
	wc.TmuxController = tr
	wc.PortController = pc
	wc.WorkspaceRepository = NewWorkspaceRepository(topics, storePath)
	return wc
}

func (wc *WorkspaceController) PeriodValidation(name string) error {
	if strings.ContainsRune(name, '.') {
		return errors.New("workspace name cannot contain '.'")
	}
	return nil
}

func (wc *WorkspaceController) CreateWorkspace(name string, topic *Topic) (*Workspace, error) {
	if err := wc.PeriodValidation(name); err != nil {
		return nil, err
	}

	workspace := NewWorkspace(name, topic)

	if err := wc.WorkspaceRepository.SetSelectedWorkspace(workspace); err != nil {
		return nil, err
	}

	events.EmitEvent(constants.WorkspaceChangeEventName)
	events.EmitEvent(constants.TopicChangeEventName)

	return workspace, nil
}

func (wc *WorkspaceController) DeleteWorkspace(w *Workspace) error {
	pl := wc.GetPortsByWorkspace(w)
	refreshPorts := pl.Len() > 0
	if refreshPorts {
		defer events.EmitEvent(constants.PortChangeEventName)
	}

	wc.DeleteWorkspaceTmuxSession(w)

	if err := wc.WorkspaceRepository.Delete(w); err != nil {
		return err
	}

	events.EmitEvent(constants.WorkspaceChangeEventName)
	events.EmitEvent(constants.TopicChangeEventName)
	return nil
}

func (wc *WorkspaceController) RenameWorkspace(w *Workspace, newName string) error {
	if err := wc.PeriodValidation(newName); err != nil {
		return err
	}

	s := wc.TmuxController.GetTmuxSessionByName(w.Path)

	if err := wc.WorkspaceRepository.Rename(w, newName); err != nil {
		return err
	}

	if s != nil {
		wc.TmuxController.RenameTmuxSession(s, w.Path)
	}

	events.EmitEvent(constants.WorkspaceChangeEventName)
	return nil
}

func (wc *WorkspaceController) GetWorkspaceTmuxSessionCount() int {
	out := 0
	for _, w := range wc.WorkspaceRepository.Container {
		if wc.TmuxController.GetTmuxSessionByName(w.Path) != nil {
			out++
		}
	}
	return out
}

func (wc *WorkspaceController) SetDescription(description string, w *Workspace) {
	w.Metadata.Description = description
	wc.WorkspaceRepository.SetSelectedWorkspace(w)
	events.EmitEvent(constants.WorkspaceChangeEventName)
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
		events.EmitEvent(constants.WorkspaceChangeEventName)
		return nil
	}

	wc.WorkspaceRepository.SetSelectedWorkspace(w)
	if err := wc.TmuxController.CreateAndAttachTmuxSession(w.Path, w.Path); err != nil {
		return err
	}

	events.EmitEvent(constants.WorkspaceChangeEventName)
	return nil
}

func (wc *WorkspaceController) DeleteWorkspaceTmuxSession(w *Workspace) {
	if ts := wc.TmuxController.GetTmuxSessionByName(w.Path); ts != nil {
		wc.TmuxController.DeleteTmuxSession(ts)
		events.EmitEvent(constants.TmuxSessionChangeEventName)
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

	events.EmitEvent(constants.TmuxSessionChangeEventName)
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
	return len(wc.WorkspaceRepository.Container)
}

func (wc *WorkspaceController) GetWorkspacesByTopicCount(t *Topic) int {
	return wc.GetWorkspaces().ByTopic(t).Len()
}

func (wc *WorkspaceController) GetWorkspaces() Workspaces {
	return wc.WorkspaceRepository.Container.ToList()
}

func (wc *WorkspaceController) DeleteWorkspacesByTopic(t *Topic) error {
	for _, w := range wc.WorkspaceRepository.Container.ToList().ByTopic(t) {
		if err := wc.WorkspaceRepository.Delete(w); err != nil {
			return err
		}
	}

	events.EmitEvent(constants.WorkspaceChangeEventName)
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
	events.EmitEvent(constants.WorkspaceChangeEventName)
	return nil
}
