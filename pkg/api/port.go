package api

import "path/filepath"

type Port struct {
	TmuxSession *TmuxSession
	Exe         string
	Port        int
	Pid         int
}

func NewPort(port int, pid int, tm *TmuxSession, exe string) *Port {
	return &Port{
		Port:        port,
		Pid:         pid,
		TmuxSession: tm,
		Exe:         exe,
	}
}

func (p *Port) GetExeName() string {
	return filepath.Base(p.Exe)
}

type Ports map[int]*Port
