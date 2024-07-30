package tmux

import (
	"mynav/pkg/constants"
	"mynav/pkg/events"
	"mynav/pkg/system"
	"mynav/pkg/tasks"
	"sync"
)

type TmuxController struct {
	TmuxRepository   *TmuxRepository
	TmuxCommunicator *TmuxCommunicator
	PortController   *system.PortController
}

func NewTmuxController(pc *system.PortController) *TmuxController {
	tc := NewTmuxCommunicator()
	tr := NewTmuxRepository(tc)

	tmc := &TmuxController{
		TmuxRepository:   tr,
		TmuxCommunicator: tc,
		PortController:   pc,
	}

	// TODO: improve this
	events.AddEventListener(constants.PortSyncNeededEventName, func(_ string) {
		tmc.syncPorts()
	})

	return tmc
}

func (t *TmuxController) GetTmuxSessionChildProcesses(ts *TmuxSession) []*TmuxPaneProcess {
	panes := t.GetTmuxPanesBySession(ts)
	out := make([]*TmuxPaneProcess, 0)
	for _, p := range panes {
		processes, err := system.GetChildProcesses(p.Pid)
		if err != nil {
			continue
		}

		for _, process := range processes {
			paneProcess := &TmuxPaneProcess{
				Process: process,
				Pane:    p,
			}

			out = append(out, paneProcess)
		}
	}

	return out
}

func (tc *TmuxController) RenameTmuxSession(s *TmuxSession, newName string) error {
	if err := tc.TmuxRepository.RenameSession(s, newName); err != nil {
		return err
	}

	events.Emit(constants.TmuxSessionChangeEventName)
	return nil
}

func (tc *TmuxController) CreateAndAttachTmuxSession(session string, path string) error {
	tc.TmuxCommunicator.CreateAndAttachTmuxSession(session, path)
	tc.TmuxRepository.LoadSessions()
	tc.syncPorts()

	events.Emit(constants.TmuxSessionChangeEventName)
	return nil
}

func (tc *TmuxController) AttachTmuxSession(s *TmuxSession) error {
	tc.TmuxCommunicator.AttachTmuxSession(s.Name)
	tc.TmuxRepository.LoadSessions()
	tc.syncPorts()

	events.Emit(constants.TmuxSessionChangeEventName)
	return nil
}

func (tc *TmuxController) GetTmuxSessionCount() int {
	return tc.TmuxRepository.SessionCount()
}

func (tc *TmuxController) GetTmuxSessions() TmuxSessions {
	return tc.TmuxRepository.GetSessions()
}

func (tc *TmuxController) GetTmuxSessionByName(name string) *TmuxSession {
	return tc.TmuxRepository.GetSessionByName(name)
}

func (tc *TmuxController) DeleteTmuxSession(s *TmuxSession) error {
	refreshPorts := len(s.Ports.ToList()) > 0
	if refreshPorts {
		defer tc.syncPorts()
	}

	if err := tc.TmuxRepository.DeleteSession(s); err != nil {
		return err
	}

	events.Emit(constants.TmuxSessionChangeEventName)
	return nil
}

func (tc *TmuxController) GetTmuxSessionPidMap() map[int]*TmuxSession {
	out := map[int]*TmuxSession{}

	for _, session := range tc.GetTmuxSessions() {
		panes := tc.GetTmuxPanesBySession(session)
		for _, pane := range panes {
			out[pane.Pid] = session
		}
	}

	return out
}

func (tc *TmuxController) GetTmuxPanesBySession(ts *TmuxSession) []*TmuxPane {
	return tc.TmuxCommunicator.GetSessionPanes(ts)
}

func (tc *TmuxController) DeleteAllTmuxSessions() error {
	for _, s := range tc.TmuxRepository.GetSessions() {
		if err := tc.TmuxRepository.DeleteSession(s); err != nil {
			return err
		}
	}

	events.Emit(constants.TmuxSessionChangeEventName)
	return nil
}

func (tc *TmuxController) GetTmuxStats() (sessionCount int, windowCount int) {
	sessionCount = 0
	windowCount = 0

	for _, s := range tc.TmuxRepository.GetSessions() {
		sessionCount++
		windowCount += s.NumWindows
	}

	return
}

func (tc *TmuxController) GetTmuxSessionByPort(port *system.Port) *TmuxSession {
	for _, session := range tc.GetTmuxSessions() {
		for _, p := range session.Ports {
			if p.Pid == port.Pid {
				return session
			}
		}
	}

	return nil
}

func (tc *TmuxController) syncPorts() {
	tasks.QueueTask(func() {
		tmap := tc.GetTmuxSessionPidMap()

		tc.PortController.InitPorts()
		ports := tc.PortController.GetPorts()

		var wg sync.WaitGroup

		for _, session := range tmap {
			session.Ports = make(system.Ports)
		}

		for _, port := range ports {
			prt := port
			wg.Add(1)
			go func() {
				defer wg.Done()
				for pid, ts := range tmap {
					if system.IsProcessChildOf(prt.Pid, pid) {
						ts.Ports.AddPort(prt)
					}
				}
			}()
		}

		wg.Wait()

		events.Emit(constants.PortChangeEventName)
	})
}
