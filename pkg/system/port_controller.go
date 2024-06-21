package system

import "github.com/shirou/gopsutil/net"

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
		port.Exe = pi.Exe
		out.AddPort(port)
	}

	pc.ports = out
}

func (pc *PortController) KillPort(p *Port) error {
	if err := KillProcess(p.Pid); err != nil {
		return err
	}

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
