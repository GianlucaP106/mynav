package api

import (
	"path/filepath"
	"sort"
	"strconv"
)

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

func (p *Port) GetPortStr() string {
	return strconv.Itoa(p.Port)
}

type Ports map[int]*Port

func (p Ports) ToList() PortList {
	out := make(PortList, 0)

	for _, p := range p {
		out = append(out, p)
	}

	return out
}

type PortList []*Port

func (p PortList) Len() int { return len(p) }

func (p PortList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p PortList) Less(i, j int) bool {
	return p[i].Port < p[j].Port
}

func (p PortList) Sorted() PortList {
	sort.Sort(p)
	return p
}
