package api

import (
	"mynav/pkg/utils"
	"sync"
)

type PortController struct {
	TmuxController *TmuxController
	ports          Ports
}

func NewPortController(tc *TmuxController) *PortController {
	pc := &PortController{
		TmuxController: tc,
		ports:          make(Ports),
	}

	return pc
}

func (pc *PortController) InitPorts() {
	allActivePorts, err := utils.GetRunningPorts()
	if err != nil {
		return
	}

	sessionPidMap := pc.TmuxController.GetTmuxSessionPidMap()
	out := map[int]*Port{}

	for _, ap := range allActivePorts {
		var session *TmuxSession
		for pid, ts := range sessionPidMap {
			if utils.IsProcessChildOf(ap.Pid, pid) {
				session = ts
				break
			}
		}

		pi := utils.GetProcessInfo(ap.Pid)
		port := NewPort(ap.Port, ap.Pid, session, pi.Exe)
		out[port.Port] = port
	}

	pc.ports = out
}

func (pc *PortController) InitPortsAsync() {
	allActivePorts, err := utils.GetRunningPorts()
	if err != nil {
		return
	}

	sessionPidMap := pc.TmuxController.GetTmuxSessionPidMap()

	resultChan := make(chan *Port, len(allActivePorts))
	var wg sync.WaitGroup

	for _, ap := range allActivePorts {
		prt := ap
		wg.Add(1)
		go func() {
			defer wg.Done()
			var session *TmuxSession
			for pid, ts := range sessionPidMap {
				if utils.IsProcessChildOf(prt.Pid, pid) {
					session = ts
					break
				}
			}

			pi := utils.GetProcessInfo(prt.Pid)
			port := NewPort(prt.Port, prt.Pid, session, pi.Exe)
			resultChan <- port
		}()
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	out := map[int]*Port{}
	for port := range resultChan {
		out[port.Port] = port
	}

	pc.ports = out
}

func (pc *PortController) GetPorts() Ports {
	if len(pc.ports) == 0 {
		pc.InitPortsAsync()
		// pc.InitPorts()
	}
	return pc.ports
}
