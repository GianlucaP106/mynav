package api

type TmuxSessionController struct {
	TmuxSessionRepository *TmuxSessionRepository
}

func NewTmuxSessionController() *TmuxSessionController {
	return &TmuxSessionController{
		TmuxSessionRepository: NewTmuxSessionRepository(),
	}
}

func (tc *TmuxSessionController) RenameTmuxSession(s *TmuxSession, newName string) error {
	if err := tc.TmuxSessionRepository.RenameSession(s, newName); err != nil {
		return err
	}

	return nil
}

func (tc *TmuxSessionController) GetTmuxSessions() TmuxSessions {
	return tc.TmuxSessionRepository.TmuxSessionContainer.ToList()
}

func (tc *TmuxSessionController) GetTmuxSessionByWorkspace(w *Workspace) *TmuxSession {
	return tc.TmuxSessionRepository.GetSessionByName(w.Path)
}

func (tc *TmuxSessionController) DeleteTmuxSession(s *TmuxSession) error {
	if err := tc.TmuxSessionRepository.DeleteSession(s); err != nil {
		return err
	}

	return nil
}

func (tc *TmuxSessionController) GetTmuxStats() (sessionCount int, windowCount int) {
	sessionCount = 0
	windowCount = 0

	for _, s := range tc.TmuxSessionRepository.TmuxSessionContainer {
		sessionCount++
		windowCount += s.NumWindows
	}

	return
}
