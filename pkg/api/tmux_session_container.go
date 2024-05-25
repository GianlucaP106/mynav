package api

type TmuxSessionContainer map[string]*TmuxSession

func NewTmuxSessionContainer() TmuxSessionContainer {
	return make(TmuxSessionContainer)
}

func (tc TmuxSessionContainer) Set(t *TmuxSession) {
	tc[t.Name] = t
}

func (ws TmuxSessionContainer) Get(id string) *TmuxSession {
	return ws[id]
}

func (ws TmuxSessionContainer) Exists(id string) bool {
	return ws[id] != nil
}
