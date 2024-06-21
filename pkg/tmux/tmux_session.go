package tmux

import "mynav/pkg/system"

type TmuxSession struct {
	Ports      system.Ports
	Name       string `json:"name"`
	NumWindows int    `json:"num-windows"`
}

func NewTmuxSession(name string) *TmuxSession {
	session := &TmuxSession{
		Name:       name,
		NumWindows: 0,
		Ports:      make(system.Ports),
	}
	return session
}

type TmuxSessions []*TmuxSession
