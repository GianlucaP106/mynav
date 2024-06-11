package api

type TmuxController struct {
	TmuxRepository   *TmuxRepository
	TmuxCommunicator *TmuxCommunicator
	PortController   *PortController
}

func NewTmuxController(pc *PortController) *TmuxController {
	tc := NewTmuxCommunicator()
	return &TmuxController{
		TmuxRepository:   NewTmuxRepository(tc),
		TmuxCommunicator: tc,
		PortController:   pc,
	}
}

func (tc *TmuxController) RenameTmuxSession(s *TmuxSession, newName string) error {
	if err := tc.TmuxRepository.RenameSession(s, newName); err != nil {
		return err
	}

	return nil
}

func (tc *TmuxController) GetTmuxSessionCount() int {
	return len(tc.TmuxRepository.TmuxSessionContainer)
}

func (tc *TmuxController) GetTmuxSessions() TmuxSessions {
	return tc.TmuxRepository.TmuxSessionContainer.ToList()
}

func (tc *TmuxController) GetTmuxSessionByWorkspace(w *Workspace) *TmuxSession {
	return tc.TmuxRepository.GetSessionByName(w.Path)
}

func (tc *TmuxController) GetPortsByTmuxSession(session *TmuxSession) PortList {
	out := make(PortList, 0)
	for _, port := range tc.PortController.GetPorts() {
		if port.TmuxSession == nil {
			continue
		}

		if port.TmuxSession == session {
			out = append(out, port)
		}
	}

	return out
}

func (tc *TmuxController) DeleteTmuxSession(s *TmuxSession) error {
	pl := tc.GetPortsByTmuxSession(s)
	refreshPorts := pl.Len() > 0
	if refreshPorts {
		defer tc.PortController.InitPortsAsync()
	}

	if err := tc.TmuxRepository.DeleteSession(s); err != nil {
		return err
	}

	return nil
}

func (tc *TmuxController) GetTmuxSessionPidMap() map[int]*TmuxSession {
	out := map[int]*TmuxSession{}

	for _, session := range tc.GetTmuxSessions() {
		panes := tc.GetTmuxPanesBySession(session)
		for _, pane := range panes {
			out[pane.Pid] = session
		}
	}

	return out
}

func (tc *TmuxController) GetTmuxPanesBySession(ts *TmuxSession) []*TmuxPane {
	return tc.TmuxCommunicator.GetSessionPanes(ts)
}

func (tc *TmuxController) DeleteAllTmuxSessions() error {
	for _, s := range tc.TmuxRepository.TmuxSessionContainer {
		if err := tc.DeleteTmuxSession(s); err != nil {
			return err
		}
	}
	return nil
}

func (tc *TmuxController) GetTmuxStats() (sessionCount int, windowCount int) {
	sessionCount = 0
	windowCount = 0

	for _, s := range tc.TmuxRepository.TmuxSessionContainer {
		sessionCount++
		windowCount += s.NumWindows
	}

	return
}
