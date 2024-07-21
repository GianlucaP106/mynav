package tmux

import "mynav/pkg/persistence"

type TmuxRepository struct {
	TmuxSessionContainer *persistence.Container[TmuxSession]
	TmuxCommunicator     *TmuxCommunicator
}

func NewTmuxRepository(tc *TmuxCommunicator) *TmuxRepository {
	tr := &TmuxRepository{
		TmuxCommunicator: tc,
	}
	tr.LoadSessions()
	return tr
}

func (tr *TmuxRepository) GetSessions() TmuxSessions {
	return tr.TmuxSessionContainer.All()
}

func (tr *TmuxRepository) LoadSessions() {
	tr.TmuxSessionContainer = persistence.NewContainer[TmuxSession]()
	sessions := tr.TmuxCommunicator.GetSessions()
	for _, s := range sessions {
		tr.TmuxSessionContainer.Set(s.Name, s)
	}
}

func (tr *TmuxRepository) DeleteSession(ts *TmuxSession) error {
	if err := tr.TmuxCommunicator.KillSession(ts.Name); err != nil {
		return err
	}

	tr.TmuxSessionContainer.Delete(ts.Name)
	return nil
}

func (tr *TmuxRepository) GetSessionByName(name string) *TmuxSession {
	return tr.TmuxSessionContainer.Get(name)
}

func (tr *TmuxRepository) RenameSession(ts *TmuxSession, newName string) error {
	if err := tr.TmuxCommunicator.RenameSession(ts.Name, newName); err != nil {
		return err
	}

	tr.TmuxSessionContainer.Delete(ts.Name)
	ts.Name = newName
	tr.TmuxSessionContainer.Set(ts.Name, ts)

	return nil
}

func (tr *TmuxRepository) SessionCount() int {
	return tr.TmuxSessionContainer.Size()
}
