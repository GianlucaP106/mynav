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

type TmuxSessions []*TmuxSession
