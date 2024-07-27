package system

import (
	"mynav/pkg/constants"
	"mynav/pkg/events"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/shirou/gopsutil/net"
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
	ports Ports
}

func NewPortController() *PortController {
	pc := &PortController{
		ports: make(Ports),
	}

	return pc
}

func (pc *PortController) InitPorts() {
	allActivePorts, err := pc.GetRunningPorts()
	if err != nil {
		return
	}

	out := make(Ports)
	for _, port := range allActivePorts {
		pi := GetProcessInfo(port.Pid)
		port.Exe = ""
		if pi != nil {
			port.Exe = pi.Exe
		}
		out.AddPort(port)
	}

	pc.ports = out
}

func (pc *PortController) KillPort(p *Port) error {
	if err := KillProcess(p.Pid); err != nil {
		return err
	}

	events.Emit(constants.PortSyncNeededEventName)
	return nil
}

func (pc *PortController) GetPorts() Ports {
	return pc.ports
}

func (pc *PortController) GetRunningPorts() (PortList, error) {
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
