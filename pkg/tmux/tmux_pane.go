package tmux

import "github.com/shirou/gopsutil/process"

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

type TmuxPaneProcess struct {
	Process *process.Process
	Pane    *TmuxPane
}
