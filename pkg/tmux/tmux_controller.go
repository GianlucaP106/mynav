package tmux

import (
	"mynav/pkg/system"
	"os"
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

	return tmc
}

func (tc *TmuxController) RenameTmuxSession(s *TmuxSession, newName string) error {
	if err := tc.TmuxRepository.RenameSession(s, newName); err != nil {
		return err
	}

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
		defer tc.SyncPorts()
	}

	if err := tc.TmuxRepository.DeleteSession(s); err != nil {
		return err
	}

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
	for _, s := range tc.TmuxRepository.GetSessionContainer() {
		if err := tc.DeleteTmuxSession(s); err != nil {
			return err
		}
	}
	return nil
}

func (tc *TmuxController) GetTmuxStats() (sessionCount int, windowCount int) {
	sessionCount = 0
	windowCount = 0

	for _, s := range tc.TmuxRepository.GetSessionContainer() {
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

func (tc *TmuxController) SyncPorts() {
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
}

func GetNewTmuxSessionCmd(session string, path string) []string {
	return []string{"tmux", "new", "-s", session, "-c", path}
}

func GetAttachTmuxSessionCmd(session string) []string {
	return []string{"tmux", "a", "-t", session}
}

func GetKillTmuxSessionCmd(sessionName string) []string {
	return []string{"tmux", "kill-session", "-t", sessionName}
}

func IsTmuxSession() bool {
	return os.Getenv("TMUX") != ""
}
