package core

import (
	"log"
	"mynav/pkg/constants"
	"mynav/pkg/events"
	"mynav/pkg/system"
	"os"

	"github.com/GianlucaP106/gotmux/gotmux"
	"github.com/shirou/gopsutil/process"
)

type TmuxController struct {
	PortController *system.PortController
	Tmux           *gotmux.Tmux
}

func NewTmuxController(pc *system.PortController) *TmuxController {
	t, err := gotmux.DefaultTmux()
	if err != nil {
		log.Panicln(err)
	}

	tmc := &TmuxController{
		PortController: pc,
		Tmux:           t,
	}

	return tmc
}

func (tc *TmuxController) RenameTmuxSession(s *gotmux.Session, newName string) error {
	err := s.Rename(newName)
	if err != nil {
		return err
	}

	events.Emit(constants.TmuxSessionChangeEventName)
	return nil
}

func (tc *TmuxController) CreateAndAttachTmuxSession(session string, path string) error {
	s, err := tc.Tmux.NewSession(&gotmux.SessionOptions{
		Name:           session,
		StartDirectory: path,
	})
	if err != nil {
		return err
	}

	err = s.Attach()
	if err != nil {
		return err
	}

	events.Emit(constants.PortChangeEventName)
	events.Emit(constants.TmuxSessionChangeEventName)
	return nil
}

func (tc *TmuxController) AttachTmuxSession(s *gotmux.Session) error {
	err := s.Attach()
	if err != nil {
		return err
	}

	events.Emit(constants.PortChangeEventName)
	events.Emit(constants.TmuxSessionChangeEventName)
	return nil
}

func (tc *TmuxController) GetTmuxSessionCount() int {
	s, _ := tc.Tmux.ListSessions()
	return len(s)
}

func (tc *TmuxController) GetTmuxSessions() []*gotmux.Session {
	sessions, err := tc.Tmux.ListSessions()
	if err != nil {
		return []*gotmux.Session{}
	}

	return sessions
}

func (tc *TmuxController) GetTmuxSessionByName(name string) *gotmux.Session {
	session, _ := tc.Tmux.GetSessionByName(name)
	return session
}

func (tc *TmuxController) DeleteTmuxSession(s *gotmux.Session) error {
	// TODO:
	// refreshPorts := len(s.Ports.ToList()) > 0
	// if refreshPorts {
	// 	defer tc.syncPorts()
	// }

	err := s.Kill()
	if err != nil {
		return err
	}

	events.Emit(constants.TmuxSessionChangeEventName)
	return nil
}

func (tc *TmuxController) KillTmuxServer() error {
	err := tc.Tmux.KillServer()
	if err != nil {
		return err
	}

	events.Emit(constants.TmuxSessionChangeEventName)
	return nil
}

func (tc *TmuxController) GetTmuxStats() (sessionCount int, windowCount int) {
	sessionCount = 0
	windowCount = 0

	sessions, err := tc.Tmux.ListSessions()
	if err != nil {
		return
	}

	for _, s := range sessions {
		sessionCount++
		windowCount += s.Windows
	}

	return
}

func (t *TmuxController) GetTmuxSessionChildProcesses(session *gotmux.Session) []*process.Process {
	windows, err := session.ListWindows()
	out := make([]*process.Process, 0)
	if err != nil {
		return out
	}

	for _, window := range windows {
		panes, err := window.ListPanes()
		if err != nil {
			continue
		}

		for _, p := range panes {
			proc, err := process.NewProcess(p.Pid)
			if err != nil {
				continue
			}

			out = append(out, proc)
		}
	}

	return out
}

func (tc *TmuxController) GetTmuxSessionByPort(port *system.Port) *gotmux.Session {
	return tc.GetTmuxSessionByChildPid(port.Pid)
}

func (t *TmuxController) GetTmuxSessionByChildPid(pid int) *gotmux.Session {
	for _, session := range t.GetTmuxSessions() {
		windows, err := session.ListWindows()
		if err != nil {
			continue
		}

		for _, w := range windows {
			panes, err := w.ListPanes()
			if err != nil {
				continue
			}

			for _, p := range panes {
				if p.Pid == int32(pid) {
					return session
				}
			}
		}

	}

	return nil
}

func IsTmuxSession() bool {
	return os.Getenv("TMUX") != ""
}

// TODO:
// func (tc *TmuxController) syncPorts() {
// 	sessions, err := tc.Tmux.ListSessions()
// 	if err != nil {
// 		return
// 	}
//
// 	tasks.QueueTask(func() {
// 		tmap := tc.GetTmuxSessionPidMap()
//
// 		tc.PortController.InitPorts()
// 		ports := tc.PortController.GetPorts()
//
// 		var wg sync.WaitGroup
// 		for _, port := range ports {
// 			prt := port
// 			wg.Add(1)
// 			go func() {
// 				defer wg.Done()
// 				for pid, ts := range tmap {
// 					if system.IsProcessChildOf(prt.Pid, pid) {
// 						ts.Ports.AddPort(prt)
// 					}
// 				}
// 			}()
// 		}
//
// 		wg.Wait()
//
// 		events.Emit(constants.PortChangeEventName)
// 	})
// }
