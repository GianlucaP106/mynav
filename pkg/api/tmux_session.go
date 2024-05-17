package api

type TmuxSession struct {
	Name       string `json:"name"`
	NumWindows int    `json:"num-windows"`
}

func NewTmuxSession(name string) *TmuxSession {
	session := &TmuxSession{
		Name:       name,
		NumWindows: 0,
	}
	return session
}

type TmuxSessions map[string]*TmuxSession

func (ws TmuxSessions) Get(id string) *TmuxSession {
	return ws[id]
}

func (ws TmuxSessions) Exists(id string) bool {
	return ws[id] != nil
}
