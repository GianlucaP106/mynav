package system

import (
	"mynav/pkg/events"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
)

type Port struct {
	Exe  string
	Port int
	Pid  int
}

func NewPort(port int, pid int, exe string) *Port {
	return &Port{
		Port: port,
		Pid:  pid,
		Exe:  exe,
	}
}

func (p *Port) GetExeName() string {
	return filepath.Base(p.Exe)
}

func (p *Port) GetPortStr() string {
	return strconv.Itoa(p.Port)
}

type Ports map[int]*Port

func (p Ports) AddPort(port *Port) {
	p[port.Port] = port
}

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

type PortController struct {
	processController *ProcessController
}

func NewPortController(p *ProcessController) *PortController {
	pc := &PortController{
		processController: p,
	}
	return pc
}

func (pc *PortController) GetPorts() Ports {
	allActivePorts, err := pc.getRunningPorts()
	if err != nil {
		return nil
	}

	out := make(Ports)
	for _, port := range allActivePorts {
		proc, err := process.NewProcess(int32(port.Pid))
		if err == nil {
			exe, _ := proc.Exe()
			port.Exe = exe
		}

		out.AddPort(port)
	}

	return out
}

func (pc *PortController) KillPort(p *Port) error {
	if err := pc.processController.KillProcess(p.Pid); err != nil {
		return err
	}

	events.Emit(events.PortChangeEvent)
	return nil
}

func (pc *PortController) getRunningPorts() (PortList, error) {
	connections, err := net.Connections("inet")
	if err != nil {
		return nil, err
	}

	out := make(PortList, 0)
	for _, cs := range connections {
		if cs.Status == "LISTEN" {
			port := NewPort(int(cs.Laddr.Port), int(cs.Pid), "")
			out = append(out, port)
		}
	}

	return out, nil
}
