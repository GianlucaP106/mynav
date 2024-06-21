package tmux

type TmuxPane struct {
	Session *TmuxSession
	Pid     int
	Number  int
}

func NewTmuxPane(ts *TmuxSession, pid int, number int) *TmuxPane {
	return &TmuxPane{
		Session: ts,
		Pid:     pid,
		Number:  number,
	}
}
