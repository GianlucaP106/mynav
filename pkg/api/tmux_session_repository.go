package api

type TmuxSessionRepository struct {
	TmuxSessionContainer TmuxSessionContainer
	*TmuxCommunicator
}

func NewTmuxSessionRepository() *TmuxSessionRepository {
	tr := &TmuxSessionRepository{
		TmuxCommunicator: NewTmuxCommunicator(),
	}
	tr.LoadSessions()
	return tr
}

func (tr *TmuxSessionRepository) LoadSessions() {
	tr.TmuxSessionContainer = tr.GetSessions()
}
