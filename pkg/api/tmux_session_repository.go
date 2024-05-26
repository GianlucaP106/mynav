package api

type TmuxSessionRepository struct {
	TmuxSessionContainer TmuxSessionContainer
	TmuxCommunicator     *TmuxCommunicator
}

func NewTmuxSessionRepository() *TmuxSessionRepository {
	tr := &TmuxSessionRepository{
		TmuxCommunicator: NewTmuxCommunicator(),
	}
	tr.LoadSessions()
	return tr
}

func (tr *TmuxSessionRepository) LoadSessions() {
	tr.TmuxSessionContainer = tr.TmuxCommunicator.GetSessions()
}

func (tr *TmuxSessionRepository) DeleteSession(ts *TmuxSession) error {
	if err := tr.TmuxCommunicator.KillSession(ts.Name); err != nil {
		return err
	}

	tr.TmuxSessionContainer.Delete(ts)
	return nil
}

func (tr *TmuxSessionRepository) GetSessionByName(name string) *TmuxSession {
	return tr.TmuxSessionContainer.Get(name)
}

func (tr *TmuxSessionRepository) RenameSession(ts *TmuxSession, newName string) error {
	if err := tr.TmuxCommunicator.RenameSession(ts.Name, newName); err != nil {
		return err
	}

	tr.TmuxSessionContainer.Delete(ts)
	ts.Name = newName
	tr.TmuxSessionContainer.Set(ts)

	return nil
}
