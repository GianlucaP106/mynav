package tmux

type TmuxRepository struct {
	TmuxSessionContainer TmuxSessionContainer
	TmuxCommunicator     *TmuxCommunicator
}

func NewTmuxRepository(tc *TmuxCommunicator) *TmuxRepository {
	tr := &TmuxRepository{
		TmuxCommunicator: tc,
	}
	tr.LoadSessions()
	return tr
}

func (tr *TmuxRepository) LoadSessions() {
	tr.TmuxSessionContainer = tr.TmuxCommunicator.GetSessions()
}

func (tr *TmuxRepository) DeleteSession(ts *TmuxSession) error {
	if err := tr.TmuxCommunicator.KillSession(ts.Name); err != nil {
		return err
	}

	tr.TmuxSessionContainer.Delete(ts)
	return nil
}

func (tr *TmuxRepository) GetSessionByName(name string) *TmuxSession {
	return tr.TmuxSessionContainer.Get(name)
}

func (tr *TmuxRepository) RenameSession(ts *TmuxSession, newName string) error {
	if err := tr.TmuxCommunicator.RenameSession(ts.Name, newName); err != nil {
		return err
	}

	tr.TmuxSessionContainer.Delete(ts)
	ts.Name = newName
	tr.TmuxSessionContainer.Set(ts)

	return nil
}
